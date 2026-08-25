package destinations_test

import (
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	acc "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/acc"
	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

var googleAdwordsOfflineConversionsTestConfigs = []c.TestConfig{
	{
		TerraformCreate: `
				rudder_account_id = "account-id-1"
				customer_id       = "1234567890"
			`,
		APICreate: `{
				"rudderAccountId": "account-id-1",
				"customerId": "1234567890",
				"subAccount": false,
				"UserIdentifierSource": "none",
				"conversionEnvironment": "none",
				"defaultUserIdentifier": "email",
				"hashUserIdentifier": true,
				"validateOnly": false
			}`,
		TerraformUpdate: `
				rudder_account_id          = "account-id-1"
				customer_id                = "1234567890"
				sub_account                = true
				login_customer_id          = "0987654321"
				user_identifier_source     = "THIRD_PARTY"
				conversion_environment     = "WEB"
				default_user_identifier    = "phone"
				hash_user_identifier       = false
				validate_only              = true

				events_to_offline_conversions_type_mapping = [
					{
						from = "Order Completed"
						to   = "click"
					},
					{
						from = "Call Requested"
						to   = "call"
					}
				]

				events_to_conversions_names_mapping = [
					{
						from = "Order Completed"
						to   = "Purchase Conversion"
					},
					{
						from = "Sign Up Completed"
						to   = "Signup Conversion"
					}
				]

				custom_variables = [
					{
						from = "revenue"
						to   = "cart_value"
					},
					{
						from = "category"
						to   = "product_category"
					}
				]

				connection_mode {
					android        = "cloud"
					android_kotlin = "cloud"
					ios            = "cloud"
					ios_swift      = "cloud"
					web            = "cloud"
					unity          = "cloud"
					amp            = "cloud"
					cloud          = "cloud"
					reactnative    = "cloud"
					flutter        = "cloud"
					cordova        = "cloud"
					warehouse      = "cloud"
					shopify        = "cloud"
				}

				consent_management {
					web = [
						{
							provider            = "oneTrust"
							consents            = ["one_web", "two_web", "three_web"]
							resolution_strategy = ""
						},
						{
							provider            = "ketch"
							consents            = ["one_web", "two_web", "three_web"]
							resolution_strategy = ""
						},
						{
							provider            = "custom"
							resolution_strategy = "and"
							consents            = ["one_web", "two_web", "three_web"]
						}
					]
					android = [{
						provider            = "ketch"
						consents            = ["one_android", "two_android", "three_android"]
						resolution_strategy = ""
					}]
					android_kotlin = [{
						provider            = "ketch"
						consents            = ["one_android_kotlin", "two_android_kotlin", "three_android_kotlin"]
						resolution_strategy = ""
					}]
					ios = [{
						provider            = "custom"
						resolution_strategy = "and"
						consents            = ["one_ios", "two_ios", "three_ios"]
					}]
					ios_swift = [{
						provider            = "custom"
						resolution_strategy = "and"
						consents            = ["one_ios_swift", "two_ios_swift", "three_ios_swift"]
					}]
					unity = [{
						provider            = "custom"
						resolution_strategy = "or"
						consents            = ["one_unity", "two_unity", "three_unity"]
					}]
					amp = [{
						provider            = "custom"
						resolution_strategy = "and"
						consents            = ["one_amp", "two_amp", "three_amp"]
					}]
					cloud = [{
						provider            = "custom"
						resolution_strategy = "and"
						consents            = ["one_cloud", "two_cloud", "three_cloud"]
					}]
					reactnative = [{
						provider            = "custom"
						resolution_strategy = "and"
						consents            = ["one_reactnative", "two_reactnative", "three_reactnative"]
					}]
					flutter = [{
						provider            = "custom"
						resolution_strategy = "and"
						consents            = ["one_flutter", "two_flutter", "three_flutter"]
					}]
					cordova = [{
						provider            = "custom"
						resolution_strategy = "and"
						consents            = ["one_cordova", "two_cordova", "three_cordova"]
					}]
					warehouse = [{
						provider            = "custom"
						resolution_strategy = "and"
						consents            = ["one_warehouse", "two_warehouse", "three_warehouse"]
					}]
					shopify = [{
						provider            = "custom"
						resolution_strategy = "and"
						consents            = ["one_shopify", "two_shopify", "three_shopify"]
					}]
				}
			`,
		APIUpdate: `{
				"rudderAccountId": "account-id-1",
				"customerId": "1234567890",
				"subAccount": true,
				"loginCustomerId": "0987654321",
				"eventsToOfflineConversionsTypeMapping": [
					{ "from": "Order Completed", "to": "click" },
					{ "from": "Call Requested", "to": "call" }
				],
				"eventsToConversionsNamesMapping": [
					{ "from": "Order Completed", "to": "Purchase Conversion" },
					{ "from": "Sign Up Completed", "to": "Signup Conversion" }
				],
				"customVariables": [
					{ "from": "revenue", "to": "cart_value" },
					{ "from": "category", "to": "product_category" }
				],
				"UserIdentifierSource": "THIRD_PARTY",
				"conversionEnvironment": "WEB",
				"defaultUserIdentifier": "phone",
				"hashUserIdentifier": false,
				"validateOnly": true,
				"connectionMode": {
					"android": "cloud",
					"androidKotlin": "cloud",
					"ios": "cloud",
					"iosSwift": "cloud",
					"web": "cloud",
					"unity": "cloud",
					"amp": "cloud",
					"cloud": "cloud",
					"reactnative": "cloud",
					"flutter": "cloud",
					"cordova": "cloud",
					"warehouse": "cloud",
					"shopify": "cloud"
				},
				"consentManagement": {
					"web": [
						{
							"provider": "oneTrust",
							"resolutionStrategy": "",
							"consents": [
								{"consent": "one_web"},
								{"consent": "two_web"},
								{"consent": "three_web"}
							]
						},
						{
							"provider": "ketch",
							"resolutionStrategy": "",
							"consents": [
								{"consent": "one_web"},
								{"consent": "two_web"},
								{"consent": "three_web"}
							]
						},
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{"consent": "one_web"},
								{"consent": "two_web"},
								{"consent": "three_web"}
							]
						}
					],
					"android": [
						{
							"provider": "ketch",
							"resolutionStrategy": "",
							"consents": [
								{"consent": "one_android"},
								{"consent": "two_android"},
								{"consent": "three_android"}
							]
						}
					],
					"androidKotlin": [
						{
							"provider": "ketch",
							"resolutionStrategy": "",
							"consents": [
								{"consent": "one_android_kotlin"},
								{"consent": "two_android_kotlin"},
								{"consent": "three_android_kotlin"}
							]
						}
					],
					"ios": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{"consent": "one_ios"},
								{"consent": "two_ios"},
								{"consent": "three_ios"}
							]
						}
					],
					"iosSwift": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{"consent": "one_ios_swift"},
								{"consent": "two_ios_swift"},
								{"consent": "three_ios_swift"}
							]
						}
					],
					"unity": [
						{
							"provider": "custom",
							"resolutionStrategy": "or",
							"consents": [
								{"consent": "one_unity"},
								{"consent": "two_unity"},
								{"consent": "three_unity"}
							]
						}
					],
					"amp": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{"consent": "one_amp"},
								{"consent": "two_amp"},
								{"consent": "three_amp"}
							]
						}
					],
					"cloud": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{"consent": "one_cloud"},
								{"consent": "two_cloud"},
								{"consent": "three_cloud"}
							]
						}
					],
					"reactnative": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{"consent": "one_reactnative"},
								{"consent": "two_reactnative"},
								{"consent": "three_reactnative"}
							]
						}
					],
					"flutter": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{"consent": "one_flutter"},
								{"consent": "two_flutter"},
								{"consent": "three_flutter"}
							]
						}
					],
					"cordova": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{"consent": "one_cordova"},
								{"consent": "two_cordova"},
								{"consent": "three_cordova"}
							]
						}
					],
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
					],
					"shopify": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{"consent": "one_shopify"},
								{"consent": "two_shopify"},
								{"consent": "three_shopify"}
							]
						}
					]
				}
			}`,
	},
}

func TestDestinationResourceGoogleAdwordsOfflineConversions(t *testing.T) {
	cmt.AssertDestination(t, "google_adwords_offline_conversions", googleAdwordsOfflineConversionsTestConfigs)
}

func TestDestinationResourceGoogleAdwordsOfflineConversionsRejectsEmptyConversionType(t *testing.T) {
	configSchema := c.Destinations.Entries()["google_adwords_offline_conversions"].ConfigSchema
	mappingSchema := configSchema["events_to_offline_conversions_type_mapping"].Elem.(*schema.Resource)
	toSchema := mappingSchema.Schema["to"]

	if diags := toSchema.ValidateDiagFunc("", cty.Path{}); !diags.HasError() {
		t.Fatal("expected empty offline conversion type to fail validation")
	}
}

func TestAccDestinationGoogleAdwordsOfflineConversions(t *testing.T) {
	if !acc.PlanOnly() {
		t.Skip("skipping: requires valid OAuth account in workspace")
	}
	acc.AccAssertDestination(t, "google_adwords_offline_conversions", googleAdwordsOfflineConversionsTestConfigs)
}
