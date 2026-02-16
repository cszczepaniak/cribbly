package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Database struct {
	db *sql.DB
}

type DatabaseFactory func() (*sql.DB, error)

func New(factory DatabaseFactory) (Database, error) {
	db, err := factory()
	if err != nil {
		return Database{}, err
	}
	return Database{db: db}, nil
}

type stmtHandler interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

func getDB(ctx context.Context, db *sql.DB) stmtHandler {
	tx, ok := getTx(ctx)
	if ok {
		fmt.Println("got a tx")
		return tx
	}
	fmt.Println("did not got a tx")
	return db
}

func (db Database) ExecContext(ctx context.Context, stmt string, args ...any) (sql.Result, error) {
	return getDB(ctx, db.db).ExecContext(ctx, stmt, args...)
}

func (db Database) ExecVoid(ctx context.Context, stmt string, args ...any) error {
	_, err := getDB(ctx, db.db).ExecContext(ctx, stmt, args...)
	return err
}

func (db Database) ExecOne(ctx context.Context, stmt string, args ...any) error {
	res, err := getDB(ctx, db.db).ExecContext(ctx, stmt, args...)
	if err != nil {
		return err
	}

	if n, err := res.RowsAffected(); err != nil {
		return err
	} else if n != 1 {
		return errors.New("expected exactly one row to be affected")
	}

	return nil
}

func (db Database) QueryContext(ctx context.Context, stmt string, args ...any) (*sql.Rows, error) {
	return getDB(ctx, db.db).QueryContext(ctx, stmt, args...)
}

func (db Database) QueryRowContext(ctx context.Context, stmt string, args ...any) *sql.Row {
	return getDB(ctx, db.db).QueryRowContext(ctx, stmt, args...)
}

func (db Database) WithTx(ctx context.Context, fn func(context.Context) error) error {
	return withTx(ctx, db.db, fn)
}
