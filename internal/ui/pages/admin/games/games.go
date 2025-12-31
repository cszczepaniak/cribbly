package games

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"net/http"
	"strconv"

	"github.com/starfederation/datastar-go/datastar"

	"github.com/cszczepaniak/cribbly/internal/persistence/divisions"
	"github.com/cszczepaniak/cribbly/internal/persistence/games"
	"github.com/cszczepaniak/cribbly/internal/persistence/teams"
)

type Handler struct {
	DivisionService divisions.Service
	TeamService     teams.Service
	GameService     games.Service
}

func (h Handler) Index(w http.ResponseWriter, r *http.Request) error {
	gs, err := h.getAllGames(r.Context())
	if err != nil {
		return err
	}

	return index(gs).Render(r.Context(), w)
}

func (h Handler) DeleteAll(w http.ResponseWriter, r *http.Request) error {
	err := h.GameService.DeleteAll(r.Context())
	if err != nil {
		return err
	}

	return index(nil).Render(r.Context(), w)
}

func (h Handler) Generate(w http.ResponseWriter, r *http.Request) error {
	err := h.generatePrelimGames(r.Context())
	if err != nil {
		return err
	}

	gs, err := h.getAllGames(r.Context())
	if err != nil {
		return err
	}

	sse := datastar.NewSSE(w, r)
	return sse.PatchElementTempl(gameList(gs))
}

type signals struct {
	Score string `json:"score"`
}

func (h Handler) Edit(w http.ResponseWriter, r *http.Request) error {
	gameID := r.URL.Query().Get("gameID")
	teamID := r.URL.Query().Get("teamID")

	score, err := h.GameService.GetScore(r.Context(), gameID, teamID)
	if err != nil {
		return err
	}

	sse := datastar.NewSSE(w, r)
	err = sse.MarshalAndPatchSignals(signals{Score: fmt.Sprint(score)})
	if err != nil {
		return err
	}

	return sse.PatchElementTempl(
		editScore(gameID, teamID, score),
		datastar.WithSelectorID(scoreCellID(gameID, teamID)),
		datastar.WithModeInner(),
	)
}

func (h Handler) Save(w http.ResponseWriter, r *http.Request) error {
	gameID := r.URL.Query().Get("gameID")
	teamID := r.URL.Query().Get("teamID")

	var signals signals
	err := datastar.ReadSignals(r, &signals)
	if err != nil {
		return err
	}

	score, err := strconv.Atoi(signals.Score)
	if err != nil {
		return err
	}

	err = h.GameService.UpdateScore(r.Context(), gameID, teamID, score)
	if err != nil {
		return err
	}

	sse := datastar.NewSSE(w, r)
	return sse.PatchElementTempl(
		displayScore(gameID, teamID, score),
		datastar.WithSelectorID(scoreCellID(gameID, teamID)),
		datastar.WithModeInner(),
	)
}

func scoreCellID(gameID, teamID string) string {
	return fmt.Sprintf("score-%s-%s", gameID, teamID)
}

type game struct {
	gameID     string
	division   divisions.Division
	team1      teams.Team
	team1Score int
	team2      teams.Team
	team2Score int
}

func (h Handler) generatePrelimGames(ctx context.Context) error {
	allTeams, err := h.TeamService.GetAll(ctx)
	if err != nil {
		return err
	}

	teamsByDivision := make(map[string][]teams.Team)
	for _, team := range allTeams {
		if team.DivisionID == "" {
			return errors.New("all teams must be in divisions to generate games")
		}
		teamsByDivision[team.DivisionID] = append(teamsByDivision[team.DivisionID], team)
	}

	// TODO: don't insert one-by-one!
	for _, teams := range teamsByDivision {
		pairs, err := generateMatchups(teams)
		if err != nil {
			return err
		}

		for _, pair := range pairs {
			_, err := h.GameService.Create(ctx, pair[0].ID, pair[1].ID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func cycle[T any](in []T) iter.Seq[T] {
	return func(yield func(T) bool) {
		i := 0
		for {
			idx := i % len(in)
			if !yield(in[idx]) {
				return
			}
			i++
		}
	}
}

func generateMatchups(allTeams []teams.Team) ([][2]teams.Team, error) {
	if len(allTeams) != 4 && len(allTeams) != 6 {
		return nil, errors.New("can only generate matchups with 4 or 6 teams")
	}

	s1, done1 := iter.Pull(cycle(allTeams))
	defer done1()
	s2, done2 := iter.Pull(cycle(allTeams))
	defer done2()

	nByTeam := make(map[string]int, len(allTeams))

	var allPairs [][2]teams.Team
	for {
		_, _ = s2()

		// Advance each cycle in lock step to generate a matchup for each team. Initially, t1 vs.
		// t2, t2 vs. t3, etc. (note that these games obviously don't happen in a single round).
		//
		// The next outer loop will yield t1 vs. t3, t2 vs. t4, etc.
		for range len(allTeams) {
			t1, _ := s1()
			t2, _ := s2()

			if nByTeam[t1.ID] == 3 || nByTeam[t2.ID] == 3 {
				// We can't add this game because then at least one team would have more than 3.
				continue
			}

			allPairs = append(allPairs, [2]teams.Team{t1, t2})
			nByTeam[t1.ID]++
			nByTeam[t2.ID]++
		}

		done := true
		for _, count := range nByTeam {
			if count != 3 {
				done = false
				break
			}
		}
		if done {
			return allPairs, nil
		}
	}
}

func (h Handler) getAllGames(ctx context.Context) ([]game, error) {
	// This is clearly a very inefficient way to do this (querying all of everything and then
	// mapping in-memory), but it keeps the database logic simple... let's do it this way until we
	// know we need something better?

	scores, err := h.GameService.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	scoresByGame := make(map[string][]games.Score)
	for _, s := range scores {
		scoresByGame[s.GameID] = append(scoresByGame[s.GameID], s)
	}

	daTeams, err := h.TeamService.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	teamsByID := make(map[string]teams.Team, len(daTeams))
	for _, t := range daTeams {
		teamsByID[t.ID] = t
	}

	daDivisions, err := h.DivisionService.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	divisionsByID := make(map[string]divisions.Division, len(daDivisions))
	for _, d := range daDivisions {
		divisionsByID[d.ID] = d
	}

	var daGames []game
	for gameID, scores := range scoresByGame {
		if len(scores) != 2 {
			return nil, errors.New("dev error: each game should have two scores")
		}
		team1, ok := teamsByID[scores[0].TeamID]
		if !ok {
			return nil, errors.New("team not found")
		}
		team2, ok := teamsByID[scores[1].TeamID]
		if !ok {
			return nil, errors.New("team not found")
		}

		if team1.DivisionID != team2.DivisionID {
			return nil, errors.New("should not have inter-division games")
		}

		division, ok := divisionsByID[team1.DivisionID]
		if !ok {
			return nil, errors.New("division not found")
		}

		g := game{
			gameID:     gameID,
			division:   division,
			team1:      team1,
			team1Score: scores[0].Score,
			team2:      team2,
			team2Score: scores[1].Score,
		}
		daGames = append(daGames, g)
	}

	return daGames, nil
}
