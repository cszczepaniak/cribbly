package games

import (
	"cmp"
	"context"
	"errors"
	"iter"
	"net/http"
	"slices"
	"strconv"

	"github.com/starfederation/datastar-go/datastar"

	"github.com/cszczepaniak/cribbly/internal/persistence/divisions"
	"github.com/cszczepaniak/cribbly/internal/persistence/games"
	"github.com/cszczepaniak/cribbly/internal/persistence/teams"
	"github.com/cszczepaniak/cribbly/internal/ui/components"
)

type Handler struct {
	DivisionRepo divisions.Repository
	TeamRepo     teams.Repository
	GameRepo     games.Repository
}

func (h Handler) Index(w http.ResponseWriter, r *http.Request) error {
	gs, err := h.getAllGames(r.Context())
	if err != nil {
		return err
	}

	return index(gs).Render(r.Context(), w)
}

func (h Handler) DeleteAll(w http.ResponseWriter, r *http.Request) error {
	err := h.GameRepo.DeleteAll(r.Context())
	if err != nil {
		return err
	}

	return datastar.NewSSE(w, r).Redirect("/admin/games")
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

// gameRowInfo is used to add a signal to each row so that we don't have to re-query for all of this
// information on subsequent requests.
type gameRowInfo struct {
	DivisionName string `json:"divisionName"`
	Team1ID      string `json:"team1ID"`
	Team1Name    string `json:"team1Name"`
	Team2Name    string `json:"team2Name"`
	Team2ID      string `json:"team2ID"`
}

type gamesSignal struct {
	Games []game `json:"games"`
}

func (h Handler) Edit(w http.ResponseWriter, r *http.Request) error {
	gameID := r.URL.Query().Get("gameID")

	var sigs gamesSignal
	err := datastar.ReadSignals(r, &sigs)
	if err != nil {
		return err
	}

	idx := slices.IndexFunc(sigs.Games, func(g game) bool { return g.GameID == gameID })
	game := sigs.Games[idx]

	sse := datastar.NewSSE(w, r)
	return sse.PatchElementTempl(gameRowEditing(game))
}

func (h Handler) Save(w http.ResponseWriter, r *http.Request) error {
	gameID := r.URL.Query().Get("gameID")

	var signals struct {
		gamesSignal
		Team1Score string `json:"team1Score"`
		Team2Score string `json:"team2Score"`
	}

	err := datastar.ReadSignals(r, &signals)
	if err != nil {
		return err
	}

	team1Score, err := strconv.Atoi(signals.Team1Score)
	if err != nil {
		return err
	}

	team2Score, err := strconv.Atoi(signals.Team2Score)
	if err != nil {
		return err
	}

	var isReset bool
	var losingScore int
	switch {
	case team1Score == 0 && team2Score == 0:
		isReset = true
	case team1Score == 121 && team2Score != 121:
		losingScore = team2Score
	case team2Score == 121 && team1Score != 121:
		losingScore = team1Score
	default:
		return components.ShowErrorToast(w, r, "Either both scores must be 0 or one score must be 121.")
	}

	if !isReset && (losingScore >= 121 || losingScore <= 0) {
		return components.ShowErrorToast(w, r, "The losing score must be between 0 and 121.")
	}

	idx := slices.IndexFunc(signals.Games, func(g game) bool { return g.GameID == gameID })
	game := signals.Games[idx]

	err = h.GameRepo.UpdateScores(r.Context(), gameID, game.Team1ID, team1Score, game.Team2ID, team2Score)
	if err != nil {
		return err
	}

	game.Team1Score = team1Score
	game.Team2Score = team2Score

	sse := datastar.NewSSE(w, r)
	return sse.PatchElementTempl(gameRow(game))
}

type game struct {
	GameID       string
	DivisionName string
	Team1ID      string
	Team1Name    string
	Team1Score   int
	Team2ID      string
	Team2Name    string
	Team2Score   int
}

func (h Handler) generatePrelimGames(ctx context.Context) error {
	allTeams, err := h.TeamRepo.GetAll(ctx)
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
			_, err := h.GameRepo.Create(ctx, pair[0].ID, pair[1].ID)
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

	scores, err := h.GameRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	scoresByGame := make(map[string][]games.Score)
	for _, s := range scores {
		scoresByGame[s.GameID] = append(scoresByGame[s.GameID], s)
	}

	daTeams, err := h.TeamRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	teamsByID := make(map[string]teams.Team, len(daTeams))
	for _, t := range daTeams {
		teamsByID[t.ID] = t
	}

	daDivisions, err := h.DivisionRepo.GetAll(ctx)
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
			GameID:       gameID,
			DivisionName: division.Name,
			Team1ID:      team1.ID,
			Team1Name:    team1.Name,
			Team1Score:   scores[0].Score,
			Team2ID:      team2.ID,
			Team2Name:    team2.Name,
			Team2Score:   scores[1].Score,
		}
		daGames = append(daGames, g)
	}

	slices.SortFunc(daGames, func(a, b game) int {
		return cmp.Or(
			cmp.Compare(a.DivisionName, b.DivisionName),
			cmp.Compare(a.Team1Name, b.Team1Name),
			cmp.Compare(a.Team2Name, b.Team2Name),
		)
	})

	return daGames, nil
}
