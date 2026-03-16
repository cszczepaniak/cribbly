package tournament

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/starfederation/datastar-go/datastar"

	"github.com/cszczepaniak/cribbly/internal/notifier"
	"github.com/cszczepaniak/cribbly/internal/persistence/database"
	"github.com/cszczepaniak/cribbly/internal/persistence/games"
	"github.com/cszczepaniak/cribbly/internal/persistence/teams"
)

type teamAreaProps struct {
	id            string
	round         int
	idx           int
	name          string
	winnerName    string
	top           bool
	gameReady     bool
	showRevert    bool
	revertFromIdx int
	revertToRound int
}

func (p teamAreaProps) isWinner() bool {
	return p.winnerName != "" && p.winnerName == p.name
}

func (p teamAreaProps) isLoser() bool {
	return p.winnerName != "" && p.winnerName != p.name
}

type Handler struct {
	TeamRepo           teams.Repository
	GameRepo           games.Repository
	TournamentNotifier *notifier.Notifier
	Transactor         database.Transactor
}

type signalInt int

func (i *signalInt) UnmarshalText(text []byte) error {
	n, err := strconv.Atoi(string(text))
	if err != nil {
		return err
	}
	*i = signalInt(n)
	return nil
}

func (i signalInt) N() int {
	return int(i)
}

type row struct {
	round     int
	idx       int
	team1ID   string
	team1Name string
	team2ID   string
	team2Name string
	winner    string
	winnerID  string
}

type round struct {
	games []row
}

func (h Handler) Index(w http.ResponseWriter, r *http.Request) error {
	rounds, teamCount, err := h.loadRounds(r.Context())
	if err != nil {
		return err
	}

	return index(rounds, teamCount, championFromRounds(rounds)).Render(r.Context(), w)
}

func (h Handler) Stream(w http.ResponseWriter, r *http.Request) error {
	sub, done := h.TournamentNotifier.Subscribe()
	defer done()

	sse := datastar.NewSSE(w, r)
	for {
		select {
		case <-r.Context().Done():
			return nil
		case <-sub:
			rounds, _, err := h.loadRounds(r.Context())
			if err != nil {
				return err
			}
			err = sse.PatchElementTempl(roundDisplay(rounds, 0, championFromRounds(rounds)), datastar.WithViewTransitions())
			if err != nil {
				return err
			}
		}
	}
}

func (h Handler) AdvanceTeam(w http.ResponseWriter, r *http.Request) error {
	teamID := r.PathValue("id")

	fromIdx, err := strconv.Atoi(r.URL.Query().Get("fromIdx"))
	if err != nil {
		return err
	}

	toRound, err := strconv.Atoi(r.URL.Query().Get("toRound"))
	if err != nil {
		return err
	}

	err = h.GameRepo.SetTournamentGameWinner(r.Context(), toRound-1, fromIdx, teamID)
	if err != nil {
		return err
	}

	rounds, _, err := h.loadRounds(r.Context())
	if err != nil {
		return err
	}

	// Only advance team to next round if there is one (skip for final/champion game)
	if toRound < len(rounds) {
		newIdx := fromIdx / 2
		if fromIdx%2 == 0 {
			err = h.GameRepo.PutTeam1IntoTournamentGame(r.Context(), toRound, newIdx, teamID)
		} else {
			err = h.GameRepo.PutTeam2IntoTournamentGame(r.Context(), toRound, newIdx, teamID)
		}
		if err != nil {
			return err
		}
	}

	rounds, _, err = h.loadRounds(r.Context())
	if err != nil {
		return err
	}

	h.TournamentNotifier.Notify()
	return datastar.NewSSE(w, r).PatchElementTempl(roundDisplay(rounds, 0, championFromRounds(rounds)), datastar.WithViewTransitions())
}

