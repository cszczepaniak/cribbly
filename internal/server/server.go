package server

import (
	"fmt"
	"math/rand/v2"
	"net/http"

	"github.com/cszczepaniak/cribbly/internal/ui/components"
	"github.com/cszczepaniak/cribbly/internal/ui/pages/index"
)

func Setup() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /", components.Handle(index.Index))
	mux.Handle("GET /list", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		items := []string{
			`Item 1`,
			`Item 2`,
			`Item 3`,
			`Item 4`,
		}
		rand.Shuffle(len(items), func(i, j int) {
			items[i], items[j] = items[j], items[i]
		})
		w.Write([]byte(`<ul id="items">`))
		for _, item := range items {
			fmt.Fprintf(w, `<li>%s</li>`, item)
		}
		w.Write([]byte(`</ul>`))
	}))

	return mux
}
