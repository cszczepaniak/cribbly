package profile

import (
	"errors"
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/starfederation/datastar-go/datastar"

	"github.com/cszczepaniak/cribbly/internal/persistence/users"
	"github.com/cszczepaniak/cribbly/internal/ui/pages/admin/middleware"
)

type ProfileHandler struct {
	UserRepo users.Repository
}

func (h ProfileHandler) Index(w http.ResponseWriter, r *http.Request) error {
	return index().Render(r.Context(), w)
}

func (h ProfileHandler) ChangePassword(w http.ResponseWriter, r *http.Request) error {
	var signals struct {
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}
	err := datastar.ReadSignals(r, &signals)
	if err != nil {
		return err
	}

	if signals.Password != signals.ConfirmPassword {
		// TODO better error
		return errors.New("passwords must match")
	}

	sesh := middleware.GetSession(r.Context())

	hash, err := argon2id.CreateHash(signals.Password, argon2id.DefaultParams)
	if err != nil {
		return err
	}

	err = h.UserRepo.ChangePassword(r.Context(), sesh.Username, hash)
	if err != nil {
		return err
	}

	sse := datastar.NewSSE(w, r)
	signals.Password = ""
	signals.ConfirmPassword = ""
	return sse.MarshalAndPatchSignals(signals)
}
