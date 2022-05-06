---
page_title: "Import existing resources"
---

# Import existing resources

This provider repo supports a tool that can help users with already configured pipelines in RudderStack to bootstrap their
Terraform scripts and import the state.

-> In order to generate a Terraform script from existing resources, golang 1.16 or later is required.

## Clone the provider repo

Generating Terraform scripts is possible using a tool provided in the provider's GitHub repo. Clone the repo and change to the new directory:

```sh
git clone https://github.com/rudderlabs/terraform-provider-rudderstack
cd terraform-provider-rudderstack
```

## Configure access token

Running the tool requires a RudderStack access token to be set in the `RUDDERSTACK_ACCESS_TOKEN` environmental variable. For more information about creating an access token, please check the relevant [RudderStack Documentation](https://www.rudderstack.com/docs/transformations/api-access-token/) page. Once this is set, you can generate the script using `go run ./cmd/generatetf`. The following example sets the access token and outputs the script to the `ruddestack.tf` file:

```sh
RUDDERSTACK_ACCESS_TOKEN=my_rudderstack_access_token go run ./cmd/generatetf > rudderstack.tf
```

## Generate Terraform script

The generated Terraform script will include any resources that are supported by this provider, named after each resource's unique ID. However, the rudderstack provider block needs to be configured independently. You can either add the provider block at the top of the generated script, or in another tf file, depending on your preferred structure. For more information about setting a provider block, please check the [Provider](https://registry.terraform.io/providers/rudderlabs/rudderstack/latest/docs) documentation. An example of a rudderstack provider block is:

```terraform
terraform {
  required_providers {
    rudderstack = {
      source  = "rudderlabs/rudderstack"
      version = "~> 0.3.0"
    }
  }
  required_version = "~> 1.1.0"
}

provider "rudderstack" {
  # access_token = ""
}
```

Note that, in this example, the access token is not set directly. In this case the provider will read it from the  `RUDDERSTACK_ACCESS_TOKEN` environmental variable.

## Import Terraform state

Once the scripts are setup, Terraform needs to import the state of your resource in order to be able to manage them.
This can happen by running the terraform import commands. As an example, for a postgres destination named `dst_some_id` with ID `some_id`, the terraform import command looks like:

```sh
terraform import rudderstack_destination_postgres.dst_some_id some_id
```

Since your generated script might contain many resources, of which you have to know their IDs, running the terraform commands by hand is tedious. You can use the same tool used for generating the Terraform script to list all the terraform commands for any imported resource:

```sh
RUDDERSTACK_ACCESS_TOKEN=my_rudderstack_access_token go run ./cmd/generatetf -import
```

## Add any sensitive configuration fields

RudderStack API is not exposing any sensitive credentials in resource configurations. Because of this, the generated terraform script will not include them and they need to be added manually to each resource. Please, refer to each resource's documentation for all available sensitive resource configuration fields.