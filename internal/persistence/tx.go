package persistence

import (
	"context"
	"database/sql"
)

type Transactor struct {
	db *sql.DB
}

func NewTransactor(db *sql.DB) Transactor {
	return Transactor{
		db: db,
	}
}

func (t Transactor) Begin(ctx context.Context) (*sql.Tx, func(), error) {
	ctx, cancel := context.WithCancel(ctx)
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		cancel()
		return nil, nil, err
	}
	return tx, cancel, nil
}

type Execer interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}

type Queryer interface {
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

type DB interface {
	Execer
	Queryer
}
