# E2E Testing — Operational Guide

Operational reference for the team that owns the Terraform provider's E2E suite. If you're a contributor adding tests for a new integration, see the contributor section in [`README.md`](README.md#e2e-testing) and the templates in `.claude/skills/onboard-integration/reference/e2e-testing.md`. This document covers running and debugging the suite plus the credentials model in §5; the step-by-step PAT rotation procedure lives in the platform team's runbook.

---

## 1. Scope — what "E2E" really means here

Tests in `internal/testutil/acc/` exercise:

- The **real** RudderStack control-plane API (Sources, Destinations, Connections) via `rudder-iac/api/client`.
- The **real** Terraform provider schema, lifecycle (Create → Update → Import → Destroy), and HCL parsing via `terraform-plugin-sdk/v2`'s `resource.Test` harness.
- **Subset verification** of the resource config returned by the API against the expected JSON in each integration's `[]configs.TestConfig` (see `config_verify.go` — extra fields the API adds are tolerated; missing fields fail the test).

Tests do **not** exercise:

- Downstream event delivery. Destinations are created with placeholder credentials (e.g. `webhook_url = "https://example.com/test"`); no events flow to the configured third party.
- Third-party OAuth flows / vendor APIs. Only the control-plane CRUD path is verified.
- Hard-delete semantics. The control-plane API soft-deletes sources and destinations — `Get` after destroy may still return 200, so `CheckDestroy` is best-effort for those two resource types. Connection destroy *is* verified strictly (`testAccCheckConnectionDestroy` fails if `Get` succeeds).

There are two run modes:

| Mode | Env | What runs | API calls per integration |
|---|---|---|---|
| **Plan-only** | `TF_ACC=1 TF_ACC_PLAN_ONLY=1` | HCL parse + provider schema validation | 0 (dummy token auto-set) |
| **Full CRUD** | `TF_ACC=1` | Create → Update → Import → Destroy + config-subset checks | ~10+ (Create/Update/Delete plus multiple Gets — Terraform refreshes state after each apply, and the acceptance helpers issue separate Gets for exists + config + settings + destroy verification) |

---

## 2. What's covered

Every registered destination and source has a `TestAcc*` function — enforced at build time by `internal/testutil/acc/coverage_test.go` (run as part of `make test-ci`). Connection coverage is intentionally narrow (two source→destination wiring tests); per-integration behavior is covered by the source/destination tests.

To list the current coverage:
```bash
grep -rh --include='*.go' "^func TestAcc" rudderstack/integrations \
  | sed 's/func //; s/(t.*//' | sort
```

**Plan-only exceptions.** A small number of tests `t.Skip` in full-CRUD mode because they need a vendor account that can't be provisioned generically. `TestAccDestinationLinkedinAds` is currently the only one (requires a valid LinkedIn OAuth account in the workspace) — it runs in `TF_ACC_PLAN_ONLY=1` mode but is skipped under `TF_ACC=1`. To find all such skips:
```bash
grep -rB1 -A2 --include='*_test.go' "acc.PlanOnly()" rudderstack/integrations
```

---

## 3. How a run executes

`acc.AccAssertDestination` / `AccAssertSource` / `AccAssertConnection` all follow the same shape (see `internal/testutil/acc/destinations.go`, `sources.go`, `connections.go`):

1. **Provider factory** — `TestAccProviderFactories` constructs the provider with `rudderstack.New()`. The provider reads `RUDDERSTACK_ACCESS_TOKEN` and `RUDDERSTACK_API_URL` from env.
2. **Pre-check** — In full CRUD mode, `TestAccPreCheck` aborts the test if `RUDDERSTACK_ACCESS_TOKEN` is unset. In plan-only mode, `ensureDummyToken` sets `RUDDERSTACK_ACCESS_TOKEN=plan-only-dummy-token` via `os.Setenv` (chosen over `t.Setenv` so `t.Parallel()` still works).
3. **Random name** — `RandomName("<integration>")` prefixes resources with `tf-acc-` and a random 62-bit int, so leftover resources from a failed run can be identified and cleaned up by name prefix.
4. **Test steps** — destinations and sources run the full lifecycle:
   - Step 1: Apply `TerraformCreate` HCL. Check resource exists in API, then subset-match its `Config` JSON against `APICreate`.
   - Step 2: Apply `TerraformUpdate` HCL. Re-check, subset-match against `APIUpdate`.
   - Step 3: `ImportStateVerify` — Terraform imports by ID and asserts the resulting state matches the in-memory state. For sources, `write_key` is excluded (computed, not returned on import).
   - Cleanup: `CheckDestroy` — verifies Delete handler ran (soft-delete tolerant for sources/destinations; strict for connections).

   Connection tests are narrower — they run Create + `ImportStateVerify` + Destroy only. There's no Update step and no config-subset verification; the checks just assert that `source_id`/`destination_id` are wired correctly and the connection exists in the API.
5. **API URL handling** — `newTestAPIClient()` strips a trailing `/v2` from `RUDDERSTACK_API_URL` for backward compatibility (the client adds it back internally). Pass the base URL; the `/v2` is optional.

---

## 4. Environments & triggers

### Triggers

- **PR to `main`** — `.github/workflows/e2e-tests.yml` runs `plan-only` for everything and `e2e-crud` / `e2e-source-crud` / `e2e-conn-crud` for affected integrations. Gated by `e2e-summary`.
- **Push to `main`** — no E2E run on push. `release-please` opens/updates the release PR; release artifacts publish only when the release PR is merged and the GitHub Release is created.
- **No nightly / scheduled run.** If a destination's API config drifts upstream (the control plane adds/removes a field), the next PR that touches that destination will catch it — not before. The workflow only triggers on `pull_request` (no `push` or `workflow_dispatch`); to force a fresh run, open a no-op PR to `main` or run `make testacc-all` locally against the CI `RUDDERSTACK_API_URL` / PAT.

### Affected-integration detection

`detect-changes` (in `e2e-tests.yml`) decides which integrations are affected:

- Touching any of `rudderstack/resource_destination.go`, `rudderstack/resource_source.go`, `rudderstack/client.go`, `rudderstack/provider.go`, `rudderstack/configs/`, `rudderstack/integrations/destinations/common_config_meta.go`, or `internal/testutil/acc/` ⇒ runs **all** destinations, sources, and connections.
- Otherwise — any changed file matching `rudderstack/integrations/destinations/destination_<name>...` (regardless of extension — `.go`, `.tf` examples, `.md` templates, test files) triggers CRUD for that destination. Sources and connections only run CRUD if their respective directories are touched.

### Environments

- The CRUD jobs run inside the GitHub **Environment** `e2e` (`environment: e2e` in the workflow). The environment binds the test workspace's PAT and API URL; if you need a separate staging/prod target, add a second environment and a duplicate matrix job.
- There is currently **one** environment (the dev control plane). No staging or prod target. Treat the dev workspace as the only source of truth for what these tests verify.

---

## 5. Test accounts & credentials

Two-credential model:

### (a) RudderStack control-plane PAT (per environment)

- **What it is:** A Personal Access Token minted in the RudderStack dev control plane.
- **Used by:** The provider (`RUDDERSTACK_ACCESS_TOKEN`) when running E2E CRUD against the API.
- **Source of truth:** Team dev vault, under the eng path (`integrations_team/e2e_test/terraform-provider/rudderstack-account`) for this provider's e2e PAT.
- **Mirrored to:** GitHub Actions secret `RUDDERSTACK_ACC_TEST_TOKEN` on the `e2e` environment of `rudderlabs/terraform-provider-rudderstack`.
- **Mirroring direction:** Vault → GitHub (vault is the source of truth; never edit the GitHub secret directly except for emergency invalidation).

### (b) Underlying test user account

- **What it is:** The control-plane user the PAT is minted for — the shared e2e test user for the dev workspace.
- **Used by:** Humans regenerating the PAT, or when debugging via the control-plane UI.
- **Source of truth:** Team dev vault, under the eng path (`integrations_team/e2e_test/terraform-provider/rudderstack-account`) for this provider's e2e account (login email + password).
- **Mirrored to:** Nothing — this credential is human-only. It is **not** synced to GitHub. CI never authenticates with the user/password; it only uses the PAT minted from this account.

### Per-environment table

| Env | Workspace API URL (`vars.RUDDERSTACK_ACC_TEST_API_URL` in CI; `RUDDERSTACK_API_URL` locally) | PAT vault entry | Account vault entry |
|---|---|---|---|
| dev | `https://api.dev.rudderlabs.com` | eng vault: e2e PAT for this provider | eng vault: e2e account for this provider |

> The test user's email is intentionally not published here — get it from the platform team's runbook before signing in to rotate.

---

## 6. Running locally

The Terraform CLI must be installed on `PATH` before running `make test-ci` or the `testacc-*` targets below. These commands use Terraform Plugin SDK test harnesses (`resource.UnitTest` / `resource.Test`) that invoke the Terraform CLI. Package-scoped tests that avoid those harnesses may not need Terraform.

Minimal commands. Full reference in the `Makefile`.

```bash
# Plan-only — validates HCL + schema for every integration. No token, zero API calls.
make testacc-plan

# Full CRUD — one destination. Requires .env with RUDDERSTACK_ACCESS_TOKEN.
make testacc-dest DEST=webhook

# Full CRUD — one source.
make testacc-source SRC=http

# Full CRUD — connections.
make testacc-conn

# Full CRUD — everything (60-minute timeout, hundreds of API calls).
make testacc-all
```

Create `.env` at the repo root (git-ignored). The Makefile auto-loads it via `-include .env`:

```
RUDDERSTACK_ACCESS_TOKEN=<paste from vault>
RUDDERSTACK_API_URL=https://api.dev.rudderlabs.com
```

The single-integration targets (`testacc-dest`, `testacc-source`) — and the CI destination CRUD job — use a `(?i)`-prefixed Go regex, so `DEST=webhook` matches `TestAccDestinationWebhook` and `DEST=customer_io` matches `TestAccDestinationCustomerIO` (underscores become `.*`). The bulk targets (`testacc-plan`, `testacc-all`, `testacc-conn`) and the CI source/connection jobs use a case-sensitive `TestAcc*` prefix instead, which works because every generated test name starts with `TestAcc`.

---

## 7. Troubleshooting & debugging

### What's captured

This is a Go-test suite against an HTTP API. There are no browser screenshots, DOM dumps, or MHTML snapshots — those don't apply here. The captured artifacts are:

| Artifact | Where | How to access | Useful for |
|---|---|---|---|
| GitHub Actions job log | `Actions → E2E Tests → <run> → <job>` | Right-click "Download log archive" or stream live | All failures — start here |
| Diff dump on config mismatch | Inline in the job log (printed by `compareConfig` in `config_verify.go`) | Search the log for `=== expected config ===` | The API returned a config that doesn't match the test's `APICreate`/`APIUpdate` expectations — shows both sides pretty-printed |
| `terraform-plugin-sdk` debug output | Set `TF_LOG=DEBUG` locally; not enabled in CI | `TF_LOG=DEBUG make testacc-dest DEST=...` | Inspecting Terraform plan/apply internals when a step fails before any API call |
| Leftover API resources | RudderStack dev control plane, filtered by `tf-acc-` prefix | Log in as the test account; filter by name | Cleaning up after a hard-cancelled local run that never reached `Destroy` |

### Concrete failure modes

| Symptom | First check | Then |
|---|---|---|
| `RUDDERSTACK_ACCESS_TOKEN must be set for acceptance tests` | Local: `.env` exists at repo root and has the token. CI: `e2e` environment is selected on the job (`environment: e2e`) | Re-mint the PAT if the workspace was rebuilt |
| `failed to get destination from API: 401` | PAT is expired or revoked | Rotate the PAT (mint a new one as the e2e test user, write to vault, update the `RUDDERSTACK_ACC_TEST_TOKEN` GitHub Environment secret on `e2e`, revoke the old PAT once a test run succeeds) |
| `API config verification failed: missing field "<x>"` | The integration's `APICreate`/`APIUpdate` JSON expects `<x>` but the API didn't return it | Either the upstream control plane stopped persisting that field (legitimate change — update the expected JSON) or the provider isn't sending it (bug — fix the provider's `MarshalJSON`/`flatten` path) |
| `ImportStateVerify` mismatch | A computed-only field is leaking into the state diff | Add the field to `ImportStateVerifyIgnore` in the helper (already done for `write_key` on sources) |
| `destination <id> not found in API` after Create | The Create step thinks it succeeded but the resource isn't queryable | Check the job log for an earlier 4xx/5xx that the SDK swallowed; reproduce locally with `TF_LOG=DEBUG` |
| Test passes locally, fails in CI | Different workspace, different feature flags, different rate limits | Set `RUDDERSTACK_API_URL` locally to the same value CI uses (the `vars.RUDDERSTACK_ACC_TEST_API_URL` GitHub variable); compare auth scopes |
| Flake: random transient 5xx | Control plane was deploying | Re-run the job; if it recurs, file a control-plane issue with the request ID from the log |

