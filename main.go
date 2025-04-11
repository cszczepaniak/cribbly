package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/cszczepaniak/cribbly/internal/persistence/players"
	"github.com/cszczepaniak/cribbly/internal/persistence/sqlite"
	"github.com/cszczepaniak/cribbly/internal/server"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer cancel()

	// TODO: for now we just always set up an in-memory sqlite instance. We could add configuration
	// to control this, i.e. to use a file for more persistent local development or our prod
	// database.
	db, err := sqlite.NewInMemory()
	if err != nil {
		log.Fatal(err)
	}

	playerService := players.NewService(db)
	err = playerService.Init(ctx)
	if err != nil {
		log.Fatal(err)
	}

	cfg := server.Config{
		PlayerService: playerService,
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
