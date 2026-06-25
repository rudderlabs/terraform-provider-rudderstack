# Staging smoke test — BigQuery account → rETL source → connection

This config applies a minimal chain against the RudderStack staging environment
to verify that the Terraform provider can create and link:

1. A BigQuery rETL **account** (`rudderstack_account_source_bigquery`)
2. An rETL **source table** (`rudderstack_retl_source_table`) backed by that account
3. A webhook **destination** (`rudderstack_destination_webhook`) — throwaway endpoint
4. An rETL **connection** (`rudderstack_retl_connection`) wiring source → destination

No real syncs are triggered (schedule type is `manual`).

---

## Prerequisites

- Terraform ≥ 1.0
- Go ≥ 1.21 (to build the provider locally)
- A RudderStack staging personal access token
- A GCP service-account JSON key with BigQuery access

---

## 1. Build and install the provider locally

```sh
cd /path/to/terraform-provider-rudderstack
go build -o /tmp/terraform-provider-rudderstack .
```

Create a dev-override Terraform CLI config at `/tmp/dev.tfrc`:

```hcl
provider_installation {
  dev_overrides {
    "rudderstack.com/rudderlabs/rudderstack" = "/tmp"
  }
  direct {}
}
```

---

## 2. Supply credentials via a git-ignored tfvars file

Create `test/e2e/staging/secret.tfvars` (this path is in `.gitignore` — never
commit it):

```hcl
# secret.tfvars — DO NOT COMMIT
access_token   = "rsa_REPLACE_ME"
bq_project     = "my-gcp-project"
bq_dataset     = "my_dataset"
bq_table       = "users"
bq_credentials = <<EOT
{
  "type": "service_account",
  "project_id": "my-gcp-project",
  ...
}
EOT
```

Optional overrides (have defaults):

```hcl
api_url     = "https://api.staging.rudderlabs.com"
bq_location = "US"
```

---

## 3. Run

`run.sh` builds the provider locally into `.bin/` (git-ignored), wires a
Terraform dev-override pointing at that directory, runs `apply`, asserts no
drift, prints the created resource IDs, and always destroys on exit.

```sh
# Standard smoke run — no staging creds needed in the env:
./run.sh

# Override the tfvars file path:
./run.sh path/to/other.tfvars

# Or via env:
TFVARS_FILE=path/to/other.tfvars ./run.sh
```

The script does **not** require `TF_CLI_CONFIG_FILE` to be set beforehand; it
writes a temporary dev-override config and exports the variable itself.

### `--backfill` flag (opt-in, not yet wired)

```sh
./run.sh --backfill
```

After the drift assertion the script will attempt to trigger a manual rETL
sync on the connection and poll to completion. **This branch currently exits
with a clear error (exit 3)** because the rudder-iac client
(`api/client/retl/connections.go`, v0.17.1) does not yet expose a
sync-trigger endpoint. See the comment block in `run.sh` for exactly which
endpoint and polling contract need to be confirmed before this can be wired.

### Generated files

| Path | Description |
|------|-------------|
| `.bin/terraform-provider-rudderstack` | Provider binary built by `run.sh` — git-ignored |

To apply manually without `run.sh`:

```sh
# 1. Build first:
go build -o test/e2e/staging/.bin/terraform-provider-rudderstack .

# 2. Write a dev-override config:
cat > /tmp/dev.tfrc <<'HCL'
provider_installation {
  dev_overrides {
    "rudderstack.com/rudderlabs/rudderstack" = "test/e2e/staging/.bin"
  }
  direct {}
}
HCL

# 3. Apply:
TF_CLI_CONFIG_FILE=/tmp/dev.tfrc terraform -chdir=test/e2e/staging \
  apply -var-file=test/e2e/staging/secret.tfvars -auto-approve
```

---

## 4. What success looks like

After `apply`, Terraform prints four non-empty IDs:

```
account_id     = "<id>"
connection_id  = "<id>"
destination_id = "<id>"
retl_source_id = "<id>"
```

`run.sh` then asserts these are correct by running
`terraform plan -detailed-exitcode`.  Exit code 0 means the provider Read
path returned state that matches the config exactly — no drift.  Exit code 2
(drift) or 1 (error) fail the script loudly before destroy runs.
