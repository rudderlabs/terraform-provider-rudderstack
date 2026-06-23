package retl

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestDecodeCustomerIOAudienceID(t *testing.T) {
	cases := []struct {
		name    string
		in      string
		wantID  int
		wantErr bool
		wantNil bool // expect errCustomerIOAudienceNullConfig sentinel
	}{
		{name: "valid int audienceId", in: `{"audienceId": 7}`, wantID: 7},
		{name: "extra fields are ignored", in: `{"audienceId": 9, "name": "x"}`, wantID: 9},
		// JSON null is the server saying "no destination-specific config".
		// Must surface as the sentinel error so the caller can treat it as a
		// soft signal (clear the field) rather than a malformed payload.
		{name: "json null returns sentinel", in: `null`, wantNil: true},
		{name: "missing audienceId is unsupported", in: `{}`, wantErr: true},
		{name: "non-numeric audienceId is unsupported", in: `{"audienceId": "abc"}`, wantErr: true},
		// Reject fractional audienceId — int(42.5) would silently truncate to 42.
		{name: "fractional audienceId is rejected", in: `{"audienceId": 42.5}`, wantErr: true},
		{name: "invalid JSON", in: `not json`, wantErr: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := decodeCustomerIOAudienceID(json.RawMessage(tc.in))
			if tc.wantNil {
				if !errors.Is(err, errCustomerIOAudienceNullConfig) {
					t.Fatalf("expected errCustomerIOAudienceNullConfig, got %v", err)
				}
				return
			}
			if (err != nil) != tc.wantErr {
				t.Fatalf("err=%v, wantErr=%v", err, tc.wantErr)
			}
			if tc.wantErr {
				return
			}
			if got != tc.wantID {
				t.Errorf("audience_id = %v, want %v", got, tc.wantID)
			}
		})
	}
}

func TestDecodeCustomerIOObject(t *testing.T) {
	cases := []struct {
		name    string
		in      string
		want    string
		wantErr bool
	}{
		{name: "valid object", in: `{"object": "customers"}`, want: "customers"},
		{name: "extra fields are ignored", in: `{"object": "customers", "x": 1}`, want: "customers"},
		// A 200 with no usable object is a persistent server-side inconsistency,
		// not a transient soft signal — every shape below is a hard error so the
		// problem surfaces at refresh instead of being masked.
		{name: "json null is an error", in: `null`, wantErr: true},
		{name: "empty input is an error", in: ``, wantErr: true},
		{name: "missing object is unsupported", in: `{}`, wantErr: true},
		{name: "non-string object is unsupported", in: `{"object": 7}`, wantErr: true},
		{name: "empty object string is unsupported", in: `{"object": ""}`, wantErr: true},
		{name: "invalid JSON", in: `not json`, wantErr: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := decodeCustomerIOObject(json.RawMessage(tc.in))
			if (err != nil) != tc.wantErr {
				t.Fatalf("err=%v, wantErr=%v", err, tc.wantErr)
			}
			if tc.wantErr {
				return
			}
			if got != tc.want {
				t.Errorf("object = %q, want %q", got, tc.want)
			}
		})
	}
}
