package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceGCS(t *testing.T) {
	cmt.AssertDestination(t, "gcs", []c.TestConfig{
		{
			TerraformCreate: `
				bucket_name = "bucket"
			`,
			APICreate: `{
				"bucketName": "bucket"
			}`,
			TerraformUpdate: `
				bucket_name = "bucket"
				prefix        = "prefix"
				credentials   = "..."
				onetrust_cookie_categories = ["one", "two", "three"]
			`,
			APIUpdate: `{
				"bucketName": "bucket",
				"prefix": "prefix",
				"credentials": "...",
				"oneTrustCookieCategories": [
					{ "oneTrustCookieCategory": "one" },
					{ "oneTrustCookieCategory": "two" },
					{ "oneTrustCookieCategory": "three" }
				]
			}`,
		},
	})
}
