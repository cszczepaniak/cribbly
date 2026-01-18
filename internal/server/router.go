package server

import (
	"log/slog"
	"net/http"
	"path"
	"slices"
	"strings"
	"time"
)

type handler = func(http.ResponseWriter, *http.Request) error
type middleware = func(handler) handler

type router struct {
	m      *http.ServeMux
	prefix string
	mw     []middleware
}

func NewRouter(m *http.ServeMux) *router {
	return &router{m: m}
}

func (r *router) Handle(route string, handler handler) {
	method, route, ok := strings.Cut(route, " ")
	if !ok {
		panic("must have a method and a route")
	}

	route = path.Join(r.prefix, route)

	finalHandler := handler
	for _, mw := range slices.Backward(r.mw) {
		finalHandler = mw(finalHandler)
	}

	r.m.Handle(method+" "+route, handleWithError(finalHandler))
}

func (r *router) Group(prefix string, mw ...middleware) *router {
	return &router{
		m:      r.m,
		prefix: path.Join(r.prefix, prefix),
		mw:     slices.Concat(r.mw, mw),
	}
}

func handleWithError(fn func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t0 := time.Now()
		slog.Info("http.start", "method", r.Method, "url", r.URL)

		err := fn(w, r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			slog.Error("http.error", "error", err, "method", r.Method, "url", r.URL, "dur", time.Since(t0))
			return
		}

		slog.Info("http.done", "method", r.Method, "url", r.URL, "dur", time.Since(t0))
	}
}
