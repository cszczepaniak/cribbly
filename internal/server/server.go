package server

import (
	"log"
	"net/http"

	"github.com/cszczepaniak/cribbly/internal/ui/pages/admin"
	"github.com/cszczepaniak/cribbly/internal/ui/pages/admin/adminmiddleware"
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
)

func Setup(cfg Config) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /public/", http.StripPrefix("/public", http.FileServer(http.Dir("public"))))
	mux.Handle("GET /", handleWithError(index.Index))
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("unknown route", r.Method, r.URL)
	}))

	r := NewRouter(mux)
	setupAdminRoutes(cfg, r)

	dh := pubdiv.Handler{
		DivisionService: cfg.DivisionService,
		TeamService:     cfg.TeamService,
	}
	r.Handle("GET /divisions", dh.Index)
	r.Handle("GET /divisions/{id}", dh.GetDivisions)

	th := pubteam.Handler{
		GameService: cfg.GameService,
		TeamService: cfg.TeamService,
	}
	r.Handle("GET /teams/{id}/games", th.GetGames)

	gh := pubgame.Handler{
		GameService: cfg.GameService,
		TeamService: cfg.TeamService,
	}
	r.Handle("GET /games/{id}", gh.GetGame)
	r.Handle("PUT /games/{id}", gh.UpdateGame)
	r.Handle("GET /standings", gh.StandingsPage)
	r.Handle("GET /standings/stream", gh.StreamStandings)

	return mux
}

func setupAdminRoutes(cfg Config, r *router) {
	// NOTE: these two admin routes must be registered without using the admin router because they
	// must _not_ have the auth middleware attached to them.
	ah := admin.AdminHandler{
		UserService: cfg.UserService,
	}
	r.Handle("GET /admin/login", admin.LoginPage)
	r.Handle("POST /admin/login", ah.DoLogin)
	r.Handle("GET /admin/register", admin.RegisterPage)
	r.Handle("POST /admin/register", ah.Register)

	adminRouter := r.Group("/admin", adminmiddleware.AuthenticationMiddleware(cfg.UserService))
	adminRouter.Handle("GET /", admin.Index)

	ph := players.PlayersHandler{
		PlayerService: cfg.PlayerService,
	}
	playersRouter := adminRouter.Group("/players")
	playersRouter.Handle("GET /", ph.RegistrationPage)
	playersRouter.Handle("POST /", ph.PostPlayer)
	playersRouter.Handle("POST /random", ph.GenerateRandomPlayers)
	playersRouter.Handle("DELETE /{id}", ph.DeletePlayer)
	playersRouter.Handle("DELETE /", ph.DeleteAllPlayers)

	th := teams.TeamsHandler{
		PlayerService: cfg.PlayerService,
		TeamService:   cfg.TeamService,
	}
	teamsRouter := adminRouter.Group("/teams")
	teamsRouter.Handle("GET /", th.Index)
	teamsRouter.Handle("GET /edit/{id}", th.Edit)
	teamsRouter.Handle("POST /", th.Create)
	teamsRouter.Handle("PUT /{id}", th.Save)
	teamsRouter.Handle("DELETE /{id}", th.Delete)

	dh := divisions.DivisionsHandler{
		TeamService:     cfg.TeamService,
		DivisionService: cfg.DivisionService,
	}
	divisionsRouter := adminRouter.Group("/divisions")
	divisionsRouter.Handle("GET /", dh.Index)
	divisionsRouter.Handle("GET /edit/{id}", dh.Edit)
	divisionsRouter.Handle("POST /", dh.Create)
	divisionsRouter.Handle("PUT /{id}", dh.Save)
	divisionsRouter.Handle("DELETE /{id}", dh.Delete)

	gh := games.Handler{
		DivisionService: cfg.DivisionService,
		TeamService:     cfg.TeamService,
		GameService:     cfg.GameService,
	}
	gamesRouter := adminRouter.Group("/games")
	gamesRouter.Handle("GET /", gh.Index)
	gamesRouter.Handle("POST /generate", gh.Generate)
	gamesRouter.Handle("DELETE /", gh.DeleteAll)
	gamesRouter.Handle("PUT /scores/edit", gh.Edit)
	gamesRouter.Handle("PUT /scores/save", gh.Save)

	uh := users.UsersHandler{
		UserService: cfg.UserService,
	}
	usersRouter := adminRouter.Group("/users")
	usersRouter.Handle("GET /", uh.Index)
	usersRouter.Handle("POST /", uh.Create)
	usersRouter.Handle("DELETE /{name}", uh.Delete)

	pph := profile.ProfileHandler{
		UserService: cfg.UserService,
	}
	profileRouter := adminRouter.Group("/profile")
	profileRouter.Handle("GET /", pph.Index)
	profileRouter.Handle("POST /password", pph.ChangePassword)
}
