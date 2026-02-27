package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cszczepaniak/cribbly/internal/assert"
	"github.com/cszczepaniak/cribbly/internal/persistence/database"
	"github.com/cszczepaniak/cribbly/internal/persistence/users"
)

func newUserRepo(t *testing.T) users.Repository {
	t.Helper()

	db := database.NewInMemory(t)
	repo := users.NewRepository(db)
	assert.NoError(t, repo.Init(t.Context()))

	return repo
}

func TestAuthenticationMiddleware_NoCookieLeavesContextUnauthenticated(t *testing.T) {
	repo := newUserRepo(t)
	mw := AuthenticationMiddleware(repo)

	var isAdmin bool
	h := mw(func(w http.ResponseWriter, r *http.Request) error {
		isAdmin = IsAdmin(r.Context())
		return nil
	})

	resp := runHandler(t, h, nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.False(t, isAdmin)
}

func TestAuthenticationMiddleware_ExpiredSessionLeavesContextUnauthenticated(t *testing.T) {
	repo := newUserRepo(t)

	ctx := context.Background()
	username := t.Name() + "@example.com"
	assert.NoError(t, repo.CreateUser(ctx, username, "hash"))

	// Create an already-expired session.
	sessionID, err := repo.CreateSession(ctx, username, -1*time.Minute)
	assert.NoError(t, err)

	mw := AuthenticationMiddleware(repo)

	var isAdmin bool
	h := mw(func(w http.ResponseWriter, r *http.Request) error {
		isAdmin = IsAdmin(r.Context())
		return nil
	})

	resp := runHandler(t, h, func(req *http.Request) {
		req.AddCookie(&http.Cookie{Name: "session", Value: sessionID})
	})
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.False(t, isAdmin)
}

func TestAuthenticationMiddleware_ValidSessionMarksContextAsAdmin(t *testing.T) {
	repo := newUserRepo(t)

	ctx := context.Background()
	username := t.Name() + "@example.com"
	assert.NoError(t, repo.CreateUser(ctx, username, "hash"))

	sessionID, err := repo.CreateSession(ctx, username, time.Hour)
	assert.NoError(t, err)

	mw := AuthenticationMiddleware(repo)

	var gotSession users.Session
	h := mw(func(w http.ResponseWriter, r *http.Request) error {
		if !IsAdmin(r.Context()) {
			t.Fatalf("expected admin in context")
		}
		gotSession = GetSession(r.Context())
		return nil
	})

	resp := runHandler(t, h, func(req *http.Request) {
		req.AddCookie(&http.Cookie{Name: "session", Value: sessionID})
	})
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	if gotSession.ID != sessionID {
		t.Fatalf("expected session ID %q, got %q", sessionID, gotSession.ID)
	}
}

func TestErrorIfNotAdmin(t *testing.T) {
	mw := ErrorIfNotAdmin()

	h := mw(func(w http.ResponseWriter, r *http.Request) error {
		return nil
	})

	resp := runHandler(t, h, nil)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestRedirectToLoginIfNotAdmin(t *testing.T) {
	mw := RedirectToLoginIfNotAdmin()

	h := mw(func(w http.ResponseWriter, r *http.Request) error {
		t.Fatalf("handler should not be called for non-admin")
		return nil
	})

	resp := runHandler(t, h, nil)
	assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)
}

func runHandler(t *testing.T, h handler, mutateReq func(*http.Request)) *http.Response {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if mutateReq != nil {
			mutateReq(r)
		}
		if err := h(w, r); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	req, err := http.NewRequest(http.MethodGet, srv.URL+"/", nil)
	assert.NoError(t, err)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Return the first response without following redirects.
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	assert.NoError(t, err)

	return resp
}
