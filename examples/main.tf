terraform {
  required_providers {
    rudderstack = {
      source  = "rudderlabs/rudderstack"
      version = "~> 0.8.2"
    }
  }
  required_version = "~> 1.1.0"
}

provider "rudderstack" {
  # api_url      = "https://api.rudderstack.com/v2"
  # access_token = ""
}
