terraform {
  required_providers {
    rudderstack = {
      source  = "rudderstack.com/rudderlabs/rudderstack"
      version = "~> 0.3.0"
    }
  }
  required_version = "~> 1.1.0"
}

provider "rudderstack" {
  # access_token = ""
}
