---
page_title: "Provider: RudderStack"
subcategory: "Customer Data Platform"
description: |-
  Terraform provider for interacting with RudderStack control plane API.
---

# RudderStack Provider

-> Visit the [Call APIs with Terraform Providers](https://learn.hashicorp.com/collections/terraform/providers?utm_source=WEBSITE&utm_medium=WEB_IO&utm_offer=ARTICLE_PAGE&utm_content=DOCS) Learn tutorials for an interactive getting started experience.

The RudderStack provider is used to interact with the control plane of RudderStack Customer Data Platform product. Using this provider, one can create sources, destinations, connections and transformations([TBD]) in the Rudderstack cloud:

Underneath, this provider uses [REST API client in Golang for](https://github.com/rudderlabs/cp-client-go) for interfacing with RudderStack control plane.

Use the navigation to the left to read about the available resources.

## Example Usage

Do not keep your authentication password in HCL for production environments, use Terraform environment variables.

```terraform
provider "hashicups" {
  username = "education"
  password = "test123"
}
```

## Schema

### Optional

- **username** (String, Optional) Username to authenticate to HashiCups API
- **password** (String, Optional) Password to authenticate to HashiCups API
- **host** (String, Optional) HashiCups API address (defaults to `localhost:19090`)
