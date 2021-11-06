# Rudderstack Provider
Use the Rudderstack provider to interact with control plane API of the Ruddertack CDP.

## Resources 
   1. [Source](resources/source.md)
   1. [Destination](resources/destination.md)
   1. [Connection](resources/connection.md)

## Example Usage 
Terraform 1.1.0 and later:
```
terraform {
  required_providers {
    rudderstack = {
      version = "~> 0.0.1"
      source  = "rudderstack.com/cdp"
    }
  }
  required_version = "~> 1.1.0"
}

provider "rudderstack" {
  # Set to control plane API host. Usually "https://api.rudderlabs.com/v2/".
  # If null, falls back on env variable RUDDERSTACK_HOST.
  host = null

  # Set to access token for control plane API host.
  # The token can be generated at https://app.dev.rudderlabs.com/profile/tokens. Example "1fasdasdsdas".
  # If null, falls back on env variable RUDDERSTACK_TOKEN.
  token = null 

  # [Deprecated, Optional]
  # Set to V1 control plane API host to be used. Usually "https://api.rudderlabs.com/".
  # If null, falls back on env variable RUDDERSTACK_SCHEMA.
  schema_host = null

  # [Deprecated, Optional]
  # Set to access token for V1 control plane API host.
  # The token is available at https://app.rudderstack.com/home. Example "1lajsdlkjasdl".
  # If null, falls back on env variable RUDDERSTACK_SCHEMA_TOKEN.
  schema_token = null
}

```

### Attributes 

- **host** (String, Optional)
  Set to control plane API host. Usually "https://api.rudderlabs.com/v2/". If null, it falls back on env variable RUDDERSTACK_HOST. For the provider to work, value must be set, one way or another.
- **token** (String, Sensitive, Optional)
  Set to access token for control plane API host. The token can be generated at https://app.dev.rudderlabs.com/profile/tokens. Example "1fasdasdsdas". If null, it falls back on env variable RUDDERSTACK_TOKEN. For the provider to work, value must be set, one way or another.
- **schema_host** (String, Optional, *Deprecated*)
  Set to V1 control plane API host to be used. Usually "https://api.rudderlabs.com/". If null, falls back on env variable RUDDERSTACK_SCHEMA.
- **schema_token** (String, Sensitive, Optional, *Deprecated*)
  Set to access token for V1 control plane API host. The token is available at https://app.rudderstack.com/home. Example "1lajsdlkjasdl". If null, falls back on env variable RUDDERSTACK_SCHEMA_TOKEN.
