<p align="center"><a href="https://rudderstack.com"><img src="https://user-images.githubusercontent.com/59817155/126267034-ae9870b7-9137-4f45-be65-d621b055a972.png" alt="RudderStack - Customer Data Platform for Developers" height="50"/></a></p>
<h1 align="center"></h1>
<p align="center"><b>Customer Platform for Developers</b></p>
<br/>


# Terraform Provider RudderStack 

# Description

This repo implements RudderStack terraform provider. Use it to access RudderStack control plane API from within Terraform.  

Questions? Please join our [Slack channel](https://resources.rudderstack.com/join-rudderstack-slack) or read about us on [Product Hunt](https://www.producthunt.com/posts/rudderstack).

# Getting Started
Good place to start with Terraform is [here](https://www.terraform.io/intro/index.html). Next, checkout example configuration for RudderStack Terraform Provider [below](#example). Detailed documentation for RudderStack Terraform
Provider is available [here](docs/index.md). 

<a id="example"></a>
# Setup dev and build env 

## PreRequisites 
Make sure that following are installed.
1. bash (On Windows, consider using WSL2 Ubuntu) 
2. go
3. make

## Build provider

Run the following command to build and install the provider

```shell
$ make
$ make install
```

Next, make sure that your ~/.terraformrc has the following lines. 
```
provider_installation {
  dev_overrides {
    "rudderlabs/rudderstack" = "~/.terraform.d/plugins/rudderstack.com/rudderlabs/rudderstack/0.2.5/linux_amd64/"
  }
}
```
The above ensures that that locally built terraform provider binary is unsed instead of the one available at registry.terraform.io

## Test sample configuration
Navigate to the `examples` directory. 

```shell
$ cd examples
```

Run the following command to initialize the workspace and apply the sample configuration.

```shell
$ terraform init && terraform apply
```

# Related 
   1) https://github.com/rudderlabs/cp-client-go : This repo implements REST API client for RudderStack Control Plain in Golang.
   1) https://github.com/rudderlabs/rscp_pyclient : This repo implements REST API client for RudderStack Control Plain in Python. Few additional RudderStack related helpful methods also available.
   1) https://github.com/rudderlabs/segment-migrator : Source code for segment migrator web app. Helps migrate from
      Segment to RudderStack.
   1) https://segment-migrator.dev-rudder.rudderlabs.com/ : If you are trying to migrate from Segment to RudderStack, you can use this web app to migrate. 

# License

RudderStack Terraform Provider is released under the [MIT License][mit_license].

# Contribute

We would love to see you contribute to RudderStack. Get more information on how to contribute [here](CONTRIBUTING.md).

# Follow Us

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

