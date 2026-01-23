package tournament

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/cszczepaniak/cribbly/internal/persistence/divisions"
	"github.com/cszczepaniak/cribbly/internal/persistence/games"
	"github.com/cszczepaniak/cribbly/internal/persistence/teams"
	"github.com/starfederation/datastar-go/datastar"
)

type Handler struct {
	DivisionRepo divisions.Repository
	TeamRepo     teams.Repository
	GameRepo     games.Repository
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
}

func (h Handler) Index(w http.ResponseWriter, r *http.Request) error {
	tourney, err := h.GameRepo.LoadTournament(r.Context())
	if err != nil {
		return err
	}

	ts, err := h.TeamRepo.GetAll(r.Context())
	if err != nil {
		return err
	}
	teamNamesByID := make(map[string]string, len(ts))
	for _, t := range ts {
		teamNamesByID[t.ID] = t.Name
	}

	var rows []row
	for i, round := range tourney.Rounds {
		for j, game := range round.Games {
			rows = append(rows, row{
				round:     i,
				idx:       j,
				team1ID:   game.TeamIDs[0],
				team1Name: teamNamesByID[game.TeamIDs[0]],
				team2ID:   game.TeamIDs[1],
				team2Name: teamNamesByID[game.TeamIDs[1]],
				winner:    teamNamesByID[game.Winner],
			})
		}
	}

	return index(rows).Render(r.Context(), w)
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

	// 0,1->0; 2,3->1; etc.
	newIdx := fromIdx / 2
	err = h.GameRepo.PutTeamIntoTournamentGame(r.Context(), toRound, newIdx, teamID)
	if err != nil {
		return err
	}

	return datastar.NewSSE(w, r).Redirect("/admin/tournament")
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
		err := h.GameRepo.PutTeamIntoTournamentGame(r.Context(), 0, i, standings[idx1].TeamID)
		if err != nil {
			return err
		}
		err = h.GameRepo.PutTeamIntoTournamentGame(r.Context(), 0, i, standings[idx2].TeamID)
		if err != nil {
			return err
		}
	}

	return datastar.NewSSE(w, r).Redirect("/admin/tournament")
}

func (h Handler) Delete(w http.ResponseWriter, r *http.Request) error {
	err := h.GameRepo.DeleteTournament(r.Context())
	if err != nil {
		return err
	}

	return datastar.NewSSE(w, r).Redirect("/admin/tournament")
}
