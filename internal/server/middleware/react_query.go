package middleware

import (
	"net/http"
	"strings"
)

// ReactQueryMiddleware serves the embedded React index.html when the query
// parameter react=true is present. Otherwise it delegates to the next handler.
// Static assets under /app/ and legacy files under /public/ always pass through.
func ReactQueryMiddleware(indexHTML func() []byte, isProd bool, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			next.ServeHTTP(w, r)
			return
		}
		if r.URL.Query().Get("react") != "true" {
			next.ServeHTTP(w, r)
			return
		}
		if !shouldServeReactShell(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}
		if !isProd {
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if r.Method == http.MethodHead {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(indexHTML())
	})
}

func shouldServeReactShell(path string) bool {
	switch {
	case strings.HasPrefix(path, "/app/"):
		// Bundled JS/CSS and other Vite output live here; never serve HTML for these paths.
		return false
	case strings.HasPrefix(path, "/public/"):
		// Legacy static files (CSS, images, etc.).
		return false
	case strings.HasPrefix(path, "/api/"):
		// Future JSON / Connect RPC endpoints.
		return false
	case strings.HasSuffix(path, "/stream"):
		// Legacy SSE endpoints (Datastar).
		return false
	default:
		return true
	}
}
