package acc

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

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
		},
		Properties: []configs.ConfigProperty{
			configs.Simple("apiKey", "api_key"),
			configs.Simple("eventKey", "event_key"),
			configs.Simple("apiUrl", "api_url"),
			configs.Simple("headers", "headers"),
			configs.Simple("auth.accessKey", "auth.0.access_key"),
			configs.Simple("auth.id", "auth.0.id"),
		},
	}

	paths, err := writeOnlyAPIConfigPaths(cm)
	if err != nil {
		t.Fatalf("expected write-only path derivation to succeed, got %v", err)
	}

	expected := []string{"apiKey", "auth.accessKey", "headers"}
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
		},
	}

	paths := destinationWriteOnlyTerraformStatePaths(cm)
	expected := []string{"config.0.api_key", "config.0.headers"}
	if !reflect.DeepEqual(paths, expected) {
		t.Fatalf("expected paths %v, got %v", expected, paths)
	}
}
