package sqlite

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func New(dsn string) (*sql.DB, error) {
	filePath, ok := strings.CutPrefix(dsn, "file:")
	if ok && filePath != ":memory:" {
		err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	return sql.Open("libsql", dsn)
}
