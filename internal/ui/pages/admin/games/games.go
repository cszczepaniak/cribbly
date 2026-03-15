package games

import (
	"cmp"
	"context"
	"errors"
	"net/http"
	"slices"
	"strconv"

	"github.com/starfederation/datastar-go/datastar"

	"github.com/cszczepaniak/cribbly/internal/persistence/database"
	"github.com/cszczepaniak/cribbly/internal/persistence/divisions"
	"github.com/cszczepaniak/cribbly/internal/persistence/games"
	"github.com/cszczepaniak/cribbly/internal/persistence/teams"
	"github.com/cszczepaniak/cribbly/internal/ui/components"
)

type Handler struct {
	Transactor   database.Transactor
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

func (h Handler) ResetScores(w http.ResponseWriter, r *http.Request) error {
	gameID := r.URL.Query().Get("gameID")

	var signals gamesSignal
	err := datastar.ReadSignals(r, &signals)
	if err != nil {
		return err
	}

	idx := slices.IndexFunc(signals.Games, func(g game) bool { return g.GameID == gameID })
	game := signals.Games[idx]

	err = h.GameRepo.UpdateScores(r.Context(), gameID, game.Team1ID, 0, game.Team2ID, 0)
	if err != nil {
		return err
	}
	game.Team1Score = 0
	game.Team2Score = 0

	return datastar.NewSSE(w, r).PatchElementTempl(gameRow(game))
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
	return h.Transactor.WithTx(ctx, func(ctx context.Context) error {
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

		// TODO: don't insert one-by-one (but at least we're in a transaction now!)
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
	})
}

// generateMatchups returns pairings so every team plays the same number of games (fair for
// standings). 3 teams → 2 games each; 4 → 3 each; 5 → 4 each; 6 → 3 each. Only 3–6 teams allowed.
func generateMatchups(allTeams []teams.Team) ([][2]teams.Team, error) {
	switch n := len(allTeams); n {
	case 3, 4, 5:
		// Full round robin: all pairs once. 2, 3, or 4 games per team.
		var pairs [][2]teams.Team
		for i := range n {
			for j := i + 1; j < n; j++ {
				pairs = append(pairs, [2]teams.Team{allTeams[i], allTeams[j]})
			}
		}
		return pairs, nil
	case 6:
		// Circle method: 3 games per team = 9 games, no duplicate pairings.
		//
		// Imagine 6 teams sitting in a circle. Team 0 stays fixed; teams 1–5 sit in order around
		// the rest. Each round we pair "across": (0,5), (1,4), (2,3). Then we rotate teams 1–5
		// one position (e.g. 1->2->3->4->5->1) and repeat. After 3 rounds every team has played 3 games
		// and no pair has met twice.
		others := []int{1, 2, 3, 4, 5} // indices of the rotating teams (0 is fixed)
		var pairs [][2]teams.Team
		for r := range 3 {
			// Build this round's seating order: 0 fixed, then others rotated by r positions.
			// Round 0: [0,1,2,3,4,5]. Round 1: [0,5,1,2,3,4]. Round 2: [0,4,5,1,2,3].
			order := make([]int, 6)
			order[0] = 0
			for i := range 5 {
				order[i+1] = others[(i-r+5)%5]
			}
			// Pair across the circle: first with last, second with second-last, third with third-last.
			pairs = append(pairs,
				[2]teams.Team{allTeams[order[0]], allTeams[order[5]]},
				[2]teams.Team{allTeams[order[1]], allTeams[order[4]]},
				[2]teams.Team{allTeams[order[2]], allTeams[order[3]]},
			)
		}
		return pairs, nil
	default:
		return nil, errors.New("can only generate matchups for 3, 4, 5, or 6 teams")
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
