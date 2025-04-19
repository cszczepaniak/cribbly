package admin

import (
	"net/http"

	"github.com/a-h/templ"
)

func Index(*http.Request) (templ.Component, error) {
	return adminPage(), nil
}
