package teams

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/column"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/filter"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/sel"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/table"
	"github.com/google/uuid"

	"github.com/cszczepaniak/cribbly/internal/persistence/internal/repo"
)

var (
	ErrTeamAlreadyInDivision = errors.New("team was already assigned to a division")
)

type Team struct {
	ID         string
	Name       string
	DivisionID string
}

type Repository struct {
	repo.Base
}

func NewRepository(db *sql.DB) Repository {
	return Repository{
		Base: repo.NewBase(db),
	}
}

func (s Repository) WithTx(tx *sql.Tx) Repository {
	s.Base = s.Base.WithTx(tx)
	return s
}

func (s Repository) Init(ctx context.Context) error {
	_, err := s.Builder.CreateTable("Teams").
		IfNotExists().
		Columns(
			column.VarChar("ID", 36).PrimaryKey(),
			column.VarChar("Name", 255),
			column.VarChar("DivisionID", 36),
		).
		ExecContext(ctx, s.DB)
	return err
}

func (s Repository) Create(ctx context.Context, name string) (Team, error) {
	team := Team{
		ID:   uuid.NewString(),
		Name: name,
	}

	_, err := s.Builder.InsertIntoTable("Teams").
		Fields("ID", "Name").
		Values(team.ID, team.Name).
		ExecContext(ctx, s.DB)
	if err != nil {
		return Team{}, err
	}

	return team, nil
}

func (s Repository) Delete(ctx context.Context, id string) error {
	_, err := s.Builder.DeleteFromTable("Teams").
		Where(filter.Equals("ID", id)).
		ExecContext(ctx, s.DB)
	return err
}

func (s Repository) DeleteAll(ctx context.Context) error {
	_, err := s.Builder.DeleteFromTable("Teams").
		ExecContext(ctx, s.DB)
	return err
}

func (s Repository) Rename(ctx context.Context, id, newName string) error {
	_, err := s.Builder.UpdateTable("Teams").
		SetFieldTo("Name", newName).
		Where(filter.Equals("ID", id)).
		ExecContext(ctx, s.DB)
	return err
}

func (s Repository) Get(ctx context.Context, id string) (Team, error) {
	row, err := s.Builder.SelectFrom(table.Named("Teams")).
		Columns("ID", "Name", "DivisionID").
		Where(filter.Equals("ID", id)).
		QueryRowContext(ctx, s.DB)
	if err != nil {
		return Team{}, err
	}

	var team Team
	var divisionID sql.Null[string]
	err = row.Scan(&team.ID, &team.Name, &divisionID)
	if err != nil {
		return Team{}, err
	}
	team.DivisionID = divisionID.V

	return team, nil
}

func (s Repository) GetAll(ctx context.Context) ([]Team, error) {
	return scanTeams(
		s.selectTeams().
			QueryContext(ctx, s.DB),
	)
}

// GetWithoutDivision returns all teams that are not assigned to a division.
func (s Repository) GetWithoutDivision(ctx context.Context) ([]Team, error) {
	return scanTeams(
		s.selectTeams().
			Where(filter.IsNull("DivisionID")).
			QueryContext(ctx, s.DB),
	)
}

// GetForDivision returns the teams in the given division.
func (s Repository) GetForDivision(ctx context.Context, divisionID string) ([]Team, error) {
	return scanTeams(
		s.selectTeams().
			Where(filter.Equals("DivisionID", divisionID)).
			QueryContext(ctx, s.DB),
	)
}

// AssignToDivision assigns the given team to the given division.
func (s Repository) AssignToDivision(ctx context.Context, teamID, divisionID string) error {
	res, err := s.Builder.UpdateTable("Teams").
		SetFieldTo("DivisionID", divisionID).
		WhereAll(
			filter.Equals("ID", teamID),
			filter.IsNull("DivisionID"),
		).
		ExecContext(ctx, s.DB)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	switch rowsAffected {
	case 0:
		return ErrTeamAlreadyInDivision
	case 1:
	// This is the expected case.
	default:
		// This should never happen, but we'll include it as a sanity check.
		return fmt.Errorf("unexpected: assigning team to division updated %d teams", rowsAffected)
	}

	return nil
}

func (s Repository) selectTeams() *sel.Builder {
	return s.Builder.SelectFrom(table.Named("Teams")).
		Columns("ID", "Name", "DivisionID")
}

func (s Repository) UnassignFromDivision(ctx context.Context, id string) error {
	_, err := s.Builder.UpdateTable("Teams").
		SetFieldToNull("DivisionID").
		Where(filter.Equals("ID", id)).
		ExecContext(ctx, s.DB)
	return err
}

func scanTeams(rows *sql.Rows, err error) ([]Team, error) {
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []Team
	for rows.Next() {
		var t Team
		var divisionID sql.Null[string]
		err := rows.Scan(&t.ID, &t.Name, &divisionID)
		if err != nil {
			return nil, err
		}

		// If the team ID was null, V will be an empty string.
		t.DivisionID = divisionID.V

		teams = append(teams, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return teams, nil
}
