package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceGoogleAds(t *testing.T) {
	cmt.AssertDestination(t, "google_ads", []c.TestConfig{
		{
			TerraformCreate: `
				conversion_id = "AW-00000000"
			`,
			APICreate: `{
				"conversionID": "AW-00000000"
			}`,
			TerraformUpdate: `
				conversion_id = "AW-00000000"

				default_page_conversion = "..."
			
				page_load_conversions = [
					{
						"label" = "..."
						"name"  = "..."
					}
				]
			
				click_event_conversions = [
					{
						"label" = "..."
						"name"  = "..."
					}
				]
			
				dynamic_remarketing {
					web = true
				}
			
				conversion_linker          = true
				send_page_view             = true
				disable_ad_personalization = true
			
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
				"conversionID": "AW-00000000",
				"pageLoadConversions": [
				  {
					"conversionLabel": "...",
					"name": "..."
				  }
				],
				"clickEventConversions": [
				  {
					"conversionLabel": "...",
					"name": "..."
				  }
				],
				"defaultPageConversion": "...",
				"dynamicRemarketing": {
				  "web": true
				},
				"conversionLinker": true,
				"sendPageView": true,
				"disableAdPersonalization": true,
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
					{
					  "oneTrustCookieCategory": "one"
					},
					{
					  "oneTrustCookieCategory": "two"
					},
					{
					  "oneTrustCookieCategory": "three"
					}
				  ]
				}
			}`,
		},
	})
}
