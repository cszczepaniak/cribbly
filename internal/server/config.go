package server

import (
	"github.com/cszczepaniak/cribbly/internal/persistence/divisions"
	"github.com/cszczepaniak/cribbly/internal/persistence/games"
	"github.com/cszczepaniak/cribbly/internal/persistence/players"
	"github.com/cszczepaniak/cribbly/internal/persistence/teams"
	"github.com/cszczepaniak/cribbly/internal/persistence/users"
)

type Config struct {
	PlayerService   players.Service
	TeamService     teams.Service
	DivisionService divisions.Service
	GameService     games.Service
	UserService     users.Service
}
