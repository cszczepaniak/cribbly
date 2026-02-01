package users

import (
	"net/http"
	"net/mail"

	"github.com/starfederation/datastar-go/datastar"

	"github.com/cszczepaniak/cribbly/internal/persistence/users"
)

type UsersHandler struct {
	UserRepo users.Repository
}

func (h UsersHandler) Index(w http.ResponseWriter, r *http.Request) error {
	users, err := h.UserRepo.GetAll(r.Context())
	if err != nil {
		return err
	}

	return index(users).Render(r.Context(), w)
}

func (h UsersHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var signals struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := datastar.ReadSignals(r, &signals)
	if err != nil {
		return err
	}

	var userErr, pwErr string
	if signals.Username == "" {
		userErr = "Username is required."
	}
	if signals.Password == "" {
		pwErr = "Password is required."
	}

	if signals.Username != "" {
		_, err = mail.ParseAddress(signals.Username)
		if err != nil {
			userErr = "Username must be a valid email address."
		}
	}

	if userErr != "" || pwErr != "" {
		return datastar.NewSSE(w, r).PatchElementTempl(newUserForm(userErr, pwErr))
	}

	err = h.UserRepo.CreateUser(r.Context(), signals.Username, signals.Password)
	if err != nil {
		return err
	}

	users, err := h.UserRepo.GetAll(r.Context())
	if err != nil {
		return err
	}

	sse := datastar.NewSSE(w, r)
	// Clear the form
	signals.Username = ""
	signals.Password = ""
	err = sse.MarshalAndPatchSignals(signals)
	if err != nil {
		return err
	}
	err = sse.PatchElementTempl(newUserForm("", ""))
	if err != nil {
		return err
	}
	return sse.PatchElementTempl(userTable(users))
}

func (h UsersHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	name := r.PathValue("name")
	err := h.UserRepo.DeleteUser(r.Context(), name)
	if err != nil {
		return err
	}

	users, err := h.UserRepo.GetAll(r.Context())
	if err != nil {
		return err
	}

	sse := datastar.NewSSE(w, r)
	return sse.PatchElementTempl(userTable(users))
}
