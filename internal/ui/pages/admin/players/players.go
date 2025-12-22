package players

import (
	"net/http"

	"github.com/starfederation/datastar-go/datastar"

	"github.com/cszczepaniak/cribbly/internal/persistence/players"
)

const (
	nameFormKey = "name"
)

type PlayersHandler struct {
	PlayerService players.Service
}

func (h PlayersHandler) RegistrationPage(w http.ResponseWriter, r *http.Request) error {
	players, err := h.PlayerService.GetAll(r.Context())
	if err != nil {
		return err
	}

	tm := playerRegistrationPage(players)
	return tm.Render(r.Context(), w)
}

func (h PlayersHandler) PostPlayer(w http.ResponseWriter, r *http.Request) error {
	var signals struct {
		Name string `json:"name"`
	}
	if err := datastar.ReadSignals(r, &signals); err != nil {
		return err
	}

	_, err := h.PlayerService.Create(r.Context(), signals.Name)
	if err != nil {
		return err
	}

	players, err := h.PlayerService.GetAll(r.Context())
	if err != nil {
		return err
	}

	sse := datastar.NewSSE(w, r)
	return sse.PatchElementTempl(playerTable(players))
}

func (h PlayersHandler) DeletePlayer(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")

	err := h.PlayerService.UnassignFromTeam(r.Context(), id)
	if err != nil {
		return err
	}

	err = h.PlayerService.Delete(r.Context(), id)
	if err != nil {
		return err
	}

	players, err := h.PlayerService.GetAll(r.Context())
	if err != nil {
		return err
	}

	sse := datastar.NewSSE(w, r)
	return sse.PatchElementTempl(playerTable(players))
}
