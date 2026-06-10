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

## DEX-377 — Narrow service seam for pre-client resources

- When resource work must proceed before `rudderstack.Client` grows a concrete field, use a narrow interface seam: define only the CRUD methods needed by the resource, expose it through a package-local accessor on `*Client`, and keep the shim isolated in a deletable file.
- Keep test injection local to the seam (for example, a test-only setter keyed by client instance) so tests can provide fakes without widening provider wiring or importing not-yet-stable client symbols.
- Reuse the existing RETL seam style (`rudderstack/retl/common.go` pattern): resource code depends on the minimal interface, while provider construction remains unchanged.
