package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceLinkedinInsightTag(t *testing.T) {
	cmt.AssertDestination(t, "LINKEDIN_INSIGHT_TAG", []c.TestConfig{
		{
			TerraformCreate: `
				partner_id = "p-id"
			`,
			APICreate: `{
				"partnerId": "p-id"
			}`,
			TerraformUpdate: `
				partner_id = "p-id"
				event_to_conversion_id_map = [
				{
					from = "a1"
					to   = "b1"
				}, 
				{
					from = "a2"
					to   = "b2"
				}]
				use_native_sdk {
					web = true
				}
				event_filtering {
					whitelist = ["one", "two", "three"]
				}
				onetrust_cookie_categories = ["one", "two", "three"]
			`,
			APIUpdate: `
			{
				"partnerId": "p-id",
				"eventToConversionIdMap": [
				  { "from": "a1", "to": "b1" },
				  { "from": "a2", "to": "b2" }
				],
				"useNativeSDK": {
					"web": true
				},
				"eventFilteringOption": "whitelistedEvents",
				"whitelistedEvents": [{
						"eventName": "one"
					},
					{
						"eventName": "two"
					},
					{
						"eventName": "three"
					}
				],
				"oneTrustCookieCategories": [
					{ "oneTrustCookieCategory": "one" },
					{ "oneTrustCookieCategory": "two" },
					{ "oneTrustCookieCategory": "three" }
				]
			}			
			`,
		},
	})
}
