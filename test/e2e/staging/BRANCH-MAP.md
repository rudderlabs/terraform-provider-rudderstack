# Branch / PR Map — TF rETL Account Management verification

Branch & PR layout for the
[TF — rETL Account Management (Lovable Phase 1)](https://linear.app/rudderstack/project/tf-retl-account-management-lovable-phase-1-6076c97fd07f)
verification work. Render with any Mermaid-aware viewer (GitHub, Notion, VS Code, mermaid.live).

## Git tree

```mermaid
%%{init: {'gitGraph': {'mainBranchName': 'main'}}}%%
gitGraph
   commit id: "4ab4efc (4.6.0)"
   branch accounts-stack
   commit id: "#260–266 dex-376..382 (accounts)" tag: "stack tip e1b8c9c"
   branch pr1-237
   commit id: "Alexandros: rETL acc + CI" type: HIGHLIGHT
   commit id: "Alexandros: account-id update" type: HIGHLIGHT
   commit id: "enabled + cio-excluded doc"
   commit id: "→ BigQuery"
   branch pr2-client-staging
   commit id: "wire Accounts client [#617]"
   commit id: "404 fix"
   commit id: "staging smoke + run.sh"
   branch pr3-docs
   commit id: "HANDOFF + BRANCH-MAP"
   branch pr4-cio
   commit id: "customerio_audience acc coverage"
```

*HIGHLIGHT = Alexandros Milaios's preserved commits.*

## PR stack (bottom-up merge order)

```mermaid
flowchart TD
  main["origin/main"]:::remote
  stack["accounts stack · PR#260–266 · dex-376…382"]:::remote
  p237["PR#237 · CLOSED (superseded)"]:::closed
  pr1["PR#271 · feature/pro-5676… · rETL acc → BigQuery (rebased from #237)"]:::remote
  pr2["PR#272 · feature/retl-e2e-account-client-and-staging · Accounts client + staging smoke ⚠️#617"]:::remote
  pr3["PR3 · feature/retl-e2e-docs · HANDOFF + BRANCH-MAP (this PR)"]:::remote
  pr4["PR4 · feature/retl-customerio-audience-acc · CIO audience coverage"]:::remote
  verify["verify/accounts-retl-e2e · integrated reference (local only, all layers)"]:::localonly

  main --> stack --> pr1 --> pr2 --> pr3 --> pr4
  p237 -. "force-push blocked reopen → new PR" .-> pr1
  stack -. "complete integrated copy" .-> verify

  classDef remote fill:#0d4429,stroke:#3fb950,color:#e6edf3;
  classDef localonly fill:#5a3b00,stroke:#d29922,color:#e6edf3;
  classDef closed fill:#3d3d3d,stroke:#8b949e,color:#e6edf3,stroke-dasharray:4 3;
```

🟩 remote/pushed · 🟧 local-only · ⬛ closed

| PR | Branch | Base | Contents | Merge gate |
|----|--------|------|----------|-----------|
| **#271** | `feature/pro-5676…terraform` | `feature/dex-382` (#266) | rETL acceptance tests (Alexandros, authorship preserved) → BigQuery; CIO documented-excluded | after accounts stack |
| **#272** | `feature/retl-e2e-account-client-and-staging` | #271 | real Accounts client wiring + 404 fix + `test/e2e/staging` smoke | **rudder-iac #617** |
| **PR3** | `feature/retl-e2e-docs` | #272 | `HANDOFF.md` + `BRANCH-MAP.md` | with #272 |
| **PR4** | `feature/retl-customerio-audience-acc` | PR3 | `retl_connection_customerio_audience` acc coverage (flips gate exclusion→enforced) | with stack |

## Notes

- **#237 could not be reopened** — its branch was force-pushed after close, which GitHub blocks. #271 carries the identical commits (Alexandros's authorship intact) as a new PR.
- **`verify/accounts-retl-e2e`** is kept locally as the complete integrated reference (all four layers in one branch) and as the runnable target for the [HANDOFF](HANDOFF.md) test runbook.
- The accounts stack #260–266 should merge bottom-up into `main` first; then this 4-PR stack rebases onto `main` and #272's `go.mod` pin flips to a released rudder-iac version once #617 lands.
