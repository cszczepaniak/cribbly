package server

import (
	"log"
	"net/http"

	"github.com/cszczepaniak/cribbly/internal/ui/components"
	"github.com/cszczepaniak/cribbly/internal/ui/pages/admin"
	"github.com/cszczepaniak/cribbly/internal/ui/pages/admin/divisions"
	"github.com/cszczepaniak/cribbly/internal/ui/pages/admin/games"
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
	mux.Handle("GET /admin/players", handleWithError(ph.RegistrationPage))
	mux.Handle("POST /admin/players", handleWithError(ph.PostPlayer))
	mux.Handle("POST /admin/players/random", handleWithError(ph.GenerateRandomPlayers))
	mux.Handle("DELETE /admin/players/{id}", handleWithError(ph.DeletePlayer))
	mux.Handle("DELETE /admin/players", handleWithError(ph.DeleteAllPlayers))

	th := teams.TeamsHandler{
		PlayerService: cfg.PlayerService,
		TeamService:   cfg.TeamService,
	}
	mux.Handle("GET /admin/teams", handleWithError(th.Index))
	mux.Handle("GET /admin/teams/edit/{id}", handleWithError(th.Edit))
	mux.Handle("POST /admin/teams/edit/cancel", handleWithError(th.CancelEdit))
	mux.Handle("POST /admin/teams", handleWithError(th.Create))
	mux.Handle("PUT /admin/teams/{id}", handleWithError(th.Save))
	mux.Handle("DELETE /admin/teams/{id}", handleWithError(th.Delete))

	dh := divisions.DivisionsHandler{
		TeamService:     cfg.TeamService,
		DivisionService: cfg.DivisionService,
	}
	mux.Handle("GET /admin/divisions", handleWithError(dh.Index))
	mux.Handle("GET /admin/divisions/edit/{id}", handleWithError(dh.Edit))
	mux.Handle("POST /admin/divisions/edit/cancel", handleWithError(dh.CancelEdit))
	mux.Handle("POST /admin/divisions", handleWithError(dh.Create))
	mux.Handle("PUT /admin/divisions/{id}", handleWithError(dh.Save))
	mux.Handle("DELETE /admin/divisions/{id}", handleWithError(dh.Delete))

	gh := games.Handler{}
	mux.Handle("GET /admin/games", handleWithError(gh.Index))

	return mux
}

func handleWithError(fn func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL)

		err := fn(w, r)
		if err != nil {
			log.Println(r.Method, r.URL, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Println(r.Method, r.URL, "COMPLETE")
	}
}
