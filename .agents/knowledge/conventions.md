# Conventions

> Coding conventions and naming schemes - things a linter can't catch.
> Append-only. Agent-authored sections may optionally carry an HTML-comment tag
> (e.g., `<!-- pr:<id> -->`) identifying the writer/PR/run; human-authored
> sections are conventionally left untouched by automated runs.

## Naming and layout

- Provider resource names are `rudderstack_<kind>_<integration>`; the resource
  registry derives these from `configs.Sources.Entries()` and
  `configs.Destinations.Entries()`.
- Terraform-facing fields use `snake_case`, while API payload fields stay in
  the config metadata layer and are translated through `ConfigProperty`.
- Integration definitions live in `rudderstack/integrations/sources` and
  `.../destinations`; the package-level `init()` pattern is intentional and is
  what makes the blank import in `rudderstack/integrations/integrations.go`
  work.
- Generated documentation lives in `docs/resources/` and the source templates in
  `templates/resources/`; `make docs` is the expected regeneration path.
- Acceptance tests are named `TestAcc*` and are grouped by resource type in
  `rudderstack/integrations/{sources,destinations,connections}`.

## Provider behavior

- The provider keeps `created_at` and `updated_at` as RFC3339 strings in state.
- Provider configuration is environment-aware: `api_url` and `access_token`
  read from `RUDDERSTACK_API_URL` and `RUDDERSTACK_ACCESS_TOKEN` when unset.
- `configureClient` strips a trailing `/v2` from the base URL for backwards
  compatibility with older user configurations.
- Release versioning is synchronized across `Makefile`, `rudderstack/provider.go`,
  `examples/main.tf`, and the generated docs via release-please.

## INT-6562 — Additive Snowflake Parity for Backward Compatibility

- For Snowflake destination parity changes, prefer additive schema/mapping updates when config metadata still carries legacy consent-adjacent keys.
- Keep existing `consent_management` Terraform structure stable while reintroducing or retaining legacy fields (`one_trust_cookie_categories`, `ketch_consent_purposes`) if they remain in upstream destination schema, to avoid breaking existing configurations.
