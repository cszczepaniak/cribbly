package tournament

import (
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

func (h Handler) Index(w http.ResponseWriter, r *http.Request) error {
	tourney, err := h.GameRepo.LoadTournament(r.Context())
	if err != nil {
		return err
	}

	return index(tourney).Render(r.Context(), w)
}

func (h Handler) Generate(w http.ResponseWriter, r *http.Request) error {
	var signals struct {
		Size signalInt `json:"size"`
	}
	err := datastar.ReadSignals(r, &signals)
	if err != nil {
		return err
	}

	err = h.GameRepo.InitializeTournament(r.Context(), signals.Size.N())
	if err != nil {
		return err
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
