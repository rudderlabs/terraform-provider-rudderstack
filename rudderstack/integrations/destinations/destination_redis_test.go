package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceRedis(t *testing.T) {
	cmt.AssertDestination(t, "redis", []c.TestConfig{
		{
			TerraformCreate: `
				address = "https://some-url"
			`,
			APICreate: `{
				"address": "https://some-url",
				"clusterMode": false,
				"secure": false,
				"skipVerify": false
			}`,
			TerraformUpdate: `
				address = "1.2.3.4"

				password      = "..."
				cluster_mode  = true
				secure        = true
				prefix        = "..."
				database      = "..."
				ca_certificate = "..."
				skip_verify   = true
				onetrust_cookie_categories = ["one", "two", "three"]
			`,
			APIUpdate: `{
				"address": "1.2.3.4",
				"password": "...",
				"clusterMode": true,
				"secure": true,
				"prefix": "...",
				"database": "...",
				"caCertificate": "...",
				"skipVerify": true,
				"oneTrustCookieCategories": [
					{ "oneTrustCookieCategory": "one" },
					{ "oneTrustCookieCategory": "two" },
					{ "oneTrustCookieCategory": "three" }
				]
			}`,
		},
	})
}
