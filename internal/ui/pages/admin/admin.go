package admin

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/starfederation/datastar-go/datastar"

	"github.com/cszczepaniak/cribbly/internal/persistence/users"
)

type AdminHandler struct {
	UserRepo users.Repository
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

		return sse.PatchElementTempl(loginError("Invalid credentials"))
	}

	pwHash, err := h.UserRepo.GetPassword(r.Context(), signals.Username)
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
	sessionID, err := h.UserRepo.CreateSession(r.Context(), signals.Username, 24*time.Hour)
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

	badLogin := func(msg string, clearUser bool) error {
		sse := datastar.NewSSE(w, r)

		if clearUser {
			signals.Username = ""
		}
		signals.Password = ""
		signals.RepeatPassword = ""
		err := sse.MarshalAndPatchSignals(signals)
		if err != nil {
			return err
		}

		return sse.PatchElementTempl(registerError(msg))
	}

	if signals.Password != signals.RepeatPassword {
		return badLogin("Passwords don't match.", false)
	}

	hash, err := argon2id.CreateHash(signals.Password, argon2id.DefaultParams)
	if err != nil {
		return err
	}

	err = h.UserRepo.CreateUser(r.Context(), signals.Username, hash)
	if err != nil {
		log.Println("registration error:", err)
		return badLogin("Error with registration. Contact an admin for help.", true)
	}

	http.Redirect(w, r, "/admin/login", http.StatusFound)
	return nil
}
