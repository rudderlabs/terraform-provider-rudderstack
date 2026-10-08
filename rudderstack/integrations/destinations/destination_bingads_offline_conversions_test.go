package destinations_test

import (
	"testing"

	acc "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/acc"
	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

var bingadsOfflineConversionsTestConfigs = []c.TestConfig{
	{
		TerraformCreate: `
				rudder_account_id   = "__ACCOUNT_ID__"
				customer_account_id = "53212345"
				customer_id         = "343598"
			`,
		APICreate: `{
				"rudderAccountId": "__ACCOUNT_ID__",
				"customerAccountId": "53212345",
				"customerId": "343598",
				"isHashRequired": false
			}`,
		TerraformUpdate: `
				rudder_account_id   = "__ACCOUNT_ID__"
				customer_account_id = "53212345"
				customer_id         = "343598"
				is_hash_required    = true

				connection_mode {
					warehouse = "cloud"
				}

				consent_management {
					warehouse = [
						{
							provider            = "custom"
							resolution_strategy = "and"
							consents            = ["one_warehouse", "two_warehouse", "three_warehouse"]
						},
						{
							provider            = "ketch"
							resolution_strategy = ""
							consents            = ["ketch_warehouse"]
						}
					]
				}
			`,
		APIUpdate: `{
				"rudderAccountId": "__ACCOUNT_ID__",
				"customerAccountId": "53212345",
				"customerId": "343598",
				"isHashRequired": true,
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
						},
						{
							"provider": "ketch",
							"resolutionStrategy": "",
							"consents": [
								{"consent": "ketch_warehouse"}
							]
						}
					]
				}
			}`,
	},
}

func TestDestinationResourceBingadsOfflineConversions(t *testing.T) {
	cmt.AssertDestination(t, "bingads_offline_conversions", bingadsOfflineConversionsTestConfigs)
}

func TestAccDestinationBingadsOfflineConversions(t *testing.T) {
	acc.AccAssertOAuthDestination(t, "bingads_offline_conversions", bingadsOfflineConversionsTestConfigs)
}
