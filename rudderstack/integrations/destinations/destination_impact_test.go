package destinations_test

import (
	"testing"

	acc "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/acc"
	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

const impactConsentManagementTF = `
				consent_management {
					android = [{
						provider            = "oneTrust"
						consents            = ["one_android"]
						resolution_strategy = ""
					}]
					android_kotlin = [{
						provider            = "oneTrust"
						consents            = ["one_android_kotlin"]
						resolution_strategy = ""
					}]
					ios = [{
						provider            = "oneTrust"
						consents            = ["one_ios"]
						resolution_strategy = ""
					}]
					ios_swift = [{
						provider            = "oneTrust"
						consents            = ["one_ios_swift"]
						resolution_strategy = ""
					}]
					web = [{
						provider            = "oneTrust"
						consents            = ["one_web"]
						resolution_strategy = ""
					}]
					unity = [{
						provider            = "oneTrust"
						consents            = ["one_unity"]
						resolution_strategy = ""
					}]
					amp = [{
						provider            = "oneTrust"
						consents            = ["one_amp"]
						resolution_strategy = ""
					}]
					cloud = [{
						provider            = "oneTrust"
						consents            = ["one_cloud"]
						resolution_strategy = ""
					}]
					warehouse = [{
						provider            = "oneTrust"
						consents            = ["one_warehouse"]
						resolution_strategy = ""
					}]
					reactnative = [{
						provider            = "oneTrust"
						consents            = ["one_reactnative"]
						resolution_strategy = ""
					}]
					flutter = [{
						provider            = "oneTrust"
						consents            = ["one_flutter"]
						resolution_strategy = ""
					}]
					cordova = [{
						provider            = "oneTrust"
						consents            = ["one_cordova"]
						resolution_strategy = ""
					}]
					shopify = [{
						provider            = "oneTrust"
						consents            = ["one_shopify"]
						resolution_strategy = ""
					}]
				}
`

const impactConsentManagementAPI = `
				"consentManagement": {
					"android": [{"provider": "oneTrust", "resolutionStrategy": "", "consents": [{"consent": "one_android"}]}],
					"androidKotlin": [{"provider": "oneTrust", "resolutionStrategy": "", "consents": [{"consent": "one_android_kotlin"}]}],
					"ios": [{"provider": "oneTrust", "resolutionStrategy": "", "consents": [{"consent": "one_ios"}]}],
					"iosSwift": [{"provider": "oneTrust", "resolutionStrategy": "", "consents": [{"consent": "one_ios_swift"}]}],
					"web": [{"provider": "oneTrust", "resolutionStrategy": "", "consents": [{"consent": "one_web"}]}],
					"unity": [{"provider": "oneTrust", "resolutionStrategy": "", "consents": [{"consent": "one_unity"}]}],
					"amp": [{"provider": "oneTrust", "resolutionStrategy": "", "consents": [{"consent": "one_amp"}]}],
					"cloud": [{"provider": "oneTrust", "resolutionStrategy": "", "consents": [{"consent": "one_cloud"}]}],
					"warehouse": [{"provider": "oneTrust", "resolutionStrategy": "", "consents": [{"consent": "one_warehouse"}]}],
					"reactnative": [{"provider": "oneTrust", "resolutionStrategy": "", "consents": [{"consent": "one_reactnative"}]}],
					"flutter": [{"provider": "oneTrust", "resolutionStrategy": "", "consents": [{"consent": "one_flutter"}]}],
					"cordova": [{"provider": "oneTrust", "resolutionStrategy": "", "consents": [{"consent": "one_cordova"}]}],
					"shopify": [{"provider": "oneTrust", "resolutionStrategy": "", "consents": [{"consent": "one_shopify"}]}]
				}
`

var impactTestConfigs = []c.TestConfig{
	{
		TerraformCreate: `
				account_sid = "account_sid"
				api_key     = "api_key"
				campaign_id = "12345"
			`,
		APICreate: `{
				"accountSID": "account_sid",
				"apiKey": "api_key",
				"campaignId": "12345"
			}`,
		TerraformUpdate: `
				account_sid              = "updated_account_sid"
				api_key                  = "updated_api_key"
				campaign_id              = "23456"
				event_type_id            = "34567"
				impact_app_id            = "45678"
				enable_email_hashing     = true
				enable_identify_events   = true
				enable_page_events       = true
				enable_screen_events     = true
				rudder_to_impact_property = [{
					from = "email"
					to   = "CustomerEmail"
				}, {
					from = "order_id"
					to   = "OrderId"
				}]
				products_mapping = [{
					from = "products.$.sku"
					to   = "Sku"
				}, {
					from = "products.$.price"
					to   = "Price"
				}]
				action_event_names  = ["Order Completed", "Product Added"]
				install_event_names = ["Application Installed", "App Opened"]
				connection_mode {
					android        = "cloud"
					android_kotlin = "cloud"
					ios            = "cloud"
					ios_swift      = "cloud"
					web            = "cloud"
					unity          = "cloud"
					amp            = "cloud"
					cloud          = "cloud"
					warehouse      = "cloud"
					reactnative    = "cloud"
					flutter        = "cloud"
					cordova        = "cloud"
					shopify        = "cloud"
				}
` + impactConsentManagementTF,
		APIUpdate: `{
				"accountSID": "updated_account_sid",
				"apiKey": "updated_api_key",
				"campaignId": "23456",
				"eventTypeId": "34567",
				"impactAppId": "45678",
				"enableEmailHashing": true,
				"enableIdentifyEvents": true,
				"enablePageEvents": true,
				"enableScreenEvents": true,
				"rudderToImpactProperty": [
					{"from": "email", "to": "CustomerEmail"},
					{"from": "order_id", "to": "OrderId"}
				],
				"productsMapping": [
					{"from": "products.$.sku", "to": "Sku"},
					{"from": "products.$.price", "to": "Price"}
				],
				"actionEventNames": [
					{"eventName": "Order Completed"},
					{"eventName": "Product Added"}
				],
				"installEventNames": [
					{"eventName": "Application Installed"},
					{"eventName": "App Opened"}
				],
				"connectionMode": {
					"android": "cloud",
					"androidKotlin": "cloud",
					"ios": "cloud",
					"iosSwift": "cloud",
					"web": "cloud",
					"unity": "cloud",
					"amp": "cloud",
					"cloud": "cloud",
					"warehouse": "cloud",
					"reactnative": "cloud",
					"flutter": "cloud",
					"cordova": "cloud",
					"shopify": "cloud"
				},` + impactConsentManagementAPI + `
			}`,
	},
}

func TestDestinationResourceImpact(t *testing.T) {
	cmt.AssertDestination(t, "impact", impactTestConfigs)
}

func TestAccDestinationImpact(t *testing.T) {
	acc.AccAssertDestination(t, "impact", impactTestConfigs)
}
