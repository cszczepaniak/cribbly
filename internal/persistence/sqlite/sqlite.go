package sqlite

import (
	"database/sql"
	"os"
	"path/filepath"
)

func NewInMemory() (*sql.DB, error) {
	return sql.Open("libsql", "file::memory:")
}

func NewFromFile(path string) (*sql.DB, error) {
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return nil, err
	}

	return sql.Open("libsql", "file:"+path)
}
