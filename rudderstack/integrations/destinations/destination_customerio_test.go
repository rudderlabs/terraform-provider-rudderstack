package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceCustomerIO(t *testing.T) {
	cmt.AssertDestination(t, "customerio", []c.TestConfig{
		{
			TerraformCreate: `
				site_id = "cd820c1b31d8f2696f3b"
				api_key = "cg044d23bc1beb3031c5"
				datacenter = "US"

				use_native_sdk {
					web = true
				}
			`,
			APICreate: `{
				"siteID": "cd820c1b31d8f2696f3b",
				"apiKey": "cg044d23bc1beb3031c5",
				"datacenter": "US",
				"useNativeSDK": {
					"web": true
				}
			}`,
			TerraformUpdate: `
				site_id = "cd820c1b31d8f2696f3b"
				api_key = "cg044d23bc1beb3031c5"
				datacenter = "EU"
				device_token_event_name = "name"

				event_filtering {
					blacklist = [ "one", "two", "three" ]
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
			`,
			APIUpdate: `{
				"siteID": "cd820c1b31d8f2696f3b",
				"apiKey": "cg044d23bc1beb3031c5",
				datacenter = "EU"
				"deviceTokenEventName": "name",
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
