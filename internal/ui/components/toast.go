package components

import (
	"net/http"

	"github.com/starfederation/datastar-go/datastar"
)

func ShowErrorToast(w http.ResponseWriter, r *http.Request, message string) error {
	sse := datastar.NewSSE(w, r)
	return sse.PatchElementTempl(errorToast(message), datastar.WithModeAppend(), datastar.WithSelectorID("body"))
}
