terraform {
  required_providers {
    rudderstack = {
      source  = "rudderlabs/rudderstack"
      version = "~> 2.0.1"
    }
  }
  required_version = "~> 1.6.6"
}

provider "rudderstack" {
  # api_url      = "https://api.rudderstack.com/v2"
  # access_token = ""
}
