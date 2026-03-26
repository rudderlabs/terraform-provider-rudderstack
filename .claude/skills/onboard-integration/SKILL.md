---
name: onboard-integration
description: Onboard a new source or destination integration to the terraform-provider-rudderstack by reading config JSON files from rudder-integrations-config and generating the .go implementation, tests, example .tf, and docs template.
argument-hint: "[name] [source|destination]"
disable-model-invocation: true
allowed-tools: Read, Write, Edit, Bash, Grep, Glob, Agent
---

# Onboard Integration

Add a new source or destination integration to the terraform-provider-rudderstack.

**Integration name:** `$0`
**Type:** `$1`

---

## Step 0: Gather Inputs

Parse positional arguments: integration name = `$0`, type = `$1`. If either is missing, ask the user:

1. **Integration name** — the snake_case name (e.g., `webhook`, `slack`, `google_analytics`).
2. **Type** — `source` or `destination`.
3. **Config files from `rudder-integrations-config`** — The 3 config JSON files for the integration are required. Try to locate them in this order:

   **Option A: Local sibling repo** — Auto-detect by checking for a sibling folder:
   ```bash
   ls -d ../rudder-integrations-config 2>/dev/null
   ```
   If found, use that path and tell the user: "Found `rudder-integrations-config` at `{resolved_path}`, using it."

   **Option B: Fetch from GitHub** — If the local repo is not found, ask the user:
   "I couldn't find `rudder-integrations-config` locally. Would you like me to:
   1. Fetch the config files directly from GitHub (requires the GitHub MCP connector)
   2. Provide the local path to your `rudder-integrations-config` clone"

   If fetching from GitHub, read the files from `https://github.com/rudderlabs/rudder-integrations-config` at:
   - `src/configurations/{destinations|sources}/{name}/db-config.json`
   - `src/configurations/{destinations|sources}/{name}/schema.json`
   - `src/configurations/{destinations|sources}/{name}/ui-config.json`

   **Option C: User provides path** — If the user provides a custom path, use that.

Once you have the config files (from any option), verify all 3 exist:
- `db-config.json`
- `schema.json`
- `ui-config.json`

Read ALL THREE files. If any are missing, tell the user which are missing and stop — all three files are required.

### Check for Existing Integration

Before proceeding, check if an integration with the same or a similar name already exists in the provider:

```bash
# For destinations:
ls rudderstack/integrations/destinations/destination_*{name}*.go 2>/dev/null
# For sources — check sources.go for a matching Register call:
grep -i '{name}' rudderstack/integrations/sources/sources.go 2>/dev/null
```

If a match is found, **stop and ask the user** what they'd like to do:

- "I found an existing integration: `{matched_file_or_name}`. What would you like to do?
  1. **Add new fields** from the latest config JSON that are missing in the current implementation
  2. **Something else** (refactor, fix types, update descriptions, etc.)"

If the user chooses **add new fields**:

1. Read the current implementation files (`.go`, `_test.go`, example `.tf`, and docs template).
2. Read the config JSON files (`db-config.json`, `schema.json`, `ui-config.json`) from `rudder-integrations-config`.
3. Compare the two to identify fields that exist in the config JSON but are **missing from the current `.go` implementation**.
4. Present the diff to the user in a table:

   "I found the following new fields comparing `rudder-integrations-config` with the current implementation:

   | # | Field (API key) | Type | Status |
   |---|---|---|---|
   | 1 | `newApiField` | string | optional |
   | 2 | `anotherField` | boolean | required |

   Would you like to add **all** of these fields, or select specific ones? (Enter `all` or comma-separated numbers like `1,2`)"

5. Wait for the user to respond. Only add the fields the user selects.
6. Proceed with the implementation, adding only the selected fields to the `.go`, `_test.go`, example `.tf`, and docs template files. Then continue from Step 3 (Run Unit Tests) onwards.

If the user chooses **something else**, **stop the skill** and tell them: "This skill only supports onboarding new integrations or adding new fields to existing ones. For other changes (refactoring, fixing types, updating descriptions, etc.), please make those changes manually."

Do NOT proceed if the user asks for anything other than adding new fields.

---

## Step 1: Extract Integration Metadata from Config Files

From the 3 JSON files, extract metadata following the rules in [reference/config-extraction.md](reference/config-extraction.md). This covers:

