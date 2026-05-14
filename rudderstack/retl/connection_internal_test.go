package retl

import (
	"encoding/json"
	"testing"
)

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
