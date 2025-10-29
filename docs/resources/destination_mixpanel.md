---
page_title: "rudderstack_destination_mixpanel Resource - terraform-provider-rudderstack"
subcategory: ""
description: |-
  
---

# rudderstack_destination_mixpanel (Resource)

This resource represents a MixPanel destination. For more information check 
https://www.rudderstack.com/docs/destinations/streaming-destinations/mixpanel/

## Example Usage

```terraform
resource "rudderstack_destination_mixpanel" example{
  name = "my-mixpanel"

  config {
    token = "..."
    data_residency = "us"
    identity_merge_api = "simplified"  # Required field - can be "simplified" or "original"
    connection_mode {
      web = "cloud"
    }

    # api_secret = "..."
    # people = true
    # set_all_traits_by_default = true
    # consolidated_page_calls = true
    # track_categorized_pages = true
    # track_named_pages = true
    # source_name = "my-mixpanel"
    # cross_subdomain_cookie = true
    # secure_cookie = true
    # super_properties = ["one","two","three"]
    # people_properties = ["one","two","three"]
    # event_increments = ["one","two","three"]
    # prop_increments = ["one","two","three"]
    # group_key_settings = ["one","two","three"]
    # set_once_properties = ["one","two","three"]
    # union_properties = ["one","two","three"]
    # append_properties = ["one","two","three"]
    # gdpr_api_token = "..."
    # user_deletion_api = "engage"
    # persistence_name = "none"
    # persistence_type = "..."
    # ignore_dnt = false
    # use_user_defined_page_event_name = false
    # user_defined_page_event_template = "Viewed {{ category }} {{ name }} page"
    # use_user_defined_screen_event_name = false
    # user_defined_screen_event_template = "Viewed {{ category }} {{ name }} screen"
    # drop_traits_in_track_event = false

    # session_replay_percentage {
    #   web = "0"
    # }

    # use_native_sdk {
    #   web = true
    # }

    # event_filtering {
    #   whitelist = ["one", "two", "three"]
    #   blacklist = ["one","two","three"]
    # }

    # consent_management {
    # 	web = [
    # 		{
    # 			provider = "oneTrust"
    # 			consents = ["one_web", "two_web", "three_web"]
    # 			resolution_strategy = ""
    # 		},
    # 		{
    # 			provider = "ketch"
    # 			consents = ["one_web", "two_web", "three_web"]
    # 			resolution_strategy = ""
    # 		},
    # 		{
    # 			provider = "custom"
    # 			resolution_strategy = "and"
    # 			consents = ["one_web", "two_web", "three_web"]
    # 		}
    # 	]
    # 	android = [{
    # 		provider = "ketch"
    # 		consents = ["one_android", "two_android", "three_android"]
    # 		resolution_strategy = ""
    # 	}]
    # 	ios = [{
    # 		provider = "custom"
    # 		resolution_strategy = "and"
    # 		consents = ["one_ios", "two_ios", "three_ios"]
    # 	}]
    # 	unity = [{
    # 		provider = "custom"
    # 		resolution_strategy = "or"
    # 		consents = ["one_unity", "two_unity", "three_unity"]
    # 	}]
    # 	reactnative = [{
    # 		provider = "custom"
    # 		resolution_strategy = "and"
    # 		consents = ["one_reactnative", "two_reactnative", "three_reactnative"]
    # 	}]
    # 	flutter = [{
    # 		provider = "custom"
    # 		resolution_strategy = "and"
    # 		consents = ["one_flutter", "two_flutter", "three_flutter"]
    # 	}]
    # 	cordova = [{
    # 		provider = "custom"
    # 		resolution_strategy = "and"
    # 		consents = ["one_cordova", "two_cordova", "three_cordova"]
    # 	}]
    # 	amp = [{
    # 		provider = "custom"
    # 		resolution_strategy = "and"
    # 		consents = ["one_amp", "two_amp", "three_amp"]
    # 	}]
    # 	cloud = [{
    # 		provider = "custom"
    # 		resolution_strategy = "and"
    # 		consents = ["one_cloud", "two_cloud", "three_cloud"]
    # 	}]
    # 	warehouse = [{
    # 		provider = "custom"
    # 		resolution_strategy = "and"
    # 		consents = ["one_warehouse", "two_warehouse", "three_warehouse"]
    # 	}]
    # 	shopify = [{
    # 		provider = "custom"
    # 		resolution_strategy = "and"
    # 		consents = ["one_shopify", "two_shopify", "three_shopify"]
    # 	}]
    # }
    # use_new_mapping = true
  }
}
```

