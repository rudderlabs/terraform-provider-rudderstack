# Staging smoke tests — per-scenario rETL chains

This directory applies minimal chains against the RudderStack **staging**
environment to verify the Terraform provider can create, link, read (no drift),
and destroy real resources. Each **scenario** is an independent Terraform root
module; `run.sh` runs them one at a time (apply → assert no drift → destroy).

No real syncs are triggered (every connection uses a `manual` schedule).

## Layout

```
test/e2e/
├── run.sh                       # builds the provider, runs scenarios
├── secret.tfvars                # git-ignored, shared by all scenarios
├── modules/
│   └── bigquery_source/         # shared: BigQuery rETL account + source table
└── scenarios/
    ├── _shared/_shared.tf       # provider + ALL variable declarations (one real file)
    ├── customerio/              # BigQuery → Customer.io (VDM v2)
    │   ├── _shared.tf           # symlink → ../_shared/_shared.tf
    │   └── main.tf
    └── customerio_audience/     # BigQuery → Customer.io Audience
        ├── _shared.tf           # symlink → ../_shared/_shared.tf
        └── main.tf
```

Each scenario stands up its **own** BigQuery account + rETL source (via the
`bigquery_source` module, namespaced by `name_prefix`), then a destination + a
typed rETL connection. Every scenario exposes the same four outputs —
`account_id`, `source_id`, `destination_id`, `connection_id` — which `run.sh`
verifies are non-empty.

The `_shared.tf` symlink lets a single `secret.tfvars` feed every scenario with
no "undeclared variable" warnings (all variables are declared in one place).

## Prerequisites

- Terraform ≥ 1.0
- Go ≥ 1.21 (to build the provider locally)
- A RudderStack staging personal access token
- A GCP service-account JSON key with BigQuery access
- Vendor creds for whichever scenarios you run (see below)

## Credentials — `test/e2e/secret.tfvars`

Create `test/e2e/secret.tfvars` (git-ignored — never commit it):

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

# customerio scenario (BigQuery → Customer.io VDM v2):
customerio_site_id    = "REPLACE_ME"
customerio_api_key    = "REPLACE_ME"
customerio_datacenter = "US"   # or "EU"; optional, defaults to US

# customerio_audience scenario (BigQuery → Customer.io Audience):
customerio_audience_app_api_key = "REPLACE_ME"
customerio_audience_id          = 16     # a real Customer.io audience ID
customerio_audience_region      = "US"   # optional, defaults to US
```

Optional overrides (have defaults): `api_url` (defaults to staging),
`bq_location` (defaults to `US`).

A scenario needs the creds for its destination. Run only the scenarios whose
creds you've supplied (see the filter below).

## Run

`run.sh` builds the provider into `.bin/` (git-ignored), wires a Terraform
dev-override pointing at it, then for each selected scenario runs
`init → apply → verify outputs → assert no drift → destroy`. It always destroys
the in-flight scenario on exit (success or failure).

```sh
# All scenarios:
./run.sh

# One scenario:
./run.sh customerio
./run.sh customerio_audience

# Custom tfvars path (flag or env):
./run.sh --tfvars path/to/other.tfvars customerio
TFVARS_FILE=path/to/other.tfvars ./run.sh
```

`TF_CLI_CONFIG_FILE` does not need to be set beforehand — `run.sh` writes a
temporary dev-override config and exports it.

Keep-going: if a scenario fails, it's destroyed and the run continues; the
script exits non-zero at the end listing the failed scenarios.

### Hold resources open for inspection (`PAUSE=true`)

Apply, then block before destroy so you can inspect the live resources in the
staging UI. Press Enter to tear down. Most useful with a single scenario:

```sh
PAUSE=true ./run.sh customerio
```

### `--backfill` (opt-in, not yet wired)

```sh
./run.sh --backfill
```

Intended to trigger a manual rETL sync and poll to completion. **Currently exits
with error 3** — the rudder-iac client (`api/client/retl/connections.go`) does
not expose a sync-trigger endpoint yet. See the comment block in `run.sh`.

## Add a new scenario

1. `mkdir test/e2e/scenarios/<name>`
2. Symlink the shared config: `ln -s ../_shared/_shared.tf test/e2e/scenarios/<name>/_shared.tf` (commit the symlink).
3. Write `main.tf`: call `module "bq"` (`source = "../../modules/bigquery_source"`) with a unique `name_prefix`, add the destination + typed connection, and expose the four standardized outputs (`account_id`, `source_id`, `destination_id`, `connection_id`).
4. If the scenario needs a new credential variable, declare it once in `scenarios/_shared/_shared.tf` and add it to `secret.tfvars` (and, for CI, the workflow generator + `e2e` environment).

No `run.sh` change needed — scenarios are discovered automatically (dirs under
`scenarios/`, skipping `_`-prefixed ones).

## What success looks like

Per scenario, `run.sh` prints four non-empty IDs and asserts no drift:

```
    account_id     = "<id>"
    source_id      = "<id>"
    destination_id = "<id>"
    connection_id  = "<id>"
==> [<scenario>] ASSERT PASSED: plan reports no drift.
```

`terraform plan -detailed-exitcode` exit 0 means the provider Read path matches
config (no drift). Exit 2 (drift) or 1 (error) fail the scenario loudly before
destroy runs. At the end: `All scenarios passed: …`.

## CI

`.github/workflows/e2e-staging-smoke.yml` runs this harness (opt-in via the
`e2e-staging` PR label or manual dispatch). It generates `secret.tfvars` from the
`e2e` GitHub Environment and selects scenarios based on which vendor creds are
configured, so a missing optional cred skips that scenario rather than failing.
