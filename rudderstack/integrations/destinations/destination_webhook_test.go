package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceWebhook(t *testing.T) {
	cmt.AssertDestination(t, "webhook", []c.TestConfig{
		{
			TerraformCreate: `
				webhook_url = "https://example.com/some/path?query=a"
				webhook_method = "GET"
			`,
			APICreate: `{
				"webhookUrl": "https://example.com/some/path?query=a",
				"webhookMethod": "GET"
			}`,
			TerraformUpdate: `
				webhook_url = "https://example.com/some/path?query=a"
				webhook_method = "GET"
				headers = [
					{
						from = "a1"
						to   = "b1"
					},
					{
						from = "a2"
						to   = "b2"
					}
				]
				onetrust_cookie_categories = ["one", "two", "three"]
			`,
			APIUpdate: `{
				"webhookUrl": "https://example.com/some/path?query=a",
				"webhookMethod": "GET",
				"headers": [
					{
						"from": "a1",
						"to": "b1"
					},
					{
						"from": "a2",
						"to": "b2"
					}
				],
				"oneTrustCookieCategories": [
					{ "oneTrustCookieCategory": "one" },
					{ "oneTrustCookieCategory": "two" },
					{ "oneTrustCookieCategory": "three" }
				]
			}`,
		},
	})
}