> **:warning: Breaking Change**
> 
> Note that from the provider versions 4.0.0 and above, 
> - `api_secret` property has been removed. 
> - `persistence` property has been removed.
> - the schema of `super_properties` property has been changed.
> - the schema of `people_properties` property has been changed.
> - the schema of `event_increments` property has been changed.
> - the schema of `prop_increments` property has been changed.
> - the schema of `group_key_settings` property has been changed.
> - the schema of `set_once_properties` property has been changed.
> - the schema of `union_properties` property has been changed.
> - the schema of `append_properties` property has been changed.
> - the schema of `identity_merge_api` property has been changed. It is a required property now.
> Please refer to the example above.

> **:warning: Breaking Change**
> 
> Note that from the provider versions 3.0.0 and above, `onetrust_cookie_categories` property is replaced with `consent_management` that supports multiple consent management providers. Please refer to the example above.

> **:warning: Breaking Change**
> 
> Note that from the provider versions 1.0.0 and above, the schema of `onetrust_cookie_categories` property has been changed. Please refer to the example above.

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

- `data_residency` (String) Mixpanel Server region either us/eu
- `identity_merge_api` (String) Mixpanel Identity Merge types
- `token` (String, Sensitive) Mixpanel API Token

Optional:

- `append_properties` (List of String) Properties to append
- `connection_mode` (Block List, Max: 1) Enable this setting to send the events via the cloud mode. (see [below for nested schema](#nestedblock--config--connection_mode))
- `consent_management` (Block List, Max: 1) Allows you to specify consent configuration data for multiple providers for each source type. (see [below for nested schema](#nestedblock--config--consent_management))
- `consolidated_page_calls` (Boolean) This will track Loaded a Page events to Mixpanel for all page method calls. We enable this by default as it's how Mixpanel suggests sending these calls.
- `cross_subdomain_cookie` (Boolean) This will allow the Mixpanel cookie to persist between different pages of your application.
- `drop_traits_in_track_event` (Boolean) Boolean flag to drop traits in track event calls
- `event_filtering` (Block List, Max: 1) With this option, you can determine which events are blocked or allowed to flow through to Mixpanel. (see [below for nested schema](#nestedblock--config--event_filtering))
- `event_increments` (List of String) Events to increment in People.
- `gdpr_api_token` (String, Sensitive)
- `group_key_settings` (List of String) Group Key
- `ignore_dnt` (Boolean) Ignore 'Do Not Track' setting
- `people` (Boolean) Boolean flag to send all of your identify calls to Mixpanel's People feature
- `people_properties` (List of String) Traits to set as People Properties.
- `persistence_name` (String) Mixpanel Persistence Name
- `persistence_type` (String) Mixpanel Persistence Type
- `prop_increments` (List of String) Properties to increment in People
- `secure_cookie` (Boolean) This will mark the Mixpanel cookie as secure, meaning it will only be transmitted over https.
- `session_replay_percentage` (Block List, Max: 1) Percentage of SDK initializations that will qualify for replay data capture (see [below for nested schema](#nestedblock--config--session_replay_percentage))
- `set_all_traits_by_default` (Boolean) While this is checked, our integration automatically sets all traits on identify calls as super properties and people properties if Mixpanel People is checked as well.
- `set_once_properties` (List of String) Properties to set only once
- `source_name` (String) This value, if it's not blank, will be sent as rudderstack_source_name to Mixpanel for every event/page/screen call.
- `strict_mode` (Boolean) Boolean flag to enable strict mode
- `super_properties` (List of String) Property to send as super Properties.
- `track_categorized_pages` (Boolean) This will track events to Mixpanel for page method calls that have a category associated with them. For example page('Docs', 'Index') would translate to Viewed Docs Index Page.
- `track_named_pages` (Boolean) This will track events to Mixpanel for page method calls that have a name associated with them. For example page('Signup') would translate to Viewed Signup Page.
- `union_properties` (List of String) Properties to union
- `use_native_sdk` (Block List, Max: 1) Enable this setting to send the events via the device mode. (see [below for nested schema](#nestedblock--config--use_native_sdk))
- `use_new_mapping` (Boolean) This value is true by default and when this flag is enabled, camel case fields are mapped to snake case fields while sending to Mixpanel. Please refer to https://www.rudderstack.com/docs/destinations/streaming-destinations/mixpanel/#connection-settings for more details.
- `use_user_defined_page_event_name` (Boolean) Boolean flag to use user-defined page event names
- `use_user_defined_screen_event_name` (Boolean) Boolean flag to use user-defined screen event names
- `user_defined_page_event_template` (String) Template for user-defined page event names
- `user_defined_screen_event_template` (String) Template for user-defined screen event names
- `user_deletion_api` (String)

<a id="nestedblock--config--connection_mode"></a>
### Nested Schema for `config.connection_mode`

Optional:

- `amp` (String)
- `android` (String)
- `cloud` (String)
- `cordova` (String)
- `flutter` (String)
- `ios` (String)
- `reactnative` (String)
- `shopify` (String)
- `unity` (String)
- `warehouse` (String)
- `web` (String)


<a id="nestedblock--config--consent_management"></a>
### Nested Schema for `config.consent_management`

Optional:

- `amp` (List of Object) Allows you to specify consent configuration data for multiple providers. (see [below for nested schema](#nestedatt--config--consent_management--amp))
- `android` (List of Object) Allows you to specify consent configuration data for multiple providers. (see [below for nested schema](#nestedatt--config--consent_management--android))
- `cloud` (List of Object) Allows you to specify consent configuration data for multiple providers. (see [below for nested schema](#nestedatt--config--consent_management--cloud))
- `cordova` (List of Object) Allows you to specify consent configuration data for multiple providers. (see [below for nested schema](#nestedatt--config--consent_management--cordova))
- `flutter` (List of Object) Allows you to specify consent configuration data for multiple providers. (see [below for nested schema](#nestedatt--config--consent_management--flutter))
- `ios` (List of Object) Allows you to specify consent configuration data for multiple providers. (see [below for nested schema](#nestedatt--config--consent_management--ios))
- `reactnative` (List of Object) Allows you to specify consent configuration data for multiple providers. (see [below for nested schema](#nestedatt--config--consent_management--reactnative))
- `shopify` (List of Object) Allows you to specify consent configuration data for multiple providers. (see [below for nested schema](#nestedatt--config--consent_management--shopify))
- `unity` (List of Object) Allows you to specify consent configuration data for multiple providers. (see [below for nested schema](#nestedatt--config--consent_management--unity))
- `warehouse` (List of Object) Allows you to specify consent configuration data for multiple providers. (see [below for nested schema](#nestedatt--config--consent_management--warehouse))
- `web` (List of Object) Allows you to specify consent configuration data for multiple providers. (see [below for nested schema](#nestedatt--config--consent_management--web))

<a id="nestedatt--config--consent_management--amp"></a>
### Nested Schema for `config.consent_management.amp`

Optional:

- `consents` (List of String)
- `provider` (String)
- `resolution_strategy` (String)


<a id="nestedatt--config--consent_management--android"></a>
### Nested Schema for `config.consent_management.android`

Optional:

- `consents` (List of String)
- `provider` (String)
- `resolution_strategy` (String)


<a id="nestedatt--config--consent_management--cloud"></a>
### Nested Schema for `config.consent_management.cloud`

Optional:

- `consents` (List of String)
- `provider` (String)
- `resolution_strategy` (String)


<a id="nestedatt--config--consent_management--cordova"></a>
### Nested Schema for `config.consent_management.cordova`

Optional:

- `consents` (List of String)
- `provider` (String)
- `resolution_strategy` (String)


<a id="nestedatt--config--consent_management--flutter"></a>
### Nested Schema for `config.consent_management.flutter`

Optional:

- `consents` (List of String)
- `provider` (String)
- `resolution_strategy` (String)


<a id="nestedatt--config--consent_management--ios"></a>
### Nested Schema for `config.consent_management.ios`

Optional:

- `consents` (List of String)
- `provider` (String)
- `resolution_strategy` (String)


<a id="nestedatt--config--consent_management--reactnative"></a>
### Nested Schema for `config.consent_management.reactnative`

Optional:

- `consents` (List of String)
- `provider` (String)
- `resolution_strategy` (String)


<a id="nestedatt--config--consent_management--shopify"></a>
### Nested Schema for `config.consent_management.shopify`

Optional:

- `consents` (List of String)
- `provider` (String)
- `resolution_strategy` (String)


<a id="nestedatt--config--consent_management--unity"></a>
### Nested Schema for `config.consent_management.unity`

Optional:

- `consents` (List of String)
- `provider` (String)
- `resolution_strategy` (String)


<a id="nestedatt--config--consent_management--warehouse"></a>
### Nested Schema for `config.consent_management.warehouse`

Optional:

- `consents` (List of String)
- `provider` (String)
- `resolution_strategy` (String)


<a id="nestedatt--config--consent_management--web"></a>
### Nested Schema for `config.consent_management.web`

Optional:

- `consents` (List of String)
- `provider` (String)
- `resolution_strategy` (String)



<a id="nestedblock--config--event_filtering"></a>
### Nested Schema for `config.event_filtering`

Optional:

- `blacklist` (List of String) Enter the event names to be denylisted.
- `whitelist` (List of String) Enter the event names to be allowlisted.


<a id="nestedblock--config--session_replay_percentage"></a>
### Nested Schema for `config.session_replay_percentage`

Optional:

- `web` (String)


<a id="nestedblock--config--use_native_sdk"></a>
### Nested Schema for `config.use_native_sdk`

Optional:

- `web` (Boolean)
