package players

import (
	"context"
	"database/sql"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/column"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/filter"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/formatter"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/table"
)

type Player struct {
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
	stmt, err := s.b.CreateTable("Players").
		IfNotExists().
		Columns(
			column.VarChar("ID", 32).PrimaryKey(),
			column.VarChar("Name", 255),
		).
		Build()
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, stmt)
	return err
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

func (s Service) Create(ctx context.Context, id, name string) error {
	_, err := s.b.InsertIntoTable("Players").
		Fields(
			"ID", "Name",
		).
		Values(
			id,
			name,
		).
		ExecContext(ctx, s.db)
	return err
}
