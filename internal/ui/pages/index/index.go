package index

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/a-h/templ"
)

func Index(r *http.Request) (templ.Component, error) {
	return index(), nil
}

func Vals(r *http.Request) (templ.Component, error) {
	vals := []string{"a", "b", "c", "d"}
	rand.Shuffle(len(vals), func(i, j int) {
		vals[i], vals[j] = vals[j], vals[i]
	})
	return items(vals), nil
}

func SSE() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming unsupported!", http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		flusher.Flush()

		vals := []string{"a", "b", "c", "d"}
		for {
			rand.Shuffle(len(vals), func(i, j int) {
				vals[i], vals[j] = vals[j], vals[i]
			})
			eventRenderer := items(vals)

			fmt.Fprintln(w, "event: values")
			fmt.Fprint(w, "data: ")
			eventRenderer.Render(r.Context(), w)
			fmt.Fprintln(w)
			fmt.Fprintln(w)
			flusher.Flush()

			time.Sleep(2000 * time.Millisecond)
		}
	})
}
