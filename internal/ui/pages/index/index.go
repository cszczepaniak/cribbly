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

