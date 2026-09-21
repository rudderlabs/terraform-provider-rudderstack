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

## RUD-3134 — Missing OAuth accounts in the dev workspace

<!-- session: 2026-09-21 -->

- `client.Account.Definition.Type` uses the lowercase Terraform destination key, such as `linkedin_ads`, `bingads_offline_conversions`, or `google_adwords_offline_conversions`. Earlier `no OAuth account found` failures from `internal/testutil/acc/oauth_destinations.go::resolveOAuthAccountID` occurred because those accounts were not yet connected in the dev workspace, not because the lookup used the wrong namespace.
- `configs.TestConfig.APIResponseOptionalFields` only relaxes the explicit live API comparison; it does not bypass Terraform SDK acceptance tests' automatic post-apply empty-plan check. Live fixtures must round-trip through API reads: omit unsupported explicit options, avoid changing immutable values on update, and use the backend's read-back zero value for omitted boolean defaults when outbound serialization still needs coverage.
