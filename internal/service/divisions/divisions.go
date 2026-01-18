package divisions

import (
	"context"

	"github.com/cszczepaniak/cribbly/internal/persistence"
	"github.com/cszczepaniak/cribbly/internal/persistence/divisions"
	"github.com/cszczepaniak/cribbly/internal/persistence/teams"
)

type Division struct {
	ID    string
	Name  string
	Size  int
	Teams []teams.Team
}

type Service struct {
	txer         persistence.Transactor
	teamRepo     teams.Repository
	divisionRepo divisions.Repository
}

func New(txer persistence.Transactor, teamRepo teams.Repository, divisionRepo divisions.Repository) Service {
	return Service{
		txer:         txer,
		teamRepo:     teamRepo,
		divisionRepo: divisionRepo,
	}
}

func (s Service) Get(ctx context.Context, id string) (Division, error) {
	div, err := s.divisionRepo.Get(ctx, id)
	if err != nil {
		return Division{}, err
	}

	teams, err := s.teamRepo.GetForDivision(ctx, id)
	if err != nil {
		return Division{}, err
	}

	return Division{
		ID:    div.ID,
		Name:  div.Name,
		Size:  div.Size,
		Teams: teams,
	}, nil
}
