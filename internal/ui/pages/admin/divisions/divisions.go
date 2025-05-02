package divisions

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
	"github.com/cszczepaniak/cribbly/internal/persistence/divisions"
	"github.com/cszczepaniak/cribbly/internal/persistence/teams"
	"github.com/cszczepaniak/cribbly/internal/ui/hx"
)

type divisionsData struct {
	divisions []divisions.Division
}

type divisionData struct {
	division divisions.Division
	teams    []teams.Team
}

type editDivisionData struct {
	divisionData
	availableTeams []teams.Team
}

type DivisionsHandler struct {
	TeamService     teams.Service
	DivisionService divisions.Service
}

func (h DivisionsHandler) Index(_ http.ResponseWriter, r *http.Request) (templ.Component, error) {
	if id := r.URL.Query().Get("edit"); id != "" {
		return h.renderIndexWithEditForm(r.Context(), id)
	}

	return h.renderIndex(r.Context())
}

func (h DivisionsHandler) Create(_ http.ResponseWriter, r *http.Request) (templ.Component, error) {
	_, err := h.DivisionService.Create(r.Context())
	if err != nil {
		return nil, err
	}

	return h.renderIndex(r.Context())
}

func (h DivisionsHandler) Save(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	divisionID := r.PathValue("id")

	teamToDelete := r.FormValue("unassignTeam")
	if teamToDelete != "" {
		err := h.TeamService.UnassignFromDivision(r.Context(), teamToDelete)
		if err != nil {
			return nil, err
		}

		// If we're unassigning a team, we'll keep the modal open (by not redirecting).
		return h.renderIndexWithEditForm(r.Context(), divisionID)
	}

	teamToAssign := r.FormValue("assignTeam")
	if teamToAssign != "" {
		// TODO: validate that we're not adding too many teams to the division.
		err := h.TeamService.AssignToDivision(r.Context(), teamToAssign, divisionID)
		if err != nil {
			return nil, err
		}

		// If we're assigning a player, we'll keep the modal open (by not redirecting).
		return h.renderIndexWithEditForm(r.Context(), divisionID)
	}

	newName := r.FormValue("name")
	if newName != "" {
		err := h.DivisionService.Rename(r.Context(), divisionID, newName)
		if err != nil {
			return nil, err
		}
	}

	hx.RedirectTo(w, "/admin/divisions")

	// Since we're redirecting, the index will get loaded and we don't need to return a component.
	return nil, nil
}

func (h DivisionsHandler) renderIndex(ctx context.Context) (templ.Component, error) {
	divisions, err := h.DivisionService.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return index(divisions), nil
}

func (h DivisionsHandler) renderIndexWithEditForm(ctx context.Context, divisionID string) (templ.Component, error) {
	// TODO: we could also simply find the division from the list of all divisions we query below instead of
	// making another query to the DB.
	division, err := h.DivisionService.Get(ctx, divisionID)
	if err != nil {
		return nil, err
	}

	allDivisions, err := h.DivisionService.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	teamsInThisDivision, err := h.TeamService.GetForDivision(ctx, divisionID)
	if err != nil {
		return nil, err
	}

	availableTeams, err := h.TeamService.GetWithoutDivision(ctx)
	if err != nil {
		return nil, err
	}

	return indexEditing(allDivisions, editDivisionData{
		divisionData: divisionData{
			division: division,
			teams:    teamsInThisDivision,
		},
		availableTeams: availableTeams,
	}), nil
}
