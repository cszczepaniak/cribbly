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
	_, ok := ctx.Value(sessionKey{}).(users.Session)
	return ok
}

func GetSession(ctx context.Context) users.Session {
	return ctx.Value(sessionKey{}).(users.Session)
}

// TODO: share this type somewhere (can't be server pkg because that'd cause an import cycle)
type handler = func(http.ResponseWriter, *http.Request) error

// AuthenticationMiddleware extracts the user's session (if any) and adds it to the request's Go
// context. It always forwards the request to the next handler regardless of whether the user was
// successfully authenticated.
func AuthenticationMiddleware(userRepo users.Repository) func(next handler) handler {
	return func(next handler) handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			cookie, err := r.Cookie("session")
			if err != nil {
				if errors.Is(err, http.ErrNoCookie) {
					return next(w, r)
				}

				return fmt.Errorf("get session cookie: %w", err)
			}

			// TODO: in-memory caching of sessions to avoid a DB call for each request
			sesh, err := userRepo.GetSession(r.Context(), cookie.Value)
			if err != nil {
				if errors.Is(err, users.ErrSessionExpired) {
					return next(w, r)
				}

				return err
			}

			if !sesh.Expired() {
				ctx := context.WithValue(r.Context(), sessionKey{}, sesh)
				r = r.WithContext(ctx)
			}

			return next(w, r)
		}
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
