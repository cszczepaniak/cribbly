package divisions

import (
	"context"
	"database/sql"

	"github.com/cszczepaniak/cribbly/internal/persistence/internal/repo"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/column"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/filter"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/table"
	"github.com/google/uuid"
)

type Division struct {
	ID   string
	Name string
	Size int
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
	_, err := s.Builder.CreateTable("Divisions").
		IfNotExists().
		Columns(
			column.VarChar("ID", 36).PrimaryKey(),
			column.VarChar("Name", 255),
			column.TinyInt("Size"), // 4 or 6
		).
		ExecContext(ctx, s.DB)
	return err
}

func (s Repository) Create(ctx context.Context) (Division, error) {
	division := Division{
		ID:   uuid.NewString(),
		Name: "Unnamed Division",
		Size: 4,
	}

	_, err := s.Builder.InsertIntoTable("Divisions").
		Fields("ID", "Name", "Size").
		Values(division.ID, division.Name, 4).
		ExecContext(ctx, s.DB)
	if err != nil {
		return Division{}, err
	}

	return division, nil
}

func (s Repository) Delete(ctx context.Context, id string) error {
	_, err := s.Builder.DeleteFromTable("Divisions").
		Where(filter.Equals("ID", id)).
		ExecContext(ctx, s.DB)
	return err
}

func (s Repository) Rename(ctx context.Context, id, newName string) error {
	_, err := s.Builder.UpdateTable("Divisions").
		SetFieldTo("Name", newName).
		Where(filter.Equals("ID", id)).
		ExecContext(ctx, s.DB)
	return err
}

func (s Repository) UpdateSize(ctx context.Context, id string, size int) error {
	_, err := s.Builder.UpdateTable("Divisions").
		SetFieldTo("Size", size).
		Where(filter.Equals("ID", id)).
		ExecContext(ctx, s.DB)
	return err
}

func (s Repository) Get(ctx context.Context, id string) (Division, error) {
	row, err := s.Builder.SelectFrom(table.Named("Divisions")).
		Columns("ID", "Name", "Size").
		Where(filter.Equals("ID", id)).
		QueryRowContext(ctx, s.DB)
	if err != nil {
		return Division{}, err
	}

	var division Division
	err = row.Scan(&division.ID, &division.Name, &division.Size)
	if err != nil {
		return Division{}, err
	}

	return division, nil
}

func (s Repository) GetAll(ctx context.Context) ([]Division, error) {
	rows, err := s.Builder.SelectFrom(table.Named("Divisions")).
		Columns("ID", "Name", "Size").
		QueryContext(ctx, s.DB)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var divisions []Division
	for rows.Next() {
		var d Division
		err := rows.Scan(&d.ID, &d.Name, &d.Size)
		if err != nil {
			return nil, err
		}

		divisions = append(divisions, d)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return divisions, nil
}
