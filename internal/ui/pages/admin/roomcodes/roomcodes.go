package roomcodes

import (
	"crypto/rand"
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
	const (
		codeLength  = 6
		maxAttempts = 5
	)

	code, err := generateCode(codeLength)
	if err != nil {
		return err
	}

	expiresAt := time.Now().Add(24 * time.Hour)

	// Best-effort to avoid collisions by retrying a few times.
	for i := 0; i < maxAttempts; i++ {
		err = h.RoomCodeRepo.Create(r.Context(), code, expiresAt)
		if err == nil {
			http.Redirect(w, r, "/admin/room-codes", http.StatusFound)
			return nil
		}

		// Generate a new code and try again.
		code, err = generateCode(codeLength)
		if err != nil {
			return err
		}
	}

	return err
}

func generateCode(n int) (string, error) {
	const alphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

	buf := make([]byte, n)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	for i := range buf {
		buf[i] = alphabet[int(buf[i])%len(alphabet)]
	}

	return string(buf), nil
}
