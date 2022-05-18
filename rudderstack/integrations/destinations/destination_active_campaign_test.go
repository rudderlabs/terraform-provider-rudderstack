package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceActiveCampaign(t *testing.T) {
	cmt.AssertDestination(t, "active_campaign", []c.TestConfig{
		{
			TerraformCreate: `
				api_url = "https://some-url"
			`,
			APICreate: `{
				"apiUrl": "https://some-url"
			}`,
			TerraformUpdate: `
				api_url   = "https://some-url"
				api_key   = "api-key"
				actid     = "actid"
				event_key = "event-key"
			`,
			APIUpdate: `{
				"apiUrl": "https://some-url",
				"apiKey": "api-key",
				"actid": "actid",
				"eventKey": "event-key"
			}`,
		},
	})
}
