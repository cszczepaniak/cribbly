package assert

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Equal[T any](t testing.TB, act, exp T) {
	t.Helper()
	if diff := cmp.Diff(act, exp); diff != "" {
		t.Fatalf("values differed:\n%v", diff)
	}
}

func ElementsMatch[T any](t testing.TB, act, exp []T, compare func(a, b T) int) {
	t.Helper()
	if diff := cmp.Diff(act, exp, cmpopts.SortSlices(compare)); diff != "" {
		t.Fatalf("values differed:\n%v", diff)
	}
}

func MapHasKey[K comparable, V any, M ~map[K]V](t testing.TB, m M, k K) {
	t.Helper()

	_, ok := m[k]
	if !ok {
		t.Fatalf("map doesn't contain key %v", k)
	}
}

func SliceLen[S ~[]T, T any](t testing.TB, sl S, l int) {
	t.Helper()
	if len(sl) != l {
		t.Fatalf("expected slice to have length %v, had %v", l, len(sl))
	}
}

func MapLen[M ~map[K]V, K comparable, V any](t testing.TB, m M, l int) {
	t.Helper()
	if len(m) != l {
		t.Fatalf("expected map to have length %v, had %v", l, len(m))
	}
}

func Error(t testing.TB, err error) {
	if err == nil {
		t.Fatal("expected an error but got none")
	}
}

func NoError(t testing.TB, err error) {
	if err != nil {
		t.Fatalf("expected no error but got: %q", err)
	}
}

func ErrorIs(t testing.TB, err, exp error) {
	t.Helper()

	if !errors.Is(err, exp) {
		msg := &strings.Builder{}
		fmt.Fprintf(msg, "errors not equal\nexpected %q (%T)\n", exp, exp)

		if err == nil {
			msg.WriteString("got no error (nil)")
			t.Fatal(msg.String())
		}

		fmt.Fprintln(msg, "got chain:")
		for err != nil {
			fmt.Fprintf(msg, "\t%v (%T)\n", err, err)

			switch e := err.(type) {
			case interface{ Unwrap() error }:
				err = e.Unwrap()
			case interface{ Unwrap() []error }:
				// TODO: handle joined errors
				err = nil
			default:
				err = nil
			}
		}

		t.Fatal(msg.String())
	}
}

func False(t testing.TB, x bool) {
	t.Helper()

	if x {
		t.Fatal("expected the value to be false")
	}
}
