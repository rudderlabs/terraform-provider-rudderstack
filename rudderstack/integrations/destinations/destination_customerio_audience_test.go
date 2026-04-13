package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceCustomerioAudience(t *testing.T) {
	cmt.AssertDestination(t, "customerio_audience", []c.TestConfig{
		{
			TerraformCreate: `
				site_id     = "site-id-1"
				api_key     = "api-key-1"
				app_api_key = "app-api-key-1"
				region      = "US"
			`,
			APICreate: `{
				"siteId":    "site-id-1",
				"apiKey":    "api-key-1",
				"appApiKey": "app-api-key-1",
				"region":    "US"
			}`,
			TerraformUpdate: `
				site_id     = "site-id-1"
				api_key     = "api-key-1"
				app_api_key = "app-api-key-1"
				region      = "EU"
				connection_mode {
					warehouse = "cloud"
				}
				consent_management {
					warehouse = [{
						provider = "custom"
						resolution_strategy = "and"
						consents = ["one_warehouse", "two_warehouse", "three_warehouse"]
					}]
				}
			`,
			APIUpdate: `{
				"siteId":    "site-id-1",
				"apiKey":    "api-key-1",
				"appApiKey": "app-api-key-1",
				"region":    "EU",
				"connectionMode": {
					"warehouse": "cloud"
				},
				"consentManagement": {
					"warehouse": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{"consent": "one_warehouse"},
								{"consent": "two_warehouse"},
								{"consent": "three_warehouse"}
							]
						}
					]
				}
			}`,
		},
	})
}
