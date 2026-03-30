package roomcodes

import (
	"errors"
	"net/http"
	"time"

	"github.com/cszczepaniak/cribbly/internal/persistence/roomcodes"
)

type Handler struct {
	RoomCodeRepo roomcodes.Repository
}

type latestView struct {
	Code    string
	Expires time.Time
}

func (h Handler) Index(w http.ResponseWriter, r *http.Request) error {
	var latest *latestView

	rc, err := h.RoomCodeRepo.Latest(r.Context())
	if err != nil {
		if !errors.Is(err, roomcodes.ErrCodeNotFound) {
			return err
		}
	} else {
		latest = &latestView{
			Code:    rc.Code,
			Expires: rc.Expires,
		}
	}

	return index(latest).Render(r.Context(), w)
}

func (h Handler) Generate(w http.ResponseWriter, r *http.Request) error {
	_, err := h.RoomCodeRepo.CreateRandomCode(r.Context())
	if err != nil {
		return err
	}
	http.Redirect(w, r, "/admin/room-codes", http.StatusFound)
	return nil
}
