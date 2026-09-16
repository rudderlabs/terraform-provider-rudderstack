# Field Validation Reference

The self-verification gate that runs **before** any `.go` is generated. Its job is to prove every config field is accounted for and correctly translated, so the first draft doesn't repeat the recurring silent misses: per-source-type connection-mode enums, field-level regex patterns, enums, secret flags, and nested source-type blocks flattened to scalars.

Build on [config-extraction.md](config-extraction.md) — it owns the type mapping, RE2 regex translation, and the `GetCommonConfigMeta` skip-list. This file owns the *table* and the *completeness check*.

## How the gate runs

1. Build one table row for **every** field — the union of:
   - every property in `schema.json` → `configSchema.properties`, **and**
   - every key in `db-config.json` → `destConfig.defaultConfig` and each `destConfig.<sourceType>`.
2. Fill each cell by **reading a named source location**, not from memory.
3. Run the completeness check below.
   - **All rows reconcile →** proceed to code generation automatically and include the table in the PR description / tech spec so the reviewer can see what was generated.
   - **Any row can't be reconciled →** STOP and ask. Do not guess.

It is a *self-verification* gate: a clean table needs no human sign-off; only an irreconcilable row pulls a human in.

## Columns (one row per field)

| Column | Read it from | Notes |
|---|---|---|
| API key | `schema.json` property name / `db-config` key | camelCase |
| TF key | derived | snake_case |
| Go type | `schema.json` `type` → mapping in config-extraction.md | |
| Validator | `schema.json` `pattern` → `c.StringMatchesRegexp("<verbatim>")`; `enum` → `StringInSlice` / `^(a|b)$` | copy the pattern **verbatim**, then RE2-translate per config-extraction.md |
| Required | is it in `schema.json` `required[]`? | else Optional |
| Default | `schema.json` `default` | omit the column if none |
| Secret | is it in `db-config` `secretKeys`? | `Sensitive` **only** if yes — never from ui-config `secret: true` |
| Description | `ui-config.json` — this field's `label` / `footerNote` | becomes the TF `Description:`; this is the column read from `ui-config.json` (so the table genuinely spans all three files) |
| Source-type scoped | is the property an object keyed by source types? | if yes → one sub-row per source type + its allowed values |
| Skip | in the `GetCommonConfigMeta` skip-list? | `consentManagement` / `oneTrustCookieCategories` / `ketchConsentPurposes` |

## Per-source-type keys — the most-missed class

If a `schema.json` property is an **object whose keys are source types** (e.g. `connectionMode`, `useNativeSDK`, `match_id`), it is NOT a flat field:

- Add one sub-row **per source type**, each carrying that source type's own allowed values. For `connectionMode`, cross-check `db-config.json` → `supportedConnectionModes.<sourceType>`.
- Model it as `TypeList` + `MaxItems: 1` with one nested field per source type, mapped `c.Simple("<key>.<sourceType>", "<key>.0.<snake_sourceType>", c.SkipZeroValue)`.
- **Omit `c.SkipZeroValue` when a source-scoped field has a non-zero `schema.json` default** (e.g. a bool defaulting to `true`, like `useNativeSDK.web`). `SkipZeroValue` drops the Go zero value (`false`/`""`), so an explicit `false` would be silently discarded and revert to the backend default — the user could never turn the field off. Pair the TF-schema `Default` with a plain `c.Simple` in that case.
- **Enums can differ per source type.** e.g. `android: ["cloud","device"]` but `web: ["cloud"]` → the android validator is `^(cloud|device)$` and web's is `^(cloud)$`. Never apply one blanket enum to all source types.

## Completeness check — all must hold before codegen

- [ ] Every `schema.json` property maps to exactly one of: a field row, a source-scoped block (which itself expands to one sub-row per source type), or the skip-list — nothing unaccounted. Count *properties* covered, not table rows (a source-scoped property is many rows but one property).
- [ ] Every `db-config` `supportedSourceTypes` entry has a `connection_mode` sub-row (unless the destination has no `connectionMode` property at all).
- [ ] Every field with a `schema.json` `pattern` has a `StringMatchesRegexp` carrying that verbatim pattern (RE2-translated where needed).
- [ ] Every field with a `schema.json` `enum` has a matching validator.
- [ ] `Sensitive` is set on exactly the `db-config` `secretKeys` fields — no more, no fewer.
- [ ] Every per-source-type object property is a nested block, not a scalar.

## Worked example (excerpt)

A destination whose `schema.json` has `apiKey` (has `pattern`, listed in `secretKeys`), `campaignId` (`pattern` `^[0-9]+$`, in `required`), and `connectionMode` (object: `android: ["cloud","device"]`, `web: ["cloud"]`):

| API key | TF key | Go type | Validator | Required | Secret | Source-scoped |
|---|---|---|---|---|---|---|
| apiKey | api_key | string | `StringMatchesRegexp("<verbatim>")` | yes | **yes** (secretKeys) | no |
| campaignId | campaign_id | string | `StringMatchesRegexp("…\|^[0-9]+$")` | yes | no | no |
| connectionMode.android | connection_mode.0.android | string | `StringMatchesRegexp("^(cloud\|device)$")` | no | no | android |
| connectionMode.web | connection_mode.0.web | string | `StringMatchesRegexp("^(cloud)$")` | no | no | web |

`android` allows `device`, `web` does not — the per-source enums come from `supportedConnectionModes`, not a shared default.
