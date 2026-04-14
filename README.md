<p align="center"><a href="https://rudderstack.com"><img src="https://user-images.githubusercontent.com/59817155/126267034-ae9870b7-9137-4f45-be65-d621b055a972.png" alt="RudderStack - Customer Data Platform for Developers" height="50"/></a></p>
<h1 align="center"></h1>
<p align="center"><b>Customer Platform for Developers</b></p>
<br/>

<p align="center">
  <b>
    <a href="https://rudderstack.com">Website</a>
    ·
    <a href="https://registry.terraform.io/providers/rudderlabs/rudderstack/latest/docs">Documentation</a>
    ·
    <a href="https://rudderstack.com/join-rudderstack-slack-community">Community Slack</a>
  </b>
</p>


# RudderStack Terraform Provider

This repository implements the RudderStack Terraform provider. You can use this provider to programmatically access the [RudderStack control plane](https://www.rudderstack.com/docs/get-started/rudderstack-architecture/#control-plane) via Terraform and seamlessly manage your source-destination configurations.

Questions? Join our [Slack community](https://resources.rudderstack.com/join-rudderstack-slack) for quick support.

# Getting started

- If you new to the Terraform platform, then their [docs](https://www.terraform.io/intro) are a good place to start.
- If you are interested in enhancing the RudderStack Terraform provider, create a local build and test your environment using an example configuration listed [here](#example).
- To manage your production RudderStack resources via Terraform, get detailed documentation for the Terraform Provider [here](docs/index.md).

<a id="example"></a>
# Setting up the development and build environment

Before you get started, make sure that following dependencies are installed:

- Bash (On Windows, consider using WSL2 Ubuntu) 
- Go
- Make

## Building the Terraform provider

1. Run the following command to build and install the provider:

```shell
$ make
$ make install
```

2. Next, make sure that your `/.terraformrc` file has the following lines: 

```
provider_installation {
  dev_overrides {
    "rudderlabs/rudderstack" = "/users/abc/.terraform.d/plugins/rudderstack.com/rudderlabs/rudderstack/0.2.12/linux_amd64/"
  }
}
```

The above snippet ensures that you use the locally built Terraform provider binary instead of the one available at the [Terraform Registry](https://registry.terraform.io).

A few things to note here:

- Use the full path, that is, `/users/xyz/.terraform.d/....` depending on your home directory.
- The `linux_amd64` part depends on your system's architecture. For example, on Macs it is `darwin_amd64`. This is essentially the path that is reported after the `make install` command runs.

## Testing the sample configuration

1. Navigate to the `examples` directory: 

```shell
$ cd examples
```

2. Run the following command to initialize the workspace and apply the sample configuration:

```shell
$ terraform init && terraform apply
```

## Making a new release

### Create a new tag with properly named version number

``` shell
git tag vX.Y.Z
git push
```

### Create new release

``` shell
goreleaser release --rm-dist
``` 

# Onboarding New Integrations with Claude Code

This repo includes a Claude Code skill (`/onboard-integration`) that automates onboarding source/destination integrations to the Terraform provider.

## What can this skill do?

| Scenario | Supported? |
|---|---|
| Onboard a **brand new** integration | Yes |
| **Add new fields** to an existing integration (from latest config JSON) | Yes |
| Refactor, fix types, update descriptions on existing integrations | No — make those changes manually |

## Prerequisites

Before running the skill, make sure you have:

1. **[Claude Code CLI](https://claude.com/claude-code)** installed
2. **Integration config files** — The skill needs 3 JSON files (`db-config.json`, `schema.json`, `ui-config.json`) from the [`rudder-integrations-config`](https://github.com/rudderlabs/rudder-integrations-config) repo. You can provide them in one of three ways:
   - **Auto-detect** — Clone `rudder-integrations-config` as a sibling directory (i.e., `../rudder-integrations-config`). The skill finds it automatically.
   - **GitHub fetch** — If you have the GitHub MCP connector configured, the skill can fetch the files directly from GitHub.
   - **Manual path** — Provide the absolute path to your local clone when prompted.
3. **For E2E testing** — A `.env` file at the repo root with `RUDDERSTACK_ACCESS_TOKEN` and `RUDDERSTACK_API_URL` (see [E2E Testing](#e2e-testing) below).

## Usage

```
/onboard-integration <name> <source|destination>
```

**Examples:**
```
/onboard-integration slack destination
/onboard-integration shopify source
```

If you omit the name or type, the skill will prompt you.

## How it works

### New integration
1. **Gathers inputs** — Parses integration name, type, and locates config JSON files.
2. **Checks for duplicates** — If an integration with the same name already exists, it offers to add new fields instead (see below).
3. **Extracts metadata** — Reads `db-config.json`, `schema.json`, and `ui-config.json` to determine field names, types, validations, defaults, descriptions, consent management config, and source-type-specific fields (e.g., `connectionMode`, `useNativeSDK`).
4. **Studies a similar integration** — Reads an existing integration with similar field patterns as a reference to ensure consistency.
5. **Generates files** — Creates the `.go` implementation, unit tests, example `.tf`, and docs template following the repo's established patterns.
6. **Validates** — Runs unit tests, generates docs (`make docs`), runs the full test suite (`go test ./...`), and lints (`make lint`).
7. **E2E testing** — Runs plan-only validation (`make testacc-plan`) and, if a `RUDDERSTACK_ACCESS_TOKEN` is available, runs full CRUD acceptance tests against the live API.

### Adding new fields to an existing integration
When you run the skill for an integration that already exists, it:
1. Compares the config JSON files against the current `.go` implementation.
2. Shows you a table of **new fields** that exist in the config but are missing from the code.
3. Lets you select which fields to add (all or specific ones).
4. Updates the `.go`, test, example, and docs files, then runs the full validation and E2E pipeline.

> **Note:** This skill does not support modifying existing fields, fixing types, refactoring, or any other changes to already-implemented integrations. For those, make changes manually.

<a id="e2e-testing"></a>
# E2E Testing

The provider includes an E2E acceptance test framework that validates integrations against the real RudderStack API. Every registered source and destination has a corresponding `TestAcc*` function.

## Two Modes

| Mode | Env vars | What it does | API calls |
|---|---|---|---|
| **Plan-only** | `TF_ACC=1 TF_ACC_PLAN_ONLY=1` | Validates HCL config + provider schema | Zero |
| **Full CRUD** | `TF_ACC=1` | Create → Update → Import → Destroy | ~5 per integration |

- **Plan-only** runs for all integrations on every PR. It sets a dummy token automatically — no credentials needed.
- **Full CRUD** runs only for integrations affected by the PR, using real API credentials. After each Create and Update step, the test fetches the resource from the API and verifies its config matches the expected `APICreate`/`APIUpdate` JSON from the test configs.

## Setup

Create a `.env` file at the repo root (git-ignored):

```
RUDDERSTACK_ACCESS_TOKEN=your-access-token
RUDDERSTACK_API_URL=https://api.rudderstack.com/v2
```

The Makefile auto-loads `.env` via `-include .env` + `export`.

> **Note:** `RUDDERSTACK_API_URL` must include `/v2` — the API client uses it as the complete base URL.

## Running Locally

```bash
# Plan-only validation for all integrations (no token needed):
make testacc-plan

# Full CRUD for a single destination:
make testacc-dest DEST=webhook

# Full CRUD for a single source:
make testacc-source SRC=http

# Full CRUD for connection tests:
make testacc-conn

# Full CRUD for everything:
make testacc-all
```

The `-run` patterns use `(?i)` for case-insensitive matching, so `DEST=webhook` matches `TestAccDestinationWebhook`.

## Architecture

The framework lives in `internal/testutil/acc/`:

| File | Purpose |
|---|---|
| `provider.go` | Provider factory, `TestAccPreCheck()`, `PlanOnly()`, dummy token helper |
| `destinations.go` | `AccAssertDestination()` — plan-only or full CRUD for destinations |
| `sources.go` | `AccAssertSource()` — plan-only or full CRUD for sources |
| `connections.go` | `AccAssertConnection()` — tests source→destination wiring |
| `config_verify.go` | Shared API config comparison logic (subset check) |
| `coverage_test.go` | CI enforcement — fails if any registered integration lacks a `TestAcc*` function |

E2E tests reuse the same `[]configs.TestConfig` data as unit tests — no duplicated HCL.

## Adding E2E Tests for New Integrations

### Destinations

In `destination_{name}_test.go`, extract configs to a package-level var and add the E2E function:

```go
var exampleTestConfigs = []c.TestConfig{
    {TerraformCreate: `...`, APICreate: `...`, TerraformUpdate: `...`, APIUpdate: `...`},
}

func TestDestinationResourceExample(t *testing.T) {
    cmt.AssertDestination(t, "example", exampleTestConfigs)
}

func TestAccDestinationExample(t *testing.T) {
    acc.AccAssertDestination(t, "example", exampleTestConfigs)
}
```

### Sources

In `sources_test.go`, add the E2E function:

```go
func TestAccSourceExample(t *testing.T) {
    acc.AccAssertSource(t, "example", emptyTestConfigs)
}
```

### Connections

In `connections/connections_test.go`:

```go
func TestAccConnectionExampleToWebhook(t *testing.T) {
    acc.AccAssertConnection(t, acc.ConnectionTestConfig{
        Source:      "example",
        Destination: "webhook",
        DestConfig:  `webhook_url = "https://example.com/test"
                      webhook_method = "POST"`,
    })
}
```

## CI Workflow

The GitHub Actions workflow (`.github/workflows/e2e-tests.yml`) runs on every PR:

1. **detect-changes** — Determines which integrations are affected by the PR
2. **plan-only** — Validates all integration configs with zero API calls
3. **e2e-crud** — Runs full CRUD tests only for affected integrations (matrix strategy, max 5 parallel)
4. **e2e-summary** — Gates the PR on both plan-only and CRUD results

Core file changes (provider, configs, client, acc helpers) trigger CRUD tests for **all** integrations.

## CI Enforcement

A coverage test (`internal/testutil/acc/coverage_test.go`) runs during `make test-ci` and fails if any registered integration is missing its `TestAcc*` function. This uses case-insensitive matching — exact PascalCase is not required.

# Related

- Learn more about the [RudderStack architecture](https://www.rudderstack.com/docs/get-started/rudderstack-architecture/) to understand the difference between RudderStack control plane and data plane.
- https://github.com/rudderlabs/cp-client-go: This repository implements the REST API client for the RudderStack control plane in Golang.
- https://github.com/rudderlabs/rscp_pyclient: This repository implements the REST API client for the RudderStack control plane in Python. A few additional RudderStack-related methods are also available.

<!--
   1) https://github.com/rudderlabs/segment-migrator : Source code for segment migrator web app. Helps migrate from
      Segment to RudderStack.
   1) http://segment-migrator.dev-rudder.rudderlabs.com/ : If you are trying to migrate from Segment to RudderStack, you can use this web app to migrate. 
-->
# License

The RudderStack Terraform Provider is released under the [MIT License][mit_license].

# Contribute

We would love to see you contribute to RudderStack. Get more information on how to contribute [here](CONTRIBUTING.md).

# Follow us

- [RudderStack Blog][rudderstack-blog]
- [Slack][slack]
- [Twitter][twitter]
- [LinkedIn][linkedin]
- [dev.to][devto]
- [Medium][medium]
- [YouTube][youtube]
- [HackerNews][hackernews]
- [Product Hunt][producthunt]

<!----variables---->

[slack]: https://resources.rudderstack.com/join-rudderstack-slack
[twitter]: https://twitter.com/rudderstack
[linkedin]: https://www.linkedin.com/company/rudderlabs/
[devto]: https://dev.to/rudderstack
[medium]: https://rudderstack.medium.com/
[youtube]: https://www.youtube.com/channel/UCgV-B77bV_-LOmKYHw8jvBw
[rudderstack-blog]: https://rudderstack.com/blog/
[hackernews]: https://news.ycombinator.com/item?id=21081756
[producthunt]: https://www.producthunt.com/posts/rudderstack
[mit_license]: https://opensource.org/licenses/MIT
[agplv3_license]: https://www.gnu.org/licenses/agpl-3.0-standalone.html
[sspl_license]: https://www.mongodb.com/licensing/server-side-public-license
[config-generator]: https://github.com/rudderlabs/config-generator
[config-generator-section]: https://github.com/rudderlabs/rudder-server/blob/master/README.md#rudderstack-config-generator
[rudder-logo]: https://repository-images.githubusercontent.com/197743848/b352c900-dbc8-11e9-9d45-4deb9274101f
