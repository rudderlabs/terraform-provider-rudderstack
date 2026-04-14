# Destination Code Templates

## File 1: `rudderstack/integrations/destinations/destination_{name}.go`

Follow this exact pattern:

```go
package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{...from db-config.json...}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)

	properties := []c.ConfigProperty{
		// For simple string/bool/int fields:
		c.Simple("apiKeyName", "terraform_key_name"),
		// For optional fields that should be omitted when empty:
		c.Simple("optionalField", "optional_field", c.SkipZeroValue),
		// For arrays of objects:
		c.ArrayWithObjects("apiArrayKey", "terraform_array_key", map[string]interface{}{
			"nestedApiKey1": "nested_tf_key1",
			"nestedApiKey2": "nested_tf_key2",
		}),
		// For arrays of strings with a wrapping key:
		c.ArrayWithStrings("apiArrayKey", "wrappingKey", "terraform_list_key"),
		// connectionMode — one per supportedSourceType:
		c.Simple("connectionMode.web", "connection_mode.0.web", c.SkipZeroValue),
		c.Simple("connectionMode.android", "connection_mode.0.android", c.SkipZeroValue),
		// ... etc for each supportedSourceType
	}
	properties = append(properties, commonProperties...)

	schema := map[string]*schema.Schema{
		"terraform_key_name": {
			Type:             schema.TypeString,
			Required:         true, // or Optional: true
			Sensitive:        true, // only if in secretKeys
			Description:      "Description from ui-config.json",
			ValidateDiagFunc: c.StringMatchesRegexp("pattern from schema.json"),
		},
		// TypeList with objects needs ConfigMode:
		"terraform_array_key": {
			Type:        schema.TypeList,
			Optional:    true,
			ConfigMode:  schema.SchemaConfigModeAttr,
			Description: "...",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"nested_tf_key1": {Type: schema.TypeString, Required: true},
					"nested_tf_key2": {Type: schema.TypeString, Required: true},
				},
			},
		},
		// TypeList with strings (no ConfigMode needed):
		"terraform_list_key": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "...",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		// connectionMode block:
		"connection_mode": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Configure the connection mode.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"web": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud|device|hybrid)$"),
					},
					// ... one per supportedSourceType
				},
			},
		},
	}

	for key, value := range commonSchema {
		schema[key] = value
	}

	c.Destinations.Register("{name}", c.ConfigMeta{
		APIType:      "{API_TYPE from db-config.json}",
		Properties:   properties,
		ConfigSchema: schema,
	})
}
```

## IMPORTANT Patterns

- `ConfigMode: schema.SchemaConfigModeAttr` is REQUIRED on all `TypeList` fields that use `ArrayWithObjects`
- Fields in `secretKeys` from db-config.json get `Sensitive: true`
- Required fields from schema.json `required` array get `Required: true`, all others get `Optional: true`
- Optional fields that should be omitted when zero-valued get `c.SkipZeroValue` in the property
- **Do NOT use `c.SkipZeroValue` on fields with a non-zero `Default`.** For example, a field with `Default: true` must NOT use `SkipZeroValue`, because `SkipZeroValue` checks Go's zero value (`false` for bool, `""` for string), which would silently drop explicit user values like `false`. Only use `SkipZeroValue` on Optional fields without a `Default`, or where the `Default` matches Go's zero value (e.g., `Default: false`, `Default: ""`).
- **Default values and APICreate:** If a field has a non-zero `Default` in the Terraform schema (e.g., `Default: "v2"`, `Default: true`), Terraform will always send that field — even on create. You MUST include such fields with their default values in `APICreate` in the test, not just in `APIUpdate`. Otherwise the test will fail because the actual API payload includes the defaults but your expected JSON doesn't.

## File 2: `rudderstack/integrations/destinations/destination_{name}_test.go`

