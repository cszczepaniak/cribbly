package sqlite

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"

	"github.com/cszczepaniak/cribbly/internal/persistence/database"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

func New(dsn string) (database.Database, error) {
	filePath, ok := strings.CutPrefix(dsn, "file:")
	if ok && filePath != ":memory:" {
		err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
		if err != nil {
			return database.Database{}, err
		}
	}

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return database.Database{}, err
	}

	return database.New(db), nil
}
