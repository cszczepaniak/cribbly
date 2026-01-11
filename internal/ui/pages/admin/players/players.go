package players

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/starfederation/datastar-go/datastar"

	"github.com/cszczepaniak/cribbly/internal/persistence/players"
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

	signals.Name = ""
	err = sse.MarshalAndPatchSignals(signals)
	if err != nil {
		return err
	}

	return sse.PatchElementTempl(playerTable(players))
}

func (h PlayersHandler) GenerateRandomPlayers(w http.ResponseWriter, r *http.Request) error {
	var signals struct {
		Num int `json:"num"`
	}
	if err := datastar.ReadSignals(r, &signals); err != nil {
		return err
	}

	for range signals.Num {
		name := "Player " + uuid.NewString()
		_, err := h.PlayerService.Create(r.Context(), name)
		if err != nil {
			return err
		}
	}

	players, err := h.PlayerService.GetAll(r.Context())
	if err != nil {
		return err
	}

	sse := datastar.NewSSE(w, r)
	return sse.PatchElementTempl(playerTable(players))
}

func (h PlayersHandler) DeleteAllPlayers(w http.ResponseWriter, r *http.Request) error {
	players, err := h.PlayerService.GetAll(r.Context())
	if err != nil {
		return err
	}

	for _, p := range players {
		err := h.PlayerService.UnassignFromTeam(r.Context(), p.ID)
		if err != nil {
			return err
		}

		err = h.PlayerService.Delete(r.Context(), p.ID)
		if err != nil {
			return err
		}
	}

	sse := datastar.NewSSE(w, r)
	return sse.PatchElementTempl(playerTable(nil))
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
