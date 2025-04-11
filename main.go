package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/cszczepaniak/cribbly/internal/server"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer cancel()

	s := server.Setup()

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
