package server

import (
	"github.com/cszczepaniak/cribbly/internal/persistence/divisions"
	"github.com/cszczepaniak/cribbly/internal/persistence/players"
	"github.com/cszczepaniak/cribbly/internal/persistence/teams"
)

type Config struct {
	PlayerService   players.Service
	TeamService     teams.Service
	DivisionService divisions.Service
}
