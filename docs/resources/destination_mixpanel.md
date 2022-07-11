---
page_title: "rudderstack_destination_mixpanel Resource - terraform-provider-rudderstack"
subcategory: ""
description: |-
  
---

# rudderstack_destination_mixpanel (Resource)

This resource represents a MixPanel destination. For more information check 
https://www.rudderstack.com/docs/destinations/error-reporting/mixpanel

## Example Usage

```terraform
resource "rudderstack_destination_mixpanel" example{
  name = "my-mixpanel"

  config {
    token = "avasaffav1241"
    data_residency = "us"
    persistence = "none"
    # api_secret = ""
    people = true
    set_all_traits_by_default = true
    consolidated_page_calls = true
    track_categorized_pages = true
    track_named_pages = true
    source_name = "my-mixpanel"
    cross_subdomain_cookie = true
    secure_cookie = true
    super_properties = [
      {
        "property" = ""
      }
    ]
    people_properties = [
      {
        "property" = ""
      }
    ]
    event_increments = [
      {
        "property" = ""
      }
    ]
    prop_increments = [
      {
        "property" = ""
      }
    ]
    group_key_settings = [
      {
        "group_key" = ""
      }
    ]
    # use_native_sdk = [
    #   {
    #     "web" = true
    #   }
    # ]

    event_filtering {
      whitelist = ["one", "two", "three"]
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

- `token` (String, Sensitive) Mixpanel API Token

Optional:

- `api_secret` (String, Sensitive) Mixpanel API secret
- `consolidated_page_calls` (Boolean) This will track Loaded a Page events to Mixpanel for all page method calls. We enable this by default as it's how Mixpanel suggests sending these calls.
- `cross_subdomain_cookie` (Boolean) This will allow the Mixpanel cookie to persist between different pages of your application.
- `data_residency` (String) Mixpanel Server region either us/eu
- `event_filtering` (Block List, Max: 1) This option allows you filter the events you want to send to Amplitude. (see [below for nested schema](#nestedblock--config--event_filtering))
- `event_increments` (List of Object) (see [below for nested schema](#nestedatt--config--event_increments))
- `group_key_settings` (List of Object) (see [below for nested schema](#nestedatt--config--group_key_settings))
- `onetrust_cookie_categories` (Block List, Max: 1) Specify the OneTrust category name for mapping the OneTrust consent settings to RudderStack's consent purposes. (see [below for nested schema](#nestedblock--config--onetrust_cookie_categories))
- `people` (Boolean) Boolean flag to send all of your identify calls to Mixpanel's People feature
- `people_properties` (List of Object) (see [below for nested schema](#nestedatt--config--people_properties))
- `persistence` (String) Choose persistence for Mixpanel SDK. One of none|cookie|localStorage
- `prop_increments` (List of Object) (see [below for nested schema](#nestedatt--config--prop_increments))
- `secure_cookie` (Boolean) This will mark the Mixpanel cookie as secure, meaning it will only be transmitted over https
- `set_all_traits_by_default` (Boolean) While this is checked, our integration automatically sets all traits on identify calls as super properties and people properties if Mixpanel People is checked as well.
- `source_name` (String) This value, if it's not blank, will be sent as rudderstack_source_name to Mixpanel for every event/page/screen call.
- `super_properties` (List of Object) Property to send as super Properties (see [below for nested schema](#nestedatt--config--super_properties))
- `track_categorized_pages` (Boolean) This will track events to Mixpanel for page method calls that have a category associated with them. For example page('Docs', 'Index') would translate to Viewed Docs Index Page.
- `track_named_pages` (Boolean) This will track events to Mixpanel for page method calls that have a name associated with them. For example page('Signup') would translate to Viewed Signup Page.
- `use_native_sdk` (Block List, Max: 1) Enable this setting to send events to Mixpanel via the device mode. (see [below for nested schema](#nestedblock--config--use_native_sdk))

<a id="nestedblock--config--event_filtering"></a>
### Nested Schema for `config.event_filtering`

Optional:

- `blacklist` (List of String) Enter the event names to be blacklisted.
- `whitelist` (List of String) Enter the event names to be whitelisted.


<a id="nestedatt--config--event_increments"></a>
### Nested Schema for `config.event_increments`

Optional:

- `property` (String)


<a id="nestedatt--config--group_key_settings"></a>
### Nested Schema for `config.group_key_settings`

Optional:

- `group_key` (String)


<a id="nestedblock--config--onetrust_cookie_categories"></a>
### Nested Schema for `config.onetrust_cookie_categories`

Optional:

- `web` (List of String)


<a id="nestedatt--config--people_properties"></a>
### Nested Schema for `config.people_properties`

Optional:

- `property` (String)


<a id="nestedatt--config--prop_increments"></a>
### Nested Schema for `config.prop_increments`

Optional:

- `property` (String)


<a id="nestedatt--config--super_properties"></a>
### Nested Schema for `config.super_properties`

Optional:

- `property` (String)


<a id="nestedblock--config--use_native_sdk"></a>
### Nested Schema for `config.use_native_sdk`

Optional:

- `web` (Boolean)