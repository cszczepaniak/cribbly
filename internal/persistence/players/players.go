package players

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/cszczepaniak/cribbly/internal/persistence/internal/repo"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/column"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/filter"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/sel"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/table"
	"github.com/google/uuid"
)

var (
	ErrPlayerAlreadyOnATeam = errors.New("player was already assigned to a team")
)

type Player struct {
	ID        string
	FirstName string
	LastName  string
	TeamID    string
}

func (p Player) Name() string {
	return p.FirstName + " " + p.LastName
}

type Repository struct {
	repo.Base
}

func NewRepository(db *sql.DB) Repository {
	return Repository{
		Base: repo.NewBase(db),
	}
}

func (s Repository) Init(ctx context.Context) error {
	_, err := s.Builder.CreateTable("Players").
		IfNotExists().
		Columns(
			column.VarChar("ID", 36).PrimaryKey(),
			column.VarChar("FirstName", 255),
			column.VarChar("LastName", 255),
			column.VarChar("TeamID", 36).DefaultNull(),
		).
		ExecContext(ctx, s.DB)
	return err
}

func (s Repository) GetAll(ctx context.Context) ([]Player, error) {
	return scanPlayers(
		s.selectPlayers().
			QueryContext(ctx, s.DB),
	)
}

// GetFreeAgents returns all players who are not assigned to a team.
func (s Repository) GetFreeAgents(ctx context.Context) ([]Player, error) {
	return scanPlayers(
		s.selectPlayers().
			Where(filter.IsNull("TeamID")).
			QueryContext(ctx, s.DB),
	)
}

// GetForTeam returns the players assigned to the given team.
func (s Repository) GetForTeam(ctx context.Context, teamID string) ([]Player, error) {
	return scanPlayers(
		s.selectPlayers().
			Where(filter.Equals("TeamID", teamID)).
			QueryContext(ctx, s.DB),
	)
}

// AssignToTeam assigns the given player to the given team.
func (s Repository) AssignToTeam(ctx context.Context, playerID, teamID string) error {
	res, err := s.Builder.UpdateTable("Players").
		SetFieldTo("TeamID", teamID).
		WhereAll(
			filter.Equals("ID", playerID),
			filter.IsNull("TeamID"),
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
		return ErrPlayerAlreadyOnATeam
	case 1:
	// This is the expected case.
	default:
		// This should never happen, but we'll include it as a sanity check.
		return fmt.Errorf("unexpected: assigning player to team updated %d players", rowsAffected)
	}

	return nil
}

func (s Repository) UnassignFromTeam(ctx context.Context, id string) error {
	_, err := s.Builder.UpdateTable("Players").
		SetFieldToNull("TeamID").
		Where(filter.Equals("ID", id)).
		ExecContext(ctx, s.DB)
	return err
}

func (s Repository) Create(ctx context.Context, firstName, lastName string) (string, error) {
	id := uuid.NewString()

	if firstName == "" || lastName == "" {
		return "", errors.New("must have a first and last name")
	}

	_, err := s.Builder.InsertIntoTable("Players").
		Fields("ID", "FirstName", "LastName").
		Values(id, firstName, lastName).
		ExecContext(ctx, s.DB)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s Repository) Delete(ctx context.Context, id string) error {
	_, err := s.Builder.DeleteFromTable("Players").
		Where(filter.Equals("ID", id)).
		ExecContext(ctx, s.DB)
	return err
}

func (s Repository) selectPlayers() *sel.Builder {
	return s.Builder.SelectFrom(table.Named("Players")).
		Columns("ID", "FirstName", "LastName", "TeamID")
}

func scanPlayers(rows *sql.Rows, err error) ([]Player, error) {
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []Player
	for rows.Next() {
		var p Player
		var teamID sql.Null[string]
		err := rows.Scan(&p.ID, &p.FirstName, &p.LastName, &teamID)
		if err != nil {
			return nil, err
		}

		// If the team ID was null, V will be an empty string.
		p.TeamID = teamID.V

		players = append(players, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return players, nil
}
