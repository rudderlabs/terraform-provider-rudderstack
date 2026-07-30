package acc

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestCompareDestinationConfigIgnoresSensitiveFieldsMissingFromAPIResponse(t *testing.T) {
	cm := configs.ConfigMeta{
		ConfigSchema: map[string]*schema.Schema{
			"api_key": {
				Type:      schema.TypeString,
				Sensitive: true,
			},
			"api_url": {
				Type: schema.TypeString,
			},
		},
		Properties: []configs.ConfigProperty{
			configs.Simple("apiKey", "api_key"),
			configs.Simple("apiUrl", "api_url"),
		},
	}

	err := compareDestinationConfig(
		json.RawMessage(`{"apiUrl":"https://example.com"}`),
		`{"apiKey":"secret","apiUrl":"https://example.com"}`,
		cm,
	)

	require.NoError(t, err)
}

func TestCompareDestinationConfigIgnoresSensitiveNestedFieldsMaskedByAPIResponse(t *testing.T) {
	cm := configs.ConfigMeta{
		ConfigSchema: map[string]*schema.Schema{
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
			"url": {
				Type: schema.TypeString,
			},
		},
		Properties: []configs.ConfigProperty{
			configs.ArrayWithObjects("headers", "headers", map[string]any{
				"from": "from",
				"to":   "to",
			}),
			configs.Simple("url", "url"),
		},
	}

	err := compareDestinationConfig(
		json.RawMessage(`{"headers":[{"from":"header-name","to":""}],"url":"https://example.com"}`),
		`{"headers":[{"from":"header-name","to":"secret-value"}],"url":"https://example.com"}`,
		cm,
	)

	require.NoError(t, err)
}

func TestCompareDestinationConfigStillChecksReturnedFieldsInSensitiveBlocks(t *testing.T) {
	cm := configs.ConfigMeta{
		ConfigSchema: map[string]*schema.Schema{
			"headers": {
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
		Properties: []configs.ConfigProperty{
			configs.ArrayWithObjects("headers", "headers", map[string]any{
				"from": "from",
				"to":   "to",
			}),
		},
	}

	err := compareDestinationConfig(
		json.RawMessage(`{"headers":[{"from":"actual-header-name","to":""}]}`),
		`{"headers":[{"from":"expected-header-name","to":"secret-value"}]}`,
		cm,
	)

	require.ErrorContains(t, err, `field "headers[0].from"`)
}

func TestCompareDestinationConfigStillFailsForMissingNonSensitiveField(t *testing.T) {
	cm := configs.ConfigMeta{
		ConfigSchema: map[string]*schema.Schema{
			"api_key": {
				Type:      schema.TypeString,
				Sensitive: true,
			},
			"api_url": {
				Type: schema.TypeString,
			},
		},
		Properties: []configs.ConfigProperty{
			configs.Simple("apiKey", "api_key"),
			configs.Simple("apiUrl", "api_url"),
		},
	}

	err := compareDestinationConfig(
		json.RawMessage(`{"apiKey":"secret"}`),
		`{"apiKey":"secret","apiUrl":"https://example.com"}`,
		cm,
	)

	require.ErrorContains(t, err, `missing field "apiUrl"`)
}
