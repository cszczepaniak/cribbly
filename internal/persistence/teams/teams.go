package teams

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/column"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/filter"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/formatter"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/sel"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/table"
	"github.com/google/uuid"
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
	db *sql.DB
	b  *sqlbuilder.Builder
}

func NewRepository(db *sql.DB) Repository {
	return Repository{
		db: db,
		b:  sqlbuilder.New(formatter.Sqlite{}),
	}
}

func (s Repository) Init(ctx context.Context) error {
	_, err := s.b.CreateTable("Teams").
		IfNotExists().
		Columns(
			column.VarChar("ID", 36).PrimaryKey(),
			column.VarChar("Name", 255),
			column.VarChar("DivisionID", 36),
		).
		Exec(s.db)
	return err
}

func (s Repository) Create(ctx context.Context) (Team, error) {
	team := Team{
		ID:   uuid.NewString(),
		Name: "Unnamed Team",
	}

	_, err := s.b.InsertIntoTable("Teams").
		Fields("ID", "Name").
		Values(team.ID, team.Name).
		ExecContext(ctx, s.db)
	if err != nil {
		return Team{}, err
	}

	return team, nil
}

func (s Repository) Delete(ctx context.Context, id string) error {
	_, err := s.b.DeleteFromTable("Teams").
		Where(filter.Equals("ID", id)).
		ExecContext(ctx, s.db)
	return err
}

func (s Repository) DeleteAll(ctx context.Context) error {
	_, err := s.b.DeleteFromTable("Teams").
		ExecContext(ctx, s.db)
	return err
}

func (s Repository) Rename(ctx context.Context, id, newName string) error {
	_, err := s.b.UpdateTable("Teams").
		SetFieldTo("Name", newName).
		Where(filter.Equals("ID", id)).
		ExecContext(ctx, s.db)
	return err
}

func (s Repository) Get(ctx context.Context, id string) (Team, error) {
	row, err := s.b.SelectFrom(table.Named("Teams")).
		Columns("ID", "Name").
		Where(filter.Equals("ID", id)).
		QueryRowContext(ctx, s.db)
	if err != nil {
		return Team{}, err
	}

	var team Team
	err = row.Scan(&team.ID, &team.Name)
	if err != nil {
		return Team{}, err
	}

	return team, nil
}

func (s Repository) GetAll(ctx context.Context) ([]Team, error) {
	return scanTeams(
		s.selectTeams().
			QueryContext(ctx, s.db),
	)
}

// GetWithoutDivision returns all teams that are not assigned to a division.
func (s Repository) GetWithoutDivision(ctx context.Context) ([]Team, error) {
	return scanTeams(
		s.selectTeams().
			Where(filter.IsNull("DivisionID")).
			QueryContext(ctx, s.db),
	)
}

// GetForDivision returns the teams in the given division.
func (s Repository) GetForDivision(ctx context.Context, divisionID string) ([]Team, error) {
	return scanTeams(
		s.selectTeams().
			Where(filter.Equals("DivisionID", divisionID)).
			QueryContext(ctx, s.db),
	)
}

// AssignToDivision assigns the given team to the given division.
func (s Repository) AssignToDivision(ctx context.Context, teamID, divisionID string) error {
	res, err := s.b.UpdateTable("Teams").
		SetFieldTo("DivisionID", divisionID).
		WhereAll(
			filter.Equals("ID", teamID),
			filter.IsNull("DivisionID"),
		).
		ExecContext(ctx, s.db)
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
	return s.b.SelectFrom(table.Named("Teams")).
		Columns("ID", "Name", "DivisionID")
}

func (s Repository) UnassignFromDivision(ctx context.Context, id string) error {
	_, err := s.b.UpdateTable("Teams").
		SetFieldToNull("DivisionID").
		Where(filter.Equals("ID", id)).
		ExecContext(ctx, s.db)
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
