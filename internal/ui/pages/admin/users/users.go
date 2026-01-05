package users

import (
	"errors"
	"net/http"

	"github.com/starfederation/datastar-go/datastar"

	"github.com/cszczepaniak/cribbly/internal/persistence/users"
)

type UsersHandler struct {
	UserService users.Service
}

func (h UsersHandler) Index(w http.ResponseWriter, r *http.Request) error {
	users, err := h.UserService.GetAll(r.Context())
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

	if signals.Username == "" || signals.Password == "" {
		return errors.New("missing username or password")
	}

	err = h.UserService.CreateUser(r.Context(), signals.Username, signals.Password)
	if err != nil {
		return err
	}

	users, err := h.UserService.GetAll(r.Context())
	if err != nil {
		return err
	}

	sse := datastar.NewSSE(w, r)
	return sse.PatchElementTempl(userTable(users))
}

func (h UsersHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	name := r.PathValue("name")
	err := h.UserService.DeleteUser(r.Context(), name)
	if err != nil {
		return err
	}

	users, err := h.UserService.GetAll(r.Context())
	if err != nil {
		return err
	}

	sse := datastar.NewSSE(w, r)
	return sse.PatchElementTempl(userTable(users))
}
