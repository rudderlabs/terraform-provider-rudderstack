# Branch / PR Map — TF rETL Account Management verification

Layout for the
[TF — rETL Account Management (Lovable Phase 1)](https://linear.app/rudderstack/project/tf-retl-account-management-lovable-phase-1-6076c97fd07f)
project. The original **13 stacked PRs** have been **collapsed into 3 review units
(A / B / C)** stacked on the **#277 CI scaffold**. Render with any Mermaid-aware
viewer (GitHub, Notion, VS Code, mermaid.live).

## Current state — 3 review PRs on the #277 scaffold

```mermaid
%%{init: {'gitGraph': {'mainBranchName': 'main'}}}%%
gitGraph
   commit id: "main"
   branch "#277 ci/fix-pr-triggers (scaffold)"
   commit id: "ci-trigger-fix"
   branch "A #266 accounts"
   commit id: "DEX-376..382 + v0.18.0"
   branch "B #278 rETL tests + config-verify"
   commit id: "PRO-5676 / PRO-5768"
   branch "C #274 client+staging+docs+cio"
   commit id: "verification"
```

Linear chain: `main → #277 → A (#266) → B (#278) → C (#274)`. Reviewers see **3
PRs**; #277 is a never-merged base that lets every PR run checks pre-merge.

## How the 13 PRs collapsed

```mermaid
flowchart LR
  classDef keep fill:#0d4429,stroke:#3fb950,color:#e6edf3;
  classDef closed fill:#3d3d3d,stroke:#8b949e,color:#e6edf3,stroke-dasharray:4 3;
  classDef scaffold fill:#5a3b00,stroke:#d29922,color:#e6edf3;

  main["origin/main"] --> s["#277 ci/fix-pr-triggers — scaffold, never merged"]:::scaffold
  s --> A["A · #266 — feat(accounts) DEX-376–382"]:::keep
  A --> B["B · #278 — rETL tests + config-verify"]:::keep
  B --> C["C · #274 — client + staging + docs + customerio"]:::keep

  subgraph GA ["folded into A (#266)"]
    direction TB
    n260["#260 DEX-376"]; n261["#261 DEX-377"]; n262["#262 DEX-378"]
    n263["#263 DEX-381"]; n264["#264 DEX-379"]; n265["#265 DEX-380"]
  end
  subgraph GB ["folded into B (#278)"]
    n271["#271 rETL acc · Alexandros"]
  end
  subgraph GC ["folded into C (#274)"]
    n272["#272 client+staging"]; n273["#273 docs"]
  end
  class n260,n261,n262,n263,n264,n265,n271,n272,n273 closed
  GA -. closed, commits preserved .-> A
  GB -. closed, commits preserved .-> B
  GC -. closed, commits preserved .-> C
```

🟩 kept (review unit) · 🟧 scaffold (never merged) · ⬛ closed (folded in)

## Consolidation map

| Review PR | Branch | Base | Folds in (closed) | Contents |
|-----------|--------|------|-------------------|----------|
| **#277** scaffold | `ci/fix-pr-triggers` | `main` | — | Run e2e + unit on all PRs (drop `branches:[main]`) — PRO-5776. **Never merged**; closed at the end. |
| **A · #266** | `feature/dex-382-…` | #277 | #260, #261, #262, #263, #264, #265 | Full accounts feature (DEX-376–382): registry, generic `resource_account` CRUD, data source, `AssertAccount` helper, BigQuery `ConfigMeta`, provider + `NewAPIClient` wiring, BigQuery integration test. rudder-iac **v0.18.0**, `client.*` consumer migration, `e2e-account-crud` job. |
| **B · #278** | `feature/pro-5768-…` | A (#266) | #271 | rETL source/connection acceptance tests on BigQuery (Alexandros's commits preserved) + PRO-5768 upstream config verification. |
| **C · #274** | `feature/retl-customerio-audience-acc` | B (#278) | #272, #273 | Real Accounts client + 404 fix, `test/e2e/staging` smoke (`run.sh` + PAUSE hold-open + label-gated `e2e-staging-smoke.yml`), HANDOFF/BRANCH-MAP docs, `customerio_audience` coverage. |

Closed PRs carry a comment pointing to their keeper; their commits live on in the
kept branch (nothing abandoned). They can be reopened if a split is ever needed.

## Merge plan (never merging #277)

Review A → B → C. Then collapse upward — merge **C into B**, **B into A** (one
final check run on the combined A) — then **rebase A off #277 onto `main`** (drops
the single scaffold commit; restores `branches:[main]`), and merge A into `main`.
Close #277. The CI-trigger change never lands in `main`.

## Notes

- **#237** could not be reopened (force-pushed after close); its commits live in **B** (#278) via #271, authorship intact.
- **Secrets, not Vault.** A Vault integration was prototyped and removed — the e2e tests source creds from GitHub Environment secrets/vars on the `e2e` environment (`RUDDERSTACK_ACC_TEST_TOKEN`, `RUDDERSTACK_STAGING_*`, and `RUDDERSTACK_RETL_TEST_ACCOUNT_ID` once the dev seed account is aligned).
- **Staging smoke** runs pre-merge via the `e2e-staging` label on #274 (and `workflow_dispatch` once on `main`); it applies → asserts no drift → destroys against staging.
- **rudder-iac v0.18.0** is standardized across the chain (the old `#617` pseudo-version gate is gone).
- `verify/accounts-retl-e2e` (the old local integrated reference) has been **deleted** — the C tip / `integration/retl-accounts-full` is the current full integration.
