package dstar

import "testing"

func TestSendFormatsRequests(t *testing.T) {
	tests := []struct {
		name string
		got  string
		exp  string
	}{{
		name: "get",
		got:  SendGetf("/games/%s", "123"),
		exp:  "@get('/games/123', { requestCancellation: 'disabled' })",
	}, {
		name: "post",
		got:  SendPostf("/games/%s", "123"),
		exp:  "@post('/games/123', { requestCancellation: 'disabled' })",
	}, {
		name: "delete",
		got:  SendDeletef("/games/%s", "123"),
		exp:  "@delete('/games/123', { requestCancellation: 'disabled' })",
	}, {
		name: "put",
		got:  SendPutf("/games/%s", "123"),
		exp:  "@put('/games/123', { requestCancellation: 'disabled' })",
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.got != tc.exp {
				t.Fatalf("expected %q, got %q", tc.exp, tc.got)
			}
		})
	}
}
