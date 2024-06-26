---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "rudderstack_destination_LINKEDIN_INSIGHT_TAG Resource - terraform-provider-rudderstack"
subcategory: ""
description: |-
  
---

# rudderstack_destination_LINKEDIN_INSIGHT_TAG (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `config` (Block List, Min: 1, Max: 1) Destination specific configuration. Check the nested block documenation for more information. (see [below for nested schema](#nestedblock--config))
- `name` (String) Human readable name of the destination. The value has to be unique across all the destinations.

### Optional

- `enabled` (Boolean) An enabled destination allows data to be sent to it.

### Read-Only

- `created_at` (String) Time when the resource was created, in ISO 8601 format.
- `id` (String) The ID of this resource.
- `updated_at` (String) Time when the resource was last updated, in ISO 8601 format.

<a id="nestedblock--config"></a>
### Nested Schema for `config`

Required:

- `partner_id` (String) Enter your LinkedIn Partner ID.

Optional:

- `event_filtering` (Block List, Max: 1) With this option, you can determine which events are blocked or allowed to flow through to LinkedIn. (see [below for nested schema](#nestedblock--config--event_filtering))
- `event_to_conversion_id_map` (List of Object) Event Conversion IDs. (see [below for nested schema](#nestedatt--config--event_to_conversion_id_map))
- `onetrust_cookie_categories` (List of String) Specify the OneTrust category name for mapping the OneTrust consent settings to RudderStack's consent purposes.
- `use_native_sdk` (Block List, Max: 1) As this is a device mode destination, this setting will always be enabled. (see [below for nested schema](#nestedblock--config--use_native_sdk))

<a id="nestedblock--config--event_filtering"></a>
### Nested Schema for `config.event_filtering`

Optional:

- `blacklist` (List of String) Enter the event names to be denylisted.
- `whitelist` (List of String) Enter the event names to be allowlisted.


<a id="nestedatt--config--event_to_conversion_id_map"></a>
### Nested Schema for `config.event_to_conversion_id_map`

Optional:

- `from` (String)
- `to` (String)


<a id="nestedblock--config--use_native_sdk"></a>
### Nested Schema for `config.use_native_sdk`

Optional:

- `web` (Boolean)


