package teams

import (
	"context"
	"database/sql"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/column"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/filter"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/formatter"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/table"
	"github.com/google/uuid"
)

type Team struct {
	ID   string
	Name string
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
	_, err := s.b.CreateTable("Teams").
		IfNotExists().
		Columns(
			column.VarChar("ID", 36).PrimaryKey(),
			column.VarChar("Name", 255),
		).
		Exec(s.db)
	return err
}

func (s Service) Create(ctx context.Context) (Team, error) {
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

func (s Service) Delete(ctx context.Context, id string) error {
	_, err := s.b.DeleteFromTable("Teams").
		Where(filter.Equals("ID", id)).
		ExecContext(ctx, s.db)
	return err
}

func (s Service) Rename(ctx context.Context, id, newName string) error {
	_, err := s.b.UpdateTable("Teams").
		SetFieldTo("Name", newName).
		Where(filter.Equals("ID", id)).
		ExecContext(ctx, s.db)
	return err
}

func (s Service) Get(ctx context.Context, id string) (Team, error) {
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

func (s Service) GetAll(ctx context.Context) ([]Team, error) {
	rows, err := s.b.SelectFrom(table.Named("Teams")).
		Columns("ID", "Name").
		QueryContext(ctx, s.db)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []Team
	for rows.Next() {
		var t Team
		err := rows.Scan(&t.ID, &t.Name)
		if err != nil {
			return nil, err
		}

		teams = append(teams, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return teams, nil
}
