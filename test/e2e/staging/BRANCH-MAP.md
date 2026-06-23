# Branch / PR Map — TF rETL Account Management verification

Full branch & PR layout for the
[TF — rETL Account Management (Lovable Phase 1)](https://linear.app/rudderstack/project/tf-retl-account-management-lovable-phase-1-6076c97fd07f)
project: the 7-PR **accounts feature stack** (Stream E) plus the **e2e verification stack** built on top of it — the whole chain now rooted on the CI-enablement PR **#277**.
Render with any Mermaid-aware viewer (GitHub, Notion, VS Code, mermaid.live).

## Git tree

```mermaid
%%{init: {'gitGraph': {'mainBranchName': 'main'}}}%%
gitGraph
   commit id: "ef89fed (main)"
   branch "#277 ci/fix-pr-triggers"
   commit id: "ci-trigger-fix"
   branch "#260 DEX-376 accounts registry"
   commit id: "dex-376"
   branch "#261 DEX-377 resource_account CRUD"
   commit id: "dex-377"
   branch "#262 DEX-378 data source"
   commit id: "dex-378"
   branch "#263 DEX-381 AssertAccount helper"
   commit id: "dex-381"
   branch "#264 DEX-379 bq ConfigMeta"
   commit id: "dex-379"
   branch "#265 DEX-380 provider wiring"
   commit id: "dex-380"
   branch "#266 DEX-382 bq integration test"
   commit id: "dex-382"
   branch "#271 rETL acc to BigQuery"
   commit id: "feature/pro-5676" type: HIGHLIGHT
   branch "#278 rETL config-verify (PRO-5768)"
   commit id: "feature/pro-5768"
   branch "#272 Accounts client + staging smoke"
   commit id: "retl-e2e-account-client-and-staging"
   branch "#273 docs (HANDOFF + BRANCH-MAP)"
   commit id: "retl-e2e-docs"
   branch "#274 customerio_audience coverage"
   commit id: "retl-customerio-audience-acc"
```

*Branch lanes are labelled by PR (`#NNN` + DEX issue); commit dots carry the branch name. HIGHLIGHT (#271) = the lane carrying Alexandros Milaios's preserved commits. Lower 7 lanes = accounts feature stack (Stream E); upper 4 = this verification stack.*

## PR stack (bottom-up merge order)

```mermaid
flowchart TD
  main["origin/main @ef89fed"]:::remote

  subgraph ACC["Accounts feature stack · Stream E · #260-266"]
    direction TB
    d376["#260 · DEX-376 · dex-376 · accounts integration registry"]:::remote
    d377["#261 · DEX-377 · dex-377 · generic resource_account.go CRUD"]:::remote
    d378["#262 · DEX-378 · dex-378 · data_source_account (read-only)"]:::remote
    d381["#263 · DEX-381 · dex-381 · AssertAccount test helper"]:::remote
    d379["#264 · DEX-379 · dex-379 · BigQuery account ConfigMeta"]:::remote
    d380["#265 · DEX-380 · dex-380 · provider.go resources wiring"]:::remote
    d382["#266 · DEX-382 · dex-382 · BigQuery account integration test"]:::remote
    d376 --> d377 --> d378 --> d381 --> d379 --> d380 --> d382
  end

  subgraph VER["E2E verification stack · this work"]
    direction TB
    pr1["#271 · feature/pro-5676 · rETL acc to BigQuery (rebased from #237)"]:::remote
    pcv["#278 · feature/pro-5768 · rETL upstream config-verify (PRO-5768)"]:::remote
    pr2["#272 · feature/retl-e2e-account-client-and-staging · Accounts client + staging smoke (v0.18.0, #617 resolved)"]:::remote
    pr3["#273 · feature/retl-e2e-docs · HANDOFF + BRANCH-MAP (this PR)"]:::remote
    pr4["#274 · feature/retl-customerio-audience-acc · customerio_audience coverage"]:::remote
    pr1 --> pcv --> pr2 --> pr3 --> pr4
  end

  ci277["#277 · ci/fix-pr-triggers · run e2e+unit on all PRs (PRO-5776)"]:::remote
  main --> ci277
  ci277 --> d376
  d382 --> pr1
  p237["PR#237 · CLOSED (superseded)"]:::closed
  p237 -. "force-push blocked reopen, new PR" .-> pr1
  verify["verify/accounts-retl-e2e · integrated reference (local only, all layers)"]:::localonly
  d382 -. "complete integrated copy" .-> verify

  classDef remote fill:#0d4429,stroke:#3fb950,color:#e6edf3;
  classDef localonly fill:#5a3b00,stroke:#d29922,color:#e6edf3;
  classDef closed fill:#3d3d3d,stroke:#8b949e,color:#e6edf3,stroke-dasharray:4 3;
```

🟩 remote/pushed · 🟧 local-only · ⬛ closed

## Accounts feature stack (Stream E)

| PR | Issue | Branch (`feature/…`) | Base | Summary |
|----|-------|----------------------|------|---------|
| #260 | DEX-376 | `dex-376-…-add-accounts-registry` | `main` | Accounts integration registry |
| #261 | DEX-377 | `dex-377-…-generic-resource_accountgo-crud-handler` | #260 | generic `resource_account.go` CRUD + import |
| #262 | DEX-378 | `dex-378-…-data_source_accountgo-read-only` | #261 | read-only `rudderstack_account` data source |
| #263 | DEX-381 | `dex-381-…-build-assertaccount-test-helper` | #262 | `AssertAccount` / `AccAssertAccount` helpers |
| #264 | DEX-379 | `dex-379-…-bigquery-account-integration-file` | #263 | BigQuery account `ConfigMeta` + example |
| #265 | DEX-380 | `dex-380-…-wire-providergo-resources-datasourcesmap` | #264 | register resource + data source in `provider.go` |
| #266 | DEX-382 | `dex-382-…-bigquery-account-integration-test` | #265 | BigQuery account integration test + docs |

## Linear tree (current — bottom-up merge order)

The whole chain is now stacked on **#277** so every PR's head carries the CI
trigger-fix and runs checks. rudder-iac is standardized to **v0.18.0**
throughout (the #617 gate is dissolved — it's now a tagged release).

| PR | Branch | Base | Contents |
|----|--------|------|----------|
| **#277** | `ci/fix-pr-triggers` | `main` | run e2e + unit on all PRs (drop `branches:[main]`) — PRO-5776 |
| #260–#266 | `feature/dex-376…382` | #277 → … | accounts stack (Stream E); #261 bumps rudder-iac→v0.18.0; #262/#263 carry the `client.*` consumer migration; #266 adds the `e2e-account-crud` job |
| **#271** | `feature/pro-5676…` | #266 | rETL acceptance tests (Alexandros authorship preserved) → BigQuery |
| **#278** | `feature/pro-5768-retl-upstream-config-verify` | #271 | rETL upstream config verification — PRO-5768 |
| **#272** | `feature/retl-e2e-account-client-and-staging` | #278 | real Accounts client wiring + 404 fix + `test/e2e/staging` smoke (v0.18.0; duplicate stub-migration dropped) |
| **#273** | `feature/retl-e2e-docs` | #272 | `HANDOFF.md` + `BRANCH-MAP.md` (this PR) |
| **#274** | `feature/retl-customerio-audience-acc` | #273 | `retl_connection_customerio_audience` coverage |

## Notes

- **#237 could not be reopened** — its branch was force-pushed after close, which GitHub blocks. #271 carries the identical commits (Alexandros's authorship intact) as a new PR.
- **`verify/accounts-retl-e2e`** is kept locally as the complete integrated reference (all layers in one branch); the [HANDOFF](HANDOFF.md) runbook targets the pushed stack tip `feature/retl-customerio-audience-acc`.
- **Merge order:** merge bottom-up starting at **#277** → main; each PR below auto-retargets to `main` as the one beneath it merges, and its checks run green before merge.
- **rudder-iac v0.18.0 / #617:** the released v0.18.0 supersedes the old `#617` pseudo-version pin, so #272 is no longer gated. The stub-removal + `client.*` migration (originally duplicated across #261 and #272) now lands once, in the accounts stack.
