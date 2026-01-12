package index

import (
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) error {
	return index().Render(r.Context(), w)
}
