package database

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/cszczepaniak/cribbly/internal/assert"
)

func TestWithTransaction(t *testing.T) {
	db, err := NewSQLiteDB("file:db.sqlite")
	assert.NoError(t, err)

	t.Cleanup(func() {
		_ = os.Remove("db.sqlite")
	})

	_, err = db.db.Conn(t.Context())
	assert.NoError(t, err)

	assert.NoError(t, db.ExecVoid(t.Context(), `CREATE TABLE Test (
		A INT
	)`))

	// Insert some things transactionally
	err = db.WithTx(t.Context(), func(ctx context.Context) error {
		assert.NoError(t, db.ExecVoid(ctx, `INSERT INTO Test (A) VALUES (?), (?)`, 1, 2))
		assert.NoError(t, db.ExecVoid(ctx, `INSERT INTO Test (A) VALUES (?), (?)`, 3, 4))
		return nil
	})
	assert.NoError(t, err)

	// Should be committed
	var val int
	err = db.QueryRowContext(t.Context(), `SELECT SUM(A) FROM Test`).Scan(&val)
	assert.NoError(t, err)
	assert.Equal(t, 10, val)

	// Insert some things transactionally, but cause a rollback (with an error)
	err = db.WithTx(t.Context(), func(ctx context.Context) error {
		assert.NoError(t, db.ExecVoid(ctx, `INSERT INTO Test (A) VALUES (?), (?)`, 1, 2))
		assert.NoError(t, db.ExecVoid(ctx, `INSERT INTO Test (A) VALUES (?), (?)`, 3, 4))
		return errors.New("foo")
	})
	assert.Error(t, err)

	// Shouldn't be committed
	err = db.QueryRowContext(t.Context(), `SELECT SUM(A) FROM Test`).Scan(&val)
	assert.NoError(t, err)
	assert.Equal(t, 10, val)

	// Insert some things transactionally, but cause a rollback (with a panic)
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected a panic but got none")
			}
		}()

		_ = db.WithTx(t.Context(), func(ctx context.Context) error {
			assert.NoError(t, db.ExecVoid(ctx, `INSERT INTO Test (A) VALUES (?), (?)`, 1, 2))
			assert.NoError(t, db.ExecVoid(ctx, `INSERT INTO Test (A) VALUES (?), (?)`, 3, 4))
			panic("ahh")
		})
	}()

	// Shouldn't be committed
	err = db.QueryRowContext(t.Context(), `SELECT SUM(A) FROM Test`).Scan(&val)
	assert.NoError(t, err)
	assert.Equal(t, 10, val)

	// It should work when a transaction already exists
	tx, err := db.db.BeginTx(t.Context(), nil)
	assert.NoError(t, err)

	err = db.WithTx(ctxWithTx(t.Context(), tx), func(ctx context.Context) error {
		assert.NoError(t, db.ExecVoid(ctx, `INSERT INTO Test (A) VALUES (?), (?)`, 5, 6))
		return nil
	})
	assert.NoError(t, err)
	assert.NoError(t, tx.Commit())

	err = db.QueryRowContext(t.Context(), `SELECT SUM(A) FROM Test`).Scan(&val)
	assert.NoError(t, err)
	assert.Equal(t, 21, val)
}
