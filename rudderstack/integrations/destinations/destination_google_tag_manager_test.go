package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceGoogleTagManager(t *testing.T) {
	cmt.AssertDestination(t, "google_tag_manager", []c.TestConfig{
		{
			TerraformCreate: `
				container_id = "GTM-000000"
			`,
			APICreate: `{
				"containerID": "GTM-000000"
			}`,
			TerraformUpdate: `
				container_id = "GTM-000000"

				server_url = "https://example.com"
			
				use_native_sdk {
					web = true
				}
			
				event_filtering {
					blacklist = ["one", "two", "three"]
				}
			
				onetrust_cookie_categories {
					web = ["one", "two", "three"]
				}
			`,
			APIUpdate: `{
				"containerID": "GTM-000000",
				"serverUrl": "https://example.com",
				"blacklistedEvents": [
				  {
					"eventName": "one"
				  },
				  {
					"eventName": "two"
				  },
				  {
					"eventName": "three"
				  }
				],
				"eventFilteringOption": "blacklistedEvents",
				"useNativeSDK": {
				  "web": true
				},
				"oneTrustCookieCategories": {
					"web": [
						{ "oneTrustCookieCategory": "one" },
						{ "oneTrustCookieCategory": "two" },
						{ "oneTrustCookieCategory": "three" }
					]
				}
			}`,
		},
	})
}
