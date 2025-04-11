package components

import (
	"net/http"

	"github.com/a-h/templ"
)

type ComponentHandler func(r *http.Request) (templ.Component, error)

func Handle(handler ComponentHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := handler(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		templ.Handler(c).ServeHTTP(w, r)
	})
}
