package divisions

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

type Division struct {
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
	_, err := s.b.CreateTable("Divisions").
		IfNotExists().
		Columns(
			column.VarChar("ID", 36).PrimaryKey(),
			column.VarChar("Name", 255),
		).
		Exec(s.db)
	return err
}

func (s Service) Create(ctx context.Context) (Division, error) {
	division := Division{
		ID:   uuid.NewString(),
		Name: "Unnamed Division",
	}

	_, err := s.b.InsertIntoTable("Divisions").
		Fields("ID", "Name").
		Values(division.ID, division.Name).
		ExecContext(ctx, s.db)
	if err != nil {
		return Division{}, err
	}

	return division, nil
}

func (s Service) Delete(ctx context.Context, id string) error {
	_, err := s.b.DeleteFromTable("Divisions").
		Where(filter.Equals("ID", id)).
		ExecContext(ctx, s.db)
	return err
}

func (s Service) Rename(ctx context.Context, id, newName string) error {
	_, err := s.b.UpdateTable("Divisions").
		SetFieldTo("Name", newName).
		Where(filter.Equals("ID", id)).
		ExecContext(ctx, s.db)
	return err
}

func (s Service) Get(ctx context.Context, id string) (Division, error) {
	row, err := s.b.SelectFrom(table.Named("Divisions")).
		Columns("ID", "Name").
		Where(filter.Equals("ID", id)).
		QueryRowContext(ctx, s.db)
	if err != nil {
		return Division{}, err
	}

	var division Division
	err = row.Scan(&division.ID, &division.Name)
	if err != nil {
		return Division{}, err
	}

	return division, nil
}

func (s Service) GetAll(ctx context.Context) ([]Division, error) {
	rows, err := s.b.SelectFrom(table.Named("Divisions")).
		Columns("ID", "Name").
		QueryContext(ctx, s.db)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var divisions []Division
	for rows.Next() {
		var t Division
		err := rows.Scan(&t.ID, &t.Name)
		if err != nil {
			return nil, err
		}

		divisions = append(divisions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return divisions, nil
}
