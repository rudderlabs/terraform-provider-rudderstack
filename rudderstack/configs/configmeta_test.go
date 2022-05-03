package configs_test

import (
	"testing"

	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	// StateToAPI will check all conditionals and use that last value that exists in state
	api, err := cm.StateToAPI(`{
		"s": 123,
		"d": true,
		"c1": "condition1",
		"c2": "condition2"
	}`)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"simple": 123,
		"discriminator": "VALUE",
		"conditional": "condition2"
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
		"conditional": "condition2"
	}`)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"s": 123,
		"c2": { "v": "condition2" }
	}`, api)
}
