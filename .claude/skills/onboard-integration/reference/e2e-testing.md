# E2E Testing Steps

Automatically proceed with E2E testing. Do NOT ask the user whether to run E2E tests — just run them.

## 1. Load the access token from `.env`

```bash
cat .env 2>/dev/null | grep RUDDERSTACK_ACCESS_TOKEN || echo "NOT_FOUND"
```

- If `.env` contains `RUDDERSTACK_ACCESS_TOKEN`, use that value automatically. Tell the user: "Using access token from `.env` file."
- If `.env` is missing or doesn't contain the token, ask the user: "No `RUDDERSTACK_ACCESS_TOKEN` found in `.env`. Please provide your access token."

## 2. Check for `dev_overrides` in `~/.terraformrc`

```bash
cat ~/.terraformrc 2>/dev/null || echo "No .terraformrc found"
```

- If `dev_overrides` exist for `rudderlabs/rudderstack`, note the override path and source.
- If overrides exist, you must use the overridden source in the `required_providers` block and **skip `terraform init`** (it fails with dev_overrides). Instead, copy the built binary to the override path.

## 3. Build and install the provider locally

```bash
make install
```

If dev_overrides exist, also copy the binary to the override path:
```bash
cp $(go env GOPATH)/bin/terraform-provider-rudderstack {override_path}/
```

## 4. Create a temporary workspace

```bash
mkdir -p /tmp/tf-test-{name} && cd /tmp/tf-test-{name}
```

## 5. Write a `main.tf`

Write a `main.tf` with the provider config and a resource:

```hcl
terraform {
  required_providers {
    rudderstack = {
      source = "rudderlabs/rudderstack"  # use override source if dev_overrides exist
    }
  }
}

provider "rudderstack" {}

resource "rudderstack_{type}_{name}" "test" {
  name = "e2e-test-{name}"
  config {
    // ALL config fields (required AND optional) with realistic test values.
    // This ensures the E2E test validates every field mapping, not just required ones.
    // Use the same field values from TerraformUpdate in the unit test.
  }
}
```

## 6. Run Terraform

Use the token loaded from `.env` (or provided by the user) as an environment variable prefix for all terraform and verify commands:

```bash
# Skip `terraform init` if dev_overrides are active — it will fail.
# Only run `terraform init` if there are NO dev_overrides.
RUDDERSTACK_ACCESS_TOKEN="{token}" terraform init  # skip if dev_overrides
RUDDERSTACK_ACCESS_TOKEN="{token}" terraform plan
RUDDERSTACK_ACCESS_TOKEN="{token}" terraform apply -auto-approve
```

## 7. Verify the resource was created

```bash
terraform show
```

Check that the resource has a valid ID.

## 8. Run the verify script

Deterministically compare the .tf config against the API:

```bash
RESOURCE_ID=$(terraform show -json | jq -r '.values.root_module.resources[0].values.id')
cd <terraform-provider-rudderstack repo path>
go run ./cmd/integration-verify/ -file /tmp/tf-test-{name}/main.tf -id "$RESOURCE_ID"
```

This performs a subset comparison: every config key from the .tf file must exist and match in the API response. If the verify script reports FAIL, investigate the differences before proceeding.

## 9. Ask before cleaning up

Ask the user: "Would you like to verify the resource from the RudderStack dashboard first, or can I go ahead and delete it?" Wait for confirmation before proceeding. Do NOT proceed with destroy until the user explicitly confirms.

## 10. Clean up (only after user confirms)

```bash
terraform destroy -auto-approve
cd /
rm -rf /tmp/tf-test-{name}
```
