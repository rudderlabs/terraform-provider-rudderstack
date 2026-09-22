# Feedback

> Human direction, preferences, corrections, or review guidance.
> Append-only. Agent-authored sections may optionally carry an HTML-comment tag
> (e.g., `<!-- pr:<id> -->`) identifying the writer/PR/run; human-authored
> sections are conventionally left untouched by automated runs.

## CFD-71 — Bing Ads Offline Conversions docs URL

- For `templates/resources/destination_bingads_offline_conversions.md.tmpl`, use the RudderStack docs URL `https://www.rudderstack.com/docs/destinations/streaming-destinations/bingads-offline-conversions/`. The public destination docs slug is `bingads-offline-conversions` and it lives under `streaming-destinations`, not `reverse-etl` or `bing-ads-offline-conversions`.

## RUD-3119 — Spotify Pixel field sensitivity

- For the Spotify Pixel Terraform destination, do not mark `config.pixel_id` as sensitive unless the upstream destination `db-config.json` lists `pixelId` under `secretKeys`; UI-only `secret: true` metadata is not sufficient for provider sensitivity.

## RUD-3119 — Spotify Pixel does have a connectionMode key (correction)

- The earlier `conventions.md` entry claiming Spotify Pixel "should not expose a Terraform `connection_mode` block" because "upstream includes no `connectionMode` config key" was wrong: `schema.json` defines `connectionMode.web` (enum `["device"]`) and `db-config.json` lists `supportedConnectionModes: {"web": ["device"]}`.
- Precedent elsewhere in the provider (e.g. Mixpanel's `connection_mode.android`, validated to `^(cloud)$`) is to still expose a per-source-type `connection_mode` field even when only one enum value is valid, rather than omitting the field because the UI offers no selector. Spotify Pixel now follows that pattern: `connection_mode { web = "device" }`, validated to `^(device)$`.

## RUD-3119 — `ui-config.json` defaults are console-only, not backend defaults (context, not something Spotify Pixel acts on)

- `ui-config.json`'s `defaultConnectionModes` (and similar UI-layer defaults) are applied by the RudderStack console's form-submission logic, not by the destinations Public API itself. Terraform talks to that API directly and never goes through the console, so a destination created via Terraform without an explicit `connectionMode` will NOT get the ui-config default applied on the backend's behalf — the field is simply absent from the stored config unless the provider sends it.
- This was explored for Spotify Pixel via `c.SimpleWithDefault("connectionMode.web", "connection_mode.0.web", "device")` (paired with `Optional + Computed` on the block and `Optional + Default` on the nested field, mirroring Amplitude's `sdkVersion.web`) — but that was deliberately reverted. **Final decision: `connectionMode.web` stays `c.Simple("connectionMode.web", "connection_mode.0.web", c.SkipZeroValue)` with plain `Optional: true`, no injected default.** This matches the universal convention every other `connectionMode.*` mapping in this provider uses (254 other mappings across ~22 destinations, all plain `c.Simple` + `c.SkipZeroValue`, none inject a default) — when the block is omitted, `connectionMode` is simply absent from the API payload, same as every sibling destination.
- Do not reintroduce `SimpleWithDefault` for `connectionMode` on this or other destinations without an explicit decision to do so — the console-only-default fact above is real, but the chosen tradeoff here is consistency with the rest of the provider over auto-correcting for it.

## INT-7142 — Generator redacted sensitive fields review guidance

- Do not add destination-specific guards such as `terraformType == "impact"` in `cmd/generatetf/generator` to work around redacted Required+Sensitive fields; this is a generic generator issue across destinations.
- Any future generated-HCL/import handling for redacted Required+Sensitive fields should derive policy from `configs.ConfigMeta` / registered schema metadata and be handled separately from individual destination onboarding PRs.
- For newly onboarded Terraform destination resources, do not include `(^env[.].+)` in new `c.StringMatchesRegexp(...)` validators. Keep existing dynamic-config support for the `{{ ... || ... }}` form when required, but do not add new `env.` validator support because env-based dynamic configuration is discontinued and upstream integration-config cleanup is pending.

## RUD-3134 — OAuth acceptance account lookup review guidance

- Match `client.Account.Definition.Type` directly against the lowercase Terraform destination key (for example, `linkedin_ads`) with `Definition.Category == "destination"`; do not substitute the uppercase destination API type or add a separate API-type lookup argument.
- OAuth account names can contain vendor login emails or usernames. Acceptance-test diagnostics should log only the selected account ID and type, never `Account.Name`.
- Keep `E2E_TESTING.md`'s `Plan-only exceptions` note: identify `TestAccSourceFacebookLeadAds` as intentionally skipping full CRUD because the catalog-hidden source is blocked by the backend resource gate, and retain the grep command that discovers all `acc.PlanOnly()` skip branches.
