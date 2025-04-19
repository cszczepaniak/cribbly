package players

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/column"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/filter"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/formatter"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/table"
	"github.com/google/uuid"
)

var (
	ErrPlayerAlreadyOnATeam = errors.New("player was already assigned to a team")
)

type Player struct {
	ID     string
	Name   string
	TeamID string
}

type Service struct {
	db *sql.DB
	b  *sqlbuilder.Builder
}

func NewService(db *sql.DB) Service {
	return Service{
		db: db,
		b:  sqlbuilder.New(formatter.Sqlite{}),
	}
}

func (s Service) Init(ctx context.Context) error {
	_, err := s.b.CreateTable("Players").
		IfNotExists().
		Columns(
			column.VarChar("ID", 36).PrimaryKey(),
			column.VarChar("Name", 255),
			column.VarChar("TeamID", 36).DefaultNull(),
		).
		Exec(s.db)
	return err
}

func (s Service) GetAll(ctx context.Context) ([]Player, error) {
	rows, err := s.b.SelectFrom(table.Named("Players")).
		Columns("ID", "Name").
		QueryContext(ctx, s.db)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	players := make([]Player, 0)
	for rows.Next() {
		var p Player
		err := rows.Scan(&p.ID, &p.Name)
		if err != nil {
			return nil, err
		}

		players = append(players, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return players, nil
}

// GetFreeAgents returns all players who are not assigned to a team.
func (s Service) GetFreeAgents(ctx context.Context) ([]Player, error) {
	rows, err := s.b.SelectFrom(table.Named("Players")).
		Columns("ID", "Name", "TeamID").
		Where(filter.IsNull("TeamID")).
		QueryContext(ctx, s.db)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var players []Player
	for rows.Next() {
		var p Player
		var teamID sql.Null[string]
		err := rows.Scan(&p.ID, &p.Name, &teamID)
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

// AssignToTeam assigns the given player to the given team.
func (s Service) AssignToTeam(ctx context.Context, p Player, teamID string) error {
	res, err := s.b.UpdateTable("Players").
		SetFieldTo("TeamID", teamID).
		WhereAll(
			filter.Equals("ID", p.ID),
			filter.IsNull("TeamID"),
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
		return ErrPlayerAlreadyOnATeam
	case 1:
	// This is the expected case.
	default:
		// This should never happen, but we'll include it as a sanity check.
		return fmt.Errorf("unexpected: assigning player to team updated %d players", rowsAffected)
	}

	return nil
}

func (s Service) Get(ctx context.Context, ids ...string) ([]Player, error) {
	rows, err := s.b.SelectFrom(table.Named("Players")).
		Columns("ID", "Name").
		Where(
			filter.In("ID", ids...),
		).
		QueryContext(ctx, s.db)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	players := make([]Player, 0, len(ids))
	for rows.Next() {
		var p Player
		err := rows.Scan(&p.ID, &p.Name)
		if err != nil {
			return nil, err
		}

		players = append(players, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return players, nil
}

func (s Service) GetOne(ctx context.Context, id string) (Player, error) {
	players, err := s.Get(ctx, id)
	if err != nil {
		return Player{}, err
	}

	return players[0], nil
}

func (s Service) Create(ctx context.Context, name string) (string, error) {
	id := uuid.NewString()

	_, err := s.b.InsertIntoTable("Players").
		Fields(
			"ID", "Name",
		).
		Values(
			id,
			name,
		).
		ExecContext(ctx, s.db)
	if err != nil {
		return "", err
	}

	return id, nil
}
