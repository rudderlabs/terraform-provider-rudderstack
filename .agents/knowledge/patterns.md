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

## SDK-4941 — Default normalization for omitted API keys

- For integration fields that are `Optional` with a Terraform `Default`, avoid
  relying on `configs.Simple(...)` if control-plane `GET` responses may omit
  the API key; missing keys are not written into state by default and can cause
  repeated plan/apply churn.
- Prefer a custom `ConfigProperty.ToStateFunc` that writes the effective
  default into state when the API key is absent, while preserving explicit API
  values when present. Example applied in
  `rudderstack/integrations/destinations/destination_amplitude.go` for
  `apiVersion` <-> `api_version` with read fallback to `"v1"`.
