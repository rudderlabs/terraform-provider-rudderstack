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
			connection_mode = "cloud",
			data_center = "US-01",
			rest_api_key = "rest_api_pass"
			`,
			APICreate: `{
			connection_mode = "cloud",
			data_center = "US-01",
			rest_api_key = "rest_api_pass",
			enable_subscription_group_in_group_call = true,
			track_anonymous_user = true
			}`,
			TerraformUpdate: `
				connection_mode = "cloud",
				data_center = "US-03",
				rest_api_key = "updated_rest_api_pass"
				onetrust_cookie_categories = ["one", "two", "three"]
			`,
			APIUpdate: `{
				"connection_mode = "cloud",
				data_center = "US-03",
				rest_api_key = "updated_rest_api_pass"
				"oneTrustCookieCategories": [
					{ "oneTrustCookieCategory": "one" },
					{ "oneTrustCookieCategory": "two" },
					{ "oneTrustCookieCategory": "three" }
				]
			}`,
		},
	})
}
