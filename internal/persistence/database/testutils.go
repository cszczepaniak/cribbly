package database

import (
	"testing"

	"github.com/cszczepaniak/cribbly/internal/assert"
)

func NewInMemory(t *testing.T) Database {
	db, err := NewSQLiteDB("file::memory:")
	assert.NoError(t, err)
	return db
}
