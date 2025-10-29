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
			`,
			APICreate: `{
				"apiKey": "key"
			}`,
			TerraformUpdate: `
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
