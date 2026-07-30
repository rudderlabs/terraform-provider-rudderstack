package testutil

import (
	"strings"
	"testing"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

// RedactedAPIConfigKeys returns the set of top-level API config keys the backend
// redacts from responses: those that map from a Sensitive (secret) schema field.
//
// The API key for each Sensitive field (including nested ones) is discovered by
// setting a sentinel at its state path and running the destination's real
// state->API transform — the same path the provider uses — so no mapping is kept
// by hand. A "contains" match tolerates transforms that wrap the value (e.g. PEM).
//
// Derivation failures are fatal rather than silently returning nil: a swallowed
// error would regress the suite to a confusing "missing field" with no hint that
// redaction detection is what broke.
func RedactedAPIConfigKeys(t *testing.T, cm configs.ConfigMeta) map[string]bool {
	return redactedAPIConfigKeysFor(t, cm, cm.SensitiveConfigPaths())
}

// RedactedScalarAPIKeys is like RedactedAPIConfigKeys but only for scalar secrets
// — those the provider blanks on read. Sensitive collections (e.g. webhook
// headers) are excluded: the real backend blanks them in place (so they diff in
// e2e), but the provider itself doesn't touch them (so they don't diff against a
// mock that returns them intact).
func RedactedScalarAPIKeys(t *testing.T, cm configs.ConfigMeta) map[string]bool {
	return redactedAPIConfigKeysFor(t, cm, cm.SensitiveScalarPaths())
}

func redactedAPIConfigKeysFor(t *testing.T, cm configs.ConfigMeta, paths []string) map[string]bool {
	t.Helper()
	const sentinel = "accRedactedSecretSentinel"

	if len(paths) == 0 {
		return nil
	}

	stateJSON := "{}"
	for _, p := range paths {
		var err error
		if stateJSON, err = sjson.Set(stateJSON, p, sentinel); err != nil {
			t.Fatalf("RedactedAPIConfigKeys: sjson.Set(%q): %v", p, err)
		}
	}
	apiJSON, err := cm.StateToAPI(stateJSON)
	if err != nil {
		t.Fatalf("RedactedAPIConfigKeys: StateToAPI failed for %q: %v", cm.APIType, err)
	}

	out := map[string]bool{}
	gjson.Parse(apiJSON).ForEach(func(k, v gjson.Result) bool {
		if strings.Contains(v.String(), sentinel) {
			out[k.String()] = true
		}
		return true
	})
	return out
}

// BlankAPISecrets returns apiJSON with each present redacted key set to "".
// Mocked destination tests use it to make a fake Get response behave like the
// real backend, which redacts secrets from reads. Only scalar keys should be
// passed (collections would be corrupted by setting them to a string).
func BlankAPISecrets(apiJSON string, keys map[string]bool) string {
	for key := range keys {
		if !gjson.Get(apiJSON, key).Exists() {
			continue
		}
		if updated, err := sjson.Set(apiJSON, key, ""); err == nil {
			apiJSON = updated
		}
	}
	return apiJSON
}

// ConfigHasRedactedSecret reports whether apiJSON sets any redacted key to a
// non-empty value. This is what drives the perpetual secret diff: a step only
// plans a change when its config actually carries a redacted secret (e.g. a
// destination whose secret is set only on update has an empty plan on create).
func ConfigHasRedactedSecret(apiJSON string, redacted map[string]bool) bool {
	for key := range redacted {
		r := gjson.Get(apiJSON, key)
		if !r.Exists() {
			continue
		}
		switch r.Type {
		case gjson.String:
			if r.String() != "" {
				return true
			}
		case gjson.JSON:
			if len(r.Array()) > 0 || len(r.Map()) > 0 {
				return true
			}
		default:
			return true
		}
	}
	return false
}
