---
page_title: "Import Existing Resources"
---

# Import existing resources

The Terraform Provider lets you bootstrap your Terraform scripts and import the state of your already-configured pipelines in RudderStack.

> To generate a Terraform script from the existing resources, Golang 1.16 or later is required.

## Step 1: Clone the provider repository

Generating the Terraform scripts is possible using a tool provided in the provider's [GitHub repository](https://github.com/rudderlabs/terraform-provider-rudderstack). Clone this repo and navigate to the new directory using the following commands:

```sh
git clone https://github.com/rudderlabs/terraform-provider-rudderstack
cd terraform-provider-rudderstack
```

## Step 2: Generate your personal access token

Running the tool requires a RudderStack personal access token to be set in the `RUDDERSTACK_ACCESS_TOKEN` environemnt variable. For more information about generating this token, refer to the [RudderStack Documentation](https://www.rudderstack.com/docs/rudderstack-api/personal-access-tokens/).

Once you set this token, you can generate the script using the following command:

```sh
go run ./scripts/bootstrap-terraform.sh
```

The following example sets the personal access token and outputs the script to the `rudderstack.tf` file:

```sh
RUDDERSTACK_ACCESS_TOKEN=my_rudderstack_access_token ./scripts/bootstrap-terraform.sh > rudderstack.tf
```

## Step 3: Generate the Terraform script

The generated Terraform script will include any resources supported by this provider, named after each resource's unique ID. However, the RudderStack provider block needs to be configured independently. You can either add the provider block at the top of the generated script or in another `tf` file, depending on your preferred structure. 

For more information about setting a provider block, refer to the [Provider](https://registry.terraform.io/providers/rudderlabs/rudderstack/latest/docs) documentation. 

An example of a RudderStack provider block is shown below:

```terraform
terraform {
  required_providers {
    rudderstack = {
      source  = "rudderlabs/rudderstack"
      version = "~> 0.8.0"
    }
  }
  required_version = "~> 1.1.0"
}

provider "rudderstack" {
  # api_url      = "https://api.rudderstack.com/v2"
  # access_token = ""
}
```

> Note that in the above example, the access token is not set directly. In this case, the provider will read it from the  `RUDDERSTACK_ACCESS_TOKEN` environment variable.

## Step 4: Import Terraform state

Once the scripts are setup, Terraform needs to import the state of your resource to be able to manage them. This can be done by running the Terraform import commands.

For example, for a Redshift destination named `dest_dev` with the destination ID `id`, the Terraform import command will be as follows:

```sh
terraform import rudderstack_destination_redshift.dest_dev id
```

Since your generated script might contain many resources, remembering the IDs of all of them and running the Terraform commands manually can be very tedious. To avoid this, you can use the same tool used to generate the Terraform script that lists all the Terraform commands for any imported resource:

```sh
RUDDERSTACK_ACCESS_TOKEN=my_rudderstack_access_token ./scripts/bootstrap-terraform-import.sh
```

## Step 5: Add any sensitive configuration fields

The RudderStack API does not expose any sensitive credentials in the resource configurations. Because of this, the generated Terraform script will not include these credentials and they need to be added manually to each resource. Therefore, refer to each resource's documentation for all the relevant sensitive resource configuration fields.
