terraform {
  required_providers {
    rudderstack = {
      version = "~> 0.0.1"
      source  = "rudderstack.com/cdp/rudderstack"
    }
  }
  required_version = "~> 1.0.3"
}

provider "rudderstack" {
  # Set to default "https://api.rudderlabs.com/v2/"
  host = null

  # Set Workspace token that can be generated at https://app.dev.rudderlabs.com/profile/tokens. Example "1fasdasdsdas".
  token = null 

  # Set to default "https://api.rudderlabs.com/"
  catalog_host = null

  # Set Catalog token available at https://app.rudderstack.com/home. Example "1lajsdlkjasdl".
  catalog_token = null
}

resource "rudderstack_source" "src1" {
  name = "tfsource"
  type = "Auth0"
  config = {
    id = 0
  }
}

resource "rudderstack_destination" "dst1" {
  name = "tfdestination"
  type = "Slack"
  config = {
    id = 0
  }
}

resource "rudderstack_connection" "cnxn1" {
  source_id = "${rudderstack_source.src1.id}" 
  destination_id = "${rudderstack_destination.dst1.id}"
}

output "src1_id" {
  value = rudderstack_source.src1.id
}

output "dst1_id" {
  value = rudderstack_destination.dst1.id
}

output "cnxn1_id" {
  value = rudderstack_connection.cnxn1.id
}
