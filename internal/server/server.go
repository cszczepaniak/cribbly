package server

import (
	"net/http"

	"github.com/cszczepaniak/cribbly/internal/ui/components"
	"github.com/cszczepaniak/cribbly/internal/ui/pages/index"
)

func Setup() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /", components.Handle(index.Index))
	mux.Handle("GET /list", components.Handle(index.Vals))

	return mux
}
