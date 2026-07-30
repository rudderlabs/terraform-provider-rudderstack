package configs

import "encoding/json"

type TestConfig struct {
	TerraformCreate string
	// APICreate and APIUpdate hold the expected create/update request config JSON.
	APICreate       string
	TerraformUpdate string
	APIUpdate       string
	// APICreateResponse and APIUpdateResponse hold response-specific config JSON
	// for destination GET assertions when the API prunes request-only fields from
	// responses. Empty values fall back to APICreate/APIUpdate.
	APICreateResponse string
	APIUpdateResponse string
	// APICreatePrunedKeys and APIUpdatePrunedKeys remove request-only top-level
	// config keys from destination GET assertions while keeping APICreate/APIUpdate
	// strict for outgoing create/update payload validation.
	APICreatePrunedKeys []string
	APIUpdatePrunedKeys []string
	// APICreateSettings and APIUpdateSettings hold the expected settings JSON for
	// source-level fields (GeoEnrichmentEnabled, Transient) that are asserted separately from Config.
	APICreateSettings string
	APIUpdateSettings string
}

func (tc TestConfig) DestinationAPICreateResponse() string {
	return pruneDestinationAPIResponse(fallbackDestinationAPIResponse(tc.APICreateResponse, tc.APICreate), tc.APICreatePrunedKeys...)
}

func (tc TestConfig) DestinationAPIUpdateResponse() string {
	return pruneDestinationAPIResponse(fallbackDestinationAPIResponse(tc.APIUpdateResponse, tc.APIUpdate), tc.APIUpdatePrunedKeys...)
}

func fallbackDestinationAPIResponse(response, request string) string {
	if response != "" {
		return response
	}
	return request
}

func pruneDestinationAPIResponse(response string, keys ...string) string {
	if len(keys) == 0 {
		return response
	}
	return pruneJSONKeys(response, keys...)
}

func pruneJSONKeys(jsonString string, keys ...string) string {
	if jsonString == "" {
		return ""
	}

	var m map[string]interface{}
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