func (h Handler) RevertAdvance(w http.ResponseWriter, r *http.Request) error {
	teamID := r.PathValue("id")

	fromIdx, err := strconv.Atoi(r.URL.Query().Get("fromIdx"))
	if err != nil {
		return err
	}

	toRound, err := strconv.Atoi(r.URL.Query().Get("toRound"))
	if err != nil {
		return err
	}

	rounds, _, err := h.loadRounds(r.Context())
	if err != nil {
		return err
	}

	// Only allow revert from the team's furthest game: the game they're being removed from must have no winner yet
	gameIdx := fromIdx / 2
	if toRound >= len(rounds) {
		// Final round: clearing champion is always allowed
	} else {
		game := rounds[toRound].games[gameIdx]
		if game.winner != "" {
			return errors.New("can only revert a team from their furthest game")
		}
		// Team must be in this game
		if game.team1ID != teamID && game.team2ID != teamID {
			return errors.New("team is not in this game")
		}
	}

	err = h.Transactor.WithTx(r.Context(), func(ctx context.Context) error {
		// Clear winner on the game this team was advanced from
		err = h.GameRepo.ClearTournamentGameWinner(ctx, toRound-1, fromIdx)
		if err != nil {
			return err
		}

		// If they were advanced into a next round, clear that slot too
		if toRound < len(rounds) {
			err = h.GameRepo.ClearTeamFromTournamentGame(ctx, toRound, gameIdx, teamID)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	rounds, _, err = h.loadRounds(r.Context())
	if err != nil {
		return err
	}

	h.TournamentNotifier.Notify()
	return datastar.NewSSE(w, r).PatchElementTempl(roundDisplay(rounds, 0, championFromRounds(rounds)), datastar.WithViewTransitions())
}

func (h Handler) Generate(w http.ResponseWriter, r *http.Request) error {
	var signals struct {
		Size signalInt `json:"size"`
	}
	err := datastar.ReadSignals(r, &signals)
	if err != nil {
		return err
	}

	standings, err := h.GameRepo.GetStandings(r.Context())
	if err != nil {
		return err
	}

	if len(standings) < signals.Size.N() {
		return errors.New("not enough teams to seed the tournament")
	}

	err = h.GameRepo.InitializeTournament(r.Context(), signals.Size.N())
	if err != nil {
		return err
	}

	for i := range signals.Size.N() / 2 {
		// compute 0,15 1,14 2,13 etc. for the tournament seeds
		idx1 := i
		idx2 := signals.Size.N() - (i + 1)
		err := h.GameRepo.PutTeam1IntoTournamentGame(r.Context(), 0, i, standings[idx1].TeamID)
		if err != nil {
			return err
		}
		err = h.GameRepo.PutTeam2IntoTournamentGame(r.Context(), 0, i, standings[idx2].TeamID)
		if err != nil {
			return err
		}
	}

	rounds, _, err := h.loadRounds(r.Context())
	if err != nil {
		return err
	}

	return datastar.NewSSE(w, r).PatchElementTempl(tournamentPage(rounds, 0, championFromRounds(rounds)))
}

func (h Handler) Delete(w http.ResponseWriter, r *http.Request) error {
	err := h.GameRepo.DeleteTournament(r.Context())
	if err != nil {
		return err
	}

	allTeams, err := h.TeamRepo.GetAll(r.Context())
	if err != nil {
		return err
	}
	teamCount := len(allTeams)

	return datastar.NewSSE(w, r).PatchElementTempl(tournamentPage(nil, teamCount, champion{}))
}

func (h Handler) loadRounds(ctx context.Context) ([]round, int, error) {
	tourney, err := h.GameRepo.LoadTournament(ctx)
	if err != nil {
		return nil, 0, err
	}

	ts, err := h.TeamRepo.GetAll(ctx)
	if err != nil {
		return nil, 0, err
	}
	teamNamesByID := make(map[string]string, len(ts))
	for _, t := range ts {
		teamNamesByID[t.ID] = t.Name
	}

	var rounds []round
	for i, rnd := range tourney.Rounds {
		var games []row
		for j, game := range rnd.Games {
			games = append(games, row{
				round:     i,
				idx:       j,
				team1ID:   game.TeamIDs[0],
				team1Name: teamNamesByID[game.TeamIDs[0]],
				team2ID:   game.TeamIDs[1],
				team2Name: teamNamesByID[game.TeamIDs[1]],
				winner:    teamNamesByID[game.Winner],
				winnerID:  game.Winner,
			})
		}

		rounds = append(rounds, round{
			games: games,
		})
	}

	return rounds, len(ts), nil
}

type champion struct {
	Name string
	ID   string
}

func championFromRounds(rounds []round) champion {
	if len(rounds) == 0 {
		return champion{}
	}
	last := rounds[len(rounds)-1]
	if len(last.games) != 1 {
		return champion{}
	}
	g := last.games[0]
	if g.winner == "" {
		return champion{}
	}
	return champion{Name: g.winner, ID: g.winnerID}
}
