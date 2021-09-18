terraform {
  required_providers {
    rudderstack = {
      version = "~> 0.0.1"
      source  = "hashicorp.com/edu/rudderstack"
    }
  }
  required_version = "~> 1.0.3"
}

provider "hashicups" {
  username = "education"
  password = "test123"
  host     = "http://localhost:19090"
}

resource "hashicups_order" "edu" {
  items = [{
    coffee = {
      id = 3
    }
    quantity = 2
    }, {
    coffee = {
      id = 1
    }
    quantity = 2
    }
  ]
}

output "edu_order" {
  value = hashicups_order.edu
}
