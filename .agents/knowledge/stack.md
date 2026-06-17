# Stack

> Dependencies, frameworks, tooling.
> Append-only. Agent-authored sections may optionally carry an HTML-comment tag
> (e.g., `<!-- pr:<id> -->`) identifying the writer/PR/run; human-authored
> sections are conventionally left untouched by automated runs.

## Core stack

- Language/runtime: Go `1.25.8` from `go.mod`.
- Provider framework: `github.com/hashicorp/terraform-plugin-sdk/v2 v2.40.1`.
- RudderStack API client layer: `github.com/rudderlabs/rudder-iac v0.15.0`.
- Terraform docs generation: `github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@v0.24.0`.
- Linting: `github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.9.0`.
- HCL/state shaping for the generator: `github.com/hashicorp/hcl/v2 v2.24.0`,
  `github.com/zclconf/go-cty v1.18.1`, `github.com/tidwall/gjson v1.18.0`, and
  `github.com/tidwall/sjson v1.2.5`.

## Tooling and release

- `Makefile` is the primary task runner: `build`, `install`, `lint`, `docs`,
  `test`, `test-ci`, `testacc-*`, and `release`.
- Release automation is split across `release-please-config.json`,
  `.github/workflows/release-please.yml`, and `.github/workflows/release.yml`.
- Test automation lives in `.github/workflows/test-ci.yml`,
  `.github/workflows/e2e-tests.yml`, and `.github/workflows/lint.yml`.
- The release workflow uses `.github/actions/import-gpg/action.yml` for signing
  and GoReleaser for publish-time artifacts.
- `scripts/bootstrap-terraform.sh` and `scripts/bootstrap-terraform-import.sh`
  support local example/config bootstrapping.
- `cmd/generatetf/main.go` is the offline generator used to snapshot live API
  state into HCL or import commands.

## Repo-local automation
<!-- RUD-2790 2026-06-03 -->

- `.claude/skills/onboard-integration/SKILL.md` is the repo-local automation
  hook for integration onboarding; it complements the generated-code workflow
  described in `README.md` without changing the provider runtime stack.

## INF-2349 — test-ci requires Terraform CLI

- `make test-ci` runs the Terraform Plugin SDK unit-test suite through
  `go test ./...`, and those non-acceptance `resource.UnitTest` tests invoke
  the Terraform CLI.
- The `test-ci` target depends on a `test-ci-preflight` check that requires
  `terraform` to already be present on `PATH`. The preflight is intentionally
  diagnostic only; it does not download Terraform.
- GitHub Actions performs Terraform setup before invoking `make test-ci`, while
  local and DevBox environments must provide Terraform themselves.
