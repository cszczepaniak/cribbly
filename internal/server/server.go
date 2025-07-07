package server

import (
	"net/http"

	"github.com/cszczepaniak/cribbly/internal/ui/components"
	"github.com/cszczepaniak/cribbly/internal/ui/pages/admin"
	"github.com/cszczepaniak/cribbly/internal/ui/pages/admin/divisions"
	"github.com/cszczepaniak/cribbly/internal/ui/pages/admin/players"
	"github.com/cszczepaniak/cribbly/internal/ui/pages/admin/teams"
	"github.com/cszczepaniak/cribbly/internal/ui/pages/index"
)

func Setup(cfg Config) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /public/", http.StripPrefix("/public", http.FileServer(http.Dir("public"))))

	mux.Handle("GET /", components.Handle(index.Index))
	mux.Handle("GET /admin", components.Handle(admin.Index))

	ph := players.PlayersHandler{
		PlayerService: cfg.PlayerService,
	}
	mux.Handle("GET /admin/players", components.Handle(ph.RegistrationPage))
	mux.Handle("POST /admin/players", components.Handle(ph.PostPlayer))

	th := teams.TeamsHandler{
		PlayerService: cfg.PlayerService,
		TeamService:   cfg.TeamService,
	}
	mux.Handle("GET /admin/teams", components.Handle(th.Index))
	mux.Handle("POST /admin/teams", components.Handle(th.Create))
	mux.Handle("PUT /admin/teams/{id}", components.Handle(th.Save))
	mux.Handle("DELETE /admin/teams/{id}", components.Handle(th.Delete))

	dh := divisions.DivisionsHandler{
		TeamService:     cfg.TeamService,
		DivisionService: cfg.DivisionService,
	}
	mux.Handle("GET /admin/divisions", components.Handle(dh.Index))
	mux.Handle("POST /admin/divisions", components.Handle(dh.Create))
	mux.Handle("PUT /admin/divisions/{id}", components.Handle(dh.Save))
	mux.Handle("DELETE /admin/divisions/{id}", components.Handle(dh.Delete))

	return mux
}
