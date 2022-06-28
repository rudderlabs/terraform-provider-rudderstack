---
page_title: "rudderstack_destination_google_tag_manager Resource - terraform-provider-rudderstack"
subcategory: ""
description: |-
  
---

# rudderstack_destination_google_tag_manager (Resource)

This resource represents a Google Tag Manger destination. For more information check 
https://www.rudderstack.com/docs/destinations/tag-managers/google-tag-manager/

## Example Usage

```terraform
resource "rudderstack_destination_google_tag_manager" "example" {
  name = "my-google-tag-manager"

  config {
    container_id = "GTM-000000"

    server_url = "https://example.com"

    use_native_sdk {
      web = true
    }

    event_filtering {
      whitelist = ["one", "two", "three"]
      blacklist = ["one", "two", "three"]
    }

    onetrust_cookie_categories {
      web = ["one", "two", "three"]
    }
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `config` (Block List, Min: 1, Max: 1) Destination specific configuration. Check the nested block documenation for more information. (see [below for nested schema](#nestedblock--config))
- `name` (String) Human readable name of the destination. The value has to be unique across all destinations.

### Optional

- `enabled` (Boolean) An enabled destination allows data to be sent to it.

### Read-Only

- `created_at` (String) Time when the resource was created, in ISO 8601 format.
- `id` (String) The ID of this resource.
- `updated_at` (String) Time when the resource was last updated, in ISO 8601 format.

<a id="nestedblock--config"></a>
### Nested Schema for `config`

Required:

- `container_id` (String) Enter your Google Tag Manager container ID.

Optional:

- `event_filtering` (Block List, Max: 1) With this option, you can determine which events are blocked or allowed to flow through to GTM. (see [below for nested schema](#nestedblock--config--event_filtering))
- `onetrust_cookie_categories` (Block List, Max: 1) Specify the OneTrust category name for mapping the OneTrust consent settings to RudderStack's consent purposes. (see [below for nested schema](#nestedblock--config--onetrust_cookie_categories))
- `server_url` (String) Specify Tag Manager server container URL.
- `use_native_sdk` (Block List, Max: 1) Enable this setting to send the events via the device mode. (see [below for nested schema](#nestedblock--config--use_native_sdk))

<a id="nestedblock--config--event_filtering"></a>
### Nested Schema for `config.event_filtering`

Optional:

- `blacklist` (List of String) Enter the event names to be blacklisted.
- `whitelist` (List of String) Enter the event names to be whitelisted.


<a id="nestedblock--config--onetrust_cookie_categories"></a>
### Nested Schema for `config.onetrust_cookie_categories`

Optional:

- `web` (List of String)


<a id="nestedblock--config--use_native_sdk"></a>
### Nested Schema for `config.use_native_sdk`

Optional:

- `web` (Boolean)