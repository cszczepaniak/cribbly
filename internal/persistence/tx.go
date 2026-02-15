package persistence

import (
	"context"
	"database/sql"
	"errors"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(db *sql.DB) Database {
	return Database{db: db}
}

type txKey struct{}

func getDB(ctx context.Context, db *sql.DB) DB {
	tx, ok := ctx.Value(txKey{}).(*sql.Tx)
	if ok {
		return tx
	}
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

type Transactor struct {
	db *sql.DB
}

func NewTransactor(db *sql.DB) Transactor {
	return Transactor{
		db: db,
	}
}

func (t Transactor) WithTx(ctx context.Context, fn func(context.Context) error) error {
	return withTx(ctx, t.db, fn)
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

func withTx(
	ctx context.Context,
	db *sql.DB,
	fn func(context.Context) error,
) error {
	var tx *sql.Tx
	shouldCommit := false

	if t, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		// Context already has a transaction
		tx = t
	} else {
		shouldCommit = true
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		var err error
		tx, err = db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		ctx = context.WithValue(ctx, txKey{}, tx)
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
