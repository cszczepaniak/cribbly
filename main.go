package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/alexedwards/argon2id"

	"github.com/cszczepaniak/cribbly/internal/config"
	"github.com/cszczepaniak/cribbly/internal/persistence/database"
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

	serverCfg, err := server.SetupFromDB(ctx, db, cfg.Environment == "production")
	if err != nil {
		return err
	}

	if cfg.SeedUser.Username != "" && cfg.SeedUser.Password != "" {
		passwordHash, err := argon2id.CreateHash(cfg.SeedUser.Password, argon2id.DefaultParams)
		if err != nil {
			return err
		}
		err = serverCfg.UserRepo.CreateUser(context.Background(), cfg.SeedUser.Username, passwordHash)
		if err != nil {
			log.Println("could not seed user:", err)
		}
	}

	s := server.Setup(serverCfg)

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
