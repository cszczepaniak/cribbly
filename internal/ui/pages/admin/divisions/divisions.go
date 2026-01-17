package divisions

import (
	"fmt"
	"net/http"

	"github.com/starfederation/datastar-go/datastar"

	"github.com/cszczepaniak/cribbly/internal/persistence/divisions"
	"github.com/cszczepaniak/cribbly/internal/persistence/teams"
)

type DivisionsHandler struct {
	TeamRepo     teams.Repository
	DivisionRepo divisions.Repository
}

func (h DivisionsHandler) Index(w http.ResponseWriter, r *http.Request) error {
	divisions, err := h.DivisionRepo.GetAll(r.Context())
	if err != nil {
		return err
	}

	return index(divisions).Render(r.Context(), w)
}

func (h DivisionsHandler) Edit(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	division, err := h.DivisionRepo.Get(r.Context(), id)
	if err != nil {
		return err
	}

	availableTeams, err := h.TeamRepo.GetWithoutDivision(r.Context())
	if err != nil {
		return err
	}

	inThisDivision, err := h.TeamRepo.GetForDivision(r.Context(), id)
	if err != nil {
		return err
	}

	sse := datastar.NewSSE(w, r)
	err = sse.PatchElementTempl(editNameInput(division.Name))
	if err != nil {
		return err
	}
	err = sse.PatchElementTempl(editSaveButton(id))
	if err != nil {
		return err
	}
	return sse.PatchElementTempl(itemsListing(id, availableTeams, inThisDivision))
}

func (h DivisionsHandler) Create(w http.ResponseWriter, r *http.Request) error {
	_, err := h.DivisionRepo.Create(r.Context())
	if err != nil {
		return err
	}

	divisions, err := h.DivisionRepo.GetAll(r.Context())
	if err != nil {
		return err
	}

	sse := datastar.NewSSE(w, r)
	return sse.PatchElementTempl(divisionTable(divisions))
}

func (h DivisionsHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	err := h.DivisionRepo.Delete(r.Context(), id)
	if err != nil {
		return err
	}

	sse := datastar.NewSSE(w, r)
	return sse.RemoveElementByID(fmt.Sprintf("table-row-%s", id))
}

type signals struct {
	Name string `json:"name"`
}

func (h DivisionsHandler) Save(w http.ResponseWriter, r *http.Request) error {
	divisionID := r.PathValue("id")

	assign := r.URL.Query().Get("assign")
	unassign := r.URL.Query().Get("unassign")
	if assign != "" || unassign != "" {
		var err error
		if assign != "" {
			err = h.TeamRepo.AssignToDivision(r.Context(), assign, divisionID)
		} else {
			err = h.TeamRepo.UnassignFromDivision(r.Context(), unassign)
		}
		if err != nil {
			return err
		}

		inThisDivision, err := h.TeamRepo.GetForDivision(r.Context(), divisionID)
		if err != nil {
			return err
		}

		available, err := h.TeamRepo.GetWithoutDivision(r.Context())
		if err != nil {
			return err
		}

		// If we're unassigning a team, we'll keep the modal open (by not redirecting).
		sse := datastar.NewSSE(w, r)
		return sse.PatchElementTempl(itemsListing(divisionID, available, inThisDivision))
	}

	var sigs signals
	err := datastar.ReadSignals(r, &sigs)
	if err != nil {
		return err
	}

	if sigs.Name != "" {
		err := h.DivisionRepo.Rename(r.Context(), divisionID, sigs.Name)
		if err != nil {
			return err
		}

		sse := datastar.NewSSE(w, r)

		err = sse.PatchElementTempl(divisionName(divisionID, sigs.Name))
		if err != nil {
			return err
		}
		return sse.MarshalAndPatchSignals(signals{Name: ""})
	}

	return nil
}
