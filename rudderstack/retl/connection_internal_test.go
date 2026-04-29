package retl

import "testing"

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
