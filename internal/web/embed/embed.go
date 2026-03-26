package webembed

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed all:dist
var distFS embed.FS

// StaticHandler serves JS/CSS and other files from the Vite build under /app/.
func StaticHandler() http.Handler {
	fsys, err := fs.Sub(distFS, "dist")
	if err != nil {
		panic(err)
	}
	return http.StripPrefix("/app", http.FileServer(http.FS(fsys)))
}

// MustReadIndexHTML returns the built index.html for the React shell.
func MustReadIndexHTML() []byte {
	fsys, err := fs.Sub(distFS, "dist")
	if err != nil {
		panic(err)
	}
	data, err := fs.ReadFile(fsys, "index.html")
	if err != nil {
		panic(err)
	}
	return data
}
