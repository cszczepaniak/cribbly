package games

import (
	"net/http"

	"github.com/cszczepaniak/cribbly/internal/notifier"
	"github.com/cszczepaniak/cribbly/internal/persistence/games"
	"github.com/cszczepaniak/cribbly/internal/persistence/teams"
	"github.com/starfederation/datastar-go/datastar"
)

type Handler struct {
	GameService         games.Service
	TeamService         teams.Service
	ScoreUpdateNotifier *notifier.Notifier
}

type getGameProps struct {
	gameID       string
	complete     bool
	scores       [2]games.Score
	team1, team2 teams.Team
}

func (h Handler) GetGame(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")

	g, err := h.GameService.Get(r.Context(), id)
	if err != nil {
		return err
	}

	team1, err := h.TeamService.Get(r.Context(), g[0].TeamID)
	if err != nil {
		return err
	}

	team2, err := h.TeamService.Get(r.Context(), g[1].TeamID)
	if err != nil {
		return err
	}

	return Game(getGameProps{
		gameID:   id,
		complete: g[0].Score != 0 && g[1].Score != 0,
		scores:   g,
		team1:    team1,
		team2:    team2,
	}).Render(r.Context(), w)
}

func (h Handler) UpdateGame(w http.ResponseWriter, r *http.Request) error {
	gameID := r.PathValue("id")
	var signals struct {
		WinningTeamID string `json:"winner"`
		LoserScore    int    `json:"loserScore"`
	}
	err := datastar.ReadSignals(r, &signals)
	if err != nil {
		return err
	}

	if signals.WinningTeamID == "" {
		sse := datastar.NewSSE(w, r)
		err := sse.PatchElementTempl(
			teamInputError("Must select a winning team."),
			datastar.WithSelectorID("team-input"),
			datastar.WithModeAppend(),
		)
		if err != nil {
			return err
		}
	}

	if signals.LoserScore <= 0 || signals.LoserScore > 120 {
		sse := datastar.NewSSE(w, r)
		return sse.PatchElementTempl(scoreInput("Losing score must be between 1 and 120."))
	}

	scores, err := h.GameService.Get(r.Context(), gameID)
	if err != nil {
		return err
	}

	var losingID string
	if signals.WinningTeamID == scores[0].TeamID {
		losingID = scores[1].TeamID
	} else {
		losingID = scores[0].TeamID
	}

	err = h.GameService.UpdateScores(r.Context(), gameID, signals.WinningTeamID, 121, losingID, signals.LoserScore)
	if err != nil {
		return err
	}

	return datastar.NewSSE(w, r).Redirectf("/games/%s", gameID)
}

func (h Handler) StandingsPage(w http.ResponseWriter, r *http.Request) error {
	s, err := h.GameService.GetStandings(r.Context())
	if err != nil {
		return err
	}
	return standings(s).Render(r.Context(), w)
}

func (h Handler) StreamStandings(w http.ResponseWriter, r *http.Request) error {
	sse := datastar.NewSSE(w, r)
	notify, cancel := h.ScoreUpdateNotifier.Subscribe()
	defer cancel()

	for {
		select {
		case <-r.Context().Done():
			return nil
		case <-notify:
			s, err := h.GameService.GetStandings(r.Context())
			if err != nil {
				return err
			}

			err = sse.PatchElementTempl(standingsTable(s), datastar.WithViewTransitions())
			if err != nil {
				return err
			}
		}
	}
}
