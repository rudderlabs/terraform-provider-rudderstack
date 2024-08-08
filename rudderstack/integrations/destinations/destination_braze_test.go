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
			onetrust_cookie_categories = ["one", "two", "three"]
			`,
			APIUpdate: `{
				"connectionMode": {
					"web": "cloud",
					"ios": "cloud"
				},
				"dataCenter": "US-03",
				"restApiKey": "updated_rest_api_pass",
				"oneTrustCookieCategories": [
					{ "oneTrustCookieCategory": "one" },
					{ "oneTrustCookieCategory": "two" },
					{ "oneTrustCookieCategory": "three" }
				]
			}`,
		},
	})
}
