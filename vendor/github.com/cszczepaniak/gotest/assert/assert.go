package assert

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Equal[T any](t testing.TB, act, exp T) {
	t.Helper()
	if diff := cmp.Diff(act, exp); diff != "" {
		t.Fatalf("values differed:\n%v", diff)
	}
}

func Error(t testing.TB, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected an error but got none")
	}
}

func NoError(t testing.TB, err error) {
	t.Helper()
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

		msg.WriteString("got chain:\n")
		writeErrorChain(msg, err, "")
		t.Fatal(msg.String())
	}
}

// writeErrorChain appends the error and its full chain (including joined errors) to b.
// indent is the prefix for each line (increased for nested/joined branches).
func writeErrorChain(b *strings.Builder, err error, indent string) {
	if err == nil {
		return
	}
	fmt.Fprintf(b, "%s%v (%T)\n", indent, err, err)

	switch e := err.(type) {
	case interface{ Unwrap() error }:
		writeErrorChain(b, e.Unwrap(), indent+"  ")
	case interface{ Unwrap() []error }:
		subs := e.Unwrap()
		for i, sub := range subs {
			fmt.Fprintf(b, "%s  [%d]:\n", indent, i+1)
			writeErrorChain(b, sub, indent+"    ")
		}
	default:
		// end of chain
	}
}
