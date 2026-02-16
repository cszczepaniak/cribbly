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

type txKey struct{}

func ctxWithTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func getTx(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(*sql.Tx)
	return tx, ok
}

func withTx(
	ctx context.Context,
	db *sql.DB,
	fn func(context.Context) error,
) error {
	var tx *sql.Tx
	shouldCommit := false

	if t, ok := getTx(ctx); ok {
		// Context already has a transaction
		tx = t
	} else {
		shouldCommit = true
		ctxWithCancel, cancel := context.WithCancel(ctx)
		defer cancel()

		var err error
		tx, err = db.BeginTx(ctxWithCancel, nil)
		if err != nil {
			return err
		}

		ctx = ctxWithTx(ctxWithCancel, tx)
	}

	err := fn(ctx)
	if err != nil {
		return err
	}

	if shouldCommit {
		return tx.Commit()
	}

	return nil
}

// Transactor knows how to start a transaction but can't do anything else.
type Transactor struct {
	db Database
}

func NewTransactor(db Database) Transactor {
	return Transactor{
		db: db,
	}
}

func (t Transactor) WithTx(ctx context.Context, fn func(context.Context) error) error {
	return t.db.WithTx(ctx, fn)
}
