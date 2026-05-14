package retl

import (
	"encoding/json"
	"strings"
	"testing"
)

// stubGetter satisfies the resourceGetter interface for unit-testing helpers
// that don't need a full *schema.ResourceData / *schema.ResourceDiff.
type stubGetter map[string]interface{}

func (s stubGetter) Get(key string) interface{} { return s[key] }

func TestCustomerIOAudienceToState(t *testing.T) {
	cases := []struct {
		name      string
		in        string
		wantID    int
		wantEmpty bool // expect (nil, nil) — typed block should be cleared
		wantErr   bool
	}{
		{name: "valid int audienceId", in: `{"audienceId": 7}`, wantID: 7},
		{name: "zero audienceId", in: `{"audienceId": 0}`, wantID: 0},
		{name: "extra fields are ignored", in: `{"audienceId": 9, "name": "x"}`, wantID: 9},
		// JSON null is the server saying "no destination-specific config".
		// Must NOT error — that would conflate a valid response with the
		// "unsupported destination type" signal and break refresh on a real
		// Customer.io Audience connection that ever received a null payload.
		{name: "json null clears the typed block", in: `null`, wantEmpty: true},
		{name: "missing audienceId is unsupported", in: `{}`, wantErr: true},
		{name: "non-numeric audienceId is unsupported", in: `{"audienceId": "abc"}`, wantErr: true},
		// Reject fractional audienceId — int(42.5) would silently truncate to 42.
		{name: "fractional audienceId is rejected", in: `{"audienceId": 42.5}`, wantErr: true},
		{name: "invalid JSON", in: `not json`, wantErr: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := customerIOAudienceToState(json.RawMessage(tc.in))
			if (err != nil) != tc.wantErr {
				t.Fatalf("err=%v, wantErr=%v", err, tc.wantErr)
			}
			if tc.wantErr {
				return
			}
			if tc.wantEmpty {
				if got != nil {
					t.Fatalf("expected nil typed block, got %v", got)
				}
				return
			}
			if len(got) != 1 {
				t.Fatalf("expected 1 element, got %d", len(got))
			}
			if v, _ := got[0]["audience_id"].(int); v != tc.wantID {
				t.Errorf("audience_id = %v, want %v", v, tc.wantID)
			}
		})
	}
}

func TestValidateJSONMapperIdentifierTargets(t *testing.T) {
	cases := []struct {
		name           string
		identifiers    []interface{}
		wantErr        bool
		wantSubstrings []string
	}{
		{
			name: "user_id is accepted",
			identifiers: []interface{}{
				map[string]interface{}{"from": "email", "to": "user_id"},
			},
		},
		{
			name: "anonymous_id is accepted",
			identifiers: []interface{}{
				map[string]interface{}{"from": "anon", "to": "anonymous_id"},
			},
		},
		{
			name: "mixed valid targets accepted",
			identifiers: []interface{}{
				map[string]interface{}{"from": "a", "to": "user_id"},
				map[string]interface{}{"from": "b", "to": "anonymous_id"},
			},
		},
		{
			name: "destination column rejected with both valid values in message",
			identifiers: []interface{}{
				map[string]interface{}{"from": "email", "to": "email"},
			},
			wantErr:        true,
			wantSubstrings: []string{`identifiers[0].to`, "user_id", "anonymous_id", `got "email"`, "JSON Mapper"},
		},
		{
			name: "error reports the index of the offending entry",
			identifiers: []interface{}{
				map[string]interface{}{"from": "a", "to": "user_id"},
				map[string]interface{}{"from": "b", "to": "external_id"},
			},
			wantErr:        true,
			wantSubstrings: []string{`identifiers[1].to`, `got "external_id"`},
		},
		{
			name:        "empty identifiers list passes (presence is enforced elsewhere by MinItems)",
			identifiers: []interface{}{},
		},
		{
			name:        "nil identifiers passes (presence is enforced elsewhere)",
			identifiers: nil,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateJSONMapperIdentifierTargets(stubGetter{"identifiers": tc.identifiers})
			if (err != nil) != tc.wantErr {
				t.Fatalf("err=%v, wantErr=%v", err, tc.wantErr)
			}
			if !tc.wantErr {
				return
			}
			for _, sub := range tc.wantSubstrings {
				if !strings.Contains(err.Error(), sub) {
					t.Errorf("error %q missing expected substring %q", err.Error(), sub)
				}
			}
		})
	}
}

func TestFlowForceNewRules(t *testing.T) {
	cases := []struct {
		name                           string
		objectSet, destSpecific        bool
		wantIdentifiers, wantConstants bool
	}{
		{"json mapper (no object, no typed block)", false, false, true, false},
		{"object mapping (object set)", true, false, true, true},
		{"customerio_audience (typed block set)", false, true, false, true},
		{"object set wins over typed block", true, true, true, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotID, gotC := flowForceNewRules(tc.objectSet, tc.destSpecific)
			if gotID != tc.wantIdentifiers || gotC != tc.wantConstants {
				t.Errorf("flowForceNewRules(object=%v, destSpecific=%v) = (%v, %v), want (%v, %v)",
					tc.objectSet, tc.destSpecific, gotID, gotC, tc.wantIdentifiers, tc.wantConstants)
			}
		})
	}
}