- `db-config.json` → APIType, supportedSourceTypes, secretKeys, defaultConfig
- `schema.json` → properties, types, patterns, enums, defaults, required fields
- `ui-config.json` → descriptions and labels
- Type mapping (JSON Schema → Terraform Go types)
- Fields to skip (handled by `GetCommonConfigMeta()`)
- Source-type-specific fields (connectionMode, useNativeSDK, etc.)
- Naming conventions (camelCase API → snake_case Terraform)

---

## Step 1.5: Study a Similar Existing Integration

Before generating files, find an existing destination/source with similar field types to the new integration and read its implementation and tests as a reference. This helps catch patterns the generic instructions might miss.

- **Simple fields only?** → Study `destination_webhook.go` and its test
- **Arrays of objects?** → Study `destination_slack.go` and its test
- **Event filtering / complex nested config?** → Look for a destination with similar structure in `rudderstack/integrations/destinations/`

Read both the `.go` file and the `_test.go` file to understand the exact patterns used.

---

## Step 2: Generate Files

Follow the code templates in [reference/destination-templates.md](reference/destination-templates.md) for destinations or [reference/source-templates.md](reference/source-templates.md) for sources.

### For DESTINATIONS — create these files:
1. `rudderstack/integrations/destinations/destination_{name}.go`
2. `rudderstack/integrations/destinations/destination_{name}_test.go`
3. `examples/destination_{name}.tf`
4. `templates/resources/destination_{name}.md.tmpl`

### For SOURCES — modify + create:
1. **Modify:** `rudderstack/integrations/sources/sources.go` — add `Register()` call
2. **Modify:** `rudderstack/integrations/sources/sources_test.go` — add test function
3. **Create:** `examples/source_{name}.tf`
4. **Create:** `templates/resources/source_{name}.md.tmpl`

See the template files for exact code patterns, important rules about `ConfigMode`, `SkipZeroValue`, `Sensitive`, defaults, and consent management test patterns.

---

## Step 3: Run Unit Tests

Run the unit test for the new integration:

```bash
# For destinations:
go test ./rudderstack/integrations/destinations/ -run TestDestinationResource{PascalCaseName} -v

# For sources:
go test ./rudderstack/integrations/sources/ -run TestSourceResource{PascalCaseName} -v
```

If tests fail, analyze the error output carefully:
- Schema mismatches → check field types and required/optional
- JSON key mismatches → check camelCase API keys vs snake_case terraform keys
- Consent management errors → ensure the API JSON wraps consents as `{"consent": "value"}` objects

Fix any failures and re-run until tests pass.

---

## Step 4: Generate Documentation

Run:
```bash
make docs
```

This generates `docs/resources/{destination|source}_{name}.md` from the template. If `make docs` fails, check:
- The template file exists at the correct path
- The example .tf file is valid HCL
- `tfplugindocs` is installed (run `go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest` if needed)

---

## Step 5: Run Full Test Suite

Run the full test suite to make sure nothing is broken:

```bash
go test ./... -v
```

If there are failures in other tests, investigate — the new integration should NOT affect existing tests since it self-registers via `init()`.

---

## Step 6: Lint Check

Run the linter:

```bash
make lint
```

Fix any lint issues in the generated code.

---

## Step 7: E2E Testing

Automatically proceed with E2E testing. Do NOT ask the user whether to run E2E tests — just run them. Follow the detailed E2E steps in [reference/e2e-testing.md](reference/e2e-testing.md).

High-level flow:
1. Load access token from `.env`
2. Check for `dev_overrides` in `~/.terraformrc`
3. Build and install provider (`make install`)
4. Create temp workspace, write `main.tf` with ALL config fields
5. Run `terraform plan` and `terraform apply`
6. Verify resource with `terraform show`
7. Run the verify script (`go run ./cmd/integration-verify/`)
8. **Ask before cleaning up** — wait for user confirmation before `terraform destroy`

---

## Summary Checklist

Before finishing, verify:

- [ ] All files created/modified
- [ ] `init()` function registers the integration (self-registering, no provider.go changes needed)
- [ ] Unit tests pass
- [ ] `make docs` generates the docs file
- [ ] Full test suite passes (`go test ./...`)
- [ ] Lint passes (`make lint`)
- [ ] Code follows existing patterns exactly (imports, naming, structure)

Present the user with a summary of all files created/modified and test results.
