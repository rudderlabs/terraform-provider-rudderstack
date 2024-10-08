---
page_title: "rudderstack_destination_intercom Resource - terraform-provider-rudderstack"
subcategory: ""
description: |-
  
---

# rudderstack_destination_intercom (Resource)

This resource represents an Intercom destination. For more information check 
https://www.rudderstack.com/docs/destinations/streaming-destinations/intercom/

## Example Usage

```terraform
resource "rudderstack_destination_intercom" "example" {
  name = "my-intercom-tf"

  config {
    app_id = "app-id"
    api_key = "api-key"
#    use_native_sdk {
#      web = true
#      ios = true
#    }
#    event_filtering {
#      blacklist = [ "one", "two", "three" ]
#    }
#    mobile_api_key_android = "and-key"
#    mobile_api_key_ios = "ios-key"
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
  }
}
```

> **:warning: Breaking Change**
> 
> Note that from the provider versions 3.0.0 and above, `onetrust_cookie_categories` property is replaced with `consent_management` that supports multiple consent management providers. Please refer to the example above.

> **:warning: Breaking Change**
> 
> Note that from the provider versions 2.0.0 and above, the schema of `onetrust_cookie_categories` property has been changed. Please refer to the example above.

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

- `api_key` (String) Enter your Intercom access token.
- `app_id` (String) Enter your app ID.

Optional:

- `collect_context` (Boolean) Enable this setting to include the user context along with your identify calls.
- `consent_management` (Block List, Max: 1) Allows you to specify consent configuration data for multiple providers for each source type. (see [below for nested schema](#nestedblock--config--consent_management))
- `event_filtering` (Block List, Max: 1) Use this setting to determine which events should be blocked or allowed to flow through. (see [below for nested schema](#nestedblock--config--event_filtering))
- `mobile_api_key_android` (String) Enter the Android API Key.
- `mobile_api_key_ios` (String) Enter the iOS API Key.
- `send_anonymous_id` (Boolean) Enable this setting to send anonymousId as the secondary userId.
- `update_last_request_at` (Boolean) Enable this setting to send the last seen information with the current time.
- `use_native_sdk` (Block List, Max: 1) Enable this setting to send the events through device mode, that is, using the native SDK. (see [below for nested schema](#nestedblock--config--use_native_sdk))

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


<a id="nestedblock--config--use_native_sdk"></a>
### Nested Schema for `config.use_native_sdk`

Optional:

- `android` (Boolean)
- `ios` (Boolean)
- `web` (Boolean)
