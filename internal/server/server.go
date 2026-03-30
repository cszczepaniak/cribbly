package server

import (
	"net/http"
	"sync"

	"github.com/cszczepaniak/cribbly/internal/api/playersconnect"
	"github.com/cszczepaniak/cribbly/internal/api/roomcodeconnect"
	cribblyv1connect "github.com/cszczepaniak/cribbly/internal/gen/cribbly/v1/cribblyv1connect"
	mw "github.com/cszczepaniak/cribbly/internal/server/middleware"
	"github.com/cszczepaniak/cribbly/internal/ui/pages/admin"
	"github.com/cszczepaniak/cribbly/internal/ui/pages/admin/divisions"
	"github.com/cszczepaniak/cribbly/internal/ui/pages/admin/games"
	"github.com/cszczepaniak/cribbly/internal/ui/pages/admin/players"
	"github.com/cszczepaniak/cribbly/internal/ui/pages/admin/profile"
	"github.com/cszczepaniak/cribbly/internal/ui/pages/admin/roomcodes"
	"github.com/cszczepaniak/cribbly/internal/ui/pages/admin/teams"
	"github.com/cszczepaniak/cribbly/internal/ui/pages/admin/users"
	pubdiv "github.com/cszczepaniak/cribbly/internal/ui/pages/divisions"
	pubgame "github.com/cszczepaniak/cribbly/internal/ui/pages/games"
	"github.com/cszczepaniak/cribbly/internal/ui/pages/index"
	pubteam "github.com/cszczepaniak/cribbly/internal/ui/pages/teams"
	pubtournament "github.com/cszczepaniak/cribbly/internal/ui/pages/tournament"
	webembed "github.com/cszczepaniak/cribbly/internal/web/embed"
)

// noCacheStatic wraps a handler to set Cache-Control headers so the browser
// does not cache static assets. Use in development so CSS/JS changes are visible immediately.
func noCacheStatic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		next.ServeHTTP(w, r)
	})
}

func Setup(cfg Config) http.Handler {
	mux := http.NewServeMux()

	publicFiles := http.StripPrefix("/public", http.FileServer(http.Dir("public")))
	if !cfg.IsProd {
		publicFiles = noCacheStatic(publicFiles)
	}
	mux.Handle("GET /public/", publicFiles)

	reactStatic := webembed.StaticHandler()
	if !cfg.IsProd {
		reactStatic = noCacheStatic(reactStatic)
	}
	mux.Handle("GET /app", http.RedirectHandler("/app/", http.StatusMovedPermanently))
	mux.Handle("GET /app/", reactStatic)

	r := NewRouter(
		mux,
		mw.DevAdminBypassMiddleware(cfg.DevAdminSecret, cfg.IsProd),
		mw.AuthenticationMiddleware(cfg.UserRepo),
		mw.IsProdMiddleware(cfg.IsProd),
		mw.DevToolsQueryMiddleware(),
		mw.RoomCodeMiddleware(cfg.RoomCodeRepo),
	)

	home := index.Handler{
		RoomCodeRepo: cfg.RoomCodeRepo,
	}
	r.Handle("GET /", home.Index)
	r.Handle("POST /room-code", home.SubmitRoomCode)

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
		Transactor:         cfg.Transactor,
	}
	r.Handle("GET /tournament", tourneyHandler.Index)
	r.Handle("GET /tournament/stream", tourneyHandler.Stream)
	r.Handle("POST /tournament", tourneyHandler.Generate, mw.ErrorIfNotAdmin())
	r.Handle("DELETE /tournament", tourneyHandler.Delete, mw.ErrorIfNotAdmin())
	r.Handle("POST /tournament/team/{id}/advance", tourneyHandler.AdvanceTeam, mw.ErrorIfNotAdmin())
	r.Handle("POST /tournament/team/{id}/revert", tourneyHandler.RevertAdvance, mw.ErrorIfNotAdmin())

	rcConnect := &roomcodeconnect.Server{Repo: cfg.RoomCodeRepo, UserRepo: cfg.UserRepo}
	connectMountPath, roomCodeConnectHandler := cribblyv1connect.NewRoomCodeServiceHandler(rcConnect)

	// The generated Connect HTTP handler expects r.URL.Path to match the procedure path
	// (e.g. "/cribbly.v1.RoomCodeService/SetRoomCode") exactly. Since we expose it under
	// our own "/api" prefix, we strip that prefix before invoking the handler.
	roomCodeConnect := http.StripPrefix("/api", connectWithAdminContext(cfg, roomCodeConnectHandler))
	mux.Handle("POST /api"+connectMountPath, roomCodeConnect)

	plConnect := &playersconnect.Server{PlayerRepo: cfg.PlayerRepo}
	playerMountPath, playerConnectHandler := cribblyv1connect.NewPlayerServiceHandler(plConnect)
	mux.Handle("POST /api"+playerMountPath, http.StripPrefix("/api", connectWithAdminContext(cfg, playerConnectHandler)))

	return mw.ReactQueryMiddleware(sync.OnceValue(webembed.MustReadIndexHTML), cfg.IsProd, mux)
}

// connectWithAdminContext applies dev-admin bypass and session cookie to Connect requests (the
// mux routes are not wrapped by NewRouter's AuthenticationMiddleware).
func connectWithAdminContext(cfg Config, h http.Handler) http.Handler {
	return withDevAdminRequestContext(cfg, mw.ConnectSessionMiddleware(cfg.UserRepo)(h))
}

func withDevAdminRequestContext(cfg Config, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = mw.WithDevAdminBypassIfHeader(r, cfg.DevAdminSecret, cfg.IsProd)
		h.ServeHTTP(w, r)
	})
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
	playersRouter.Handle("POST /excel", ph.UploadExcel)
	playersRouter.Handle("POST /excel/import", ph.ImportExcel)

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

	rcHandler := roomcodes.Handler{
		RoomCodeRepo: cfg.RoomCodeRepo,
	}
	roomCodesRouter := adminRouter.Group("/room-codes")
	roomCodesRouter.Handle("GET /", rcHandler.Index)
	roomCodesRouter.Handle("POST /", rcHandler.Generate)

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
