<p align="center"><a href="https://rudderstack.com"><img src="https://user-images.githubusercontent.com/59817155/126267034-ae9870b7-9137-4f45-be65-d621b055a972.png" alt="RudderStack - Customer Data Platform for Developers" height="50"/></a></p>
<h1 align="center"></h1>
<p align="center"><b>Customer Platform for Developers</b></p>
<br/>

<p align="center">
  <b>
    <a href="https://rudderstack.com">Website</a>
    ·
    <a href="https://rudderstack.com/docs/stream-sources/rudderstack-sdk-integration-guides/rudderstack-javascript-sdk/">Documentation</a>
    ·
    <a href="https://rudderstack.com/join-rudderstack-slack-community">Community Slack</a>
  </b>
</p>


# RudderStack Terraform Provider

This repository contains all the necessary resources to help you implement the RudderStack Terraform provider. You can use this provider to programmatically access the [RudderStack control plane](https://www.rudderstack.com/docs/get-started/rudderstack-architecture/#control-plane) via Terraform and seamlessly manage your source-destination configurations.

Questions? Join our [Slack community](https://resources.rudderstack.com/join-rudderstack-slack) for quick support.

# Getting started

- If you new to the Terraform platform, then their [docs](https://www.terraform.io/intro) are a good place to start.
- If you are interested in enhancing the RudderStack Terraform provider, create a local build and test your environment using an example configuration listed [here](#example).
- To manage your production RudderStack resources via Terraform, get detailed documentation for the Terraform Provider [here](docs/index.md).

## Terraform scripting flowchart for RudderStack

To create and maintain RudderStack's resources in Terraform, you can follow the below flowchart below:

![Flowchart for building and managing RudderStack's Terraform config](docs/TerraformScriptingForRudderStackFlowchart.png)

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

2. Next, make sure that your `~/.terraformrc` has the following lines: 

```
provider_installation {
  dev_overrides {
    "rudderlabs/rudderstack" = "~/.terraform.d/plugins/rudderstack.com/rudderlabs/rudderstack/0.2.12/linux_amd64/"
  }
}
```

The above snippet ensures that you use the locally built Terraform provider binary instead of the one available at the [Terraform Registry](https://registry.terraform.io).

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
