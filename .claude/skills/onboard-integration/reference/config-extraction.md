# Config Extraction Reference

Detailed rules for extracting integration metadata from the 3 config JSON files.

## From `db-config.json`

- `name` field → this becomes the `APIType` (e.g., `"WEBHOOK"`, `"SLACK"`)
- `supportedSourceTypes` → the list for `GetCommonConfigMeta()` and for generating `connectionMode` properties (destinations only)
- `config.secretKeys` → fields that need `Sensitive: true` in the Terraform schema
- `destConfig.defaultConfig` → field names that belong to this integration

## From `schema.json` (JSON Schema Draft-07)

- `properties` → each property becomes a Terraform schema field
- `type` → maps to Terraform type (see mapping table below)
- `pattern` → becomes `ValidateDiagFunc: c.StringMatchesRegexp(pattern)`
- `enum` → becomes `ValidateDiagFunc: c.StringMatchesRegexp("^(val1|val2|...)$")`
- `default` → becomes `Default: value`
- `required` array → fields listed are `Required: true`, all others are `Optional: true`

## From `ui-config.json`

- Field labels and descriptions → becomes `Description:` in Terraform schema
- Group structure → helps organize fields logically

## Type Mapping

| `schema.json` type | Terraform Go code |
|---|---|
| `"string"` | `schema.TypeString` |
| `"boolean"` | `schema.TypeBool` |
| `"integer"` / `"number"` | `schema.TypeInt` |
| `"object"` (single nested) | `schema.TypeList` + `MaxItems: 1` + `Elem: &schema.Resource{...}` |
| `"array"` + object `items` | `schema.TypeList` + `ConfigMode: schema.SchemaConfigModeAttr` + `Elem: &schema.Resource{...}` + `c.ArrayWithObjects(...)` |
| `"array"` + single nested key | `schema.TypeList` + `Elem: &schema.Schema{Type: schema.TypeString}` + `c.ArrayWithStrings(...)` |

## Fields to SKIP (handled by `GetCommonConfigMeta()`)

`GetCommonConfigMeta()` **only** handles consent management fields. Skip these:
- `consentManagement`
- `oneTrustCookieCategories`
- `ketchConsentPurposes`

## Source Type Specific Fields (NOT handled by `GetCommonConfigMeta`)

Many config fields are **per-source-type** — in the config JSON they appear with dot-notation like `fieldName.web`, `fieldName.android`, etc. These are **not** handled by `GetCommonConfigMeta()` and must be added manually. Common examples:

- `connectionMode.{sourceType}` — connection mode (cloud/device/hybrid)
- `useNativeSDK.{sourceType}` — whether to use native SDK (boolean)
- `eventUploadPeriodMillis.{sourceType}` — event batching settings
- Other destination-specific fields (e.g., `trackSessionEvents.android`, `capturePageView.web`, `enableLocationListening.android`)

Check `schema.json` for properties that have per-source-type sub-fields. For each one, use the dot-notation pattern:

**Properties** (one entry per source type per field):
```go
// connectionMode — one per supportedSourceType:
c.Simple("connectionMode.web", "connection_mode.0.web", c.SkipZeroValue),
c.Simple("connectionMode.android", "connection_mode.0.android", c.SkipZeroValue),
c.Simple("connectionMode.androidKotlin", "connection_mode.0.android_kotlin", c.SkipZeroValue),
// ... one entry per supportedSourceType from db-config.json

// useNativeSDK — only for source types that support device mode:
c.Simple("useNativeSDK.web", "use_native_sdk.0.web", c.SkipZeroValue),
c.Simple("useNativeSDK.android", "use_native_sdk.0.android", c.SkipZeroValue),

// Other per-source-type fields follow the same pattern:
c.Simple("eventUploadPeriodMillis.web", "event_upload_period_millis.0.web", c.SkipZeroValue),
```

**Schema** (each per-source-type field becomes a `TypeList` + `MaxItems: 1` block):
```go
"connection_mode": {
    Type:        schema.TypeList,
    MaxItems:    1,
    Optional:    true,
    Description: "Configure the connection mode for {name}.",
    Elem: &schema.Resource{
        Schema: map[string]*schema.Schema{
            "web": {
                Type:             schema.TypeString,
                Optional:         true,
                ValidateDiagFunc: c.StringMatchesRegexp("^(cloud|device|hybrid)$"),
            },
            // ... one entry per supportedSourceType
        },
    },
},
"use_native_sdk": {
    Type:        schema.TypeList,
    MaxItems:    1,
    Optional:    true,
    Description: "Enable native SDK for specific source types.",
    Elem: &schema.Resource{
        Schema: map[string]*schema.Schema{
            "web": {Type: schema.TypeBool, Optional: true},
            // ... one entry per relevant sourceType
        },
    },
},
```

Look at `schema.json` to determine which source types each field applies to and the allowed values/types. Study an existing destination with similar patterns (e.g., `destination_amplitude.go` for many per-source-type fields, `destination_bqstream.go` for connectionMode only).

## Naming Conventions

- API keys are camelCase (e.g., `webhookUrl`)
- Terraform keys are snake_case (e.g., `webhook_url`)
- Convert camelCase → snake_case for all field names
