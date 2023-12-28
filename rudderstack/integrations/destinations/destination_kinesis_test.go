package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceKinesis(t *testing.T) {
	cmt.AssertDestination(t, "kinesis", []c.TestConfig{
		{
			TerraformCreate: `
				region = "usa-east"
				stream = "test"
			`,
			APICreate: `{
				"region": "usa-east",
				"stream": "test"
			}`,
			TerraformUpdate: `
				region = "usa-east"
				stream = "test"
				access_key_id = ""
				access_key    = ""
				i_am_role_arn = "arn"
				role_based_auth = true
				use_message_id = false
				onetrust_cookie_categories = ["one", "two", "three"]
			`,
			APIUpdate: `{
				"region": "usa-east"
				"stream": "test"
				"accessKeyID": "",
				"accessKey": "",
				"roleBasedAuth": true,
				"iamRoleARN": "arn",
				"useMessageId": false,
				"oneTrustCookieCategories": [
					{ "oneTrustCookieCategory": "one" },
					{ "oneTrustCookieCategory": "two" },
					{ "oneTrustCookieCategory": "three" }
				]
			}`,
		},
	})
}
