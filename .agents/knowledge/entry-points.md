# Entry points

> Key entry-point files: read these first to orient in this repo.
> Append-only. Agent-authored sections may optionally carry an HTML-comment tag
> (e.g., `<!-- pr:<id> -->`) identifying the writer/PR/run; human-authored
> sections are conventionally left untouched by automated runs.

## First reads

- [README.md](/workspace/terraform-provider-rudderstack/README.md) - repo
  purpose, build/test workflow, and release process.
- [main.go](/workspace/terraform-provider-rudderstack/main.go) - Terraform
  plugin process entry point.
- [rudderstack/provider.go](/workspace/terraform-provider-rudderstack/rudderstack/provider.go)
  - provider schema, client construction, and resource registration.
- [rudderstack/integrations/integrations.go](/workspace/terraform-provider-rudderstack/rudderstack/integrations/integrations.go)
  - blank-import hub that activates source/destination registries.
- [cmd/generatetf/main.go](/workspace/terraform-provider-rudderstack/cmd/generatetf/main.go)
  - live API snapshot/export utility.
- [docs/index.md](/workspace/terraform-provider-rudderstack/docs/index.md)
  - generated provider documentation and example usage.

## Secondary reads

- [E2E_TESTING.md](/workspace/terraform-provider-rudderstack/E2E_TESTING.md)
  - how acceptance coverage is staged and what credentials it expects.

## Repo-local skill
<!-- RUD-2790 2026-06-03 -->

- `.claude/skills/onboard-integration/SKILL.md` is the repo-local Claude Code
  skill for onboarding integrations; it is the automation entry point named in
  `README.md` and should be considered before hand-editing generated integration
  files.
