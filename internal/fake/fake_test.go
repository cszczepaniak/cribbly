package fake

import "testing"

func TestFake(t *testing.T) {
	for i, fn := range []func() string{FirstName, LastName, USState} {
		if f := fn(); f == "" {
			t.Fatalf("function should have returned something (index %d)", i)
		}
	}
}
