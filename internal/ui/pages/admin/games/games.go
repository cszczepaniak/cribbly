package games

import (
	"net/http"
)

type Handler struct{}

func (h Handler) Index(w http.ResponseWriter, r *http.Request) error {
	return index().Render(r.Context(), w)
}
