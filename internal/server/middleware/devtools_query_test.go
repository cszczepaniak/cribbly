package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cszczepaniak/gotest/assert"
)

func TestDevToolsQueryMiddleware(t *testing.T) {
	var expDevQuery bool

	h := func(_ http.ResponseWriter, r *http.Request) error {
		assert.Equal(t, expDevQuery, DevToolsFromQuery(r.Context()))
		return nil
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		err := DevToolsQueryMiddleware()(h)(w, r)
		assert.NoError(t, err)
	}

	s := httptest.NewServer(http.HandlerFunc(handler))
	t.Cleanup(s.Close)
	client := s.Client()

	expDevQuery = false
	_, err := client.Get(s.URL + "/x")
	assert.NoError(t, err)

	expDevQuery = true
	_, err = client.Get(s.URL + "/x?dev=true")
	assert.NoError(t, err)

	expDevQuery = false
	_, err = client.Get(s.URL + "/x?dev=false")
	assert.NoError(t, err)

	expDevQuery = false
	_, err = client.Get(s.URL + "/x?dev=1")
	assert.NoError(t, err)
}
