package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceKafka(t *testing.T) {
	cmt.AssertDestination(t, "kafka", []c.TestConfig{
		{
			TerraformCreate: `
				host_name = "example.com"
				port = "9092"
				topic = "example-topic"
			`,
			APICreate: `{
				"hostName": "example.com",
				"port": "9092",
				"topic": "example-topic",
				"sslEnabled": true
			}`,
			TerraformUpdate: `
				host_name = "example-updated.com"
				port = "9092"
				topic = "example-topic"
				ssl_enabled = true
				ca_certificate = "example-ca-certificate"
			`,
			APIUpdate: `{
				"hostName": "example-updated.com",
				"port": "9092",
				"topic": "example-topic",
				"sslEnabled": true,
				"caCertificate": "example-ca-certificate"
			}`,
		},
	})
}
