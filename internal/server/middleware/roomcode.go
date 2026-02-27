package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/cszczepaniak/cribbly/internal/persistence/roomcodes"
)

type roomAccessKey struct{}

func HasRoomAccess(ctx context.Context) bool {
	has, ok := ctx.Value(roomAccessKey{}).(bool)
	return ok && has
}

// RoomCodeMiddleware checks for a valid room-code cookie and:
//   - stores whether the current request has room access in the context
//   - redirects non-admin users without room access to the home page
func RoomCodeMiddleware(repo roomcodes.Repository) middleware {
	return func(next handler) handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			ctx := r.Context()

			// Admins always have access.
			hasAccess := IsAdmin(ctx)

			if !hasAccess {
				if cookie, err := r.Cookie("room_code"); err == nil {
					ok, err := repo.Validate(ctx, cookie.Value)
					if err != nil {
						return err
					}
					hasAccess = ok
				} else if !errors.Is(err, http.ErrNoCookie) {
					// Only propagate non-"no cookie" errors.
					return err
				}
			}

			ctx = context.WithValue(ctx, roomAccessKey{}, hasAccess)
			r = r.WithContext(ctx)

			// If the user already has access, allow the request through.
			if hasAccess {
				return next(w, r)
			}

			// Allow some routes without a room code.
			if shouldBypassRoomCode(r) {
				return next(w, r)
			}

			// Redirect to home page so the user can enter the room code.
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return nil
		}
	}
}

func shouldBypassRoomCode(r *http.Request) bool {
	path := r.URL.Path

	// Static assets should always be allowed.
	if strings.HasPrefix(path, "/public/") {
		return true
	}

	// Admin login/registration and admin pages are controlled separately.
	if strings.HasPrefix(path, "/admin") {
		return true
	}

	// Home page should be visible so users can enter the code.
	if path == "/" {
		return true
	}

	// Endpoint for submitting the room code.
	if path == "/room-code" && r.Method == http.MethodPost {
		return true
	}

	return false
}
