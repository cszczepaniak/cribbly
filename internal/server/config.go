package server

import (
	"github.com/cszczepaniak/cribbly/internal/notifier"
	"github.com/cszczepaniak/cribbly/internal/persistence/database"
	"github.com/cszczepaniak/cribbly/internal/persistence/divisions"
	"github.com/cszczepaniak/cribbly/internal/persistence/games"
	"github.com/cszczepaniak/cribbly/internal/persistence/players"
	"github.com/cszczepaniak/cribbly/internal/persistence/teams"
	"github.com/cszczepaniak/cribbly/internal/persistence/users"
	divisionservice "github.com/cszczepaniak/cribbly/internal/service/divisions"
	teamservice "github.com/cszczepaniak/cribbly/internal/service/teams"
)

type Config struct {
	Transactor          database.Transactor
	PlayerRepo          players.Repository
	TeamRepo            teams.Repository
	DivisionRepo        divisions.Repository
	GameRepo            games.Repository
	UserRepo            users.Repository
	ScoreUpdateNotifier *notifier.Notifier
	TournamentNotifier  *notifier.Notifier
	IsProd              bool
}

func (cfg Config) TeamService() teamservice.Service {
	return teamservice.New(cfg.Transactor, cfg.PlayerRepo, cfg.TeamRepo)
}

func (cfg Config) DivisionService() divisionservice.Service {
	return divisionservice.New(cfg.Transactor, cfg.TeamRepo, cfg.DivisionRepo)
}
