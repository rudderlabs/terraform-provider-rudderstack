package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceGoogleAnalytics4(t *testing.T) {
	cmt.AssertDestination(t, "google_analytics4", []c.TestConfig{
		{
			TerraformCreate: `
				api_secret      = "..."
				measurement_id  = "G-000000"
			`,
			APICreate: `{
				"apiSecret": "...",
				"measurementId": "G-000000"
			}`,
			TerraformUpdate: `
				api_secret = "..."

				types_of_client = "gtag"
				measurement_id  = "G-000000"
				firebase_app_id = "..."
			
				block_page_view_event   = true
				extend_page_view_params = true
				send_user_id            = true
			
				use_native_sdk {
					web = true
				}
			
				event_filtering {
					blacklist = ["one", "two", "three"]
				}
			
				onetrust_cookie_categories = ["one", "two", "three"]
			`,
			APIUpdate: `{
				"apiSecret": "...",
				"typesOfClient": "gtag",
				"measurementId": "G-000000",
				"firebaseAppId": "...",
				"blockPageViewEvent": true,
				"extendPageViewParams": true,
				"sendUserId": true,
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
				"oneTrustCookieCategories": [
					{ "oneTrustCookieCategory": "one" },
					{ "oneTrustCookieCategory": "two" },
					{ "oneTrustCookieCategory": "three" }
				]
			}`,
		},
	})
}
