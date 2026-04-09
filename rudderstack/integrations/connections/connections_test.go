package connections_test

import (
	"testing"

	acc "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/acc"
)

func TestAccConnectionJavascriptToWebhook(t *testing.T) {
	acc.AccAssertConnection(t, acc.ConnectionTestConfig{
		Source:      "javascript",
		Destination: "webhook",
		DestConfig: `
			webhook_url    = "https://example.com/test"
			webhook_method = "POST"
		`,
	})
}

func TestAccConnectionHTTPToWebhook(t *testing.T) {
	acc.AccAssertConnection(t, acc.ConnectionTestConfig{
		Source:      "http",
		Destination: "webhook",
		DestConfig: `
			webhook_url    = "https://example.com/test"
			webhook_method = "POST"
		`,
	})
}
