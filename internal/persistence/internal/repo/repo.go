package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/formatter"
)

type baseDB interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

type DB interface {
	baseDB
	ExecVoid(context.Context, string, ...any) error
	ExecN(context.Context, string, ...any) (int64, error)
	ExecOne(context.Context, string, ...any) error
}

type dbWithExtensions struct {
	baseDB
}

func (db *dbWithExtensions) ExecVoid(ctx context.Context, stmt string, args ...any) error {
	_, err := db.baseDB.ExecContext(ctx, stmt, args...)
	return err
}

func (db *dbWithExtensions) ExecN(ctx context.Context, stmt string, args ...any) (int64, error) {
	res, err := db.baseDB.ExecContext(ctx, stmt, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (db *dbWithExtensions) ExecOne(ctx context.Context, stmt string, args ...any) error {
	n, err := db.ExecN(ctx, stmt, args...)
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("expected one row to be affected")
	}
	return nil
}

type Tx interface {
	DB
	Commit() error
	Rollback() error
}

type txWithExtensions struct {
	tx *sql.Tx
	DB
}

func (tx *txWithExtensions) Commit() error {
	return tx.tx.Commit()
}

func (tx *txWithExtensions) Rollback() error {
	return tx.tx.Commit()
}

type Base struct {
	db      *sql.DB
	DB      DB
	Builder *sqlbuilder.Builder
}

func NewBase(db *sql.DB) Base {
	return Base{
		db:      db,
		DB:      &dbWithExtensions{baseDB: db},
		Builder: sqlbuilder.New(formatter.Sqlite{}),
	}
}

func (b Base) WithTx(tx *sql.Tx) Base {
	return Base{
		DB:      &dbWithExtensions{baseDB: tx},
		Builder: b.Builder,
	}
}

func (b Base) BeginTx(ctx context.Context) (Tx, error) {
	tx, err := b.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &txWithExtensions{
		tx: tx,
		DB: &dbWithExtensions{baseDB: tx},
	}, nil
}

func (b Base) EnsureTx(ctx context.Context) (Tx, func(), error) {
	if tx, ok := b.DB.(*txWithExtensions); ok {
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
