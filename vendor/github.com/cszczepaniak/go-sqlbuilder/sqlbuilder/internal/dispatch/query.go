package dispatch

import (
	"context"
	"database/sql"
)

type Queryer interface {
	Query(stmt string, args ...any) (*sql.Rows, error)
}

func Query(b builder, q Queryer) (*sql.Rows, error) {
	res, err := b.Build()
	if err != nil {
		return nil, err
	}

	return q.Query(res.Stmt, res.Args...)
}

type QueryCtxer interface {
	QueryContext(ctx context.Context, stmt string, args ...any) (*sql.Rows, error)
}

func QueryContext(ctx context.Context, b builder, q QueryCtxer) (*sql.Rows, error) {
	res, err := b.Build()
	if err != nil {
		return nil, err
	}

	return q.QueryContext(ctx, res.Stmt, res.Args...)
}

type RowQueryer interface {
	QueryRow(stmt string, args ...any) *sql.Row
}

func QueryRow(b builder, q RowQueryer) (*sql.Row, error) {
	res, err := b.Build()
	if err != nil {
		return nil, err
	}

	return q.QueryRow(res.Stmt, res.Args...), nil
}

type RowQueryCtxer interface {
	QueryRowContext(ctx context.Context, stmt string, args ...any) *sql.Row
}

func QueryRowContext(ctx context.Context, b builder, q RowQueryCtxer) (*sql.Row, error) {
	res, err := b.Build()
	if err != nil {
		return nil, err
	}

	return q.QueryRowContext(ctx, res.Stmt, res.Args...), nil
}
