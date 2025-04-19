package server

import (
	"net/http"

	"github.com/cszczepaniak/cribbly/internal/ui/components"
	"github.com/cszczepaniak/cribbly/internal/ui/pages/admin"
	"github.com/cszczepaniak/cribbly/internal/ui/pages/admin/players"
	"github.com/cszczepaniak/cribbly/internal/ui/pages/admin/teams"
	"github.com/cszczepaniak/cribbly/internal/ui/pages/index"
)

func Setup(cfg Config) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /", components.Handle(index.Index))

	mux.Handle("GET /admin", components.Handle(admin.Index))

	ph := players.PlayersHandler{
		PlayerService: cfg.PlayerService,
	}
	mux.Handle("GET /admin/players", components.Handle(ph.RegistrationPage))
	mux.Handle("POST /admin/players", components.Handle(ph.PostPlayer))

	th := teams.TeamsHandler{
		PlayerService: cfg.PlayerService,
	}
	mux.Handle("GET /admin/teams", components.Handle(th.Index))
	mux.Handle("POST /admin/teams", components.Handle(th.Create))
	mux.Handle("PUT /admin/teams/{id}", components.Handle(th.Save))

	return mux
}
