package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceStatSig(t *testing.T) {
	cmt.AssertDestination(t, "statsig", []c.TestConfig{
		{
			TerraformCreate: `
				secret_key = "key"
				connection_mode {
				 web = "cloud"
				 ios = "cloud"
				}
			`,
			APICreate: `{
				"secretKey": "key",
				"connectionMode": {
					"web": "cloud",
					"ios": "cloud"
				}
			}`,
			TerraformUpdate: `
				secret_key = "key"
				connection_mode {
				 web = "cloud"
				 ios = "cloud"
				 amp = "cloud"
				 react_native = "cloud"
				}
				onetrust_cookie_categories = ["one", "two", "three"]
			`,
			APIUpdate: `{
				"secretKey": "key",
				"connectionMode": {
					"web": "cloud",
					"ios": "cloud",
					"amp": "cloud",
					"reactnative": "cloud"
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
