//go:build !prod

package sqlite

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func NewInMemoryForTest(t testing.TB) *sql.DB {
	t.Helper()

	db, err := NewInMemory()
	require.NoError(t, err)

	return db
}
