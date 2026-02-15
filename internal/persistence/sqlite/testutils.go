//go:build !prod

package sqlite

import (
	"testing"

	"github.com/cszczepaniak/cribbly/internal/assert"
	"github.com/cszczepaniak/cribbly/internal/persistence"
)

func NewInMemoryForTest(t testing.TB) persistence.Database {
	t.Helper()

	db, err := New("file::memory:")
	assert.NoError(t, err)

	return persistence.NewDatabase(db)
}
