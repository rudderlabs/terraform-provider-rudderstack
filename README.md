<p align="center"><a href="https://rudderstack.com"><img src="https://user-images.githubusercontent.com/59817155/126267034-ae9870b7-9137-4f45-be65-d621b055a972.png" alt="RudderStack - Customer Data Platform for Developers" height="50"/></a></p>
<h1 align="center"></h1>
<p align="center"><b>Customer Platform for Developers</b></p>
<br/>


# Terraform Provider RudderStack 


# Repo description

In the repo, we have the implementation of RudderStack terraform provider as a bridge between Terraform and the RudderStack control plane API.  
|:------|

Questions? Please join our [Slack channel](https://resources.rudderstack.com/join-rudderstack-slack) or read about us on [Product Hunt](https://www.producthunt.com/posts/rudderstack).

# Why Use Terraform RudderStack provider 

\*\* Describe what software does. \*\*

# Key Features

\*\* Describe key features, if necessary. \*\*

# Get Started

## Build provider

Run the following command to build the provider

```shell
$ go build -o terraform-provider-hashicups
```

## Test sample configuration

First, build and install the provider.

```shell
$ make install
```

Then, navigate to the `examples` directory. 

```shell
$ cd examples
```

Run the following command to initialize the workspace and apply the sample configuration.

```shell
$ terraform init && terraform apply
```

# License

RudderStack \*\* software name \*\* is released under the [MIT License][mit_license].

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

# \*\* Optional \*\* :clap:  Our Supporters

\*\* Update the repo names below. \*\*<br />

[![Stargazers repo roster for @rudderlabs/rudder-server](https://reporoster.com/stars/rudderlabs/rudder-repo-template)](https://github.com/rudderlabs/rudder-repo-template/stargazers)
[![Forkers repo roster for @rudderlabs/rudder-server](https://reporoster.com/forks/rudderlabs/rudder-repo-template)](https://github.com/rudderlabs/rudder-repo-template/network/members)

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
# Terraform Provider Hashicups

This repo is a companion repo to the [Call APIs with Terraform Providers](https://learn.hashicorp.com/collections/terraform/providers) Learn collection. 

In the collection, you will use the HashiCups provider as a bridge between Terraform and the HashiCups API. Then, extend Terraform by recreating the HashiCups provider. By the end of this collection, you will be able to take these intuitions to create your own custom Terraform provider. 

