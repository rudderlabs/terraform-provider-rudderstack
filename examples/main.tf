terraform {
  required_providers {
    rudderstack = {
      source  = "rudderstack"
      version = "~> 0.2.12"
    }
  }
  required_version = "~> 1.1.0"
}

provider "rudderstack" {
  # access_token = ""
}
