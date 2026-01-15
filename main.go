package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/alexedwards/argon2id"
	"github.com/cszczepaniak/cribbly/internal/notifier"
	"github.com/cszczepaniak/cribbly/internal/persistence/divisions"
	"github.com/cszczepaniak/cribbly/internal/persistence/games"
	"github.com/cszczepaniak/cribbly/internal/persistence/players"
	"github.com/cszczepaniak/cribbly/internal/persistence/sqlite"
	"github.com/cszczepaniak/cribbly/internal/persistence/teams"
	"github.com/cszczepaniak/cribbly/internal/persistence/users"
	"github.com/cszczepaniak/cribbly/internal/server"
)

func main() {
	dbSource := flag.String("db", "file", "The source to use for the sqlite database. Valid options are 'memory' or 'file'.")
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer cancel()

	var db *sql.DB
	var err error
	switch *dbSource {
	case "file":
		db, err = sqlite.NewFromFile("data/db.sqlite")
	case "memory":
		db, err = sqlite.NewInMemory()
	default:
		err = fmt.Errorf("unknown -db value: %q", *dbSource)
	}
	if err != nil {
		log.Fatal(err)
	}

	scoreUpdateNotifier := &notifier.Notifier{}

	playerService := players.NewService(db)
	err = playerService.Init(ctx)
	if err != nil {
		log.Fatal(err)
	}

	teamService := teams.NewService(db)
	err = teamService.Init(ctx)
	if err != nil {
		log.Fatal(err)
	}

	divisionService := divisions.NewService(db)
	err = divisionService.Init(ctx)
	if err != nil {
		log.Fatal(err)
	}

	gameService := games.NewService(db, scoreUpdateNotifier)
	err = gameService.Init(ctx)
	if err != nil {
		log.Fatal(err)
	}

	userService := users.NewService(db)
	err = userService.Init(ctx)
	if err != nil {
		log.Fatal(err)
	}

	seedUser, seedPass := os.Getenv("SEED_USER"), os.Getenv("SEED_PASSWORD")
	if seedUser != "" && seedPass != "" {
		passwordHash, err := argon2id.CreateHash(seedPass, argon2id.DefaultParams)
		if err != nil {
			log.Fatal(err)
		}
		err = userService.CreateUser(context.Background(), seedUser, passwordHash)
		if err != nil {
			log.Println("could not seed user:", err)
		}
	}

	cfg := server.Config{
		PlayerService:       playerService,
		TeamService:         teamService,
		DivisionService:     divisionService,
		GameService:         gameService,
		UserService:         userService,
		ScoreUpdateNotifier: scoreUpdateNotifier,
	}

	s := server.Setup(cfg)

	errCh := make(chan error)
	go func() {
		errCh <- http.ListenAndServe(":8080", s)
	}()

	select {
	case <-ctx.Done():
		log.Println("server shutting down...")
	case err := <-errCh:
		if err != nil {
			log.Fatal(err)
		}
	}
}
