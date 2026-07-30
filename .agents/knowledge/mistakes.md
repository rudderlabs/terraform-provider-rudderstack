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

## ACT2-463 — Pruned destination response comparisons

- Destination E2E CI failed after config-backend began pruning/masking destination response config because the acceptance helper compared GET response config directly against create/update request JSON. Symptoms included broad `API config verification failed` errors for missing or masked fields such as `apiSecret`, `apiKey`, `credentials`, `password`, and webhook `headers[].to`.
- The fix pattern is to keep create/update request assertions strict while deriving GET-response expectations from destination `ConfigMeta` and ignoring only fields mapped to sensitive Terraform schema values, rather than deleting request fixture fields or loosening config comparison globally.
