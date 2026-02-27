package index

import (
	"net/http"
	"strings"
	"time"

	"github.com/cszczepaniak/cribbly/internal/persistence/roomcodes"
	"github.com/cszczepaniak/cribbly/internal/server/middleware"
)

type Handler struct {
	RoomCodeRepo roomcodes.Repository
}

func (h Handler) Index(w http.ResponseWriter, r *http.Request) error {
	return index().Render(r.Context(), w)
}

func (h Handler) SubmitRoomCode(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	code := strings.TrimSpace(r.FormValue("room_code"))
	if code == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return nil
	}

	ok, err := h.RoomCodeRepo.Validate(r.Context(), code)
	if err != nil {
		return err
	}

	if !ok {
		http.Redirect(w, r, "/", http.StatusFound)
		return nil
	}

	// Room code is valid; set a cookie so the user doesn't need to enter it again.
	http.SetCookie(w, &http.Cookie{
		Name:     "room_code",
		Value:    code,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   middleware.IsProd(r.Context()),
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, "/", http.StatusFound)
	return nil
}
