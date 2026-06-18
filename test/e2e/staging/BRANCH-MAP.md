# Branch / PR Map — TF rETL Account Management verification

Full branch & PR layout for the
[TF — rETL Account Management (Lovable Phase 1)](https://linear.app/rudderstack/project/tf-retl-account-management-lovable-phase-1-6076c97fd07f)
project: the 7-PR **accounts feature stack** (Stream E) plus the 4-PR **e2e verification stack** built on top of it.
Render with any Mermaid-aware viewer (GitHub, Notion, VS Code, mermaid.live).

## Git tree

```mermaid
%%{init: {'gitGraph': {'mainBranchName': 'main'}}}%%
gitGraph
   commit id: "4ab4efc (main)"
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
  main["origin/main @4ab4efc"]:::remote

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
    pr2["#272 · feature/retl-e2e-account-client-and-staging · Accounts client + staging smoke (gated #617)"]:::remote
    pr3["#273 · feature/retl-e2e-docs · HANDOFF + BRANCH-MAP (this PR)"]:::remote
    pr4["#274 · feature/retl-customerio-audience-acc · customerio_audience coverage"]:::remote
    pr1 --> pr2 --> pr3 --> pr4
  end

  main --> d376
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

## E2E verification stack (this work)

| PR | Branch | Base | Contents | Merge gate |
|----|--------|------|----------|-----------|
| **#271** | `feature/pro-5676…terraform` | #266 | rETL acceptance tests (Alexandros, authorship preserved) → BigQuery; CIO documented-excluded | after accounts stack |
| **#272** | `feature/retl-e2e-account-client-and-staging` | #271 | real Accounts client wiring + 404 fix + `test/e2e/staging` smoke | **rudder-iac #617** |
| **#273** | `feature/retl-e2e-docs` | #272 | `HANDOFF.md` + `BRANCH-MAP.md` | with #272 |
| **#274** | `feature/retl-customerio-audience-acc` | #273 | `retl_connection_customerio_audience` coverage (gate exclusion→enforced) | with stack |

## Notes

- **#237 could not be reopened** — its branch was force-pushed after close, which GitHub blocks. #271 carries the identical commits (Alexandros's authorship intact) as a new PR.
- **`verify/accounts-retl-e2e`** is kept locally as the complete integrated reference (all layers in one branch); the [HANDOFF](HANDOFF.md) runbook targets the pushed stack tip `feature/retl-customerio-audience-acc`.
- **Merge order:** the accounts stack #260–266 merges bottom-up into `main` first; then this 4-PR stack rebases onto `main` and #272's `go.mod` pin flips to a released rudder-iac version once #617 lands.
