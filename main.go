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

	"github.com/cszczepaniak/cribbly/internal/persistence/players"
	"github.com/cszczepaniak/cribbly/internal/persistence/sqlite"
	"github.com/cszczepaniak/cribbly/internal/persistence/teams"
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

	cfg := server.Config{
		PlayerService: playerService,
		TeamService:   teamService,
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
