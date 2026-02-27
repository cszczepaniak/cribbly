package server

import (
	"context"

	"github.com/cszczepaniak/cribbly/internal/notifier"
	"github.com/cszczepaniak/cribbly/internal/persistence/database"
	"github.com/cszczepaniak/cribbly/internal/persistence/divisions"
	"github.com/cszczepaniak/cribbly/internal/persistence/games"
	"github.com/cszczepaniak/cribbly/internal/persistence/players"
	"github.com/cszczepaniak/cribbly/internal/persistence/roomcodes"
	"github.com/cszczepaniak/cribbly/internal/persistence/teams"
	"github.com/cszczepaniak/cribbly/internal/persistence/users"
)

// SetupFromDB creates all repositories from the given database, initializes them,
// and returns a server Config and HTTP handler. Callers can use the handler to
// serve traffic and the config to access repos (e.g. for seeding or tests).
func SetupFromDB(ctx context.Context, db database.Database, isProd bool) (Config, error) {
	scoreUpdateNotifier := &notifier.Notifier{}
	tournamentNotifier := &notifier.Notifier{}

	playerRepo := players.NewRepository(db)
	if err := playerRepo.Init(ctx); err != nil {
		return Config{}, err
	}

	teamRepo := teams.NewRepository(db)
	if err := teamRepo.Init(ctx); err != nil {
		return Config{}, err
	}

	divisionRepo := divisions.NewRepository(db)
	if err := divisionRepo.Init(ctx); err != nil {
		return Config{}, err
	}

	gameRepo := games.NewRepository(db, scoreUpdateNotifier)
	if err := gameRepo.Init(ctx); err != nil {
		return Config{}, err
	}

	roomCodeRepo := roomcodes.NewRepository(db)
	if err := roomCodeRepo.Init(ctx); err != nil {
		return Config{}, err
	}

	userRepo := users.NewRepository(db)
	if err := userRepo.Init(ctx); err != nil {
		return Config{}, err
	}

	cfg := Config{
		Transactor:          database.NewTransactor(db),
		PlayerRepo:          playerRepo,
		TeamRepo:            teamRepo,
		DivisionRepo:        divisionRepo,
		GameRepo:            gameRepo,
		UserRepo:            userRepo,
		RoomCodeRepo:        roomCodeRepo,
		ScoreUpdateNotifier: scoreUpdateNotifier,
		TournamentNotifier:  tournamentNotifier,
		IsProd:              isProd,
	}

	return cfg, nil
}