### Cleaning up leftover resources

If a local run is killed mid-test, the Terraform state is gone but the API resources remain. Find them via the control-plane UI by filtering names by the `tf-acc-` prefix, then delete. There is no scripted cleanup today — if leftover resources become a recurring problem, add a sweeper following the [`terraform-plugin-sdk` sweeper pattern](https://developer.hashicorp.com/terraform/plugin/sdkv2/testing/acceptance-tests/sweepers).

---

## Where the framework lives

| Path | Purpose |
|---|---|
| `internal/testutil/acc/provider.go` | Provider factory, pre-check, dummy-token helper, `PlanOnly()` |
| `internal/testutil/acc/destinations.go` | `AccAssertDestination` — plan-only or full CRUD |
| `internal/testutil/acc/sources.go` | `AccAssertSource` — plan-only or full CRUD (+ `GeoEnrichmentEnabled`/`Transient` settings check) |
| `internal/testutil/acc/connections.go` | `AccAssertConnection` — wires a source + destination together |
| `internal/testutil/acc/config_verify.go` | `compareConfig` — subset semantics for JSON comparison |
| `internal/testutil/acc/coverage_test.go` | CI gate: every registered integration must have a `TestAcc*` |
| `.github/workflows/e2e-tests.yml` | CI orchestration (detect-changes → plan-only → matrix CRUD → summary) |
| `Makefile` (`testacc-*` targets) | Local entry points |
