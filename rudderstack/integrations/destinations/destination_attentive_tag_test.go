package destinations_test

import (
	"testing"

	acc "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/acc"
	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

var attentiveTagTestConfigs = []c.TestConfig{
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
				"connectionMode": {
					"web": "cloud",
					"android": "cloud",
					"ios": "cloud",
					"unity": "cloud",
					"reactnative": "cloud",
					"flutter": "cloud",
					"cordova": "cloud",
					"amp": "cloud",
					"cloud": "cloud",
					"warehouse": "cloud",
					"shopify": "cloud"
				},
				"signUpSourceId": "123456",
				"enableNewIdentifyFlow": true
			}`,
	},
}

func TestDestinationResourceAttentiveTag(t *testing.T) {
	cmt.AssertDestination(t, "attentive_tag", attentiveTagTestConfigs)
}

func TestAccDestinationAttentiveTag(t *testing.T) {
	acc.AccAssertDestination(t, "attentive_tag", attentiveTagTestConfigs)
}
