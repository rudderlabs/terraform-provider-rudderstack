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

## SDK-5015 — Redacted Amplitude API secret in acceptance reads

- Live Amplitude destination reads from the RudderStack API redact the top-level `apiSecret` config key, so acceptance tests must not expect it during read-back comparison; the observed failure was `TestAccDestinationAmplitude` reporting `missing field "apiSecret": expected abc123` while the API response contained `apiKey` and defaults but no `apiSecret`.
- Keep `apiSecret` in `APICreate` and `APIUpdate` for mock-backed assertions that verify the provider sends the secret, but include `apiSecret` in `APIReadRedactedFields` and ignore the Terraform state path `config.0.api_secret` during import verification.
- Destination E2E CRUD can fail broadly when the live RudderStack API redacts sensitive destination credentials on GET, including API keys, credential blocks, REST API keys, passwords, private keys, and Amplitude `apiSecret`; the symptom is refreshed Terraform state losing the sensitive value after apply.
- Prefer the generic destination read-path fix over per-destination hard-coding: merge API-returned config with prior Terraform state for schema fields marked `Sensitive`, including nested credential blocks, and derive live read/import ignored sensitive paths from the destination schema.
