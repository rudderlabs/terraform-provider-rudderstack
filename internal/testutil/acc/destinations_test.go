package acc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationImportStateVerifyIgnoreIncludesSensitiveConfigPaths(t *testing.T) {
	cm := configs.ConfigMeta{
		ConfigSchema: map[string]*schema.Schema{
			"api_secret": {
				Type:      schema.TypeString,
				Sensitive: true,
			},
			"endpoint": {
				Type: schema.TypeString,
			},
			"headers": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from": {
							Type: schema.TypeString,
						},
						"to": {
							Type:      schema.TypeString,
							Sensitive: true,
						},
					},
				},
			},
			"key_based_authentication": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_key_id": {
							Type:      schema.TypeString,
							Sensitive: true,
						},
						"access_key": {
							Type:      schema.TypeString,
							Sensitive: true,
						},
					},
				},
			},
			"sensitive_headers": {
				Type:      schema.TypeList,
				Sensitive: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from": {
							Type: schema.TypeString,
						},
						"to": {
							Type: schema.TypeString,
						},
					},
				},
			},
		},
	}

	require.ElementsMatch(t, []string{
		"config.0.api_secret",
		"config.0.headers.0.to",
		"config.0.key_based_authentication.#",
		"config.0.key_based_authentication.0.%",
		"config.0.key_based_authentication.0.access_key_id",
		"config.0.key_based_authentication.0.access_key",
		"config.0.sensitive_headers",
		"config.0.sensitive_headers.#",
		"config.0.sensitive_headers.0.%",
		"config.0.sensitive_headers.0.from",
		"config.0.sensitive_headers.0.to",
	}, destinationImportStateVerifyIgnore(cm, configs.EmptyTestConfig))
}

func TestDestinationImportStateVerifyIgnoreIncludesPrunedResponseConfigPaths(t *testing.T) {
	cm := configs.ConfigMeta{
		ConfigSchema: map[string]*schema.Schema{
			"partner_param_keys": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from": {
							Type: schema.TypeString,
						},
						"to": {
							Type: schema.TypeString,
						},
					},
				},
			},
			"request_only_flag": {
				Type: schema.TypeBool,
			},
		},
		Properties: []configs.ConfigProperty{
			configs.ArrayWithObjects("partnerParamKeys", "partner_param_keys", map[string]interface{}{
				"from": "from",
				"to":   "to",
			}),
			configs.Simple("requestOnlyFlag", "request_only_flag"),
		},
	}
	cfg := configs.TestConfig{
		APIUpdate: `{
			"partnerParamKeys": [
				{
					"from": "userId",
					"to": "user_id"
				}
			],
			"requestOnlyFlag": true,
			"kept": "value"
		}`,
		APIUpdatePrunedKeys: []string{"partnerParamKeys", "requestOnlyFlag"},
	}

	require.ElementsMatch(t, []string{
		"config.0.partner_param_keys.#",
		"config.0.partner_param_keys.0.%",
		"config.0.partner_param_keys.0.from",
		"config.0.partner_param_keys.0.to",
		"config.0.request_only_flag",
	}, destinationImportStateVerifyIgnore(cm, cfg))
}
