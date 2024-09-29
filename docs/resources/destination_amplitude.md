---
page_title: "rudderstack_destination_amplitude Resource - terraform-provider-rudderstack"
subcategory: ""
description: |-
  
---

# rudderstack_destination_amplitude (Resource)

This resource represents an Amplitude destination. For more information check 
https://www.rudderstack.com/docs/destinations/analytics/amplitude

## Example Usage

```terraform
resource "rudderstack_destination_amplitude" "example" {
  name = "my-amplitude"

  config {
    api_key    = "amplitude api key"
    api_secret = "amplitude api secret"

    # group_type_trait  = ""
    # group_value_trait = ""

    # track_all_pages           = true
    # track_categorized_pages   = true
    # track_named_pages         = true
    # track_products_once       = true
    # track_revenue_per_product = true

    # track_gclid {
    #   web = true
    # }

    # track_referrer {
    #   web = true
    # }

    # track_utm_properties {
    #   web = true
    # }

    # track_session_events {
    #   android      = true
    #   ios          = true
    #   react_native = true
    # }

    # version_name = ""

    # traits_to_increment = ["one", "two", "three"]
    # traits_to_set_once  = ["one", "two", "three"]
    # traits_to_append    = ["one", "two", "three"]
    # traits_to_prepend   = ["one", "two", "three"]

    # prefer_anonymous_id_for_device_id {
    #   web = true
    # }

    # device_id_from_url_param {
    #   web = true
    # }

    # force_https {
    #   web = true
    # }

    # save_params_referrer_once_per_session {
    #   web = true
    # }

    # unset_params_referrer_on_new_session {
    #   web = true
    # }

    # batch_events {
    #   web = true
    # }

    # map_device_brand = true

    # event_upload_period_millis {
    #   web          = "1000"
    #   ios          = "1000"
    #   android      = "1000"
    #   react_native = "1000"
    # }

    # event_upload_threshold {
    #   web          = "1000"
    #   ios          = "1000"
    #   android      = "1000"
    #   react_native = "1000"
    # }

    # enable_location_listening {
    #   android      = true
    #   react_native = true
    # }

    # use_advertising_id_for_device_id {
    #   android      = true
    #   react_native = true
    # }

    # use_idfa_as_device_id {
    #   ios          = true
    #   react_native = true
    # }

    # use_native_sdk {
    #   web          = true
    #   ios          = true
    #   android      = true
    #   react_native = true
    # }

    # event_filtering {
    #   whitelist = ["one", "two", "three"]
    #   blacklist = ["one", "two", "three"]
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

    # residency_server = "EU"
  }
}
```

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

- `api_key` (String) Enter your Amplitude API key.

Optional:

