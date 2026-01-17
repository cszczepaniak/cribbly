package persistence

import (
	"context"
	"database/sql"
)

type Transactor struct {
	db *sql.DB
}

func (t Transactor) Begin(ctx context.Context) (*sql.Tx, error) {
	return t.db.BeginTx(ctx, nil)
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
