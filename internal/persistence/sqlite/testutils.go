//go:build !prod

package sqlite

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
	_ "github.com/tursodatabase/go-libsql"
)

func NewInMemory(t testing.TB) *sql.DB {
	t.Helper()

	db, err := sql.Open("libsql", ":memory:")
	require.NoError(t, err)

	return db
}
