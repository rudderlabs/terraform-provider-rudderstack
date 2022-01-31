terraform {
  required_providers {
    rudderstack = {
      source  = "rudderlabs/rudderstack"
      version = "~> 0.2.12"
    }
  }
  required_version = "~> 1.1.0"
}

provider "rudderstack" {
  # Set to control plane API host. Usually "https://api.rudderlabs.com/v2".
  # If null, falls back on env variable RUDDERSTACK_HOST.
  host = null

  # Set to access token for control plane API host. If null, falls back on env variable RUDDERSTACK_TOKEN.
  #
  # The token can be generated using steps below:
  # 1. Click <a href="https://app.rudderstack.com/home">here</a> and login.
  # 2. Click on Settings at the bottom left.
  # 3. Select "Personal Access Tokens".
  # 4. Create a new Personal Access Token and copy it for pasting here 
  token = null

  # Set to V1 control plane API host to be used. Usually "https://api.rudderlabs.com".
  # If null, falls back on env variable RUDDERSTACK_SCHEMA_HOST.
  schema_host = null

  # Set to access token for V1 control plane API host. If null, falls back on env variable RUDDERSTACK_SCHEMA_TOKEN.
  #
  # The token can be retrieved using steps below.
  # 1. Click <a href="https://app.rudderstack.com/home">here</a> and login.
  # 2. Copy the hexadecimal token string specified above data plane URL for pasting here.
  schema_token = null
}

resource "rudderstack_source" "src1" {
  name = "tfsource"
  type = "Auth0"
  config = {
    object = { 
    }
  }
}

resource "rudderstack_destination" "pub_sub_events" {
  allow_same_name = true
  config = {
    object = {
        "eventDelivery" = {
          bool = false
        },
        "eventDeliveryTS" = {
          str = "1637088430000"
        },
        "eventToTopicMap" = {
          list = [
              {
                object = {
                    "from" = {
                      str = "*"
                    },
                    "to" = {
                      str = "rudderstack_events"
                    },
                  }
              },
            ]
        },
        "projectId" = {
          str = "<valid project id>"
        },
      }
  }
  name   = "Pub-Sub rudderstack_events"
  type   = "GOOGLEPUBSUB"
}

resource "rudderstack_destination" "dst1" {
  name = "tfdestination"
  type = "GA"
  config = {
    object = { 
      "trackingID": { str = "UA-908213012-193" },
      "doubleClick": { bool = true },
      "enhancedLinkAttribution": { bool = true },
      "includeSearch": { bool = true },
      "enableServerSideIdentify": { bool = true },
      "serverSideIdentifyEventCategory": { str = "mnd,msdnf" },
      "serverSideIdentifyEventAction": { str = ",mn,m" },
      "disableMd5": { bool = true },
      "anonymizeIp": { bool = true },
      "enhancedEcommerce": { bool = true },
      "nonInteraction": { bool = true },
      "sendUserId": { bool = true },
      "dimensions": {
        list = [
          {
             object = {
               "from": { str = "mas." },
               "to": { str = "3" },
             }
          }
        ]
      },
      "blacklistedEvents" = { list = [ { object = { } } ] },
      "whitelistedEvents" = { list = [ { object = { } } ] },
      "metrics": {
        list = [
          {
             object = {
               "from": { str = "kksad1222" },
               "to": { str = "3" },
             }
          }
        ]
      },
      "contentGroupings": {
        list = [
          {
             object = {
               "from": { str = "lkjdlkjsdf" },
               "to": { str = "lkjlkjsdf" },
             }
          }
        ]
      },
      "eventFilteringOption" = { str = "disable" }
    },
  }
}

resource "rudderstack_source" "src0" {
  name = "test-1"
  type = "Javascript"
  config = { object = {} }
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
