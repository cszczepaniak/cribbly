package assert

import (
	"cmp"
	"testing"

	gocmp "github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func SliceElemsMatch[T cmp.Ordered, S ~[]T](t testing.TB, act, exp S) {
	t.Helper()
	SliceElemsMatchFunc(t, act, exp, cmp.Compare)
}

func SliceElemsMatchFunc[T any](t testing.TB, act, exp []T, compare func(a, b T) int) {
	t.Helper()
	if diff := gocmp.Diff(act, exp, cmpopts.SortSlices(compare)); diff != "" {
		t.Fatalf("values differed:\n%v", diff)
	}
}

func SliceLen[S ~[]T, T any](t testing.TB, sl S, l int) {
	t.Helper()
	if len(sl) != l {
		t.Fatalf("expected slice to have length %v, had %v", l, len(sl))
	}
}
