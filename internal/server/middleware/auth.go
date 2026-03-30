package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/cszczepaniak/cribbly/internal/persistence/users"
)

type sessionKey struct{}

func IsAdmin(ctx context.Context) bool {
	if isDevAdminBypass(ctx) {
		return true
	}
	_, ok := ctx.Value(sessionKey{}).(users.Session)
	return ok
}

func GetSession(ctx context.Context) users.Session {
	return ctx.Value(sessionKey{}).(users.Session)
}

// TODO: share this type somewhere (can't be server pkg because that'd cause an import cycle)
type (
	handler    = func(http.ResponseWriter, *http.Request) error
	middleware = func(handler) handler
)

// requestWithSessionIfAny returns r with an admin session in context when a valid session cookie
// is present; otherwise returns r unchanged. Missing or expired sessions are not errors.
func requestWithSessionIfAny(r *http.Request, userRepo users.Repository) (*http.Request, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return r, nil
		}

		return nil, fmt.Errorf("get session cookie: %w", err)
	}

	// TODO: in-memory caching of sessions to avoid a DB call for each request
	sesh, err := userRepo.GetSession(r.Context(), cookie.Value)
	if err != nil {
		if errors.Is(err, users.ErrSessionExpired) {
			return r, nil
		}

		return nil, err
	}

	if sesh.Expired() {
		return r, nil
	}

	ctx := context.WithValue(r.Context(), sessionKey{}, sesh)
	return r.WithContext(ctx), nil
}

// AuthenticationMiddleware extracts the user's session (if any) and adds it to the request's Go
// context. It always forwards the request to the next handler regardless of whether the user was
// successfully authenticated.
func AuthenticationMiddleware(userRepo users.Repository) middleware {
	return func(next handler) handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			r2, err := requestWithSessionIfAny(r, userRepo)
			if err != nil {
				return err
			}

			return next(w, r2)
		}
	}
}

// ConnectSessionMiddleware attaches a valid session cookie to the request context (same rules as
// AuthenticationMiddleware). Use on Connect RPC mounts that are not routed through NewRouter.
func ConnectSessionMiddleware(userRepo users.Repository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r2, err := requestWithSessionIfAny(r, userRepo)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			next.ServeHTTP(w, r2)
		})
	}
}

func ErrorIfNotAdmin() func(next handler) handler {
	return func(next handler) handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			if !IsAdmin(r.Context()) {
				return errors.New("must be an admin")
			}
			return next(w, r)
		}
	}
}

func RedirectToLoginIfNotAdmin() func(next handler) handler {
	return func(next handler) handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			if !IsAdmin(r.Context()) {
				http.Redirect(w, r, "/admin/login", http.StatusTemporaryRedirect)
				return nil
			}
			return next(w, r)
		}
	}
}
