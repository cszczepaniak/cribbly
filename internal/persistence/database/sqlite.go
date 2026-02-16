package database

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

func NewSQLiteDB(dsn string) (Database, error) {
	filePath, ok := strings.CutPrefix(dsn, "file:")
	if ok && filePath != ":memory:" {
		err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
		if err != nil {
			return Database{}, err
		}
	}

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return Database{}, err
	}

	return New(db), nil
}
