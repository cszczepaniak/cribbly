package index

import (
	"net/http"

	"github.com/a-h/templ"
)

func Index(r *http.Request) (templ.Component, error) {
	return index(), nil
}
