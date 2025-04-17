package server

import "github.com/cszczepaniak/cribbly/internal/persistence/players"

type Config struct {
	PlayerService players.Service
}
