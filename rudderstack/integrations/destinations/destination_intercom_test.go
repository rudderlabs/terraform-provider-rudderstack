package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceIntercom(t *testing.T) {
	cmt.AssertDestination(t, "intercom", []c.TestConfig{
		{
			TerraformCreate: `
				app_id = "app-id"
				api_key = "api-key"
				use_native_sdk {
					web = true
					ios = true
				}
			`,
			APICreate: `{
				"appId": "app-id",
				"apiKey": "api-key",
				"useNativeSDK": {
					"web": true,
					"ios": true
				},				
				"updateLastRequestAt": true
			}`,
			TerraformUpdate: `
				app_id = "app-id"
				api_key = "api-key"
				use_native_sdk {
					android = true
				}
				event_filtering {
					blacklist = [ "one", "two", "three" ]
				}
				collect_context = true
				send_anonymous_id = true
				update_last_request_at = false
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
				mobile_api_key_android = "and-key"
				mobile_api_key_ios = "ios-key"
			`,
			APIUpdate: `{
				"appId": "app-id",
				"apiKey": "api-key",
				"useNativeSDK": {
					"android": true
				},				
				"mobileApiKeyAndroid": {
				  "android": "and-key"
				},
				"mobileApiKeyIOS": {
				  "ios": "ios-key"
				},
				"collectContext":true,
				"sendAnonymousId":true,
				"eventFilteringOption": "blacklistedEvents",
				"blacklistedEvents": [{
					"eventName": "one"
				}, {
					"eventName": "two"
				}, {
					"eventName": "three"
				}],
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
				}
			}`,
		},
	})
}
