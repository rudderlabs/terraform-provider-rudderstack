# Staging E2E Verification — Handoff Runbook

**Goal:** verify the new Terraform **Accounts** feature + **rETL** end-to-end against staging (`api.staging.rudderlabs.com`): create a BigQuery account, an rETL source that uses it, and an rETL connection — via both the Go acceptance harness and a real `terraform apply` smoke.

**Branch:** `feature/retl-customerio-audience-acc` (in `terraform-provider-rudderstack`) — the tip of the 4-PR verification stack (#271 → #272 → docs → CIO), so it contains every layer. Builds clean; all plan-only/unit tests pass. (See [BRANCH-MAP.md](BRANCH-MAP.md) for the full stack.)

**Current blocker (your task):** staging does **not** yet have the `SOURCE_BIGQUERY` account definition registered. The provider, client wiring, auth, and `/v2/accounts` API are all verified working — the only missing piece is the account definition. Everything below is ready to run the moment that lands.

---

## Step 1 — Create the `SOURCE_BIGQUERY` account definition on staging

This is the foundation gap (Stream A / milestone **M2 — Foundation Deployed**): the BigQuery source `accountDefinition` files per [DEX-374](https://linear.app/rudderstack/issue/DEX-374), [APL-399](https://linear.app/rudderstack/issue/APL-399), [APL-398](https://linear.app/rudderstack/issue/APL-398), against `CONTRACT-ACCT-V1`. It must be deployed/registered in the staging control plane so `accountDefinitionName: "SOURCE_BIGQUERY"` resolves.

**Symptom today** (confirms the gap):
```bash
curl -sS -X POST "https://api.staging.rudderlabs.com/v2/accounts" \
  -H "Authorization: Bearer $STAGING_PAT" -H 'Content-Type: application/json' \
  -d '{"name":"defn-probe","accountDefinitionName":"SOURCE_BIGQUERY","options":{"projectId":"x"},"secret":{"credentials":"{}"}}'
# → 404  "Account definition with name \"SOURCE_BIGQUERY\" does not exist"
```
Once you've registered the definition, the same probe returns a created account (with an `id`) instead of the 404. Delete that probe account afterward, or ignore it.

---

## Step 2 — Prerequisites (one-time local setup)

```bash
# 1. Tooling
brew install hashicorp/tap/terraform   # terraform CLI (the Go acc tests + run.sh need it)
go version                              # 1.23+; repo builds with the toolchain in go.mod

# 2. Get the branch
cd terraform-provider-rudderstack
git fetch && git switch feature/retl-customerio-audience-acc   # tip of the verification stack (all layers)
go build ./...                          # sanity: must succeed
```
> Note: `go.mod` pins `rudder-iac` to the **unmerged** PR #617 (`...31f63ee269cf`) for the account client. If that fails to fetch, ensure your git auth can reach `github.com/rudderlabs/rudder-iac`. Re-pin to a tagged release once #617 merges.

---

## Step 3 — Secrets & fixtures

All of these are git-ignored (`.env`, `*.tfvars`, `sa.json`) — they will **not** be committed.

**3a. GCP service-account key** → `test/e2e/staging/sa.json`
Paste the JSON key for a service account with BigQuery access in your project. (Reference setup used project `big-query-integration-poc`, SA `rudder-warehouse-dev@…`.)

**3b. BigQuery fixture table** — the rETL source points at `rudder_tf_e2e.users` (cols `user_id`, `email`, `created_at`). It already exists in `big-query-integration-poc`. If you use a **different** project, create it once:
```bash
python3 -m venv /tmp/bqvenv && /tmp/bqvenv/bin/pip install -q google-cloud-bigquery
GOOGLE_APPLICATION_CREDENTIALS=test/e2e/staging/sa.json /tmp/bqvenv/bin/python - <<'PY'
from google.cloud import bigquery
c = bigquery.Client()
ds = "rudder_tf_e2e"
c.create_dataset(bigquery.Dataset(f"{c.project}.{ds}"), exists_ok=True, timeout=30)
t = bigquery.Table(f"{c.project}.{ds}.users", schema=[
    bigquery.SchemaField("user_id","STRING",mode="REQUIRED"),
    bigquery.SchemaField("email","STRING"),
    bigquery.SchemaField("created_at","TIMESTAMP")])
c.create_table(t, exists_ok=True, timeout=30)
print("fixture ready:", t.full_table_id)
PY
```

**3c. `test/e2e/staging/secret.tfvars`** (for the Terraform smoke):
```hcl
access_token = "<STAGING_PAT>"
bq_project   = "big-query-integration-poc"   # your project
bq_dataset   = "rudder_tf_e2e"
bq_table     = "users"
# bq_location = "US"   # uncomment if not US
# bq_credentials is fed from sa.json at runtime — do NOT put it here.
```

**3d. `.env`** at repo root (for the Go acceptance tests):
```
RUDDERSTACK_ACCESS_TOKEN=<STAGING_PAT>
RUDDERSTACK_API_URL=https://api.staging.rudderlabs.com
```

---

## Step 4 — Mint a BigQuery account and capture its ID (for the rETL Go tests)

The rETL acceptance tests reference an existing BigQuery account by ID. Create one and export its id:
```bash
set -a; . ./.env; set +a
ACC_ID=$(jq -n --arg creds "$(cat test/e2e/staging/sa.json)" \
  '{name:"tf-e2e-retl-acct",accountDefinitionName:"SOURCE_BIGQUERY",
    options:{projectId:"big-query-integration-poc",location:"US"},
    secret:{credentials:$creds}}' \
  | curl -sS -X POST "$RUDDERSTACK_API_URL/v2/accounts" \
      -H "Authorization: Bearer $RUDDERSTACK_ACCESS_TOKEN" -H 'Content-Type: application/json' -d @- \
  | jq -r '.id')
echo "account id: $ACC_ID"
export RUDDERSTACK_RETL_TEST_ACCOUNT_ID="$ACC_ID"
```
(Keep this account around for test runs; delete it when fully done.)

---

## Step 5 — Run the tests

### 5a. Account CRUD (Go acc test)
```bash
set -a; . ./.env; set +a
TF_ACC=1 go test ./rudderstack/integrations/accounts/ -run TestAcc -v -count=1 -timeout 15m
```
**Expect:** `TestAccAccountBigQuery` PASS — real Create→Update→Import→Destroy.

### 5b. rETL source + connection (Go acc test, BigQuery)
```bash
set -a; . ./.env; set +a
export RUDDERSTACK_RETL_TEST_ACCOUNT_ID="$ACC_ID"   # from Step 4
make testacc-retl
```
**Expect:** `TestAccRETLSourceModel_BigQuery`, `TestAccRETLSourceTable_BigQuery`, and the connection tests PASS (Create→Update→Import→Destroy against staging). Without the account ID they **skip** (not fail).

### 5c. Full-chain smoke (real `terraform apply`)
```bash
cd test/e2e/staging
export TF_VAR_bq_credentials="$(cat sa.json)"
./run.sh
```
**Expect:** applies `bigquery account → retl_source_table → retl_connection`, asserts via `terraform plan -detailed-exitcode` (0 = all exist, no drift), then destroys everything via the cleanup trap. Exit 0.

- **Backfill:** `./run.sh --backfill` currently **exits 3 by design** — there is no rETL sync/backfill-trigger endpoint in the rudder-iac client yet. Confirm the control-plane sync-trigger path before wiring this; it's the one piece of F3 still open.

---

## Troubleshooting

| Symptom | Cause / fix |
|---|---|
| `Account definition with name "SOURCE_BIGQUERY" does not exist` | Step 1 not done — definition not registered on staging |
| rETL tests SKIP | `RUDDERSTACK_RETL_TEST_ACCOUNT_ID` not exported (Step 4) |
| `terraform: command not found` | `brew install hashicorp/tap/terraform` |
| go build fails on rudder-iac | git auth to `github.com/rudderlabs/rudder-iac` (unmerged #617 pin) |
| source create fails on missing table | the rETL source points at `rudder_tf_e2e.users` in the account's project — create the fixture (Step 3b) |
| `run.sh --backfill` exits 3 | expected — no sync-trigger endpoint yet |

## Cleanup
Go tests self-destroy their resources; `run.sh` destroys via trap. Delete the Step-4 account and the `rudder_tf_e2e` dataset when finished.
