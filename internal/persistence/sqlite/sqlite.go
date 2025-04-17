package sqlite

import (
	"database/sql"
)

func NewInMemory() (*sql.DB, error) {
	return sql.Open("libsql", "file::memory:")
}
