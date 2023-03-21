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
				datacenter_eu = true

				use_native_sdk {
					web = true
				}

				
				onetrust_cookie_categories = ["one", "two", "three"]
			`,
			APICreate: `{
				"siteID": "cd820c1b31d8f2696f3b",
				"apiKey": "cg044d23bc1beb3031c5",
				"datacenterEU": true,
				"useNativeSDK": {
					"web": true
				},
				"oneTrustCookieCategories": [
					{ "oneTrustCookieCategory": "one" },
					{ "oneTrustCookieCategory": "two" },
					{ "oneTrustCookieCategory": "three" }
				]
			}`,
			TerraformUpdate: `
				site_id = "cd820c1b31d8f2696f3b"
				api_key = "cg044d23bc1beb3031c5"
				datacenter_eu = false
				device_token_event_name = "name"

				event_filtering {
					blacklist = [ "one", "two", "three" ]
				}
			`,
			APIUpdate: `{
				"siteID": "cd820c1b31d8f2696f3b",
				"apiKey": "cg044d23bc1beb3031c5",
				"deviceTokenEventName": "name",
				"eventFilteringOption": "blacklistedEvents",
				"blacklistedEvents": [{
					"eventName": "one"
				}, {
					"eventName": "two"
				}, {
					"eventName": "three"
				}]
			}`,
		},
	})
}
