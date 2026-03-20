package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceHttpWebhook(t *testing.T) {
	cmt.AssertDestination(t, "http_webhook", []c.TestConfig{
		{
			TerraformCreate: `
				webhook_url = "https://example.com/webhook"
			`,
			APICreate: `{
				"webhookUrl": "https://example.com/webhook"
			}`,
			TerraformUpdate: `
				webhook_url    = "https://example.com/webhook"
				webhook_method = "POST"
				headers = [
					{
						from = "Content-Type"
						to   = "application/json"
					},
					{
						from = "X-Custom-Header"
						to   = "my-value"
					}
				]
			`,
			APIUpdate: `{
				"webhookUrl": "https://example.com/webhook",
				"webhookMethod": "POST",
				"headers": [
					{"from": "Content-Type", "to": "application/json"},
					{"from": "X-Custom-Header", "to": "my-value"}
				]
			}`,
		},
	})
}
