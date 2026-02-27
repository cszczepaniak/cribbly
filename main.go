package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/alexedwards/argon2id"

	"github.com/cszczepaniak/cribbly/internal/config"
	"github.com/cszczepaniak/cribbly/internal/notifier"
	"github.com/cszczepaniak/cribbly/internal/persistence/database"
	"github.com/cszczepaniak/cribbly/internal/persistence/divisions"
	"github.com/cszczepaniak/cribbly/internal/persistence/games"
	"github.com/cszczepaniak/cribbly/internal/persistence/players"
	"github.com/cszczepaniak/cribbly/internal/persistence/roomcodes"
	"github.com/cszczepaniak/cribbly/internal/persistence/teams"
	"github.com/cszczepaniak/cribbly/internal/persistence/users"
	"github.com/cszczepaniak/cribbly/internal/server"
)

func main() {
	err := runMain()
	if err != nil {
		log.Fatal(err)
	}
}

func runMain() error {
	cfg := config.Config{}
	err := config.Load(&cfg)
	if err != nil {
		return err
	}

	if cfg.DSN == "" {
		cfg.DSN = "file:data/db.sqlite"
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer cancel()

	db, err := database.NewSQLiteDB(cfg.DSN)
	if err != nil {
		return err
	}

	scoreUpdateNotifier := &notifier.Notifier{}
	tournamentNotifier := &notifier.Notifier{}

	playerRepo := players.NewRepository(db)
	err = playerRepo.Init(ctx)
	if err != nil {
		return err
	}

	teamRepo := teams.NewRepository(db)
	err = teamRepo.Init(ctx)
	if err != nil {
		return err
	}

	divisionRepo := divisions.NewRepository(db)
	err = divisionRepo.Init(ctx)
	if err != nil {
		return err
	}

	gameRepo := games.NewRepository(db, scoreUpdateNotifier)
	err = gameRepo.Init(ctx)
	if err != nil {
		return err
	}

	roomCodeRepo := roomcodes.NewRepository(db)
	err = roomCodeRepo.Init(ctx)
	if err != nil {
		return err
	}

	userRepo := users.NewRepository(db)
	err = userRepo.Init(ctx)
	if err != nil {
		return err
	}

	if cfg.SeedUser.Username != "" && cfg.SeedUser.Password != "" {
		passwordHash, err := argon2id.CreateHash(cfg.SeedUser.Password, argon2id.DefaultParams)
		if err != nil {
			return err
		}
		err = userRepo.CreateUser(context.Background(), cfg.SeedUser.Username, passwordHash)
		if err != nil {
			log.Println("could not seed user:", err)
		}
	}

	scfg := server.Config{
		Transactor:          database.NewTransactor(db),
		PlayerRepo:          playerRepo,
		TeamRepo:            teamRepo,
		DivisionRepo:        divisionRepo,
		GameRepo:            gameRepo,
		UserRepo:            userRepo,
		RoomCodeRepo:        roomCodeRepo,
		ScoreUpdateNotifier: scoreUpdateNotifier,
		TournamentNotifier:  tournamentNotifier,
		IsProd:              cfg.Environment == "production",
	}

	s := server.Setup(scfg)

	errCh := make(chan error)
	go func() {
		errCh <- http.ListenAndServe(":8080", s)
	}()

	select {
	case <-ctx.Done():
		log.Println("server shutting down...")
		return nil
	case err := <-errCh:
		return err
	}
}
