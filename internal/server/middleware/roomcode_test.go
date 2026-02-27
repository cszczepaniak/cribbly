package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cszczepaniak/cribbly/internal/assert"
	"github.com/cszczepaniak/cribbly/internal/persistence/database"
	"github.com/cszczepaniak/cribbly/internal/persistence/roomcodes"
)

func newRoomCodeRepo(t *testing.T) roomcodes.Repository {
	t.Helper()

	db := database.NewInMemory(t)
	repo := roomcodes.NewRepository(db)
	assert.NoError(t, repo.Init(t.Context()))

	return repo
}

func TestRoomCodeMiddleware_AllowsValidRoomCode(t *testing.T) {
	repo := newRoomCodeRepo(t)

	// Create a valid, non-expired room code.
	expires := time.Now().Add(time.Hour)
	assert.NoError(t, repo.Create(t.Context(), "GOOD", expires))

	mw := RoomCodeMiddleware(repo)

	called := false
	var hasAccess bool
	h := mw(func(w http.ResponseWriter, r *http.Request) error {
		called = true
		hasAccess = HasRoomAccess(r.Context())
		w.WriteHeader(http.StatusOK)
		return nil
	})

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	t.Cleanup(srv.Close)

	req, err := http.NewRequest(http.MethodGet, srv.URL+"/protected", nil)
	assert.NoError(t, err)
	req.AddCookie(&http.Cookie{Name: "room_code", Value: "GOOD"})

	resp, err := srv.Client().Do(req)
	assert.NoError(t, err)
	resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	if !called {
		t.Fatalf("expected handler to be called")
	}
	if !hasAccess {
		t.Fatalf("expected HasRoomAccess to be true in handler context")
	}
}

func TestRoomCodeMiddleware_InvalidCodeRedirects(t *testing.T) {
	repo := newRoomCodeRepo(t)
	mw := RoomCodeMiddleware(repo)

	called := false
	h := mw(func(w http.ResponseWriter, r *http.Request) error {
		called = true
		return nil
	})

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	t.Cleanup(srv.Close)

	req, err := http.NewRequest(http.MethodGet, srv.URL+"/protected", nil)
	assert.NoError(t, err)
	req.AddCookie(&http.Cookie{Name: "room_code", Value: "BAD"})

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Return the first response without following redirects.
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	assert.NoError(t, err)
	resp.Body.Close()

	// Should redirect to "/" with temporary redirect.
	assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)
	assert.False(t, called)
}

func TestShouldBypassRoomCode(t *testing.T) {
	tests := []struct {
		path   string
		method string
		want   bool
	}{
		{path: "/public/app.js", method: http.MethodGet, want: true},
		{path: "/admin/login", method: http.MethodGet, want: true},
		{path: "/", method: http.MethodGet, want: true},
		{path: "/room-code", method: http.MethodPost, want: true},
		{path: "/room-code", method: http.MethodGet, want: false},
		{path: "/divisions", method: http.MethodGet, want: false},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, tc.path, nil)
			assert.NoError(t, err)

			got := shouldBypassRoomCode(req)
			assert.Equal(t, tc.want, got)
		})
	}
}

