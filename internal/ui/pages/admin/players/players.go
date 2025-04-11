package players

import (
	"net/http"

	"github.com/a-h/templ"

	"github.com/cszczepaniak/cribbly/internal/persistence/players"
)

type PlayersHandler struct {
	PlayerService players.Service
}

func (h PlayersHandler) RegistrationPage(r *http.Request) (templ.Component, error) {
	return h.renderAllPlayers(r, playerRegistrationPage)
}

func (h PlayersHandler) renderAllPlayers(
	r *http.Request,
	fn func([]string) templ.Component,
) (templ.Component, error) {
	players, err := h.PlayerService.GetAll(r.Context())
	if err != nil {
		return nil, err
	}

	playerNames := make([]string, 0, len(players))
	for _, p := range players {
		playerNames = append(playerNames, p.Name)
	}

	return fn(playerNames), nil
}
