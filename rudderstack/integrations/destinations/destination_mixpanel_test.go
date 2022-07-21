package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceMixpanel(t *testing.T) {
	cmt.AssertDestination(t, "mixpanel", []c.TestConfig{
		{
			TerraformCreate: `
				token = "..."
				data_residency = "us"
				persistence = "none"
				consolidated_page_calls = false
			`,
			APICreate: `{
				"token": "...",
				"dataResidency": "us",
				"consolidatedPageCalls": false,
				"persistence": "none"
			}`,
			TerraformUpdate: `
				token = "..."
				data_residency = "eu"
				persistence = "localStorage"
				api_secret = "..."
				people = true
				set_all_traits_by_default = true
				consolidated_page_calls = true
				track_categorized_pages = true
				track_named_pages = true
				source_name = "my-mixpanel"
				cross_subdomain_cookie = true
				secure_cookie = true
				super_properties = ["one","two","three"]
				people_properties = ["one","two","three"]
				event_increments = ["one","two","three"]
				prop_increments = ["one","two","three"]
				group_key_settings = ["one","two","three"]
				
				use_native_sdk {
					web = true
				}
		
				event_filtering {
					whitelist = ["one", "two", "three"]
				}
		
				onetrust_cookie_categories {
					web = ["one", "two", "three"]
				}
				use_new_mapping = true
			`,
			APIUpdate: `
			{
				"token": "...",
				"dataResidency": "eu",
				"persistence": "localStorage",
				"apiSecret": "...",
				"people": true,
				"setAllTraitsByDefault": true,
				"consolidatedPageCalls": true,
				"trackCategorizedPages": true,
				"trackNamedPages": true,
				"sourceName": "my-mixpanel",
				"crossSubdomainCookie": true,
				"secureCookie": true,
				"superProperties": [{
						"property": "one"
					},
					{
						"property": "two"
					},
					{
						"property": "three"
					}
				],
				"peopleProperties": [{
						"property": "one"
					},
					{
						"property": "two"
					},
					{
						"property": "three"
					}
				],
				"eventIncrements": [{
						"property": "one"
					},
					{
						"property": "two"
					},
					{
						"property": "three"
					}
				],
				"propIncrements": [{
						"property": "one"
					},
					{
						"property": "two"
					},
					{
						"property": "three"
					}
				],
				"groupKeySettings": [{
						"groupKey": "one"
					},
					{
						"groupKey": "two"
					},
					{
						"groupKey": "three"
					}
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
				"oneTrustCookieCategories": {
					"web": [{
							"oneTrustCookieCategory": "one"
						},
						{
							"oneTrustCookieCategory": "two"
						},
						{
							"oneTrustCookieCategory": "three"
						}
					]
				},
				"useNewMapping": true
			}			
			`,
		},
	})
}
