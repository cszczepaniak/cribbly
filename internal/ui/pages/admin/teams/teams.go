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
	teams              []team
	currentlyEditingID string
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

func Index(_ http.ResponseWriter, r *http.Request) (templ.Component, error) {
	data := teamsData{
		teams:              teams,
		currentlyEditingID: r.URL.Query().Get("edit"),
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
	teamIdx := slices.IndexFunc(teams, func(t team) bool {
		return t.id == id
	})
	team := teams[teamIdx]

	newName := r.FormValue("name")
	if newName != "" {
		team.name = newName
	}

	teams[teamIdx] = team

	hx.RedirectTo(w, "/admin/teams")

	// Since we're redirecting, the index will get loaded and we don't need to return a component.
	return nil, nil
}
