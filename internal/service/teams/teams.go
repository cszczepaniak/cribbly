package teams

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/cszczepaniak/cribbly/internal/moreiter"
	"github.com/cszczepaniak/cribbly/internal/persistence/database"
	"github.com/cszczepaniak/cribbly/internal/persistence/players"
	"github.com/cszczepaniak/cribbly/internal/persistence/teams"
)

type Service struct {
	txer       database.Transactor
	playerRepo players.Repository
	teamRepo   teams.Repository
}

func New(txer database.Transactor, playerRepo players.Repository, teamRepo teams.Repository) Service {
	return Service{
		txer:       txer,
		playerRepo: playerRepo,
		teamRepo:   teamRepo,
	}
}

type Team struct {
	ID      string
	Name    string
	Players []players.Player
}

func (s Service) GetTeam(ctx context.Context, teamID string) (Team, error) {
	team, err := s.teamRepo.Get(ctx, teamID)
	if err != nil {
		return Team{}, err
	}

	players, err := s.playerRepo.GetForTeam(ctx, team.ID)
	if err != nil {
		return Team{}, err
	}

	return Team{
		ID:      teamID,
		Name:    team.Name,
		Players: players,
	}, nil
}

func (s Service) CreateTeam(ctx context.Context) (Team, error) {
	team, err := s.teamRepo.Create(ctx, "Unnamed Team")
	if err != nil {
		return Team{}, err
	}

	return Team{
		ID:   team.ID,
		Name: team.Name,
	}, nil
}

func (s Service) DeleteTeam(ctx context.Context, id string) error {
	return s.txer.WithTx(ctx, func(ctx context.Context) error {
		playersOnTeam, err := s.playerRepo.GetForTeam(ctx, id)
		if err != nil {
			return err
		}

		err = s.playerRepo.UnassignFromTeam(
			ctx,
			id,
			moreiter.Map(slices.Values(playersOnTeam), func(p players.Player) string {
				return p.ID
			}),
		)
		if err != nil {
			return err
		}

		return s.teamRepo.Delete(ctx, id)
	})
}

func (s Service) AssignPlayerToTeam(ctx context.Context, playerID, teamID string) (Team, error) {
	var team Team
	err := s.txer.WithTx(ctx, func(ctx context.Context) error {
		var err error
		team, err = s.GetTeam(ctx, teamID)
		if err != nil {
			return err
		}

		if len(team.Players) > 1 {
			return errors.New("too many players on team")
		}

		err = s.playerRepo.AssignToTeam(ctx, playerID, teamID)
		if err != nil {
			return err
		}

		playersOnTeam, err := s.playerRepo.GetForTeam(ctx, teamID)
		if err != nil {
			return err
		}
		team.Players = playersOnTeam

		switch len(playersOnTeam) {
		case 0:
			// Should be unreachable
			return errors.New("no players on team")
		case 1:
			team.Name = getPlayerName(playersOnTeam[0])
		case 2:
			team.Name = getPlayerName(playersOnTeam[0]) + " / " + getPlayerName(playersOnTeam[1])
		default:
			// Should be unreachable
			return errors.New("too many players on team")
		}

		return s.teamRepo.Rename(ctx, teamID, team.Name)
	})
	if err != nil {
		return Team{}, err
	}

	return team, nil
}

func (s Service) UnassignPlayerFromTeam(ctx context.Context, playerID, teamID string) (Team, error) {
	var team Team
	err := s.txer.WithTx(ctx, func(ctx context.Context) error {
		var err error
		team, err = s.GetTeam(ctx, teamID)
		if err != nil {
			return err
		}

		err = s.playerRepo.UnassignFromTeam(ctx, teamID, moreiter.Of(playerID))
		if err != nil {
			return err
		}

		playersOnTeam, err := s.playerRepo.GetForTeam(ctx, teamID)
		if err != nil {
			return err
		}
		team.Players = playersOnTeam

		switch len(playersOnTeam) {
		case 0:
			team.Name = "Unnamed Team"
		case 1:
			// The team name is based on the other player's name now.
			team.Name = getPlayerName(playersOnTeam[0])
		default:
			return errors.New("too many players on team")
		}

		return s.teamRepo.Rename(ctx, teamID, team.Name)
	})
	if err != nil {
		return Team{}, err
	}

	return team, nil
}

func getPlayerName(p players.Player) string {
	return fmt.Sprintf("%c %s", p.FirstName[0], p.LastName)
}
