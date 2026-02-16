package server

import (
	"log"
	"net/http"

	mw "github.com/cszczepaniak/cribbly/internal/server/middleware"
	"github.com/cszczepaniak/cribbly/internal/ui/pages/admin"
	"github.com/cszczepaniak/cribbly/internal/ui/pages/admin/divisions"
	"github.com/cszczepaniak/cribbly/internal/ui/pages/admin/games"
	"github.com/cszczepaniak/cribbly/internal/ui/pages/admin/players"
	"github.com/cszczepaniak/cribbly/internal/ui/pages/admin/profile"
	"github.com/cszczepaniak/cribbly/internal/ui/pages/admin/teams"
	"github.com/cszczepaniak/cribbly/internal/ui/pages/admin/users"
	pubdiv "github.com/cszczepaniak/cribbly/internal/ui/pages/divisions"
	pubgame "github.com/cszczepaniak/cribbly/internal/ui/pages/games"
	"github.com/cszczepaniak/cribbly/internal/ui/pages/index"
	pubteam "github.com/cszczepaniak/cribbly/internal/ui/pages/teams"
	pubtournament "github.com/cszczepaniak/cribbly/internal/ui/pages/tournament"
)

func Setup(cfg Config) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /public/", http.StripPrefix("/public", http.FileServer(http.Dir("public"))))
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("unknown route", r.Method, r.URL)
	}))

	r := NewRouter(mux, mw.AuthenticationMiddleware(cfg.UserRepo), mw.IsProdMiddleware(cfg.IsProd))
	r.Handle("GET /", index.Index)

	setupAdminRoutes(cfg, r)

	dh := pubdiv.Handler{
		DivisionRepo: cfg.DivisionRepo,
		TeamRepo:     cfg.TeamRepo,
	}
	r.Handle("GET /divisions", dh.Index)
	r.Handle("GET /divisions/{id}", dh.GetDivisions)

	th := pubteam.Handler{
		GameRepo: cfg.GameRepo,
		TeamRepo: cfg.TeamRepo,
	}
	r.Handle("GET /teams/{id}/games", th.GetGames)

	gh := pubgame.Handler{
		GameRepo:            cfg.GameRepo,
		TeamRepo:            cfg.TeamRepo,
		ScoreUpdateNotifier: cfg.ScoreUpdateNotifier,
	}
	r.Handle("GET /games/{id}", gh.GetGame)
	r.Handle("PUT /games/{id}", gh.UpdateGame)

	r.Handle("GET /standings", gh.StandingsPage)
	r.Handle("GET /standings/stream", gh.StreamStandings)

	tourneyHandler := pubtournament.Handler{
		GameRepo:           cfg.GameRepo,
		TeamRepo:           cfg.TeamRepo,
		TournamentNotifier: cfg.TournamentNotifier,
	}
	r.Handle("GET /tournament", tourneyHandler.Index)
	r.Handle("GET /tournament/stream", tourneyHandler.Stream)
	r.Handle("POST /tournament", tourneyHandler.Generate, mw.ErrorIfNotAdmin())
	r.Handle("DELETE /tournament", tourneyHandler.Delete, mw.ErrorIfNotAdmin())
	r.Handle("POST /tournament/team/{id}/advance", tourneyHandler.AdvanceTeam, mw.ErrorIfNotAdmin())

	return mux
}

