package teams

import (
	"cmp"
	"net/http"
	"slices"

	"github.com/a-h/templ"
	"github.com/cszczepaniak/cribbly/internal/persistence/players"
	"github.com/cszczepaniak/cribbly/internal/ui/hx"
	"github.com/google/uuid"
)

type teamsData struct {
	teams []team
}

type editTeamData struct {
	team             team
	availablePlayers []players.Player
}

type team struct {
	id      string
	name    string
	players []players.Player
}

var teams = []team{{
	id:   uuid.NewString(),
	name: "super mario bros",
	players: []players.Player{{
		ID:   uuid.NewString(),
		Name: "mario",
	}, {
		ID:   uuid.NewString(),
		Name: "luigi",
	}},
}}

func getTeam(id string) (team, int) {
	teamIdx := slices.IndexFunc(teams, func(t team) bool {
		return t.id == id
	})

	if teamIdx == -1 {
		return team{}, -1
	}

	return teams[teamIdx], teamIdx
}

type TeamsHandler struct {
	PlayerService players.Service
}

func (h TeamsHandler) Index(_ http.ResponseWriter, r *http.Request) (templ.Component, error) {
	var editing team
	var isEditing bool
	if id := r.URL.Query().Get("edit"); id != "" {
		editing, _ = getTeam(id)
		isEditing = true
	}

	// If the team is being edited, we'll need to load the list of available players as well.
	if isEditing {
		availablePlayers, err := h.PlayerService.GetFreeAgents(r.Context())
		if err != nil {
			return nil, err
		}

		return indexEditing(
			teams,
			editTeamData{
				team:             editing,
				availablePlayers: availablePlayers,
			},
		), nil
	}

	return index(teams), nil
}

func (h TeamsHandler) Create(_ http.ResponseWriter, r *http.Request) (templ.Component, error) {
	team := team{
		id:   uuid.NewString(),
		name: cmp.Or(r.FormValue("name"), "Unnamed Team"),
	}
	teams = append(teams, team)
	return index(teams), nil
}

func (h TeamsHandler) Save(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	id := r.PathValue("id")
	team, teamIdx := getTeam(id)

	playerToDelete := r.FormValue("delete")
	if playerToDelete != "" {
		idx := slices.IndexFunc(team.players, func(p players.Player) bool {
			return p.ID == playerToDelete
		})
		if idx != -1 {
			team.players = slices.Delete(team.players, idx, idx+1)
		}

		teams[teamIdx] = team

		availablePlayers, err := h.PlayerService.GetFreeAgents(r.Context())
		if err != nil {
			return nil, err
		}

		// If we're deleting a player, we'll keep the modal open (by not redirecting).
		return indexEditing(
			teams,
			editTeamData{
				team:             team,
				availablePlayers: availablePlayers,
			},
		), nil
	}

	newName := r.FormValue("name")
	if newName != "" {
		team.name = newName
	}

	teams[teamIdx] = team

	hx.RedirectTo(w, "/admin/teams")

	// Since we're redirecting, the index will get loaded and we don't need to return a component.
	return nil, nil
}
