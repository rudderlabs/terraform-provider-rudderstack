package generator

import (
	"encoding/json"
	"testing"
)

func TestCustomerIOAudienceConfigBlock(t *testing.T) {
	cases := []struct {
		name    string
		in      string
		wantErr bool
	}{
		{name: "valid int audienceId", in: `{"audienceId": 16}`, wantErr: false},
		{name: "zero audienceId", in: `{"audienceId": 0}`, wantErr: false},
		{name: "missing audienceId", in: `{}`, wantErr: true},
		{name: "string audienceId", in: `{"audienceId": "abc"}`, wantErr: true},
		// Reject fractional audienceId — int64(42.5) would silently truncate.
		{name: "fractional audienceId", in: `{"audienceId": 42.5}`, wantErr: true},
		{name: "invalid JSON", in: `not json`, wantErr: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := customerIOAudienceConfigBlock(json.RawMessage(tc.in))
			if (err != nil) != tc.wantErr {
				t.Fatalf("err=%v, wantErr=%v", err, tc.wantErr)
			}
			if tc.wantErr {
				return
			}
			if b == nil {
				t.Fatal("expected non-nil block")
			}
		})
	}
}
