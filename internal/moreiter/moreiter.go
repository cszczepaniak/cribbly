package moreiter

import (
	"iter"
	"slices"
)

func Of[T any](vals ...T) iter.Seq[T] {
	return slices.Values(vals)
}

func Map[T, U any](in iter.Seq[T], fn func(T) U) iter.Seq[U] {
	return func(yield func(U) bool) {
		for val := range in {
			if !yield(fn(val)) {
				return
			}
		}
	}
}
