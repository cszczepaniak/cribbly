package database

import (
	"database/sql"
	"testing"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/ncruces/go-sqlite3/vfs/memdb"

	"github.com/cszczepaniak/cribbly/internal/assert"
)

func NewInMemory(t *testing.T) Database {
	// See https://pkg.go.dev/github.com/ncruces/go-sqlite3/vfs/memdb#example-package
	memdb.Create("test.db", nil)
	db, err := sql.Open("sqlite3", "file:/test.db?vfs=memdb")
	assert.NoError(t, err)
	return New(db)
}
