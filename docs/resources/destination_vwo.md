---
page_title: "rudderstack_destination_vwo Resource - terraform-provider-rudderstack"
subcategory: ""
description: |-
  
---

# rudderstack_destination_vwo (Resource)

This resource represents a VWO destination. For more information check 
https://www.rudderstack.com/docs/destinations/testing-and-personalization/vwo-beta-visual-website-optimizer

## Example Usage

```terraform
resource "rudderstack_destination_vwo" "example" {
  name = "my-vwo"

  config {
    account_id = "..."

    # is_spa                   = false
    # send_experiment_track    = false
    # send_experiment_identify = false

    # library_tolerance  = "2000"
    # settings_tolerance = "2000"

    # use_existing_jquery = false

    # use_native_sdk {
    #   web = true
    # }

    # event_filtering {
    #   whitelist = ["one", "two", "three"]
    #   blacklist = ["one", "two", "three"]
    # }
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

- `account_id` (String) Enter your VWO account ID.

Optional:

- `event_filtering` (Block List, Max: 1) Specify which events should be blocked or allowed to flow through to VWO. (see [below for nested schema](#nestedblock--config--event_filtering))
- `is_spa` (Boolean) Enable this setting if the page is a single page application (SPA).
- `library_tolerance` (String) Enter the value for the library tolerance setting.
- `send_experiment_identify` (Boolean) Enable this setting to send the experiments viewed as `identify` traits.
- `send_experiment_track` (Boolean) Enable this setting to send the experiment data as `track` events.
- `settings_tolerance` (String) Enter the value for the setting tolerance.
- `use_existing_jquery` (Boolean) Enable this setting to use the existing jQuery.
- `use_native_sdk` (Block List, Max: 1) Enable this setting to send the events via the device mode. (see [below for nested schema](#nestedblock--config--use_native_sdk))

<a id="nestedblock--config--event_filtering"></a>
### Nested Schema for `config.event_filtering`

Optional:

- `blacklist` (List of String) Enter the event names to be blacklisted.
- `whitelist` (List of String) Enter the event nams to be whitelisted.


<a id="nestedblock--config--use_native_sdk"></a>
### Nested Schema for `config.use_native_sdk`

Optional:

- `web` (Boolean)