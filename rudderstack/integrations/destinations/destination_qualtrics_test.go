package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceQualtrics(t *testing.T) {
	cmt.AssertDestination(t, "qualtrics", []c.TestConfig{
		{
			TerraformCreate: `
				project_id = "p-id"
				brand_id = "b-id"
				use_native_sdk {
					web = true
				}
			`,
			APICreate: `{
				"projectId": "p-id",
				"brandId": "b-id",
				"useNativeSDK": {
					"web": true
				}
			}`,
			TerraformUpdate: `
				project_id = "p-id"
				brand_id = "b-id"

				enable_generic_page_title = true
				use_native_sdk {
					ios = true
				}
           
				event_filtering {
					blacklist = [ "one", "two", "three" ]
				}
				onetrust_cookie_categories {
					web = ["one", "two", "three"]
					android = ["one", "two", "three"]
					ios = ["one", "two", "three"]
				}
			`,
			APIUpdate: `{
				"projectId": "p-id",
				"brandId": "b-id",
				"eventFilteringOption": "blacklistedEvents",
				"blacklistedEvents": [
					{"eventName": "one"},
					{"eventName": "two"},
					{"eventName": "three"}
				],
				"useNativeSDK":{"ios":true},
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
					]
				},
				"enableGenericPageTitle":{"web":true}
			}`,
		},
	})
}
