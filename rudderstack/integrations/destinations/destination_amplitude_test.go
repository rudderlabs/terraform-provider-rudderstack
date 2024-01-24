package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceAmplitude(t *testing.T) {
	cmt.AssertDestination(t, "amplitude", []c.TestConfig{
		{
			TerraformCreate: `
				api_key = "123abc"
				api_secret = "abc123"
			`,
			APICreate: `{
				"apiKey": "123abc",
				"apiSecret": "abc123"
			}`,
			TerraformUpdate: `
				api_key = "123abc"
				api_secret = "abc123"

				group_type_trait  = "type"
				group_value_trait = "value"
			
				track_all_pages           = true
				track_categorized_pages   = true
				track_named_pages         = true
				track_products_once       = true
				track_revenue_per_product = true
			
				track_gclid {
				  web = true
				}
			
				track_referrer {
				  web = true
				}
			
				track_utm_properties {
				  web = true
				}
			
				track_session_events {
				  android      = true
				  ios          = true
				  react_native = true
				}
			
				version_name = "name"
			
				traits_to_increment = ["one", "two", "three"]
				traits_to_set_once  = ["one", "two", "three"]
				traits_to_append    = ["one", "two", "three"]
				traits_to_prepend   = ["one", "two", "three"]
			
				prefer_anonymous_id_for_device_id {
				  web = true
				}
			
				device_id_from_url_param {
				  web = true
				}
			
				force_https {
				  web = true
				}
			
				save_params_referrer_once_per_session {
				  web = true
				}
			
				unset_params_referrer_on_new_session {
				  web = true
				}
			
				batch_events {
				  web = true
				}
							
				map_device_brand = true
			
				event_upload_period_millis {
				  web          = "1000"
				  ios          = "1000"
				  android      = "1000"
				  react_native = "1000"
				}
			
				event_upload_threshold {
				  web          = "1000"
				  ios          = "1000"
				  android      = "1000"
				  react_native = "1000"
				}
			
				enable_location_listening {
				  android      = true
				  react_native = true
				}
			
				use_advertising_id_for_device_id {
				  android      = true
				  react_native = true
				}
			
				use_idfa_as_device_id {
				  ios          = true
				  react_native = true
				}
			
				use_native_sdk {
				  web          = true
				  ios          = true
				  android      = true
				  react_native = true
				}
			
				event_filtering {
				  blacklist = ["one", "two", "three"]
				}
			
				onetrust_cookie_categories {
					web = ["one", "two", "three"]
					android = ["one", "two", "three"]
					ios = ["one", "two", "three"]
					unity = ["one", "two", "three"]
					reactnative = ["one", "two", "three"]
					flutter = ["one", "two", "three"]
					cordova = ["one", "two", "three"]
					amp = ["one", "two", "three"]
					cloud = ["one", "two", "three"]
					warehouse = ["one", "two", "three"]
					shopify = ["one", "two", "three"]
				}

				residency_server = "EU"
			`,
			APIUpdate: `{
				"apiKey": "123abc",
				"apiSecret": "abc123",
				"groupTypeTrait": "type",
				"groupValueTrait": "value",
				"trackAllPages": true,
				"trackCategorizedPages": true,
				"trackNamedPages": true,
				"trackProductsOnce": true,
				"trackRevenuePerProduct": true,
				"versionName": "name",
				"traitsToIncrement": [
				  { "traits": "one" },
				  { "traits": "two" },
				  { "traits": "three" }
				],
				"traitsToSetOnce": [
				  { "traits": "one" },
				  { "traits": "two" },
				  { "traits": "three" }
				],
				"traitsToAppend": [
				  { "traits": "one" },
				  { "traits": "two" },
				  { "traits": "three" }
				],
				"traitsToPrepend": [
				  { "traits": "one" },
				  { "traits": "two" },
				  { "traits": "three" }
				],
				"useNativeSDK": {
				  "web": true,
				  "ios": true,
				  "android": true,
				  "reactnative": true
				},
				"preferAnonymousIdForDeviceId": { "web": true },
				"deviceIdFromUrlParam": { "web": true },
				"forceHttps": { "web": true },
				"trackGclid": { "web": true },
				"trackReferrer": { "web": true },
				"saveParamsReferrerOncePerSession": { "web": true },
				"trackUtmProperties": { "web": true },
				"unsetParamsReferrerOnNewSession": { "web": true },
				"batchEvents": { "web": true },
				"eventFilteringOption": "blacklistedEvents",
				"blacklistedEvents": [
				  { "eventName": "one" },
				  { "eventName": "two" },
				  { "eventName": "three" }
				],
				"eventUploadPeriodMillis": {
				  "web": "1000",
				  "android": "1000",
				  "ios": "1000",
				  "reactnative": "1000"
				},
				"eventUploadThreshold": {
				  "web": "1000",
				  "android": "1000",
				  "ios": "1000",
				  "reactnative": "1000"
				},
				"mapDeviceBrand": true,
				"enableLocationListening": { "android": true, "reactnative": true },
				"trackSessionEvents": { "android": true, "ios": true, "reactnative": true },
				"useAdvertisingIdForDeviceId": { "android": true, "reactnative": true },
				"useIdfaAsDeviceId": { "ios": true, "reactnative": true },
				"oneTrustCookieCategories": {
					"web": [
						{ "oneTrustCookieCategory": "one" },
						{ "oneTrustCookieCategory": "two" },
						{ "oneTrustCookieCategory": "three" }
					],
					"android": [
						{ "oneTrustCookieCategory": "one" },
						{ "oneTrustCookieCategory": "two" },
						{ "oneTrustCookieCategory": "three" }
					],
					"ios": [
						{ "oneTrustCookieCategory": "one" },
						{ "oneTrustCookieCategory": "two" },
						{ "oneTrustCookieCategory": "three" }
					],
					"unity": [
						{ "oneTrustCookieCategory": "one" },
						{ "oneTrustCookieCategory": "two" },
						{ "oneTrustCookieCategory": "three" }
					],
					"reactnative": [
						{ "oneTrustCookieCategory": "one" },
						{ "oneTrustCookieCategory": "two" },
						{ "oneTrustCookieCategory": "three" }
					],
					"flutter": [
						{ "oneTrustCookieCategory": "one" },
						{ "oneTrustCookieCategory": "two" },
						{ "oneTrustCookieCategory": "three" }
					],
					"cordova": [
						{ "oneTrustCookieCategory": "one" },
						{ "oneTrustCookieCategory": "two" },
						{ "oneTrustCookieCategory": "three" }
					],
					"amp": [
						{ "oneTrustCookieCategory": "one" },
						{ "oneTrustCookieCategory": "two" },
						{ "oneTrustCookieCategory": "three" }
					],
					"cloud": [
						{ "oneTrustCookieCategory": "one" },
						{ "oneTrustCookieCategory": "two" },
						{ "oneTrustCookieCategory": "three" }
					],
					"warehouse": [
						{ "oneTrustCookieCategory": "one" },
						{ "oneTrustCookieCategory": "two" },
						{ "oneTrustCookieCategory": "three" }
					],
					"shopify": [
						{ "oneTrustCookieCategory": "one" },
						{ "oneTrustCookieCategory": "two" },
						{ "oneTrustCookieCategory": "three" }
					]
				},
				"residencyServer": "EU"
			}`,
		},
	})
}