```go
package destinations_test

import (
	"testing"

	acc "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/acc"
	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

// Extract test configs to a package-level var so both unit and E2E tests can reuse them.
var {camelCaseName}TestConfigs = []c.TestConfig{
	{
		TerraformCreate: `
			required_field = "value"
		`,
		APICreate: `{
			"requiredField": "value"
		}`,
		// NOTE: If any field has a non-zero Default (e.g. Default: "v2"),
		// include it in APICreate too — Terraform always sends defaults.
		// Example: APICreate: `{"requiredField": "value", "version": "v2"}`
		TerraformUpdate: `
			required_field = "value"
			optional_field = "updated"
			// ... all fields with values ...
			consent_management {
					web = [
						{
							provider = "oneTrust"
							consents = ["one_web", "two_web", "three_web"]
							resolution_strategy = ""
						},
						{
							provider = "ketch"
							consents = ["one_web", "two_web", "three_web"]
							resolution_strategy = ""
						},
						{
							provider = "custom"
							resolution_strategy = "and"
							consents = ["one_web", "two_web", "three_web"]
						}
					]
					android = [{
						provider = "ketch"
						consents = ["one_android", "two_android", "three_android"]
						resolution_strategy = ""
					}]
					ios = [{
						provider = "custom"
						resolution_strategy = "and"
						consents = ["one_ios", "two_ios", "three_ios"]
					}]
					// ... one block per supportedSourceType ...
					// pattern: web gets oneTrust+ketch+custom, android gets ketch, all others get custom
				}
			`,
			APIUpdate: `{
				"requiredField": "value",
				"optionalField": "updated",
				"consentManagement": {
					"web": [
						{
							"provider": "oneTrust",
							"resolutionStrategy": "",
							"consents": [
								{"consent": "one_web"},
								{"consent": "two_web"},
								{"consent": "three_web"}
							]
						},
						{
							"provider": "ketch",
							"resolutionStrategy": "",
							"consents": [
								{"consent": "one_web"},
								{"consent": "two_web"},
								{"consent": "three_web"}
							]
						},
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{"consent": "one_web"},
								{"consent": "two_web"},
								{"consent": "three_web"}
							]
						}
					],
					"android": [
						{
							"provider": "ketch",
							"resolutionStrategy": "",
							"consents": [
								{"consent": "one_android"},
								{"consent": "two_android"},
								{"consent": "three_android"}
							]
						}
					],
					"ios": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{"consent": "one_ios"},
								{"consent": "two_ios"},
								{"consent": "three_ios"}
							]
						}
					]
					// ... one entry per supportedSourceType with same pattern ...
				}
			}`,
	},
}

// Unit test — validates Terraform state ↔ API JSON conversion using mocked client.
func TestDestinationResource{PascalCaseName}(t *testing.T) {
	cmt.AssertDestination(t, "{name}", {camelCaseName}TestConfigs)
}

// E2E acceptance test — reuses the same test configs for real API validation.
func TestAccDestination{PascalCaseName}(t *testing.T) {
	acc.AccAssertDestination(t, "{name}", {camelCaseName}TestConfigs)
}
```

### Consent Management Test Pattern per Source Type

- `web` → 3 providers: oneTrust (resolutionStrategy=""), ketch (resolutionStrategy=""), custom (resolutionStrategy="and")
- `android` → 1 provider: ketch (resolutionStrategy="")
- `androidKotlin` → follows the `android` pattern: ketch (resolutionStrategy="")
- `ios` → 1 provider: custom (resolutionStrategy="and")
- `iosSwift` → follows the `ios` pattern: custom (resolutionStrategy="and")
- `unity` → 1 provider: custom (resolutionStrategy="or")
- All other source types → 1 provider: custom (resolutionStrategy="and")

### iosSwift / androidKotlin Naming

- In Terraform HCL: use `ios_swift` and `android_kotlin` (snake_case)
- In API JSON keys: use `iosSwift` and `androidKotlin` (camelCase)
- Consent value suffixes follow their base types: `"one_ios_swift"`, `"one_android_kotlin"`

**CRITICAL:** In the API JSON, consent strings are wrapped as `{"consent": "value"}` objects. In Terraform HCL, they are plain string lists.

## File 3: `examples/destination_{name}.tf`

```hcl
resource "rudderstack_destination_{name}" "example" {
  name = "my-{name}"

  config {
    required_field = "example-value"
    // show key optional fields too
  }
}
```

## File 4: `templates/resources/destination_{name}.md.tmpl`

```markdown
---
page_title: "rudderstack_destination_{name} Resource - terraform-provider-rudderstack"
subcategory: ""
description: |-

---

# rudderstack_destination_{name} (Resource)

This resource represents a {Display Name} destination. For more information check
https://www.rudderstack.com/docs/destinations/{name-with-hyphens}/

## Example Usage

{{tffile "examples/destination_{name}.tf"}}

{{ .SchemaMarkdown | trimspace }}
```
