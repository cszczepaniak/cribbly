package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWithDevAdminBypassIfHeader_DisabledInProd(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(DevAdminHeader, "secret")
	got := WithDevAdminBypassIfHeader(req, "secret", true)
	if IsAdmin(got.Context()) {
		t.Fatal("expected no bypass in prod")
	}
}

func TestWithDevAdminBypassIfHeader_RequiresMatchingSecret(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(DevAdminHeader, "wrong")
	got := WithDevAdminBypassIfHeader(req, "right", false)
	if IsAdmin(got.Context()) {
		t.Fatal("expected mismatch to not grant admin")
	}
}

func TestWithDevAdminBypassIfHeader_GrantsAdminWhenSecretMatches(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(DevAdminHeader, "dev-secret")
	got := WithDevAdminBypassIfHeader(req, "dev-secret", false)
	if !IsAdmin(got.Context()) {
		t.Fatal("expected dev header to grant admin")
	}
}

func TestDevAdminBypassMiddleware_NoOpWhenSecretEmpty(t *testing.T) {
	mw := DevAdminBypassMiddleware("", false)
	var ran bool
	h := mw(func(w http.ResponseWriter, r *http.Request) error {
		ran = true
		if IsAdmin(r.Context()) {
			t.Fatal("unexpected admin")
		}
		return nil
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(DevAdminHeader, "x")
	if err := h(w, r); err != nil {
		t.Fatal(err)
	}
	if !ran {
		t.Fatal("expected handler to run")
	}
}
