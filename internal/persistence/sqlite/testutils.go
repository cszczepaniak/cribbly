//go:build !prod

package sqlite

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func NewInMemoryForTest(t testing.TB) *sql.DB {
	t.Helper()

	db, err := New("file::memory:")
	require.NoError(t, err)

	return db
}
