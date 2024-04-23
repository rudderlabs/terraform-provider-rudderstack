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
				onetrust_cookie_categories = ["one", "two", "three"]
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
				"oneTrustCookieCategories": [
					{ "oneTrustCookieCategory": "one" },
					{ "oneTrustCookieCategory": "two" },
					{ "oneTrustCookieCategory": "three" }
				]
			}`,
		},
	})
}
