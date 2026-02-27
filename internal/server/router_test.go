package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouterAppliesMiddlewareInOrder(t *testing.T) {
	var calls []string

	baseMW := func(next handler) handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			calls = append(calls, "base")
			return next(w, r)
		}
	}

	routeMW := func(next handler) handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			calls = append(calls, "route")
			return next(w, r)
		}
	}

	mux := http.NewServeMux()
	r := NewRouter(mux, baseMW)

	r.Handle("GET /test", func(w http.ResponseWriter, r *http.Request) error {
		calls = append(calls, "handler")
		return nil
	}, routeMW)

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	resp, err := srv.Client().Get(srv.URL + "/test")
	if err != nil {
		t.Fatalf("unexpected error performing request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	want := []string{"base", "route", "handler"}
	if len(calls) != len(want) {
		t.Fatalf("expected %d calls, got %d", len(want), len(calls))
	}
	for i := range want {
		if calls[i] != want[i] {
			t.Fatalf("expected call %d to be %q, got %q", i, want[i], calls[i])
		}
	}
}

func TestRouterGroupPrefixesRoutes(t *testing.T) {
	mux := http.NewServeMux()
	r := NewRouter(mux)

	group := r.Group("/admin")

	called := false
	group.Handle("GET /dashboard", func(w http.ResponseWriter, r *http.Request) error {
		called = true
		return nil
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	resp, err := srv.Client().Get(srv.URL + "/admin/dashboard")
	if err != nil {
		t.Fatalf("unexpected error performing request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if !called {
		t.Fatalf("expected grouped handler to be called")
	}
}

func TestRouterHandlerErrorResultsIn500(t *testing.T) {
	mux := http.NewServeMux()
	r := NewRouter(mux)

	r.Handle("GET /err", func(w http.ResponseWriter, r *http.Request) error {
		return assertError("boom")
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	resp, err := srv.Client().Get(srv.URL + "/err")
	if err != nil {
		t.Fatalf("unexpected error performing request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, resp.StatusCode)
	}
}

type assertError string

func (e assertError) Error() string { return string(e) }

