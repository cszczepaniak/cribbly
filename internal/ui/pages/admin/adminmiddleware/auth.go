package adminmiddleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/cszczepaniak/cribbly/internal/persistence/users"
)

type sessionKey struct{}

func GetSession(ctx context.Context) users.Session {
	return ctx.Value(sessionKey{}).(users.Session)
}

// TODO: share this type somewhere (can't be server pkg because that'd cause an import cycle)
type handler = func(http.ResponseWriter, *http.Request) error

func AuthenticationMiddleware(userRepo users.Repository) func(next handler) handler {
	return func(next handler) handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			cookie, err := r.Cookie("session")
			if err != nil {
				if errors.Is(err, http.ErrNoCookie) {
					http.Redirect(w, r, "/admin/login", http.StatusTemporaryRedirect)
					return nil
				}

				return fmt.Errorf("get session cookie: %w", err)
			}

			clearCookieAndRedirect := func() {
				http.SetCookie(w, &http.Cookie{
					Name:    "session",
					Value:   "",
					MaxAge:  -1,
					Expires: time.Now().Add(-time.Minute),
				})
				http.Redirect(w, r, "/admin/login", http.StatusTemporaryRedirect)
			}

			// TODO: in-memory caching of sessions to avoid a DB call for each request
			sesh, err := userRepo.GetSession(r.Context(), cookie.Value)
			if err != nil {
				if errors.Is(err, users.ErrSessionExpired) {
					clearCookieAndRedirect()
					return nil
				}
				return err
			}

			if sesh.Expired() {
				clearCookieAndRedirect()
				return nil
			}

			ctx := context.WithValue(r.Context(), sessionKey{}, sesh)
			r = r.WithContext(ctx)

			return next(w, r)
		}
	}
}
