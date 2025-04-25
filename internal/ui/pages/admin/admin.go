package admin

import (
	"net/http"

	"github.com/a-h/templ"
)

func Index(http.ResponseWriter, *http.Request) (templ.Component, error) {
	return adminPage(), nil
}
