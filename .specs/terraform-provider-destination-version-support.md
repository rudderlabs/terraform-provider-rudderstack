# Technical Refinement Spec — [INT-6494](https://linear.app/rudderstack/issue/INT-6494/terraform-provider-rudderstack-codegen-resolve-by-apitype-version)

Refinement of the original story. All decisions below were validated against the current codebase (`91b4919`) and rudder-iac `feat/rud-2864-destination-resource-handler-updated` (PR #625). Open items listed at the bottom.

## 1\. Type & wire format

* `ConfigMeta.Version` is `int` (aligns with `client.Destination.Version int` in rudder-iac).
* Wire sends `version: 1` (numeric, not string). Matches ` Destination.Version` `json:"version,omitempty"`.
* Backfill value across the existing fleet is `1` (int), not `"1"`.

## 2\. State schema — no `version` field

* The destination Terraform resource schema is **unchanged**. No `version` attribute is added to state.
* Version is encoded purely by the resource type name (`rudderstack_destination_braze` vs `rudderstack_destination_braze_v2`).
* `populateDestinationFromState` writes `destination.Version = cm.Version` on create/update (wire-only).
* `storeDestinationToState` does **not** read `destination.Version` back — there is no state field to write into.
* Rationale: adding any version field (Required / Optional / Computed) causes a one-time plan diff across the entire existing fleet on upgrade. The story's "behaviorally a no-op, no plan diff" guarantee only holds if the schema is untouched.
* Known gap (deferred to rollout milestone): if the backend ever migrates a v1 instance to v2 in place, the provider cannot detect it (no state field to diff). Out of scope here.

## 3\. Backfill strategy — explicit, all 47 destinations

* Every existing bare (suffix-less) `Destinations.Register` call site in `rudderstack/integrations/destinations/destination_*.go` is edited to set `Version: 1` on its `ConfigMeta` literal.
* 47 files touched (one registration per file).
* Sources and Accounts are **not** versioned in this story. Their `ConfigMeta.Version` stays at the zero value; their registries keep the current name-only uniqueness check.

## 4\. Registry guards

* `Destinations.Register` gains two new panic guards (sources/accounts unaffected):
  1. **Presence**: panic if `cm.Version == 0` — forces every destination registration to declare an explicit version.
  2. **Uniqueness on** `(apiType, version)`: panic if another entry with the same `(APIType, Version)` is already registered. Requires a secondary index in `Registry` keyed by `(apiType, version)`, scoped to the destinations registry.
* These guards replace the "ambiguous match" risk inside the codegen resolver: failures happen at provider init (loud, early) rather than at codegen lookup time (silent, late).

## 5\. `generatetf` codegen resolver — two functions

* `configMeta(entries, apiType)` — **unchanged**, first-match on `APIType`. Used by source and account callers (lines 59, 143 — sources; accounts have no codegen path).
* `configMetaByVersion(entries, apiType, version int)` — **new**, mandatory version filter. Used by destination callers (lines 67, 162). Signature: `func configMetaByVersion(entries map[string]configs.ConfigMeta, apiType string, version int) (string, *configs.ConfigMeta)`.
* `configMetaByVersion` matches strictly on `e.APIType == apiType && e.Version == version`. **No coercion of** `version == 0` — see §6.
* Destination call sites become `configMetaByVersion(destinationConfigs, dst.Type, dst.Version)`.
* Source call sites stay `configMeta(sourceConfigs, src.Type)`.

## 6\. `absent version on the wire` — owned by [INT-6489](https://linear.app/rudderstack/issue/INT-6489/public-api-contract-version-on-rudder-api-dtos-migrate-proxy-and), not this story

* Q4 decision: this provider **relies on the backend (**[INT-6489](https://linear.app/rudderstack/issue/INT-6489/public-api-contract-version-on-rudder-api-dtos-migrate-proxy-and)**) to always send a real** `version` on every destination response.
* `configMetaByVersion` does **not** coerce `0 → 1`. If `dst.Version == 0`, no entry matches → the destination is skipped in codegen output.
* Consequence: the original story AC bullet *"absent* `version` *round-trips as v1"* and the corresponding test *"a* `null`*/absent version resolves to the v1 entry"* are **moved to** [INT-6489](https://linear.app/rudderstack/issue/INT-6489/public-api-contract-version-on-rudder-api-dtos-migrate-proxy-and) (backend concern: backfill `version: 1` on all existing records, default on create).
* This story's provider-side AC for absent version is replaced with: *"provider assumes backend always sends* `version`*; verified via mocked-API unit tests, gated on* [INT-6489](https://linear.app/rudderstack/issue/INT-6489/public-api-contract-version-on-rudder-api-dtos-migrate-proxy-and) *backend rollout for the live acceptance test."*

## 7\. Shared-base composition helper — scaffolding

* **Location**: `rudderstack/integrations/destinations/compose.go` (next to `common_config_meta.go`). Destinations-only for now.
* **Delta struct** (typed, compile-checked):

  ```go
  type Delta struct {
      Renamed           map[string]string          // old schema key -> new schema key
      Added             map[string]*schema.Schema   // new top-level schema fields
      Removed           []string                    // dropped schema keys
      AddedProperties   []c.ConfigProperty          // appended to Properties
      RemovedProperties []string                    // removed from Properties (by name)
  }
  ```
* **Helper signature**:

  ```go
  func ComposeConfigMeta(base c.ConfigMeta, delta Delta, version int) c.ConfigMeta
  ```
  * Clones `base.ConfigSchema` (shallow map copy) and `base.Properties` (slice copy).
  * Applies `Renamed` (delete old key, insert under new key — safe on shallow-copied map), `Added` (insert new `*schema.Schema`), `Removed` (delete keys) to the schema map.
  * Applies `AddedProperties` / `RemovedProperties` to the Properties slice.
  * Sets `result.Version = version`.
  * Returns `ConfigMeta` by value (matches `Registry.Register` signature).
* **Shallow-copy constraint (documented)**: `Renamed` / `Removed` operate by map-key delete+insert only. Callers must **not** mutate nested `*schema.Schema` fields in place — those pointers are shared with the base. Documented at the helper and enforced by the unit test (assert base unchanged after composing a v2).
* No real v2 destination is registered in this story. The helper exists + is unit-tested + sits ready for the rollout milestone.

## 8\. Dependency sequencing — rudder-iac SHA pin + split blockers

* `go.mod` currently pins `github.com/rudderlabs/rudder-iac v0.18.0`, which does **not** have `Destination.Version`. Confirmed by `git show v0.18.0:api/client/destinations.go`.
* The `Version` field lands in rudder-iac PR #625 (branch `feat/rud-2864-destination-resource-handler-updated`).
* **This story's PR pins rudder-iac to PR #625's merge SHA** via a pseudo-version: `go get github.com/rudderlabs/rudder-iac@<merge-sha> && go mod tidy` → `require github.com/rudderlabs/rudder-iac v0.18.1-0.<timestamp>-<sha12>`. Only valid after PR #625 **merges** (not while open — SHA must be reachable from a merged branch).
* [INT-6489](https://linear.app/rudderstack/issue/INT-6489/public-api-contract-version-on-rudder-api-dtos-migrate-proxy-and) is split into two blockers for this story:
  1. **Blocker A (rudder-iac client struct)**: PR #625 merge. Unblocks provider code + unit tests. This story can merge on this.
  2. **Blocker B (backend wire guarantee)**: backend always sends `version` on every destination response. Unblocks the live env-gated acceptance test (braze no-plan-diff). The acceptance test is written in this PR but skipped in CI until Blocker B is rolled out.
* Follow-up trivial PR (separate): once rudder-iac cuts a real tag, swap the pseudo-version → tag. No code change.

## 9\. Tests

| Test | File | Type | Gated on |
| -- | -- | -- | -- |
| `configMetaByVersion` resolution: v1, v2, mismatch → skip | `cmd/generatetf/generator/generator_version_test.go` (new) | Unit, table-driven | Blocker A |
| Existing `generator_test.go` destination fixtures bumped to `Version: 1` | `cmd/generatetf/generator/generator_test.go` | Unit (edit existing) | Blocker A |
| `populateDestinationFromState` sets `destination.Version = cm.Version` | `rudderstack/resource_destination_test.go` (new case) | Unit | Blocker A |
| `Destinations.Register` panics on `Version == 0` and on duplicate `(apiType, version)` | `rudderstack/configs/registries_test.go` (new cases) | Unit | Blocker A |
| `ComposeConfigMeta`: base + delta → expected ConfigMeta; base unchanged after | `rudderstack/integrations/destinations/compose_test.go` (new) | Unit | Blocker A |
| Braze applies with `version: 1` in request body, no plan diff before/after provider upgrade | `rudderstack/resource_destination_test.go` | Acceptance, `RUDDERSTACK_ACCESS_TOKEN`-gated | Blocker B |

## 10\. Updated acceptance criteria

- [ ] `Version int` added to `ConfigMeta`; `populateDestinationFromState` sets `destination.Version = cm.Version` (depends on rudder-iac PR #625 / [INT-6489](https://linear.app/rudderstack/issue/INT-6489/public-api-contract-version-on-rudder-api-dtos-migrate-proxy-and) Blocker A).
- [ ] All 47 existing bare destination registrations carry `ConfigMeta.Version = 1`.
- [ ] `Destinations.Register` panics on `Version == 0` and on duplicate `(apiType, version)`.
- [ ] `generatetf` has a new `configMetaByVersion(entries, apiType, version int)`; destination callers thread `dst.Version`; source/account callers unchanged.
- [ ] `generatetf` `configMeta(entries, apiType)` (first-match) retained for sources.
- [ ] Shared-base composition helper `ComposeConfigMeta(base, delta, version)` available in `integrations/destinations/compose.go` with typed `Delta` struct; unit-tested.
- [ ] No `_vX` resource type / registry double-registration (deferred to rollout milestone).
- [ ] Destination resource schema unchanged — no `version` field in state.
- [ ] "absent `version` ⇒ v1" handling moved to [INT-6489](https://linear.app/rudderstack/issue/INT-6489/public-api-contract-version-on-rudder-api-dtos-migrate-proxy-and) (backend); this provider relies on backend always sending `version`.
- [ ] `go.mod` pinned to rudder-iac PR #625 merge SHA via pseudo-version.
- [ ] All unit tests above pass on Blocker A; env-gated acceptance test written and skipped until Blocker B.
- [ ] Existing single-version HCL (`rudderstack_destination_braze`) plans/applies with no diff.

## 11\. Out of scope (explicit)

* Sources and Accounts versioning.
* The actual `_v2` destination registration (e.g. `rudderstack_destination_braze_v2`) — deferred to rollout milestone.
* Registry double-registration of `<name>_v2` resource types.
* Backend-driven v1→v2 in-place migration detection.
* rudder-iac tag release (follow-up trivial PR swaps SHA → tag).
* Backfill verification mechanism (see Open items).

## 12\. Open items

* **Backfill verification** (deferred): how to prove none of the 47 `Version: 1` edits were missed — options were (a) a unit test iterating `configs.Destinations.Entries()` asserting `Version == 1` for all bare names, (b) PR review + grep, (c) a static check. Decision deferred. Recommend (a) as a cheap belt-and-suspenders test regardless of review.
* `common_config_meta.go` **interaction**: the existing consent-management helper returns `([]ConfigProperty, map[string]*schema.Schema})` fragments composed by callers into a `ConfigMeta`. It does not need to know about `Version` — callers set `Version` on the assembled `ConfigMeta`. To be confirmed during implementation that no `common_config_meta` caller breaks.
* **PR #625 merge target**: confirm it merges into a branch the provider can consume (main or a release branch) before pinning the SHA.