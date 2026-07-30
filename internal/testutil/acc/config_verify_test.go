package acc

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestCompareConfigIgnoresConfiguredMissingPaths(t *testing.T) {
	actual := json.RawMessage(`{
		"apiUrl": "https://example.com",
		"nested": { "name": "kept" },
		"items": [{ "name": "one" }]
	}`)

	expected := `{
		"apiUrl": "https://example.com",
		"apiKey": "secret",
		"nested": { "name": "kept", "secret": "nested-secret" },
		"items": [{ "name": "one", "token": "item-secret" }]
	}`

	err := compareConfig(actual, expected, "apiKey", "nested.secret", "items[0].token")
	if err != nil {
		t.Fatalf("expected ignored paths to pass comparison, got %v", err)
	}
}

func TestCompareConfigStillFailsForMissingNonIgnoredPath(t *testing.T) {
	actual := json.RawMessage(`{ "apiUrl": "https://example.com" }`)
	expected := `{ "apiUrl": "https://example.com", "nonSecret": "required" }`

	err := compareConfig(actual, expected, "apiKey")
	if err == nil || !strings.Contains(err.Error(), `missing field "nonSecret"`) {
		t.Fatalf("expected missing non-ignored field error, got %v", err)
	}
}

func TestWriteOnlyAPIConfigPaths(t *testing.T) {
	cm := configs.ConfigMeta{
		ConfigSchema: map[string]*schema.Schema{
			"api_key": {
				Type: schema.TypeString,
			},
			"event_key": {
				Type: schema.TypeString,
			},
			"event_key_map": {
				Type: schema.TypeList,
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"from": {
						Type: schema.TypeString,
					},
					"to": {
						Type: schema.TypeString,
					},
				}},
			},
			"api_url": {
				Type: schema.TypeString,
			},
			"headers": {
				Type:      schema.TypeList,
				Sensitive: true,
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"from": {
						Type: schema.TypeString,
					},
					"to": {
						Type: schema.TypeString,
					},
				}},
			},
			"auth": {
				Type: schema.TypeList,
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"access_key": {
						Type: schema.TypeString,
					},
					"id": {
						Type: schema.TypeString,
					},
				}},
			},
			"private_key": {
				Type:      schema.TypeString,
				Sensitive: true,
			},
		},
		Properties: []configs.ConfigProperty{
			configs.Simple("apiKey", "api_key"),
			configs.Simple("eventKey", "event_key"),
			configs.ArrayWithObjects("eventKeyMap", "event_key_map", map[string]interface{}{
				"from": "from",
				"to":   "to",
			}),
			configs.Simple("apiUrl", "api_url"),
			configs.Simple("headers", "headers"),
			configs.Simple("auth.accessKey", "auth.0.access_key"),
			configs.Simple("auth.id", "auth.0.id"),
			{
				FromStateFunc: func(config, state string) (string, error) {
					v := gjson.Get(state, "private_key")
					if !v.Exists() {
						return config, nil
					}
					return sjson.Set(config, "privateKey", "-----BEGIN PRIVATE KEY-----\n"+v.String()+"\n-----END PRIVATE KEY-----")
				},
				ToStateFunc: func(state, config string) (string, error) {
					return state, nil
				},
			},
		},
	}

	paths, err := writeOnlyAPIConfigPaths(cm)
	if err != nil {
		t.Fatalf("expected write-only path derivation to succeed, got %v", err)
	}

	expected := []string{"apiKey", "auth.accessKey", "eventKey", "headers", "privateKey"}
	if !reflect.DeepEqual(paths, expected) {
		t.Fatalf("expected paths %v, got %v", expected, paths)
	}
}

func TestDestinationWriteOnlyTerraformStatePaths(t *testing.T) {
	cm := configs.ConfigMeta{
		ConfigSchema: map[string]*schema.Schema{
			"api_key": {
				Type: schema.TypeString,
			},
			"event_key": {
				Type: schema.TypeString,
			},
			"headers": {
				Type:      schema.TypeList,
				Sensitive: true,
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"from": {
						Type: schema.TypeString,
					},
					"to": {
						Type: schema.TypeString,
					},
				}},
			},
			"key_based_authentication": {
				Type: schema.TypeList,
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"access_key_id": {
						Type:      schema.TypeString,
						Sensitive: true,
					},
					"access_key": {
						Type:      schema.TypeString,
						Sensitive: true,
					},
				}},
			},
		},
	}

	paths := destinationWriteOnlyTerraformStatePaths(cm)
	expected := []string{
		"config.0.api_key",
		"config.0.event_key",
		"config.0.headers",
		"config.0.key_based_authentication",
		"config.0.key_based_authentication.#",
		"config.0.key_based_authentication.0.%",
		"config.0.key_based_authentication.0.access_key",
		"config.0.key_based_authentication.0.access_key_id",
	}
	if !reflect.DeepEqual(paths, expected) {
		t.Fatalf("expected paths %v, got %v", expected, paths)
	}
}
