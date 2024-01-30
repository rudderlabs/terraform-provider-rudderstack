package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceSalesforce(t *testing.T) {
	cmt.AssertDestination(t, "salesforce", []c.TestConfig{
		{
			TerraformCreate: `
				user_name = "user"
				password = "pwd"
				initial_access_token = "token"
			`,
			APICreate: `{
				"userName": "user",
				"password": "pwd",
				"initialAccessToken": "token",
				"mapProperties":true
			}`,
			TerraformUpdate: `
				user_name = "bucket"
				password = "pwd"
				initial_access_token = "token"
				map_properties = true
				use_contact_id = true
				onetrust_cookie_categories = ["one", "two", "three"]
			`,
			APIUpdate: `{
				"userName": "bucket",
				"password": "pwd",
				"initialAccessToken": "token",
				"useContactId": true,
				"mapProperties": true,
				"oneTrustCookieCategories": [
					{ "oneTrustCookieCategory": "one" },
					{ "oneTrustCookieCategory": "two" },
					{ "oneTrustCookieCategory": "three" }
				]
			}`,
		},
	})
}
