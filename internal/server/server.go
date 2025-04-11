package server

import (
	"net/http"

	"github.com/cszczepaniak/cribbly/internal/ui/components"
	"github.com/cszczepaniak/cribbly/internal/ui/pages/admin/players"
	"github.com/cszczepaniak/cribbly/internal/ui/pages/index"
)

func Setup(cfg Config) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /", components.Handle(index.Index))

	ph := players.PlayersHandler{
		PlayerService: cfg.PlayerService,
	}
	mux.Handle("GET /admin/players", components.Handle(ph.RegistrationPage))

	return mux
}
