# Onboard Integration

Add a new source or destination integration to the terraform-provider-rudderstack.

**Usage:** `/onboard-integration [name] [source|destination]`

Arguments from `$ARGUMENTS`: $ARGUMENTS

---

## Step 0: Gather Inputs

Parse `$ARGUMENTS` for the integration name and type (source/destination). If either is missing, ask the user:

1. **Integration name** — the snake_case name (e.g., `webhook`, `slack`, `google_analytics`).
2. **Type** — `source` or `destination`.
3. **Path to `rudder-integrations-config` repo** — Before asking the user, auto-detect by checking for a sibling folder:
   ```bash
   ls -d ../rudder-integrations-config 2>/dev/null
   ```
   If found, use that path and tell the user: "Found `rudder-integrations-config` at `{resolved_path}`, using it."
   If NOT found, ask: "What is the absolute path to your local `rudder-integrations-config` repo?" If the user says they don't have it, skip to the **Fallback** section below.

Once you have the path, check for these 3 JSON config files:
- `{path}/src/configurations/{destinations|sources}/{name}/db-config.json`
- `{path}/src/configurations/{destinations|sources}/{name}/schema.json`
- `{path}/src/configurations/{destinations|sources}/{name}/ui-config.json`

Read ALL THREE files. If any are missing, tell the user which are missing and fall back to interactive Q&A.

### Check for Existing Integration with Similar Name

Before proceeding, check if an integration with the same or a similar name already exists in the provider:

```bash
# For destinations:
ls rudderstack/integrations/destinations/destination_*{name}*.go 2>/dev/null
# For sources — check sources.go for a matching Register call:
grep -i '{name}' rudderstack/integrations/sources/sources.go 2>/dev/null
```

If a match is found, **stop and ask the user** what they'd like to do:

- "I found an existing integration that looks similar: `{matched_file_or_name}`. What would you like to do?
  1. **Add new fields** to the existing integration (update the `.go`, test, example, and docs files)
  2. **Refactor/fix** the existing integration (e.g., fix types, update descriptions, correct test data)
  3. **Create a brand new integration** — the name is similar but it's a different integration
  4. **Something else** — describe what you need"

Do NOT proceed until the user confirms.

If the user chooses to **add new fields** or **update** the existing integration:

1. Read the current implementation files (`.go`, `_test.go`, example `.tf`, and docs template).
2. Read the config JSON files (`db-config.json`, `schema.json`, `ui-config.json`) from `rudder-integrations-config`.
3. Compare the two to identify fields that exist in the config JSON but are **missing from the current `.go` implementation**. Also note fields where the type, required/optional status, default value, or validation has changed.
4. Present the diff to the user in a clear table or list, e.g.:

   "I found the following new/changed fields comparing `rudder-integrations-config` with the current implementation:

   | # | Field (API key) | Type | Status | Change |
   |---|---|---|---|---|
   | 1 | `newApiField` | string | optional | **New** — not in current implementation |
   | 2 | `existingField` | boolean → string | required | **Changed** — type changed |
   | 3 | `anotherField` | string | optional | **New** — not in current implementation |

   Would you like to add **all** of these fields, or select specific ones? (Enter `all` or comma-separated numbers like `1,3`)"

5. Wait for the user to respond. Only add/modify the fields the user selects.
6. Proceed with the implementation, modifying only the selected fields in the `.go`, `_test.go`, example `.tf`, and docs template files.

---

## Step 1: Extract Integration Metadata from Config Files

From the 3 JSON files, extract:

### From `db-config.json`:
- `name` field → this becomes the `APIType` (e.g., `"WEBHOOK"`, `"SLACK"`)
- `supportedSourceTypes` → the list for `GetCommonConfigMeta()` (destinations only)
- `config.secretKeys` → fields that need `Sensitive: true` in the Terraform schema
- `destConfig.defaultConfig` → field names that belong to this integration

### From `schema.json` (JSON Schema Draft-07):
- `properties` → each property becomes a Terraform schema field
- `type` → maps to Terraform type (see mapping table below)
- `pattern` → becomes `ValidateDiagFunc: c.StringMatchesRegexp(pattern)`
- `enum` → becomes `ValidateDiagFunc: c.StringMatchesRegexp("^(val1|val2|...)$")`
- `default` → becomes `Default: value`
- `required` array → fields listed are `Required: true`, all others are `Optional: true`

