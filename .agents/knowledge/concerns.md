# Concerns

> Technical debt, TODOs, FIXMEs, security concerns, architectural issues.
> Append-only. Agent-authored sections may optionally carry an HTML-comment tag
> (e.g., `<!-- pr:<id> -->`) identifying the writer/PR/run; human-authored
> sections are conventionally left untouched by automated runs.
> Top-5-8 highest-signal items per category, not exhaustive.

## TODO and fixme density

- `cmd/generatetf/main.go::getAPIRetlSources` has a TODO noting the RETL SDK
  list call is not paginated yet, which can silently truncate large workspaces.
- No other strong TODO/FIXME cluster surfaced in the sampled provider code; the
  outstanding work is concentrated in the generator path.

## Security

- `rudderstack/provider.go::NewWithConfigureClientFunc` exposes
  `access_token` as a plain optional provider field without `Sensitive: true`,
  so users who set it in HCL risk persisting the secret in Terraform state.
- README and generated docs show the token in example provider blocks, which
  reinforces the HCL-as-secret-storage risk if copied verbatim.
- `cmd/generatetf/main.go::setupClient` reads the access token from
  `RUDDERSTACK_ACCESS_TOKEN`, which is the safer path; keep that preference
  visible in future docs and examples.
- Acceptance and bootstrap scripts consume API credentials from the environment;
  that is better than inline HCL, but it still warrants careful secret handling
  in CI logs and local shell history.

## Architectural smells

- `rudderstack/configs/registries.go::Registry.Register` panics on duplicate
  names, so a bad integration registration fails at init time instead of as a
  typed configuration error.
- `rudderstack/configs/validators.go::StringMatchesRegexp` and
  `::StringNotMatchesRegexp` also panic on invalid regex input, which makes
  bad validation constants a process-level failure.
- `resourceConnection.go::resourceConnectionUpdate` returns
  `could not create source` on update failure, which is a copy/paste error that
  weakens diagnostics.
- `resourceSource.go::resourceSourceImportState`,
  `resourceDestination.go::resourceDestinationImportState`, and
  `resourceConnection.go::resourceConnectionImportState` all report
  `could not import connection` even when importing a source or destination.
- `cmd/generatetf/generator/generator.go::GenerateTerraform` carries a large
  amount of integration-specific filtering logic, so the generator is tightly
  coupled to current resource semantics.

## RUD-2790 — Provider token not marked sensitive

- `rudderstack/provider.go` exposes `access_token` as an optional provider
  schema field but does not mark it `Sensitive: true`, so tokens entered in HCL
  can land in Terraform state and plan output.

## Stale code and deps

- No clearly stale dependency cluster surfaced in the sampled manifest files.
- The main stale-code signal is the pagination TODO in
  `cmd/generatetf/main.go::getAPIRetlSources`; that is a known functional gap,
  not just dead code.
