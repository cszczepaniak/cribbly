package sqlite

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

func Factory(dsn string) func() (*sql.DB, error) {
	return func() (*sql.DB, error) {
		filePath, ok := strings.CutPrefix(dsn, "file:")
		if ok && filePath != ":memory:" {
			err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
			if err != nil {
				return nil, err
			}
		}

		return sql.Open("sqlite3", dsn)
	}
}

func MemoryFactory() func() (*sql.DB, error) {
	return Factory("file::memory:")
}