### From `ui-config.json`:
- Field labels and descriptions → becomes `Description:` in Terraform schema
- Group structure → helps organize fields logically

### Type Mapping:

| `schema.json` type | Terraform Go code |
|---|---|
| `"string"` | `schema.TypeString` |
| `"boolean"` | `schema.TypeBool` |
| `"integer"` / `"number"` | `schema.TypeInt` |
| `"object"` (single nested) | `schema.TypeList` + `MaxItems: 1` + `Elem: &schema.Resource{...}` |
| `"array"` + object `items` | `schema.TypeList` + `ConfigMode: schema.SchemaConfigModeAttr` + `Elem: &schema.Resource{...}` + `c.ArrayWithObjects(...)` |
| `"array"` + single nested key | `schema.TypeList` + `Elem: &schema.Schema{Type: schema.TypeString}` + `c.ArrayWithStrings(...)` |

### Fields to SKIP (handled by `GetCommonConfigMeta()`):
- `consentManagement`
- `oneTrustCookieCategories`
- `ketchConsentPurposes`
- `connectionMode`

### Naming Conventions:
- API keys are camelCase (e.g., `webhookUrl`)
- Terraform keys are snake_case (e.g., `webhook_url`)
- Convert camelCase → snake_case for all field names

---

## Step 1.5: Study a Similar Existing Integration

Before generating files, find an existing destination/source with similar field types to the new integration and read its implementation and tests as a reference. This helps catch patterns the generic instructions might miss.

- **Simple fields only?** → Study `destination_webhook.go` and its test
- **Arrays of objects?** → Study `destination_slack.go` and its test
- **Event filtering / complex nested config?** → Look for a destination with similar structure in `rudderstack/integrations/destinations/`

Read both the `.go` file and the `_test.go` file to understand the exact patterns used.

---

## Step 2: Generate Files

### For DESTINATIONS — create these files:

#### File 1: `rudderstack/integrations/destinations/destination_{name}.go`

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

**IMPORTANT patterns:**
- `ConfigMode: schema.SchemaConfigModeAttr` is REQUIRED on all `TypeList` fields that use `ArrayWithObjects`
- Fields in `secretKeys` from db-config.json get `Sensitive: true`
- Required fields from schema.json `required` array get `Required: true`, all others get `Optional: true`
- Optional fields that should be omitted when zero-valued get `c.SkipZeroValue` in the property
- **Do NOT use `c.SkipZeroValue` on fields with a non-zero `Default`.** For example, a field with `Default: true` must NOT use `SkipZeroValue`, because `SkipZeroValue` checks Go's zero value (`false` for bool, `""` for string), which would silently drop explicit user values like `false`. Only use `SkipZeroValue` on Optional fields without a `Default`, or where the `Default` matches Go's zero value (e.g., `Default: false`, `Default: ""`).
- **Default values and APICreate:** If a field has a non-zero `Default` in the Terraform schema (e.g., `Default: "v2"`, `Default: true`), Terraform will always send that field — even on create. You MUST include such fields with their default values in `APICreate` in the test, not just in `APIUpdate`. Otherwise the test will fail because the actual API payload includes the defaults but your expected JSON doesn't.

#### File 2: `rudderstack/integrations/destinations/destination_{name}_test.go`

```go
package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResource{PascalCaseName}(t *testing.T) {
	cmt.AssertDestination(t, "{name}", []c.TestConfig{
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
	})
}
```

**Consent management test pattern per source type:**
- `web` → 3 providers: oneTrust (resolutionStrategy=""), ketch (resolutionStrategy=""), custom (resolutionStrategy="and")
- `android` → 1 provider: ketch (resolutionStrategy="")
- `androidKotlin` → follows the `android` pattern: ketch (resolutionStrategy="")
- `ios` → 1 provider: custom (resolutionStrategy="and")
- `iosSwift` → follows the `ios` pattern: custom (resolutionStrategy="and")
- `unity` → 1 provider: custom (resolutionStrategy="or")
- All other source types → 1 provider: custom (resolutionStrategy="and")

