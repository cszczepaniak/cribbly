package roomcodes

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cszczepaniak/cribbly/internal/assert"
	"github.com/cszczepaniak/cribbly/internal/persistence/database"
	"github.com/cszczepaniak/cribbly/internal/persistence/roomcodes"
)

func newHandler(t *testing.T) (Handler, roomcodes.Repository) {
	t.Helper()

	db := database.NewInMemory(t)
	repo := roomcodes.NewRepository(db)
	assert.NoError(t, repo.Init(t.Context()))

	return Handler{RoomCodeRepo: repo}, repo
}

func TestIndexShowsLatestIfPresent(t *testing.T) {
	h, repo := newHandler(t)

	expires := time.Now().Add(time.Hour)
	assert.NoError(t, repo.Create(t.Context(), "ABC123", expires))

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := h.Index(w, r); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	t.Cleanup(srv.Close)

	resp, err := srv.Client().Get(srv.URL + "/admin/room-codes")
	assert.NoError(t, err)
	resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGenerateCreatesRoomCodeAndRedirects(t *testing.T) {
	h, repo := newHandler(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := h.Generate(w, r); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	t.Cleanup(srv.Close)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Return the first response without following redirects.
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest(http.MethodPost, srv.URL+"/admin/room-codes", nil)
	assert.NoError(t, err)

	resp, err := client.Do(req)
	assert.NoError(t, err)
	resp.Body.Close()

	assert.Equal(t, http.StatusFound, resp.StatusCode)

	// Ensure a room code was created.
	_, err = repo.Latest(t.Context())
	assert.NoError(t, err)
}
