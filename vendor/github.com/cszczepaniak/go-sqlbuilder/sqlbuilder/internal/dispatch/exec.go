package dispatch

import (
	"context"
	"database/sql"
)

type Execer interface {
	Exec(stmt string, args ...any) (sql.Result, error)
}

func Exec(b builder, e Execer) (sql.Result, error) {
	res, err := b.Build()
	if err != nil {
		return nil, err
	}

	return e.Exec(res.Stmt, res.Args...)
}

type ExecCtxer interface {
	ExecContext(ctx context.Context, stmt string, args ...any) (sql.Result, error)
}

func ExecContext(ctx context.Context, b builder, e ExecCtxer) (sql.Result, error) {
	res, err := b.Build()
	if err != nil {
		return nil, err
	}

	return e.ExecContext(ctx, res.Stmt, res.Args...)
}
