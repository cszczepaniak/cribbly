package teams

import (
	"context"
	"net/http"

	"github.com/a-h/templ"

	"github.com/cszczepaniak/cribbly/internal/persistence/players"
	"github.com/cszczepaniak/cribbly/internal/persistence/teams"
	"github.com/cszczepaniak/cribbly/internal/ui/hx"
)

type teamsData struct {
	teams []teams.Team
}

type teamData struct {
	team    teams.Team
	players []players.Player
}

type editTeamData struct {
	teamData
	availablePlayers []players.Player
}

type TeamsHandler struct {
	PlayerService players.Service
	TeamService   teams.Service
}

func (h TeamsHandler) Index(_ http.ResponseWriter, r *http.Request) (templ.Component, error) {
	teams, err := h.TeamService.GetAll(r.Context())
	if err != nil {
		return nil, err
	}

	if id := r.URL.Query().Get("edit"); id != "" {
		editData, err := h.getEditData(r.Context(), id)
		if err != nil {
			return nil, err
		}

		return indexEditing(teams, editData), nil
	}

	return index(teams), nil
}

func (h TeamsHandler) Create(_ http.ResponseWriter, r *http.Request) (templ.Component, error) {
	_, err := h.TeamService.Create(r.Context())
	if err != nil {
		return nil, err
	}

	teams, err := h.TeamService.GetAll(r.Context())
	if err != nil {
		return nil, err
	}

	return fullIndexPage(teams, editTeamData{}, false), nil
}

func (h TeamsHandler) Delete(_ http.ResponseWriter, r *http.Request) (templ.Component, error) {
	id := r.PathValue("id")
	err := h.TeamService.Delete(r.Context(), id)
	if err != nil {
		return nil, err
	}

	teams, err := h.TeamService.GetAll(r.Context())
	if err != nil {
		return nil, err
	}

	return fullIndexPage(teams, editTeamData{}, false), nil
}

func (h TeamsHandler) Save(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	teamID := r.PathValue("id")

	playerToDelete := r.FormValue("unassignPlayer")
	if playerToDelete != "" {
		err := h.PlayerService.UnassignFromTeam(r.Context(), playerToDelete)
		if err != nil {
			return nil, err
		}

		// If we're unassigning a player, we'll keep the modal open (by not redirecting).
		return h.renderIndexWithEditForm(r.Context(), teamID)
	}

	playerToAssign := r.FormValue("assignPlayer")
	if playerToAssign != "" {
		// TODO: validate that we're not adding too many players to the team.
		err := h.PlayerService.AssignToTeam(r.Context(), playerToAssign, teamID)
		if err != nil {
			return nil, err
		}

		// If we're assigning a player, we'll keep the modal open (by not redirecting).
		return h.renderIndexWithEditForm(r.Context(), teamID)
	}

	newName := r.FormValue("name")
	if newName != "" {
		err := h.TeamService.Rename(r.Context(), teamID, newName)
		if err != nil {
			return nil, err
		}
	}

	hx.RedirectTo(w, "/admin/teams")

	// Since we're redirecting, the index will get loaded and we don't need to return a component.
	return nil, nil
}

func (h TeamsHandler) renderIndexWithEditForm(ctx context.Context, teamID string) (templ.Component, error) {
	allTeams, err := h.TeamService.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// TODO: we could also simply find the team from the list of all teams we query below instead of
	// making another query to the DB.
	editData, err := h.getEditData(ctx, teamID)
	if err != nil {
		return nil, err
	}

	return fullIndexPage(allTeams, editData, true), nil
}

func (h TeamsHandler) getEditData(ctx context.Context, teamID string) (editTeamData, error) {
	team, err := h.TeamService.Get(ctx, teamID)
	if err != nil {
		return editTeamData{}, err
	}

	playersOnThisTeam, err := h.PlayerService.GetForTeam(ctx, teamID)
	if err != nil {
		return editTeamData{}, err
	}

	availablePlayers, err := h.PlayerService.GetFreeAgents(ctx)
	if err != nil {
		return editTeamData{}, err
	}

	return editTeamData{
		teamData: teamData{
			team:    team,
			players: playersOnThisTeam,
		},
		availablePlayers: availablePlayers,
	}, nil
}
