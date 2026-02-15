package players

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"iter"
	"slices"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/column"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/filter"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/formatter"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/sel"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/table"
	"github.com/google/uuid"

	"github.com/cszczepaniak/cribbly/internal/persistence"
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
	db persistence.DB
	b  *sqlbuilder.Builder
}

func NewRepository(db persistence.Database) Repository {
	return Repository{
		db: db,
		b:  sqlbuilder.New(formatter.Sqlite{}),
	}
}

func (s Repository) Init(ctx context.Context) error {
	_, err := s.b.CreateTable("Players").
		IfNotExists().
		Columns(
			column.VarChar("ID", 36).PrimaryKey(),
			column.VarChar("FirstName", 255),
			column.VarChar("LastName", 255),
			column.VarChar("TeamID", 36).DefaultNull(),
		).
		ExecContext(ctx, s.db)
	return err
}

func (s Repository) GetAll(ctx context.Context) ([]Player, error) {
	return scanPlayers(
		s.selectPlayers().
			QueryContext(ctx, s.db),
	)
}

func (s Repository) Get(ctx context.Context, id string) (Player, error) {
	row, err := s.selectPlayers().
		WhereAll(
			filter.Equals("ID", id),
		).
		QueryRowContext(ctx, s.db)
	if err != nil {
		return Player{}, err
	}

	return scanPlayer(row)
}

func (s Repository) GetFreeAgent(ctx context.Context, id string) (Player, error) {
	row, err := s.selectPlayers().
		WhereAll(
			filter.Equals("ID", id),
			filter.IsNull("TeamID"),
		).
		QueryRowContext(ctx, s.db)
	if err != nil {
		return Player{}, err
	}

	return scanPlayer(row)
}

// GetFreeAgents returns all players who are not assigned to a team.
func (s Repository) GetFreeAgents(ctx context.Context) ([]Player, error) {
	return scanPlayers(
		s.selectPlayers().
			Where(filter.IsNull("TeamID")).
			QueryContext(ctx, s.db),
	)
}

// GetForTeam returns the players assigned to the given team.
func (s Repository) GetForTeam(ctx context.Context, teamID string) ([]Player, error) {
	return scanPlayers(
		s.selectPlayers().
			Where(filter.Equals("TeamID", teamID)).
			QueryContext(ctx, s.db),
	)
}

// AssignToTeam assigns the given player to the given team.
func (s Repository) AssignToTeam(ctx context.Context, playerID, teamID string) error {
	res, err := s.b.UpdateTable("Players").
		SetFieldTo("TeamID", teamID).
		WhereAll(
			filter.Equals("ID", playerID),
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

func (s Repository) UnassignFromTeam(ctx context.Context, teamID string, playerIDs iter.Seq[string]) error {
	players := slices.Collect(playerIDs)
	if len(players) == 0 {
		return nil
	}

	res, err := s.b.UpdateTable("Players").
		SetFieldToNull("TeamID").
		WhereAll(
			filter.In("ID", players...),
			filter.Equals("TeamID", teamID),
		).
		ExecContext(ctx, s.db)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if int(n) != len(players) {
		// TODO: maybe we should just be okay with removing players that aren't on the team (which
		// is a no-op).
		return errors.New("all given players must be on the team")
	}

	return nil
}

func (s Repository) Create(ctx context.Context, firstName, lastName string) (string, error) {
	id := uuid.NewString()

	if firstName == "" || lastName == "" {
		return "", errors.New("must have a first and last name")
	}

	_, err := s.b.InsertIntoTable("Players").
		Fields("ID", "FirstName", "LastName").
		Values(id, firstName, lastName).
		ExecContext(ctx, s.db)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s Repository) Delete(ctx context.Context, id string) error {
	_, err := s.b.DeleteFromTable("Players").
		Where(filter.Equals("ID", id)).
		ExecContext(ctx, s.db)
	return err
}

func (s Repository) selectPlayers() *sel.Builder {
	return s.b.SelectFrom(table.Named("Players")).
		Columns("ID", "FirstName", "LastName", "TeamID")
}

func scanPlayers(rows *sql.Rows, err error) ([]Player, error) {
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []Player
	for rows.Next() {
		p, err := scanPlayer(rows)
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

func scanPlayer(scanner interface{ Scan(...any) error }) (Player, error) {
	var p Player
	var teamID sql.Null[string]
	err := scanner.Scan(&p.ID, &p.FirstName, &p.LastName, &teamID)
	if err != nil {
		return Player{}, err
	}

	// If the team ID was null, return empty string.
	p.TeamID = teamID.V
	return p, nil
}