func setupAdminRoutes(cfg Config, r *router) {
	ah := admin.AdminHandler{
		UserRepo: cfg.UserRepo,
	}

	// NOTE: these admin routes must be registered without using the admin router because they
	// must _not_ have the redirect middleware attached to them.
	r.Handle("GET /admin/login", admin.LoginPage)
	r.Handle("POST /admin/login", ah.DoLogin)
	r.Handle("POST /admin/logout", ah.DoLogout)
	r.Handle("GET /admin/register", admin.RegisterPage)
	r.Handle("POST /admin/register", ah.Register)

	adminRouter := r.Group("/admin", mw.RedirectToLoginIfNotAdmin())
	adminRouter.Handle("GET /", admin.Index)

	ph := players.PlayersHandler{
		PlayerRepo: cfg.PlayerRepo,
	}
	playersRouter := adminRouter.Group("/players")
	playersRouter.Handle("GET /", ph.RegistrationPage)
	playersRouter.Handle("POST /", ph.PostPlayer)
	playersRouter.Handle("POST /random", ph.GenerateRandomPlayers)
	playersRouter.Handle("DELETE /{id}", ph.DeletePlayer)
	playersRouter.Handle("DELETE /", ph.DeleteAllPlayers)

	th := teams.TeamsHandler{
		PlayerRepo:  cfg.PlayerRepo,
		TeamRepo:    cfg.TeamRepo,
		TeamService: cfg.TeamService(),
	}
	teamsRouter := adminRouter.Group("/teams")
	teamsRouter.Handle("GET /", th.Index)
	teamsRouter.Handle("GET /{id}", th.EditPage)
	teamsRouter.Handle("GET /delete/{id}", th.ConfirmDelete)
	teamsRouter.Handle("POST /", th.Create)
	teamsRouter.Handle("POST /generate", th.Generate)
	teamsRouter.Handle("PUT /{id}", th.Save)
	teamsRouter.Handle("DELETE /{id}", th.Delete)
	teamsRouter.Handle("DELETE /", th.DeleteAll)

	dh := divisions.DivisionsHandler{
		Transactor:      cfg.Transactor,
		TeamRepo:        cfg.TeamRepo,
		DivisionRepo:    cfg.DivisionRepo,
		DivisionService: cfg.DivisionService(),
	}
	divisionsRouter := adminRouter.Group("/divisions")
	divisionsRouter.Handle("GET /", dh.Index)
	divisionsRouter.Handle("GET /{id}", dh.EditPage)
	divisionsRouter.Handle("POST /", dh.Create)
	divisionsRouter.Handle("PUT /{id}", dh.Save)
	divisionsRouter.Handle("GET /{id}/delete", dh.Delete)
	divisionsRouter.Handle("DELETE /{id}", dh.ConfirmDelete)
	divisionsRouter.Handle("GET /{id}/editname", dh.EditName)
	divisionsRouter.Handle("PUT /{id}/savename", dh.SaveName)
	divisionsRouter.Handle("PUT /{id}/savesize", dh.SaveSize)
	divisionsRouter.Handle("POST /generate", dh.Generate)
	divisionsRouter.Handle("DELETE /", dh.DeleteAll)
	divisionsRouter.Handle("GET /qrs", dh.GenerateQRCodes)

	gh := games.Handler{
		Transactor:   cfg.Transactor,
		DivisionRepo: cfg.DivisionRepo,
		TeamRepo:     cfg.TeamRepo,
		GameRepo:     cfg.GameRepo,
	}
	gamesRouter := adminRouter.Group("/games")
	gamesRouter.Handle("GET /", gh.Index)
	gamesRouter.Handle("POST /generate", gh.Generate)
	gamesRouter.Handle("DELETE /", gh.DeleteAll)
	gamesRouter.Handle("GET /scores/edit", gh.Edit)
	gamesRouter.Handle("PUT /scores/save", gh.Save)
	gamesRouter.Handle("PUT /scores/reset", gh.ResetScores)

	uh := users.UsersHandler{
		UserRepo: cfg.UserRepo,
	}
	usersRouter := adminRouter.Group("/users")
	usersRouter.Handle("GET /", uh.Index)
	usersRouter.Handle("POST /", uh.Create)
	usersRouter.Handle("DELETE /{name}", uh.Delete)

	pph := profile.ProfileHandler{
		UserRepo: cfg.UserRepo,
	}
	profileRouter := adminRouter.Group("/profile")
	profileRouter.Handle("GET /", pph.Index)
	profileRouter.Handle("POST /password", pph.ChangePassword)
}
