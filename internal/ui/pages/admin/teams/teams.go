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
	teams            []team
	currentlyEditing team
	isEditing        bool
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

func Index(_ http.ResponseWriter, r *http.Request) (templ.Component, error) {
	var editing team
	var isEditing bool
	if id := r.URL.Query().Get("edit"); id != "" {
		editing, _ = getTeam(id)
		isEditing = true
	}

	data := teamsData{
		teams:            teams,
		currentlyEditing: editing,
		isEditing:        isEditing,
	}
	return index(data), nil
}

func Create(_ http.ResponseWriter, r *http.Request) (templ.Component, error) {
	team := team{
		id:   uuid.NewString(),
		name: cmp.Or(r.FormValue("name"), "Unnamed Team"),
	}
	teams = append(teams, team)
	return index(teamsData{teams: teams}), nil
}

func Save(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
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

		// If we're deleting a player, we'll keep the modal open (by not redirecting).
		return index(teamsData{
			teams:            teams,
			currentlyEditing: team,
			isEditing:        true,
		}), nil
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
