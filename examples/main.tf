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

resource "rudderstack_source" "cdp" {
  name = "tfsource"
  type = "Auth0"
  config = {
    id = 0
  }
}

output "edu_order" {
  value = rudderstack_source.cdp
}
