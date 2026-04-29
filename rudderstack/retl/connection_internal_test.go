package retl

import "testing"

func TestSuppressEquivalentJSON(t *testing.T) {
	cases := []struct {
		name     string
		oldVal   string
		newVal   string
		expected bool
	}{
		{"identical strings", `{"a":1}`, `{"a":1}`, true},
		{"different key order", `{"a":1,"b":2}`, `{"b":2,"a":1}`, true},
		{"different whitespace", `{"a":1}`, ` {  "a" : 1 } `, true},
		{"semantic difference", `{"a":1}`, `{"a":2}`, false},
		{"different shape", `{"a":1}`, `{"b":1}`, false},
		{"nested key reorder", `{"o":{"x":1,"y":2}}`, `{"o":{"y":2,"x":1}}`, true},
		{"array order matters", `[1,2,3]`, `[3,2,1]`, false},
		{"both invalid JSON", "{not json", "{not json", true},
		{"only one invalid", `{"a":1}`, "{not json", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := suppressEquivalentJSON("", tc.oldVal, tc.newVal, nil)
			if got != tc.expected {
				t.Errorf("suppressEquivalentJSON(%q, %q) = %v, want %v", tc.oldVal, tc.newVal, got, tc.expected)
			}
		})
	}
}

func TestNormalizeJSON(t *testing.T) {
	cases := []struct {
		name string
		in   string
		out  string
	}{
		{"empty", "", ""},
		{"canonicalizes key order", `{"b":2,"a":1}`, `{"a":1,"b":2}`},
		{"strips whitespace", ` { "a" : 1 } `, `{"a":1}`},
		{"invalid JSON returned as-is", "{not json", "{not json"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := normalizeJSON(tc.in); got != tc.out {
				t.Errorf("normalizeJSON(%q) = %q, want %q", tc.in, got, tc.out)
			}
		})
	}
}

func TestFlowForceNewRules(t *testing.T) {
	cases := []struct {
		name                              string
		objectSet, destConfigSet          bool
		wantIdentifiers, wantConstants    bool
	}{
		{"json mapper (no object, no dest_config)", false, false, true, false},
		{"object mapping (object set)", true, false, true, true},
		{"destination-specific (dest_config set)", false, true, false, true},
		{"object set wins over dest_config", true, true, true, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotID, gotC := flowForceNewRules(tc.objectSet, tc.destConfigSet)
			if gotID != tc.wantIdentifiers || gotC != tc.wantConstants {
				t.Errorf("flowForceNewRules(object=%v, destConfig=%v) = (%v, %v), want (%v, %v)",
					tc.objectSet, tc.destConfigSet, gotID, gotC, tc.wantIdentifiers, tc.wantConstants)
			}
		})
	}
}
