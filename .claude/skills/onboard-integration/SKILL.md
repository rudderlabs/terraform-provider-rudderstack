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
3. **Config files from `rudder-integrations-config`** — the 3 config JSON files are required. **Default: fetch the latest from GitHub.** The repo is public, so no auth or MCP connector is needed, and fetching `main` guarantees you validate against the currently published config rather than a stale local clone (a stale clone is the main cause of missed or removed fields):

   ```bash
   # {kind} = destinations | sources ; {name} = the integration's config folder
   base="https://raw.githubusercontent.com/rudderlabs/rudder-integrations-config/main/src/configurations/{kind}/{name}"
   for f in db-config schema ui-config; do
     curl -fsSL "$base/$f.json" -o "$f.json" || echo "MISSING: $f.json"
   done
   ```

   **Opt-in override — local clone:** use a local checkout *only* if the user explicitly provides a path or asks for it, reading the 3 files from `<path>/src/configurations/{kind}/{name}/`.

Once you have the config files (from either source), verify all 3 exist:
- `db-config.json`
- `schema.json`
- `ui-config.json`

Read ALL THREE files. If any are missing (a `curl` printed `MISSING:`, or the name/kind is wrong), tell the user which are missing and stop — all three files are required.

### Check for Existing Integration

Before proceeding, check if an integration with the same or a similar name already exists in the provider:

```bash
# For destinations:
ls rudderstack/integrations/destinations/destination_*{name}*.go 2>/dev/null
# For sources — check sources.go for a matching Register call:
grep -i '{name}' rudderstack/integrations/sources/sources.go 2>/dev/null
```

If **no match** is found, this is a new integration — continue to Step 1.

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

## Step 1.6: Field-by-Field Validation Table (self-verification gate)

Before generating any `.go`, build an exhaustive field-by-field validation table following [reference/field-validation.md](reference/field-validation.md). One row per field — the union of every `schema.json` property and every `db-config` `defaultConfig`/`destConfig` key — with each cell copied from a **named source location** (Go type, verbatim `pattern`, `enum`, `required`, `default`, secret-from-`secretKeys`, per-source-type scoping + allowed values).

Then run the completeness check in that reference. **This is a hard gate:**
- If every field reconciles, proceed to Step 2 automatically and include the table in the PR description / tech spec so the reviewer can verify what was generated.
- If any field cannot be reconciled across the three files (present in one but not another, ambiguous type, a `pattern` you cannot faithfully translate to RE2), STOP and ask — do not guess.

The most common first-draft errors are silent omissions the model doesn't notice: a per-source-type `connection_mode` enum that differs between source types, a field-level regex left off, a nested source-type object flattened to a scalar. Enumerating every field against its source file — rather than summarizing from memory — is what catches them.

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

## Step 7: E2E Acceptance Testing

The E2E test function was already added in Step 2 (see [reference/e2e-testing.md](reference/e2e-testing.md)). Now validate it works.

Run plan-only validation (no API calls, no token needed):

```bash
# For destinations:
TF_ACC=1 TF_ACC_PLAN_ONLY=1 go test ./rudderstack/integrations/destinations/ -run "(?i)TestAccDestination.*{name}" -v -count=1

# For sources:
TF_ACC=1 TF_ACC_PLAN_ONLY=1 go test ./rudderstack/integrations/sources/ -run "(?i)TestAccSource.*{name}" -v -count=1
```

If plan-only passes and `RUDDERSTACK_ACCESS_TOKEN` is available (from `.env`), also run the full CRUD test:

```bash
# For destinations:
make testacc-dest DEST={name}

# For sources:
make testacc-source SRC={name}
```

If the full CRUD test fails, analyze the error and fix. Common issues:
- Config field type mismatches between Terraform schema and API
- Missing required fields in TerraformCreate
- API validation errors for placeholder values

---

## Summary Checklist

Before finishing, verify:

- [ ] Field-by-field validation table completed and reconciled against `schema.json` + `db-config.json` + `ui-config.json` (Step 1.6) — every property accounted for
- [ ] All files created/modified
- [ ] `init()` function registers the integration (self-registering, no provider.go changes needed)
- [ ] Unit tests pass
- [ ] E2E acceptance test function added (`TestAccDestination*` or `TestAccSource*`)
- [ ] E2E plan-only validation passes
- [ ] `make docs` generates the docs file
- [ ] Full test suite passes (`go test ./...`)
- [ ] Lint passes (`make lint`)
- [ ] Code follows existing patterns exactly (imports, naming, structure)

Present the user with a summary of all files created/modified and test results.
