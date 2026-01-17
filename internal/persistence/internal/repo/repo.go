package repo

import (
	"context"
	"database/sql"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/formatter"
)

type DB interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

type Tx interface {
	DB
	Commit() error
	Rollback() error
}

type Base struct {
	db      *sql.DB
	DB      DB
	Builder *sqlbuilder.Builder
}

func NewBase(db *sql.DB) Base {
	return Base{
		db:      db,
		DB:      db,
		Builder: sqlbuilder.New(formatter.Sqlite{}),
	}
}

func (b Base) WithTx(tx *sql.Tx) Base {
	return Base{
		DB:      tx,
		Builder: b.Builder,
	}
}

func (b Base) BeginTx(ctx context.Context) (Tx, error) {
	return b.db.BeginTx(ctx, nil)
}

func (b Base) EnsureTx(ctx context.Context) (Tx, func(), error) {
	if tx, ok := b.DB.(*sql.Tx); ok {
		return tx, func() {}, nil
	}

	ctx, cancel := context.WithCancel(ctx)
	tx, err := b.BeginTx(ctx)
	if err != nil {
		cancel()
		return nil, nil, err
	}

	return tx, cancel, nil
}
