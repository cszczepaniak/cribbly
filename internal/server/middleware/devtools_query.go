package middleware

import (
	"context"
	"net/http"
)

type devToolsQueryKey struct{}

// DevToolsFromQuery reports whether the request included ?dev=true. Use with IsAdmin; it does not
// bypass the admin check.
func DevToolsFromQuery(ctx context.Context) bool {
	v, ok := ctx.Value(devToolsQueryKey{}).(bool)
	return ok && v
}

func DevToolsQueryMiddleware() middleware {
	return func(next handler) handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			ctx := r.Context()
			if r.URL.Query().Get("dev") == "true" {
				ctx = context.WithValue(ctx, devToolsQueryKey{}, true)
			}
			return next(w, r.WithContext(ctx))
		}
	}
}
