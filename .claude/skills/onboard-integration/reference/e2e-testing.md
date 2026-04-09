# E2E Acceptance Testing

When onboarding a new integration, add an E2E acceptance test alongside the unit test. The E2E framework reuses the same `TestConfig` data as unit tests.

## For Destinations

The E2E test function is added to the same `_test.go` file as the unit test. It reuses the extracted `var` containing `[]c.TestConfig`.

Add this function to `destination_{name}_test.go`:

```go
func TestAccDestination{PascalCaseName}(t *testing.T) {
	acc.AccAssertDestination(t, "{name}", {camelCaseName}TestConfigs)
}
```

This requires the `acc` import:
```go
acc "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/acc"
```

## For Sources

Add this function to `sources_test.go`:

```go
func TestAccSource{PascalCaseName}(t *testing.T) {
	acc.AccAssertSource(t, "{name}", emptyTestConfigs)
}
```

## What the E2E test does

The `AccAssertDestination` / `AccAssertSource` helpers run in two modes:

### Plan-only mode (`TF_ACC_PLAN_ONLY=1`)
- Validates HCL config + provider schema
- Zero API calls — uses a dummy token
- Runs on every PR for all integrations

### Full CRUD mode (`TF_ACC=1`, no `TF_ACC_PLAN_ONLY`)
- Creates the resource via real API, then verifies API config matches `APICreate` JSON
- Updates it with the `TerraformUpdate` config, then verifies API config matches `APIUpdate` JSON
- Imports it by ID (`ImportStateVerify` checks state consistency)
- Destroys it
- Runs only for affected integrations in CI

## Running locally

```bash
# Plan-only (no API calls, no token needed):
make testacc-plan

# Full CRUD for a single destination (requires .env with token):
make testacc-dest DEST={name}

# Full CRUD for a single source:
make testacc-source SRC={name}
```

## CI enforcement

A coverage test in `internal/testutil/acc/coverage_test.go` verifies every registered integration has a `TestAccDestination*` or `TestAccSource*` function. If you forget to add the E2E test, `make test-ci` will fail.
