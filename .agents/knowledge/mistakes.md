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

## INT-7014 — Customer.io new config keys depend on backend rollout

- Customer.io live acceptance and staging-smoke fixtures may only include `api_version = "v2"` and `user_id_identifier_type` once the target backend's integrations-config deployment persists and echoes API keys `apiVersion` and `userIdIdentifierType`. If the backend drops them on create/read, the post-apply drift assertion (`terraform plan -detailed-exitcode`) fails with a perpetual diff. As of the INT-7014 rename these keys are on rudder-integrations-config `develop` but not yet released to `main`, so an environment tracking released config can still fail; fall back to mock/unit coverage for serialization there.
- Keep mock/unit coverage for Terraform serialization of explicit `api_version` and `user_id_identifier_type`. `api_version` is Optional with schema `Default: "v1"`, so an omitted attribute still resolves to `v1` and is written to the API config as `apiVersion`; the round-trip stays symmetric because the backend echoes that value back. `user_id_identifier_type` is Optional with no default and uses `c.SkipZeroValue`, so it is omitted from the payload entirely when unset.