**iosSwift / androidKotlin naming:**
- In Terraform HCL: use `ios_swift` and `android_kotlin` (snake_case)
- In API JSON keys: use `iosSwift` and `androidKotlin` (camelCase)
- Consent value suffixes follow their base types: `"one_ios_swift"`, `"one_android_kotlin"`

**CRITICAL:** In the API JSON, consent strings are wrapped as `{"consent": "value"}` objects. In Terraform HCL, they are plain string lists.

#### File 3: `examples/destination_{name}.tf`

```hcl
resource "rudderstack_destination_{name}" "example" {
  name = "my-{name}"

  config {
    required_field = "example-value"
    // show key optional fields too
  }
}
```

#### File 4: `templates/resources/destination_{name}.md.tmpl`

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

### For SOURCES — modify existing files + create new ones:

#### Modify: `rudderstack/integrations/sources/sources.go`

Add a new `c.Sources.Register(...)` block at the end of the `init()` function. Most sources use `SkipConfig: true` with empty properties:

```go
c.Sources.Register("{name}", c.ConfigMeta{
	APIType:    "{APIType from db-config.json}",
	Properties: []c.ConfigProperty{},
	SkipConfig: true,
})
```

#### Modify: `rudderstack/integrations/sources/sources_test.go`

Add a new test function at the end of the file. Note: sources use `configs` (NOT aliased as `c`):

```go
func TestSourceResource{PascalCaseName}(t *testing.T) {
	cmt.AssertSource(t, "{name}", []configs.TestConfig{configs.EmptyTestConfig})
}
```

#### Create: `examples/source_{name}.tf`

```hcl
resource "rudderstack_source_{name}" "example" {
  name = "example-{name}"
}
```

#### Create: `templates/resources/source_{name}.md.tmpl`

```markdown
---
page_title: "rudderstack_source_{name} Resource - terraform-provider-rudderstack"
subcategory: ""
description: |-

---

# rudderstack_source_{name} (Resource)

This resource represents a {Display Name} event stream source. For more information check https://www.rudderstack.com/docs/sources/{name-with-hyphens}/

## Example Usage

{{tffile "examples/source_{name}.tf"}}

{{ .SchemaMarkdown | trimspace }}
```

---

## Step 3: Run Unit Tests

Run the unit test for the new integration:

```bash
# For destinations:
go test ./rudderstack/integrations/destinations/ -run TestDestinationResource{PascalCaseName} -v

# For sources:
go test ./rudderstack/integrations/sources/ -run TestSourceResource{PascalCaseName} -v
```

If tests fail, analyze the error output carefully:
- Schema mismatches → check field types and required/optional
- JSON key mismatches → check camelCase API keys vs snake_case terraform keys
- Consent management errors → ensure the API JSON wraps consents as `{"consent": "value"}` objects

Fix any failures and re-run until tests pass.

---

## Step 4: Generate Documentation

Run:
```bash
make docs
```

This generates `docs/resources/{destination|source}_{name}.md` from the template. If `make docs` fails, check:
- The template file exists at the correct path
- The example .tf file is valid HCL
- `tfplugindocs` is installed (run `go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest` if needed)

---

## Step 5: Run Full Test Suite

Run the full test suite to make sure nothing is broken:

```bash
go test ./... -v
```

If there are failures in other tests, investigate — the new integration should NOT affect existing tests since it self-registers via `init()`.

---

## Step 6: Lint Check

Run the linter:

```bash
make lint
```

Fix any lint issues in the generated code.

---

## Step 7: E2E Testing

Automatically proceed with E2E testing. Do NOT ask the user whether to run E2E tests — just run them.

1. **Load the access token from `.env`:**
   ```bash
   cat .env 2>/dev/null | grep RUDDERSTACK_ACCESS_TOKEN || echo "NOT_FOUND"
   ```
   - If `.env` contains `RUDDERSTACK_ACCESS_TOKEN`, use that value automatically. Tell the user: "Using access token from `.env` file."
   - If `.env` is missing or doesn't contain the token, ask the user: "No `RUDDERSTACK_ACCESS_TOKEN` found in `.env`. Please provide your access token."
