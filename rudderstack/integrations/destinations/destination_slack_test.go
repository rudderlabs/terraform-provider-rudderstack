package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceSlack(t *testing.T) {
	cmt.AssertDestination(t, "slack", []c.TestConfig{
		{
			TerraformCreate: `
				webhook_url = "https://some-url"
			`,
			APICreate: `{
				"webhookUrl": "https://some-url"
			}`,
			TerraformUpdate: `
				webhook_url = "https://some-url"
				identify_template = "it"

				event_channel_settings = [
					{
						name    = "n1"
						channel = "c1"
						regex   = false
					}
				]
			  
				event_template_settings = [
					{
					  name     = "n2"
					  template = "t2"
					  regex    = true
					}
				]
			  
				whitelisted_trait_settings = ["one", "two", "three"]
			`,
			APIUpdate: `{
				"webhookUrl": "https://some-url",
				"identifyTemplate": "it",
				"eventChannelSettings": [
					{
						"eventName": "n1",
						"eventChannel": "c1",
						"eventRegex": false
					}
				],
				"eventTemplateSettings": [
					{
						"eventName": "n2",
						"eventTemplate": "t2",
						"eventRegex": true
					}
				],
				"whitelistedTraitsSettings": [
					{ "trait": "one" },
					{ "trait": "two" },
					{ "trait": "three" }
				]
			}`,
		},
	})
}
