package configs

import "encoding/json"

type TestConfig struct {
	TerraformCreate string
	// APICreate and APIUpdate hold the exact config JSON of the outgoing
	// create/update request. They stay strict: the provider must send them
	// verbatim.
	APICreate       string
	TerraformUpdate string
	APIUpdate       string
	// APICreatePrunedKeys and APIUpdatePrunedKeys list top-level API config keys
	// that config-backend drops from the GET response even though the provider
	// sent them (it prunes keys not declared for the destination definition —
	// e.g. Amplitude's Browser SDK v1 settings once sdkVersion is 2). Response
	// assertions must not require them; the request assertion above still does.
	APICreatePrunedKeys []string
	APIUpdatePrunedKeys []string
	// APICreateSettings and APIUpdateSettings hold the expected settings JSON for
	// source-level fields (GeoEnrichmentEnabled, Transient) that are asserted separately from Config.
	APICreateSettings string
	APIUpdateSettings string
}

// APICreateResponse and APIUpdateResponse return the create/update config JSON
// as the backend is expected to echo it back: the request JSON minus the keys it
// prunes. Used to drive mock GETs so unit tests exercise the pruned read path.
func (tc TestConfig) APICreateResponse() string {
	return pruneJSONKeys(tc.APICreate, tc.APICreatePrunedKeys)
}

func (tc TestConfig) APIUpdateResponse() string {
	return pruneJSONKeys(tc.APIUpdate, tc.APIUpdatePrunedKeys)
}

func pruneJSONKeys(jsonString string, keys []string) string {
	if jsonString == "" || len(keys) == 0 {
		return jsonString
	}

	var m map[string]any
	if err := json.Unmarshal([]byte(jsonString), &m); err != nil {
		return jsonString
	}
	for _, key := range keys {
		delete(m, key)
	}

	b, err := json.Marshal(m)
	if err != nil {
		return jsonString
	}
	return string(b)
}

var EmptyTestConfig = TestConfig{TerraformCreate: "", APICreate: "{}", TerraformUpdate: "", APIUpdate: "{}"}