- `api_secret` (String, Sensitive) Enter the Amplitude API Secret key required for user deletion.
- `batch_events` (Block List, Max: 1) If this setting is enabled, the events are batched together and uploaded by the Amplitude SDK. (see [below for nested schema](#nestedblock--config--batch_events))
- `consent_management` (Block List, Max: 1) Allows you to specify consent configuration data for multiple providers for each source type. (see [below for nested schema](#nestedblock--config--consent_management))
- `device_id_from_url_param` (Block List, Max: 1) If this setting is enabled, the Amplitude SDK will parse the URL parameter and set the device ID from `amp_device_id`. (see [below for nested schema](#nestedblock--config--device_id_from_url_param))
- `enable_location_listening` (Block List, Max: 1) Enable this setting to activate location listening. (see [below for nested schema](#nestedblock--config--enable_location_listening))
- `event_filtering` (Block List, Max: 1) This option allows you filter the events you want to send to Amplitude. (see [below for nested schema](#nestedblock--config--event_filtering))
- `event_upload_period_millis` (Block List, Max: 1) If the batch events settings is enabled, this is the amount of time that the SDK waits to upload the events. (see [below for nested schema](#nestedblock--config--event_upload_period_millis))
- `event_upload_threshold` (Block List, Max: 1) If the batch events settings is enabled, this is the minimum number of events to batch together by the Amplitude SDK. (see [below for nested schema](#nestedblock--config--event_upload_threshold))
- `force_https` (Block List, Max: 1) If this setting is enabled, the events will always be uploaded by the Amplitude SDK to the HTTPS endpoint, otherwise it will use the embedding site's protocol. (see [below for nested schema](#nestedblock--config--force_https))
- `group_type_trait` (String) RudderStack will use this value as `groupType` in the `group` calls.
- `group_value_trait` (String) RudderStack will use this value as `groupValue` in the `group` calls.
- `map_device_brand` (Boolean) Enable this setting for RudderStack to send the device brand information (`context.device.brand`) to Amplitude.
- `prefer_anonymous_id_for_device_id` (Block List, Max: 1) If this setting is enabled, the device ID will be set as the `anonymousId` generated by RudderStack SDK or by the `anonymousId` set via RudderStack's `setAnonymousId()` method. (see [below for nested schema](#nestedblock--config--prefer_anonymous_id_for_device_id))
- `residency_server` (String)
- `save_params_referrer_once_per_session` (Block List, Max: 1) If this setting is enabled, the corresponding tracking of `gclid`, referrer, UTM parameters will be done once per session. (see [below for nested schema](#nestedblock--config--save_params_referrer_once_per_session))
- `track_all_pages` (Boolean) If this setting is enabled, RudderStack sends an event named `Loaded a page` / `Loaded a Screen` to Amplitude.
- `track_categorized_pages` (Boolean) If this setting is enabled and if `category` is present in a `page` / `screen` call, then an event named `Viewed {category} page` / `Viewed {category} Screen` will be sent to Amplitude.
- `track_gclid` (Block List, Max: 1) If this setting is enabled, the Amplitude SDK will capture the `gclid` URL parameters along with the user's `initial_gclid` parameters. (see [below for nested schema](#nestedblock--config--track_gclid))
- `track_named_pages` (Boolean) If this setting is enabled and `name` is present in a `page` call, then an event named `Viewed {name} page` will be sent to Amplitude.
- `track_products_once` (Boolean) If this setting is enabled and if the event payload contains an array of products, then the event is tracked with the original event name and all the products as its property. Otherwise, each product is tracked with event as `Product purchased`.
- `track_referrer` (Block List, Max: 1) If this setting is enabled, the Amplitude SDK will capture the `referrer` and `referring_domain` for each session along with the user's `initial_referrer` and `initial_referring_domain`. (see [below for nested schema](#nestedblock--config--track_referrer))
- `track_revenue_per_product` (Boolean) If this setting is enabled and if the event payload contains multiple products, each product's revenue is tracked individually.
- `track_session_events` (Block List, Max: 1) Enable this setting to track the session events. (see [below for nested schema](#nestedblock--config--track_session_events))
- `track_utm_properties` (Block List, Max: 1) If this setting is enabled, the Amplitude SDK parses the UTM parameters in the query string or `_utmz` cookie and includes them as user properties in all uploaded events. (see [below for nested schema](#nestedblock--config--track_utm_properties))
- `traits_to_append` (List of String) If this setting is enabled, the value of the corresponding trait will be appended to the corresponding trait array at Amplitude.
- `traits_to_increment` (List of String) If this setting is enabled, the value of the corresponding trait will be incremented at Amplitude, with the value provided against the trait in an `identify` call.
- `traits_to_prepend` (List of String) If this setting is enabled, the value of the corresponding trait will be prepended to the corresponding trait array at Amplitude.
- `traits_to_set_once` (List of String) If this setting is enabled, the value of the corresponding trait will be set once at Amplitude with the value provided against the trait in an `identify` call.
- `unset_params_referrer_on_new_session` (Block List, Max: 1) If this setting is disabled, the existing `referrer` and `utm_parameter` values will be passed to each new session. If enabled, `referrer` and `utm_parameter` properties will be set to `null` upon instantiating a new session. (see [below for nested schema](#nestedblock--config--unset_params_referrer_on_new_session))
- `use_advertising_id_for_device_id` (Block List, Max: 1) Enable this setting to set the advertising ID as the device ID. (see [below for nested schema](#nestedblock--config--use_advertising_id_for_device_id))
- `use_idfa_as_device_id` (Block List, Max: 1) Enable this setting to set the IDFA as the device ID. (see [below for nested schema](#nestedblock--config--use_idfa_as_device_id))
- `use_native_sdk` (Block List, Max: 1) Enable this setting to send events to Amplitude via the device mode. (see [below for nested schema](#nestedblock--config--use_native_sdk))
- `version_name` (String) The value of this field is set as the `versionName` of the Amplitude SDK.

<a id="nestedblock--config--batch_events"></a>
### Nested Schema for `config.batch_events`

Optional:

- `web` (Boolean)


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



<a id="nestedblock--config--device_id_from_url_param"></a>
### Nested Schema for `config.device_id_from_url_param`

Optional:

- `web` (Boolean)


<a id="nestedblock--config--enable_location_listening"></a>
### Nested Schema for `config.enable_location_listening`

Optional:

- `android` (Boolean)
- `react_native` (Boolean)


<a id="nestedblock--config--event_filtering"></a>
### Nested Schema for `config.event_filtering`

Optional:

- `blacklist` (List of String) Enter the event names to be denylisted.
- `whitelist` (List of String) Enter the event names to be allowlisted.


<a id="nestedblock--config--event_upload_period_millis"></a>
### Nested Schema for `config.event_upload_period_millis`

Optional:

- `android` (String)
- `ios` (String)
- `react_native` (String)
- `web` (String)


<a id="nestedblock--config--event_upload_threshold"></a>
### Nested Schema for `config.event_upload_threshold`

Optional:

- `android` (String)
- `ios` (String)
- `react_native` (String)
- `web` (String)


<a id="nestedblock--config--force_https"></a>
### Nested Schema for `config.force_https`

Optional:

- `web` (Boolean)


<a id="nestedblock--config--prefer_anonymous_id_for_device_id"></a>
### Nested Schema for `config.prefer_anonymous_id_for_device_id`

Optional:

- `web` (Boolean)


<a id="nestedblock--config--save_params_referrer_once_per_session"></a>
### Nested Schema for `config.save_params_referrer_once_per_session`

Optional:

- `web` (Boolean)


<a id="nestedblock--config--track_gclid"></a>
### Nested Schema for `config.track_gclid`

Optional:

- `web` (Boolean)


<a id="nestedblock--config--track_referrer"></a>
### Nested Schema for `config.track_referrer`

Optional:

- `web` (Boolean)


<a id="nestedblock--config--track_session_events"></a>
### Nested Schema for `config.track_session_events`

Optional:

- `android` (Boolean)
- `ios` (Boolean)
- `react_native` (Boolean)


<a id="nestedblock--config--track_utm_properties"></a>
### Nested Schema for `config.track_utm_properties`

Optional:

- `web` (Boolean)


<a id="nestedblock--config--unset_params_referrer_on_new_session"></a>
### Nested Schema for `config.unset_params_referrer_on_new_session`

Optional:

- `web` (Boolean)


<a id="nestedblock--config--use_advertising_id_for_device_id"></a>
### Nested Schema for `config.use_advertising_id_for_device_id`

Optional:

- `android` (Boolean)
- `react_native` (Boolean)


<a id="nestedblock--config--use_idfa_as_device_id"></a>
### Nested Schema for `config.use_idfa_as_device_id`

Optional:

- `ios` (Boolean)
- `react_native` (Boolean)


<a id="nestedblock--config--use_native_sdk"></a>
### Nested Schema for `config.use_native_sdk`

Optional:

- `android` (Boolean)
- `ios` (Boolean)
- `react_native` (Boolean)
- `web` (Boolean)
