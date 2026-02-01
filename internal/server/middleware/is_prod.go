package middleware

import (
	"context"
	"net/http"
)

type isProdKey struct{}

func IsProd(ctx context.Context) bool {
	is, ok := ctx.Value(isProdKey{}).(bool)
	return ok && is
}

func IsProdMiddleware(isProd bool) middleware {
	return func(next handler) handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			return next(
				w,
				r.WithContext(context.WithValue(
					r.Context(),
					isProdKey{},
					isProd,
				)),
			)
		}
	}
}
