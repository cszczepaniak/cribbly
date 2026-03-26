package webembed

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed all:dist
var distFS embed.FS

// Handler serves the SPA and static assets under /app/.
func Handler() http.Handler {
	fsys, err := fs.Sub(distFS, "dist")
	if err != nil {
		panic(err)
	}
	return http.StripPrefix("/app", SPAHandler(fsys))
}
