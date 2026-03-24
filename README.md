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
3. **For E2E testing** — A `RUDDERSTACK_ACCESS_TOKEN` in a `.env` file at the repo root (the skill reads it automatically). Also ensure `~/.terraformrc` has `dev_overrides` configured for the local provider binary (see [Setting up the development environment](#example) above).

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
7. **E2E testing** — Builds the provider, creates a resource via `terraform apply` with all config fields, and runs the verify script (`cmd/integration-verify/`) to compare against the live RudderStack API.

### Adding new fields to an existing integration
When you run the skill for an integration that already exists, it:
1. Compares the config JSON files against the current `.go` implementation.
2. Shows you a table of **new fields** that exist in the config but are missing from the code.
3. Lets you select which fields to add (all or specific ones).
4. Updates the `.go`, test, example, and docs files, then runs the full validation and E2E pipeline.

> **Note:** This skill does not support modifying existing fields, fixing types, refactoring, or any other changes to already-implemented integrations. For those, make changes manually.

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
