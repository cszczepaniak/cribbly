package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cszczepaniak/cribbly/internal/assert"
)

func TestIsProdMiddleware(t *testing.T) {
	var expIsProd bool
	var isProd bool

	assertIsProdHandler := func(_ http.ResponseWriter, r *http.Request) error {
		assert.Equal(t, expIsProd, IsProd(r.Context()))
		return nil
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		err := IsProdMiddleware(isProd)(assertIsProdHandler)(w, r)
		assert.NoError(t, err)
	}

	s := httptest.NewServer(http.HandlerFunc(handler))
	t.Cleanup(s.Close)

	isProd = true
	expIsProd = true

	_, err := s.Client().Get(s.URL)
	assert.NoError(t, err)

	isProd = false
	expIsProd = false

	_, err = s.Client().Get(s.URL)
	assert.NoError(t, err)
}
