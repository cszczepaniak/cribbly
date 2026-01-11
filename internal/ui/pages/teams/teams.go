package teams

import (
	"net/http"

	"github.com/cszczepaniak/cribbly/internal/persistence/games"
	"github.com/cszczepaniak/cribbly/internal/persistence/teams"
)

type Handler struct {
	GameService games.Service
	TeamService teams.Service
}

type team struct {
	name  string
	score int
}

type game struct {
	id    string
	team1 team
	team2 team
}

func (g game) rankedTeams() [2]team {
	if g.team1.score > g.team2.score {
		return [...]team{g.team1, g.team2}
	} else if g.team1.score < g.team2.score {
		return [...]team{g.team2, g.team1}
	} else {
		panic("ties can't happen")
	}
}

func (g game) won(team string) bool {
	return (g.team1.score > g.team2.score && g.team1.name == team) || (g.team2.score > g.team1.score && g.team2.name == team)
}

func (g game) complete() bool {
	return g.team1.score != 0 || g.team2.score != 0
}

func (h Handler) GetGames(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")

	t, err := h.TeamService.Get(r.Context(), id)
	if err != nil {
		return err
	}

	gs, err := h.GameService.GetForTeam(r.Context(), id)
	if err != nil {
		return err
	}

	var games []game
	for gameID, scores := range gs {
		team1, err := h.TeamService.Get(r.Context(), scores[0].TeamID)
		if err != nil {
			return err
		}

		team2, err := h.TeamService.Get(r.Context(), scores[1].TeamID)
		if err != nil {
			return err
		}

		games = append(games, game{
			id: gameID,
			team1: team{
				name:  team1.Name,
				score: scores[0].Score,
			},
			team2: team{
				name:  team2.Name,
				score: scores[1].Score,
			},
		})
	}

	return Games(t, games).Render(r.Context(), w)
}
