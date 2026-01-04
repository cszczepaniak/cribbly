package admin

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/starfederation/datastar-go/datastar"

	"github.com/cszczepaniak/cribbly/internal/persistence/users"
)

type AdminHandler struct {
	UserService users.Service
}

func Index(w http.ResponseWriter, r *http.Request) error {
	return adminPage().Render(r.Context(), w)
}

func LoginPage(w http.ResponseWriter, r *http.Request) error {
	return loginPage().Render(r.Context(), w)
}

func RegisterPage(w http.ResponseWriter, r *http.Request) error {
	return registerPage().Render(r.Context(), w)
}

func (h AdminHandler) DoLogin(w http.ResponseWriter, r *http.Request) error {
	var signals struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := datastar.ReadSignals(r, &signals)
	if err != nil {
		return err
	}

	badLogin := func() error {
		sse := datastar.NewSSE(w, r)

		signals.Username = ""
		signals.Password = ""
		err := sse.MarshalAndPatchSignals(signals)
		if err != nil {
			return err
		}

		// TODO: do a real error display somehow
		return sse.ExecuteScript("alert('bad login')")
	}

	pwHash, err := h.UserService.GetPassword(r.Context(), signals.Username)
	if err != nil {
		if errors.Is(err, users.ErrUnknownUser) {
			return badLogin()
		}

		return err
	}

	match, err := argon2id.ComparePasswordAndHash(signals.Password, pwHash)
	if err != nil {
		return err
	}

	if !match {
		return badLogin()
	}

	// TODO: decide on session lifetime
	sessionID, err := h.UserService.CreateSession(r.Context(), signals.Username, 24*time.Hour)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // TODO: only in dev
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, "/admin", http.StatusFound)
	return nil
}

func (h AdminHandler) Register(w http.ResponseWriter, r *http.Request) error {
	var signals struct {
		Username       string `json:"username"`
		Password       string `json:"password"`
		RepeatPassword string `json:"repeat_password"`
	}

	err := datastar.ReadSignals(r, &signals)
	if err != nil {
		return err
	}

	badLogin := func(msg string) error {
		sse := datastar.NewSSE(w, r)

		signals.Username = ""
		signals.Password = ""
		signals.RepeatPassword = ""
		err := sse.MarshalAndPatchSignals(signals)
		if err != nil {
			return err
		}

		// TODO: do a real error display somehow
		return sse.ExecuteScript(fmt.Sprintf("alert(%q)", msg))
	}

	if signals.Password != signals.RepeatPassword {
		badLogin("passwords don't match")
	}

	hash, err := argon2id.CreateHash(signals.Password, argon2id.DefaultParams)
	if err != nil {
		return err
	}

	err = h.UserService.CreateUser(r.Context(), signals.Username, hash)
	if err != nil {
		return err
	}

	http.Redirect(w, r, "/admin/login", http.StatusFound)
	return nil
}

// TODO: share this type somewhere (can't be server pkg because that'd cause an import cycle)
type handler = func(http.ResponseWriter, *http.Request) error

func (h AdminHandler) AuthenticationMiddleware(next handler) handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		cookie, err := r.Cookie("session")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				http.Redirect(w, r, "/admin/login", http.StatusTemporaryRedirect)
				return nil
			}

			return fmt.Errorf("get session cookie: %w", err)
		}

		// TODO: in-memory caching of sessions to avoid a DB call for each request
		expires, err := h.UserService.GetSession(r.Context(), cookie.Value)
		if err != nil {
			if errors.Is(err, users.ErrSessionExpired) {
				http.Redirect(w, r, "/admin/login", http.StatusTemporaryRedirect)
				return nil
			}
			return err
		}

		if time.Now().After(expires) {
			// Clear out the old session cookie
			http.SetCookie(w, &http.Cookie{
				Name:   "session",
				Value:  "",
				MaxAge: -1,
			})
			http.Redirect(w, r, "/admin/login", http.StatusTemporaryRedirect)
			return nil
		}

		return next(w, r)
	}
}
