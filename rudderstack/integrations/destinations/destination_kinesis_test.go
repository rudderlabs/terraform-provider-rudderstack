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
				role_based_authentication {
                  i_am_role_arn = "arn"
				}
			`,
			APICreate: `{
				"region":"usa-east",
				"stream":"test",
				"roleBasedAuth":true,
				"iamRoleARN":"arn"
			}`,
			TerraformUpdate: `
				region = "usa-east"
				stream = "test"
				key_based_authentication {
			     access_key = "key"
				 access_key_id = "id"
				}
				use_message_id = true
				onetrust_cookie_categories = ["one", "two", "three"]
			`,
			APIUpdate: `{
				"region": "usa-east",
				"stream": "test",
				"roleBasedAuth": false,
				"accessKeyID": "id",
				"accessKey": "key",
				"useMessageId": true,
				"oneTrustCookieCategories": [
					{ "oneTrustCookieCategory": "one" },
					{ "oneTrustCookieCategory": "two" },
					{ "oneTrustCookieCategory": "three" }
				]
			}`,
		},
	})
}
