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

## DEX-376 — Inert Registry First Pattern

- `rudderstack/configs/registries.go` can accept new exported registry globals before runtime behavior changes; this repo allows "catalog-first" additions.
- Runtime effects require two separate follow-ons after a new registry appears: blank-import activation in `rudderstack/integrations/integrations.go` and resource-map expansion in `rudderstack/provider.go`.
- Scope slicing guideline from DEX-376: introduce registry surface area first, then layer integration registrations and provider exposure in later tickets to keep PR blast radius low.
