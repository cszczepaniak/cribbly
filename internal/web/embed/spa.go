package webembed

import (
	"io/fs"
	"net/http"
	"path"
	"strings"
)

// SPAHandler serves static files from root, falling back to index.html for
// client-side routes (HTML5 history).
func SPAHandler(root fs.FS) http.Handler {
	fileServer := http.FileServer(http.FS(root))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upath := path.Clean(r.URL.Path)
		if !strings.HasPrefix(upath, "/") {
			upath = "/" + upath
		}
		name := strings.TrimPrefix(upath, "/")
		if name == "" || name == "." {
			name = "index.html"
		}
		f, err := root.Open(name)
		if err == nil {
			_ = f.Close()
			fileServer.ServeHTTP(w, r)
			return
		}
		data, err := fs.ReadFile(root, "index.html")
		if err != nil {
			http.Error(w, "index.html not found", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(data)
	})
}
