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
	TeamService     teams.Service
	DivisionService divisions.Service
}

func (h DivisionsHandler) Index(w http.ResponseWriter, r *http.Request) error {
	divisions, err := h.DivisionService.GetAll(r.Context())
	if err != nil {
		return err
	}

	return index(divisions).Render(r.Context(), w)
}

func (h DivisionsHandler) Edit(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	division, err := h.DivisionService.Get(r.Context(), id)
	if err != nil {
		return err
	}

	availableTeams, err := h.TeamService.GetWithoutDivision(r.Context())
	if err != nil {
		return err
	}

	inThisDivision, err := h.TeamService.GetForDivision(r.Context(), id)
	if err != nil {
		return err
	}

	sse := datastar.NewSSE(w, r)
	return sse.PatchElementTempl(
		admincomponents.EditTeamOrDivisionModal(division, inThisDivision, availableTeams),
		datastar.WithSelectorID("divisions"),
		datastar.WithModeAppend(),
	)
}

func (h DivisionsHandler) CancelEdit(w http.ResponseWriter, r *http.Request) error {
	sse := datastar.NewSSE(w, r)
	return resetEdit(sse)
}

func (h DivisionsHandler) Create(w http.ResponseWriter, r *http.Request) error {
	_, err := h.DivisionService.Create(r.Context())
	if err != nil {
		return err
	}

	teams, err := h.DivisionService.GetAll(r.Context())
	if err != nil {
		return err
	}

	sse := datastar.NewSSE(w, r)
	return sse.PatchElementTempl(admincomponents.TeamOrDivisionTable(teams))
}

func (h DivisionsHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	err := h.DivisionService.Delete(r.Context(), id)
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
			err = h.TeamService.AssignToDivision(r.Context(), a.Assign, divisionID)
		} else {
			err = h.TeamService.UnassignFromDivision(r.Context(), a.Unassign)
		}
		if err != nil {
			return err
		}

		inThisDivision, err := h.TeamService.GetForDivision(r.Context(), divisionID)
		if err != nil {
			return err
		}

		available, err := h.TeamService.GetWithoutDivision(r.Context())
		if err != nil {
			return err
		}

		// If we're unassigning a team, we'll keep the modal open (by not redirecting).
		sse := datastar.NewSSE(w, r)
		return sse.PatchElementTempl(admincomponents.ItemsListing[divisions.Division](divisionID, available, inThisDivision))
	}

	var signals signals
	err := datastar.ReadSignals(r, &signals)
	if err != nil {
		return err
	}

	if signals.Name != "" {
		err := h.DivisionService.Rename(r.Context(), divisionID, signals.Name)
		if err != nil {
			return err
		}

		sse := datastar.NewSSE(w, r)

		err = sse.PatchElementTempl(admincomponents.TeamOrDivisionName(divisionID, signals.Name))
		if err != nil {
			return err
		}
		return resetEdit(sse)
	}

	return nil
}

func resetEdit(sse *datastar.ServerSentEventGenerator) error {
	err := sse.MarshalAndPatchSignals(signals{Name: ""})
	if err != nil {
		return err
	}
	return sse.RemoveElementByID(admincomponents.EditModalID)
}
