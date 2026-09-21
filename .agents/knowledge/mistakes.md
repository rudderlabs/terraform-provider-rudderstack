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

## INT-7142 — Singular and Impact live CRUD acceptance failures

- Singular and impact.com destination acceptance tests failed in the shared E2E workspace on Step 1 create with `could not create destination: http status code: 500` when using dummy destination credentials.
- Reviewer guidance clarified that transient shared-workspace 5xx/API create failures should be handled by rerunning CI, not by adding permanent `if !acc.PlanOnly() { t.Skip(...) }` guards. Reserve plan-only live-CRUD skips for destinations with genuine documented prerequisites the shared workspace cannot satisfy, such as vendor OAuth account requirements.
- Preserve mock `AssertDestination` converter coverage and `TF_ACC_PLAN_ONLY=1` schema validation when live CRUD is unreliable.

## RUD-3134 — OAuth account type namespace mismatch

<!-- session: 2026-09-21 -->

- The OAuth acceptance tests initially passed lowercase Terraform registry keys such as `linkedin_ads` to `internal/testutil/acc/oauth_destinations.go::resolveOAuthAccountID`; CI then failed with `no OAuth account found` because `client.Account.Definition.Type` uses the destination's uppercase `ConfigMeta.APIType`. Pass `LINKEDIN_ADS`, `BINGADS_OFFLINE_CONVERSIONS`, and `GOOGLE_ADWORDS_OFFLINE_CONVERSIONS` while keeping the lowercase destination key as the separate resource argument.
- `configs.TestConfig.APIResponseOptionalFields` only permits an omitted field during `acc.compareConfig`; Terraform SDK acceptance tests still run an automatic post-apply plan and fail if configured state does not round-trip. Live fixtures in `rudderstack/integrations/destinations/*_test.go` must omit backend-discarded options, keep immutable fields unchanged between create/update, or explicitly configure an omitted boolean default to the API read-back value.