2. **Check for `dev_overrides` in `~/.terraformrc`:**
   ```bash
   cat ~/.terraformrc 2>/dev/null || echo "No .terraformrc found"
   ```
   - If `dev_overrides` exist for `rudderlabs/rudderstack`, note the override path and source.
   - If overrides exist, you must use the overridden source in the `required_providers` block and **skip `terraform init`** (it fails with dev_overrides). Instead, copy the built binary to the override path.
3. Build and install the provider locally:
   ```bash
   make install
   ```
   If dev_overrides exist, also copy the binary to the override path:
   ```bash
   cp $(go env GOPATH)/bin/terraform-provider-rudderstack {override_path}/
   ```
4. Create a temporary workspace:
   ```bash
   mkdir -p /tmp/tf-test-{name} && cd /tmp/tf-test-{name}
   ```
5. Write a `main.tf` with the provider config and a resource:
   ```hcl
   terraform {
     required_providers {
       rudderstack = {
         source = "rudderlabs/rudderstack"  # use override source if dev_overrides exist
       }
     }
   }

   provider "rudderstack" {}

   resource "rudderstack_{type}_{name}" "test" {
     name = "e2e-test-{name}"
     config {
       // ALL config fields (required AND optional) with realistic test values.
       // This ensures the E2E test validates every field mapping, not just required ones.
       // Use the same field values from TerraformUpdate in the unit test.
     }
   }
   ```
6. Run:
   Use the token loaded from `.env` (or provided by the user) as an environment variable prefix for all terraform and verify commands:
   ```bash
   # Skip `terraform init` if dev_overrides are active — it will fail.
   # Only run `terraform init` if there are NO dev_overrides.
   RUDDERSTACK_ACCESS_TOKEN="{token}" terraform init  # skip if dev_overrides
   RUDDERSTACK_ACCESS_TOKEN="{token}" terraform plan
   RUDDERSTACK_ACCESS_TOKEN="{token}" terraform apply -auto-approve
   ```
7. Verify the resource was created:
   ```bash
   terraform show
   ```
   Check that the resource has a valid ID.
8. **Run the verify script** to deterministically compare the .tf config against the API:
   ```bash
   RESOURCE_ID=$(terraform show -json | jq -r '.values.root_module.resources[0].values.id')
   cd <terraform-provider-rudderstack repo path>
   go run ./cmd/integration-verify/ -file /tmp/tf-test-{name}/main.tf -id "$RESOURCE_ID"
   ```
   This performs a subset comparison: every config key from the .tf file must exist and match in the API response. If the verify script reports FAIL, investigate the differences before proceeding.
9. **Ask before cleaning up:** Ask the user: "Would you like to verify the resource from the RudderStack dashboard first, or can I go ahead and delete it?" Wait for confirmation before proceeding. Do NOT proceed with destroy until the user explicitly confirms.
10. Clean up (only after user confirms):
   ```bash
   terraform destroy -auto-approve
   cd /
   rm -rf /tmp/tf-test-{name}
   ```

---

## Fallback: Interactive Q&A (when config files are unavailable)

If the `rudder-integrations-config` repo is not available or config files are missing, gather information interactively:

1. **API Type**: "What is the API type name? (e.g., `WEBHOOK`, `SLACK`, `GOOGLE_ANALYTICS`)"
2. **Supported Source Types**: "Which source types are supported? Common set: `web, android, ios, unity, reactnative, flutter, cordova, amp, cloud, warehouse, shopify`"
3. **Fields**: For each field, ask:
   - Field name (camelCase API name)
   - Type (string, boolean, integer, array of strings, array of objects)
   - Required or optional?
   - Sensitive? (contains secrets/keys)
   - Description
   - Validation pattern (regex) or allowed values (enum)
   - Default value (if any)
   - For arrays: what are the nested field names and types?
4. Generate the files following the same patterns above.

---

## Summary Checklist

Before finishing, verify:

- [ ] All files created/modified
- [ ] `init()` function registers the integration (self-registering, no provider.go changes needed)
- [ ] Unit tests pass
- [ ] `make docs` generates the docs file
- [ ] Full test suite passes (`go test ./...`)
- [ ] Lint passes (`make lint`)
- [ ] Code follows existing patterns exactly (imports, naming, structure)

Present the user with a summary of all files created/modified and test results.
