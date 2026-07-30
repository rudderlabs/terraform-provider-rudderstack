package configs_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestConfigMetaStateToAPI(t *testing.T) {
	cm := configs.ConfigMeta{
		Properties: []configs.ConfigProperty{
			configs.Simple("simple", "s"),
			configs.Discriminator("discriminator", map[string]interface{}{
				"d": "VALUE",
			}),
			configs.Conditional("conditional", "c1", configs.Equals("f", "FOO")),
			configs.Conditional("conditional", "c2", configs.Equals("f", "BAR")),
		},
	}

	// StateToAPI will check all conditionals and use that last value that exists in state.
	// Unmapped state fields are intentionally ignored so stale config values cannot be
	// round-tripped back to the Public API.
	api, err := cm.StateToAPI(`{
		"s": 123,
		"d": true,
		"c1": "condition1",
		"c2": "condition2",
		"stale": "should not be emitted",
		"nested": { "stale": "should not be emitted" }
	}`)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"simple": 123,
		"discriminator": "VALUE",
		"conditional": "condition2"
	}`, api)
}

func TestConfigMetaStateToAPIIgnoresUnmappedNestedFields(t *testing.T) {
	cm := configs.ConfigMeta{
		Properties: []configs.ConfigProperty{
			configs.Simple("nested.value", "nested.0.value"),
			configs.ArrayWithObjects("objects", "objects", map[string]interface{}{
				"known": "known",
				"nestedValues": configs.APINestedObject{
					TerraformKey: "nested_values",
					NestedKey:    "value",
				},
			}),
		},
	}

	api, err := cm.StateToAPI(`{
		"nested": [{ "value": "kept", "stale": "should not be emitted" }],
		"objects": [{
			"known": "kept-object",
			"nested_values": ["kept-nested"],
			"stale_object_field": "should not be emitted"
		}],
		"stale_top_level": "should not be emitted"
	}`)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"nested": { "value": "kept" },
		"objects": [{
			"known": "kept-object",
			"nestedValues": [{ "value": "kept-nested" }]
		}]
	}`, api)
}

func TestConfigMetaAPIToState(t *testing.T) {
	cm := configs.ConfigMeta{
		Properties: []configs.ConfigProperty{
			configs.Simple("simple", "s"),
			configs.Discriminator("discriminator", map[string]interface{}{
				"c1": "FOO",
				"c2": "BAR",
			}),
			configs.Conditional("conditional", "c1.v", configs.Equals("discriminator", "FOO")),
			configs.Conditional("conditional", "c2.v", configs.Equals("discriminator", "BAR")),
		},
	}

	api, err := cm.APIToState(`{
		"simple": 123,
		"discriminator": "BAR",
		"conditional": "condition2",
		"stale": "should not enter state",
		"nested": { "stale": "should not enter state" }
	}`)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"s": 123,
		"c2": { "v": "condition2" }
	}`, api)
}

func TestConfigMetaAPIToStateIgnoresUnmappedNestedFields(t *testing.T) {
	cm := configs.ConfigMeta{
		Properties: []configs.ConfigProperty{
			configs.Simple("nested.value", "nested.0.value"),
			configs.ArrayWithObjects("objects", "objects", map[string]interface{}{
				"known": "known",
				"nestedValues": configs.APINestedObject{
					TerraformKey: "nested_values",
					NestedKey:    "value",
				},
			}),
		},
	}

	state, err := cm.APIToState(`{
		"nested": { "value": "kept", "stale": "should not enter state" },
		"objects": [{
			"known": "kept-object",
			"staleObjectField": "should not enter state",
			"nestedValues": [
				{ "value": "kept-nested", "staleNestedField": "should not enter state" }
			]
		}],
		"staleTopLevel": "should not enter state"
	}`)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"nested": [{ "value": "kept" }],
		"objects": [{
			"known": "kept-object",
			"nested_values": ["kept-nested"]
		}]
	}`, state)
}
