package integration

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"

	"github.com/cszczepaniak/cribbly/internal/persistence/database"
	"github.com/cszczepaniak/cribbly/internal/persistence/roomcodes"
	"github.com/cszczepaniak/cribbly/internal/persistence/users"
	"github.com/cszczepaniak/cribbly/internal/server"
)

type testApp struct {
	BaseURL      string
	Browser      *rod.Browser
	RoomCodeRepo roomcodes.Repository
	UserRepo     users.Repository
}

func newTestApp(t *testing.T) *testApp {
	t.Helper()

	ctx := t.Context()

	db := database.NewInMemory(t)

	serverCfg, err := server.SetupFromDB(ctx, db, false)
	if err != nil {
		t.Fatalf("server setup: %v", err)
	}

	handler := server.Setup(serverCfg)

	const listenAddr = "127.0.0.1:18456"
	httpServer := &http.Server{
		Addr:    listenAddr,
		Handler: handler,
	}

	go func() {
		_ = httpServer.ListenAndServe()
	}()

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = httpServer.Shutdown(ctx)
	})

	u := launcher.New().NoSandbox(true).MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()
	t.Cleanup(func() {
		_ = browser.Close()
	})

	return &testApp{
		BaseURL:      "http://" + listenAddr,
		Browser:      browser,
		RoomCodeRepo: serverCfg.RoomCodeRepo,
		UserRepo:     serverCfg.UserRepo,
	}
}
