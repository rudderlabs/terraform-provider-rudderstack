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

## INT-7014 — Customer.io live API drops new config keys

- Customer.io live acceptance and staging-smoke fixtures should not include `api_version = "v2"` or `user_id_mapping` until the backend integrations-config deployment persists and echoes API keys `apiVersion` and `userIdMapping`; otherwise CRUD tests can fail when create/read drops those keys from the API response.
- Keep mock/unit coverage for Terraform serialization of explicit `api_version = "v2"` and `user_id_mapping`, while relying on the default `api_version = "v1"` path to skip API writes and restore the default in Terraform state on reads.
