package divisions

import (
	"fmt"
	"net/http"

	"github.com/starfederation/datastar-go/datastar"

	"github.com/cszczepaniak/cribbly/internal/persistence/divisions"
	"github.com/cszczepaniak/cribbly/internal/persistence/teams"
	"github.com/cszczepaniak/cribbly/internal/ui/pages/admin/admincomponents"
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
	err = sse.PatchElementTempl(admincomponents.EditNameInput(division.Name))
	if err != nil {
		return err
	}
	err = sse.PatchElementTempl(admincomponents.EditSaveButton[divisions.Division](id))
	if err != nil {
		return err
	}
	return sse.PatchElementTempl(admincomponents.ItemsListing[teams.Team](id, availableTeams, inThisDivision))
}

func (h DivisionsHandler) Create(w http.ResponseWriter, r *http.Request) error {
	_, err := h.DivisionRepo.Create(r.Context())
	if err != nil {
		return err
	}

	teams, err := h.DivisionRepo.GetAll(r.Context())
	if err != nil {
		return err
	}

	sse := datastar.NewSSE(w, r)
	return sse.PatchElementTempl(admincomponents.TeamOrDivisionTable(teams))
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

	a := admincomponents.GetAssignmentForEdit(r)
	if a != (admincomponents.Assignment{}) {
		var err error
		if a.Assign != "" {
			err = h.TeamRepo.AssignToDivision(r.Context(), a.Assign, divisionID)
		} else {
			err = h.TeamRepo.UnassignFromDivision(r.Context(), a.Unassign)
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
		return sse.PatchElementTempl(admincomponents.ItemsListing[divisions.Division](divisionID, available, inThisDivision))
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

		err = sse.PatchElementTempl(admincomponents.TeamOrDivisionName(divisionID, sigs.Name))
		if err != nil {
			return err
		}
		return sse.MarshalAndPatchSignals(signals{Name: ""})
	}

	return nil
}
