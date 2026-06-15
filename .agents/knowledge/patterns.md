# Patterns

> Recurring idioms specific to this repo (error handling, state management,
> retries, logging, DI, request lifecycle).
> Append-only. Agent-authored sections may optionally carry an HTML-comment tag
> (e.g., `<!-- pr:<id> -->`) identifying the writer/PR/run; human-authored
> sections are conventionally left untouched by automated runs.
> Every observed idiom includes a `file:line` reference.

## Terraform lifecycle

- CRUD handlers follow the same shape across source, destination, and
  connection resources: assert `m.(*Client)`, build an API model from state,
  call the shared client, set the new ID, then read back into state.
- Importers reuse the read path and convert any read diagnostic into a single
  import error, which keeps imported state consistent with normal refresh.
- `resourceSource.go`, `resourceDestination.go`, and
  `resourceConnection.go` all wrap API failures with `fmt.Errorf(...)` and then
  return `diag.FromErr`, so Terraform sees a consistent diagnostic surface.
- `resourceDestination.go::resourceDestinationCustomizeDiff` uses a diff-time
  validation hook instead of pushing all schema rules into the API client.

## State translation

- `configs.ConfigMeta::StateToAPI` and `::APIToState` apply an ordered chain of
  `ConfigProperty` transforms, so schema mapping is declarative rather than
  hand-written per field.
- `configs/configproperty.go::Simple`, `::Negated`, `::ArrayWithObjects`, and
  `::ArrayWithStrings` are the main mapping primitives; they cover direct
  copies, boolean inversion, and nested array/object reshaping.
- `configs/configproperty.go::SkipZeroValue` is used to omit empty API fields
  instead of serializing zero values into request payloads.
- `rudderstack/configs/validators.go::ValidateAll` composes multiple schema
  validators into one Terraform diag function.

## SDK-4941 — Amplitude web SDK version field

- Web (JavaScript) source-scoped destination settings are stored nested per
  source type in the control plane (`{ "web": "2" }`) and validated against that
  nested shape, even though config-backend flattens them to a top-level value
  when it delivers the source config to the JS SDK. Terraform writes the stored
  (nested) shape, so the field must map to the nested API key, not a flat one.
- The Amplitude destination already follows this for every web-scoped field
  (`forceHttps.web`, `trackGclid.web`, `preferAnonymousIdForDeviceId.web`, …).
  The SDK version field follows the same convention:
  `c.Simple("sdkVersion.web", "sdk_version.0.web", c.SkipZeroValue)` with a
  `TypeList`/`MaxItems: 1` block exposing a `web` string (`1` | `2`).
- API key is `sdkVersion`, HCL key is `sdk_version`, and the accepted values
  are `"1"` / `"2"`, matching rudder-integrations-config and rudder-sdk-js.
- No custom default-normalization `ToStateFunc` is needed: an omitted block
  resolves to version 1 via the integration schema default + the SDK fallback,
  and `SkipZeroValue` keeps an absent value out of the API payload.
- The inner `web` is `Required` (unlike the sibling *optional* web-scoped
  blocks): an empty `sdk_version {}` would otherwise serialize to nothing and
  cause a perpetual `+ sdk_version {}` plan diff. Requiring `web` turns that into
  a clear plan-time error. (Siblings share the footgun but are left as-is for
  back-compat.)

## INT-6562 — Snowflake Config Mapping Gates

- `GetCommonConfigMeta` contributes only `consent_management` mappings for declared source types; warehouse controls and other destination-level fields must be mapped in the destination integration itself.
- Snowflake `supportedSourceTypes` directly gates which nested `consent_management` source-type blocks are exposed in Terraform; omitted source types (for example Android Kotlin / iOS Swift) silently drop corresponding nested keys.
