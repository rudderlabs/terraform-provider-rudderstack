# Mistakes

> Post-mortem entries from observed failures: CI failures, reverts on prior PRs,
> prod incidents. Accrues over time - bootstrap leaves this empty.
> Append-only. Agent-authored sections may optionally carry an HTML-comment tag
> (e.g., `<!-- pr:<id> -->`) identifying the writer/PR/run; human-authored
> sections are conventionally left untouched by automated runs.

## RUD-2790 — Terraform auto-install blocked tests

- `make lint` passes, but `go test ./...` currently fails in
  `rudderstack`, `rudderstack/integrations/destinations`,
  `rudderstack/integrations/sources`, and `rudderstack/retl` because Terraform
  CLI auto-install returns HTTP 403 before provider tests execute.

## DAW-4034 — Sensitive destination fields filtered from API read-back

- Destination acceptance checks compare API read-back config against APICreate/APIUpdate fixtures; after public API filtering, secret destination config fields such as `apiKey`, `credentials`, and `apiToken` are intentionally absent from API responses and can trigger false failures like `API config verification failed: missing field "apiKey"`.
- Keep read-back comparisons strict for non-sensitive fields, but derive sensitive API config paths from each destination's Terraform schema and ignore only those sensitive paths during destination acceptance verification.
