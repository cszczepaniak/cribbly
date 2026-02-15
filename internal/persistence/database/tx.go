package database

import (
	"context"
	"database/sql"
)

type txKey struct{}

func ctxWithTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func getTx(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(*sql.Tx)
	return tx, ok
}

type Transactor struct {
	db *sql.DB
}

func NewTransactor(db Database) Transactor {
	return Transactor{
		db: db.db,
	}
}

func (t Transactor) WithTx(ctx context.Context, fn func(context.Context) error) error {
	return withTx(ctx, t.db, fn)
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
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		var err error
		tx, err = db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		ctx = ctxWithTx(ctx, tx)
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
