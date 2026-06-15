# Architecture

> Component layout, internal relationships, data flow.
> Append-only. Agent-authored sections may optionally carry an HTML-comment tag
> (e.g., `<!-- pr:<id> -->`) identifying the writer/PR/run; human-authored
> sections are conventionally left untouched by automated runs.

## Runtime topology

- `main.go::main` is a thin Terraform plugin entry point: it switches between
  normal `plugin.Serve` and debug mode, then hands provider construction to
  `rudderstack.New()`.
- `rudderstack/provider.go::NewWithConfigureClientFunc` owns provider wiring:
  provider-level config, the resource map, and the client factory used by all
  CRUD handlers.
- `rudderstack/integrations/integrations.go` only exists for blank imports so
  package `init()` functions can register source and destination metadata.
- `rudderstack/configs/registries.go::Registry` is the central catalog for
  integration metadata; `configs.Sources` and `configs.Destinations` are the
  authoritative lookups used to expand provider resources.
- `rudderstack/client.go::NewAPIClient` adapts the shared `rudder-iac` client
  into the provider-specific wrapper and adds RETL support.
- `rudderstack/resource_source.go::resourceSource`,
  `rudderstack/resource_destination.go::resourceDestination`, and
  `rudderstack/resource_connection.go::resourceConnection` implement the core
  Terraform resources and translate plan/apply into Public API CRUD calls.

## Generation path

- `cmd/generatetf/main.go::main` pulls live sources, destinations, connections,
  and RETL resources from the API and prints either HCL or import commands.
- `cmd/generatetf/generator/generator.go::GenerateTerraform` and
  `::GenerateImportScript` keep generated HCL aligned with the provider schema
  and skip unsupported integration shapes.
- `rudderstack/configs/configmeta.go::ConfigMeta` and
  `rudderstack/configs/configproperty.go::*` define the state/API translation
  layer that both the provider and generator reuse.
- `templates/resources/*.md.tmpl` and `docs/resources/*.md` are generated-doc
  surfaces for the provider schema and examples; `Makefile::docs` keeps them in
  sync.

## Cross-cutting

- The same registry metadata drives provider resources, generated docs, and the
  HCL snapshot generator, so integration additions usually require touching
  `rudderstack/integrations/*`, `rudderstack/configs/*`, and the generated docs
  together.
- Secret handling is split between provider config and environment variables:
  `access_token` is accepted at the provider level, while the generator and
  bootstrap scripts prefer `RUDDERSTACK_ACCESS_TOKEN`; see `patterns.md` and
  `concerns.md`.
- State/API translation is the shared abstraction boundary across CRUD and
  generation, which is why `ConfigMeta` and `ConfigProperty` are more central
  than the resource handlers themselves.

## INT-6562 — Snowflake vs Snowflake Streaming Scope Boundary

- The provider maintains separate destination integrations for standard Snowflake and Snowflake Streaming, each with its own implementation/tests/docs surfaces.
- When task input points to `src/configurations/destinations/snowflake` (non-streaming), changes should be scoped to the standard Snowflake destination unless explicitly requested otherwise.
