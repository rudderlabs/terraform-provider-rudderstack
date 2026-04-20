terraform {
  required_providers {
    rudderstack = {
      source  = "rudderlabs/rudderstack"
      version = "~> 4.4.0"
    }
  }
  required_version = "~> 1.10.5"
}

provider "rudderstack" {
  # api_url      = "https://api.rudderstack.com"
  # access_token = ""
}
