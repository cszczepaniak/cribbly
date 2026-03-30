package middleware

import (
	"context"
	"crypto/subtle"
	"net/http"
)

// DevAdminHeader is sent by local dev tooling so the backend can treat the request as an admin
// without a session. Honored only when the server is not in production and CRIBBLY_DEV_ADMIN_SECRET
// is set to a non-empty value; the header value must match that secret exactly.
const DevAdminHeader = "X-Cribbly-Dev-Admin"

type devAdminBypassKey struct{}

// WithDevAdminBypassIfHeader returns a clone of r whose context marks the request as a dev admin
// when the header matches secret. No-op when isProd or secret is empty.
func WithDevAdminBypassIfHeader(r *http.Request, secret string, isProd bool) *http.Request {
	if isProd || secret == "" {
		return r
	}
	if devAdminSecretsMatch(r.Header.Get(DevAdminHeader), secret) {
		return r.WithContext(context.WithValue(r.Context(), devAdminBypassKey{}, true))
	}
	return r
}

// DevAdminBypassMiddleware applies WithDevAdminBypassIfHeader for router handlers. Disabled when
// isProd is true or secret is empty.
func DevAdminBypassMiddleware(secret string, isProd bool) middleware {
	if isProd || secret == "" {
		return func(next handler) handler { return next }
	}
	return func(next handler) handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			return next(w, WithDevAdminBypassIfHeader(r, secret, false))
		}
	}
}

func devAdminSecretsMatch(got, want string) bool {
	if len(got) != len(want) {
		return false
	}
	// Constant-time compare only when lengths match (subtle requires equal length).
	return subtle.ConstantTimeCompare([]byte(got), []byte(want)) == 1
}

func isDevAdminBypass(ctx context.Context) bool {
	v, ok := ctx.Value(devAdminBypassKey{}).(bool)
	return ok && v
}

// WithDevAdminContext returns a context that satisfies IsAdmin via the dev-bypass branch.
// Intended for tests (simulates a valid X-Cribbly-Dev-Admin header without HTTP).
func WithDevAdminContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, devAdminBypassKey{}, true)
}
