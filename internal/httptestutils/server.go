package httptestutils

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cszczepaniak/cribbly/internal/ui/components"
)

type Server struct {
	s *httptest.Server
}

// NewServer returns a testing server that can be used to make real HTTP requests against the given
// [http.Handler]. The server will close itself on test cleanup.
func NewServer(t *testing.T, h http.Handler) Server {
	s := httptest.NewServer(h)
	t.Cleanup(s.Close)

	return Server{
		s: s,
	}
}

// NewServerForComponent returns a testing server that can be used to make real HTTP requests
// against the given [components.ComponentHandler]. The server will close itself on test cleanup.
func NewServerForComponent(t *testing.T, h components.ComponentHandler) Server {
	return NewServer(t, components.Handle(h))
}

func (s Server) Get(t *testing.T, path string) Response {
	t.Helper()

	return s.Do(t, http.MethodGet, path, nil)
}

func (s Server) Post(t *testing.T, path string, body io.Reader) Response {
	t.Helper()

	return s.Do(t, http.MethodPost, path, body)
}

func (s Server) Do(t *testing.T, method, path string, body io.Reader) Response {
	t.Helper()

	req, err := http.NewRequest(method, s.s.URL+path, body)
	require.NoError(t, err)

	resp, err := s.s.Client().Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return Response{
		bytes: bs,
	}
}
