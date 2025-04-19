// Package hx contains HTMX helpers to interact with HTMX via request/response headers.
package hx

import "net/http"

// RedirectTo adds the HX-Redirect header to an HTTP response, indicating to HTMX to redirect to the
// indicated location on the client side.
func RedirectTo(w http.ResponseWriter, path string) {
	w.Header().Add("HX-Redirect", path)
}
