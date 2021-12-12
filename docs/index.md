# RudderStack Provider
Use the RudderStack's Terraform provider plugin to interact with control plane API of the RudderStack CDP from within Terraform.

## Supported Resources 
   1. [Source](resources/source.md)
   1. [Destination](resources/destination.md)
   1. [Connection](resources/connection.md)

## Example Usage 
Terraform 1.0.9 and later:
```
terraform {
  required_providers {
    rudderstack = {
      source  = "rudderlabs/rudderstack"
      version = "~> 0.2.7"
    }
  }
  required_version = "~> 1.0.9"
}

provider "rudderstack" {
  # Set to control plane API host. Usually "https://api.rudderlabs.com/v2/".
  # If null, falls back on env variable RUDDERSTACK_HOST.
  host = null

  # Set to access token for control plane API host. If null, falls back on env variable RUDDERSTACK_TOKEN.
  #
  # The token can be generated using steps below:
  # 1. Click <a href="https://app.rudderstack.com/home">here</a> and login.
  # 2. Click on Settings at the bottom left.
  # 3. Select "Personal Access Tokens".
  # 4. Create a new Personal Access Token and copy it for pasting here.
  token = null 

  # [Deprecated, Optional]
  # Set to V1 control plane API host to be used. Usually "https://api.rudderlabs.com/".
  # If null, falls back on env variable RUDDERSTACK_SCHEMA_HOST.
  schema_host = null

  # [Deprecated, Optional]
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
        objects_list = [
          {
             object = {
               "from": { str = "mas." },
               "to": { str = "3" },
             }
          }
        ]
      },
      "metrics": {
        objects_list = [
          {
             object = {
               "from": { str = "kksad1222" },
               "to": { str = "2" },
             }
          }
        ]
      },
      "contentGroupings": {
        objects_list = [
          {
             object = {
               "from": { str = "lkjdlkjsdf" },
               "to": { str = "lkjlkjsdf" },
             }
          }
        ]
      },
    },
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

```

### Attributes 

- **host** (String, Optional)
  Set to control plane API host. Usually "https://api.rudderlabs.com/v2/". If null, it falls back on env variable RUDDERSTACK_HOST. For the provider to work, value must be set, one way or another.
- **token** (String, Sensitive, Optional)
  Set to access token for control plane API host. The token can be generated at https://app.rudderlabs.com/profile/tokens. Example "1fasdasdsdas". If null, it falls back on env variable RUDDERSTACK_TOKEN. For the provider to work, value must be set, one way or another.
- **schema_host** (String, Optional, *Deprecated*)
  Set to V1 control plane API host to be used. Usually "https://api.rudderlabs.com/". If null, falls back on env variable RUDDERSTACK_SCHEMA_HOST.
- **schema_token** (String, Sensitive, Optional, *Deprecated*)
  Set to access token for V1 control plane API host. The token is available at https://app.rudderstack.com/home. Example "1lajsdlkjasdl". If null, falls back on env variable RUDDERSTACK_SCHEMA_TOKEN.

