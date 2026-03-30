package assert

import "testing"

func MapHasKey[K comparable, V any, M ~map[K]V](t testing.TB, m M, k K) {
	t.Helper()

	_, ok := m[k]
	if !ok {
		t.Fatalf("map doesn't contain key %v", k)
	}
}

func MapLen[M ~map[K]V, K comparable, V any](t testing.TB, m M, l int) {
	t.Helper()
	if len(m) != l {
		t.Fatalf("expected map to have length %v, had %v", l, len(m))
	}
}
