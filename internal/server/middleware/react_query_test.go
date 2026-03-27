package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func indexHTML() []byte {
	return []byte("<html></html>")
}

func TestReactQueryMiddleware_servesReactShell(t *testing.T) {
	var nextCalled bool
	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		nextCalled = true
	})
	h := ReactQueryMiddleware(indexHTML, true, next)

	req := httptest.NewRequest(http.MethodGet, "/?react=true", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if nextCalled {
		t.Fatal("expected React shell, next handler should not run")
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("code: got %d want %d", rec.Code, http.StatusOK)
	}
	if got := rec.Body.String(); got != string(indexHTML()) {
		t.Fatalf("body: got %q want %q", got, string(indexHTML()))
	}
}

func TestReactQueryMiddleware_delegatesWithoutQuery(t *testing.T) {
	var nextCalled bool
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusTeapot)
	})
	h := ReactQueryMiddleware(indexHTML, true, next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if !nextCalled {
		t.Fatal("expected next handler without react=true")
	}
	if rec.Code != http.StatusTeapot {
		t.Fatalf("code: got %d want %d", rec.Code, http.StatusTeapot)
	}
}

func TestReactQueryMiddleware_delegatesAppAssets(t *testing.T) {
	var nextCalled bool
	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		nextCalled = true
	})
	h := ReactQueryMiddleware(indexHTML, true, next)

	req := httptest.NewRequest(http.MethodGet, "/app/assets/foo.js?react=true", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if !nextCalled {
		t.Fatal("expected /app/assets to bypass React shell even with react=true")
	}
}

func TestReactQueryMiddleware_delegatesStream(t *testing.T) {
	var nextCalled bool
	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		nextCalled = true
	})
	h := ReactQueryMiddleware(indexHTML, true, next)

	req := httptest.NewRequest(http.MethodGet, "/standings/stream?react=true", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if !nextCalled {
		t.Fatal("expected SSE stream to bypass React shell")
	}
}
