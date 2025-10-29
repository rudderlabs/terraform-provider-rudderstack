package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceAttentiveTag(t *testing.T) {
	cmt.AssertDestination(t, "attentive_tag", []c.TestConfig{
		{
			TerraformCreate: `
				api_key = "key"
				connection_mode {
					web = "cloud"
				}
			`,
			APICreate: `{
				"apiKey": "key",
				"connectionMode": {
					"web": "cloud"
				}
			}`,
			TerraformUpdate: `
				connection_mode {
					web = "cloud"
					android = "cloud"
					ios = "cloud"
					unity = "cloud"
					reactnative = "cloud"
					flutter = "cloud"
					cordova = "cloud"
					amp = "cloud"
					cloud = "cloud"
					warehouse = "cloud"
					shopify = "cloud"
				}
				api_key = "key"
				sign_up_source_id = "123456"
				enable_new_identify_flow = true
			`,
			APIUpdate: `{
				"apiKey": "key",
				"signUpSourceId": "123456",
				"enableNewIdentifyFlow": true
			}`,
		},
	})
}
