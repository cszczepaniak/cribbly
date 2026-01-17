package server

import (
	"github.com/cszczepaniak/cribbly/internal/notifier"
	"github.com/cszczepaniak/cribbly/internal/persistence"
	"github.com/cszczepaniak/cribbly/internal/persistence/divisions"
	"github.com/cszczepaniak/cribbly/internal/persistence/games"
	"github.com/cszczepaniak/cribbly/internal/persistence/players"
	"github.com/cszczepaniak/cribbly/internal/persistence/teams"
	"github.com/cszczepaniak/cribbly/internal/persistence/users"
	teamservice "github.com/cszczepaniak/cribbly/internal/service/teams"
)

type Config struct {
	Transactor          persistence.Transactor
	PlayerRepo          players.Repository
	TeamRepo            teams.Repository
	DivisionRepo        divisions.Repository
	GameRepo            games.Repository
	UserRepo            users.Repository
	ScoreUpdateNotifier *notifier.Notifier
}

func (cfg Config) TeamService() teamservice.Service {
	return teamservice.New(cfg.Transactor, cfg.PlayerRepo, cfg.TeamRepo)
}
