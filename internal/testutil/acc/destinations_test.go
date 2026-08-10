package acc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestWithPrunedKeys(t *testing.T) {
	redacted := map[string]bool{"apiSecret": true}

	got := withPrunedKeys(redacted, []string{"collectContext"})
	assert.True(t, got["apiSecret"], "redacted keys are kept")
	assert.True(t, got["collectContext"], "pruned keys are added")
	assert.False(t, redacted["collectContext"], "the shared redacted map must not be mutated")

	assert.Equal(t, redacted, withPrunedKeys(redacted, nil), "no pruned keys is a passthrough")
}

// The state path of a pruned key is derived through the real API->state
// transform, not guessed from the key name.
func TestPrunedStateAttrPaths(t *testing.T) {
	cm := configs.ConfigMeta{
		APIType: "TEST",
		ConfigSchema: map[string]*schema.Schema{
			"api_key":                  {Type: schema.TypeString},
			"device_id_from_url_param": {Type: schema.TypeBool},
			"collect_context":          {Type: schema.TypeBool},
		},
		Properties: []configs.ConfigProperty{
			configs.Simple("apiKey", "api_key"),
			configs.Simple("deviceIdFromUrlParam", "device_id_from_url_param"),
			configs.Simple("collectContext", "collect_context"),
		},
	}
	cfg := configs.TestConfig{
		APICreate:           `{"apiKey":"abc","deviceIdFromUrlParam":true,"collectContext":true}`,
		APICreatePrunedKeys: []string{"deviceIdFromUrlParam"},
		APIUpdate:           `{"apiKey":"xyz","deviceIdFromUrlParam":true,"collectContext":true}`,
		APIUpdatePrunedKeys: []string{"deviceIdFromUrlParam", "collectContext"},
	}

	got := prunedStateAttrPaths(t, cm, cfg)

	assert.ElementsMatch(t, []string{"config.0.device_id_from_url_param", "config.0.collect_context"}, got,
		"both steps' pruned keys map to their state paths, deduplicated")
	assert.NotContains(t, got, "config.0.api_key", "non-pruned keys stay verified on import")
}

func TestPrunedStateAttrPaths_NoPrunedKeys(t *testing.T) {
	cm := configs.ConfigMeta{
		APIType:      "TEST",
		ConfigSchema: map[string]*schema.Schema{"api_key": {Type: schema.TypeString}},
		Properties:   []configs.ConfigProperty{configs.Simple("apiKey", "api_key")},
	}
	assert.Empty(t, prunedStateAttrPaths(t, cm, configs.TestConfig{APICreate: `{"apiKey":"abc"}`}))
}
