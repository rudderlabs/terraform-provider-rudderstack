package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceBraze(t *testing.T) {
	cmt.AssertDestination(t, "braze", []c.TestConfig{
		{
			TerraformCreate: `
			connection_mode {
				web = "cloud"
				ios = "cloud"
			}
			data_center = "US-03"
			rest_api_key = "rest_api_pass"
			`,
			APICreate: `{
				"connectionMode": {
					"web": "cloud",
					"ios": "cloud"
				},
				"dataCenter": "US-03",
				"restApiKey": "rest_api_pass"
			}`,
			TerraformUpdate: `
			connection_mode {
				web = "cloud"
				ios = "cloud"
			}
			data_center = "US-03"
			rest_api_key = "updated_rest_api_pass"
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
				"connectionMode": {
					"web": "cloud",
					"ios": "cloud"
				},
				"dataCenter": "US-03",
				"restApiKey": "updated_rest_api_pass",
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
