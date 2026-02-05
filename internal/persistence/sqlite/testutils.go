//go:build !prod

package sqlite

import (
	"database/sql"
	"testing"

	"github.com/cszczepaniak/cribbly/internal/assert"
)

func NewInMemoryForTest(t testing.TB) *sql.DB {
	t.Helper()

	db, err := New("file::memory:")
	assert.NoError(t, err)

	return db
}
