package destinations_test

import (
	"testing"

	acc "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/acc"
	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

var adjustTestConfigs = []c.TestConfig{
	{
		TerraformCreate: `
			app_token = "test_app_token"
			environment = true
			connection_mode {
				android_kotlin = "device"
				android = "device"
				ios_swift = "device"
				ios = "device"
			}
			`,
		APICreate: `{
				"appToken": "test_app_token",
				"environment": true,
				"connectionMode": {
					"androidKotlin": "device",
					"android": "device",
					"iosSwift": "device",
					"ios": "device"
				}
			}`,
		TerraformUpdate: `
			app_token = "updated_app_token"
			delay = "5"
			environment = false
			custom_mappings = [
				{
					from = "Product Purchased"
					to = "abc123"
				},
				{
					from = "Signup"
					to = "def456"
				}
			]
			partner_param_keys = [
				{
					from = "userId"
					to = "user_id"
				}
			]
			enable_install_attribution_tracking {
				android = true
				ios = true
			}
			event_filtering {
				whitelist = [ "one", "two", "three" ]
			}
			connection_mode {
				android_kotlin = "cloud"
				android = "cloud"
				ios_swift = "cloud"
				ios = "cloud"
				unity = "cloud"
				reactnative = "cloud"
				flutter = "cloud"
				cordova = "cloud"
				shopify = "cloud"
				cloud = "cloud"
				warehouse = "cloud"
			}
			consent_management {
				android_kotlin = [{
					provider = "ketch"
					consents = ["one_android_kotlin", "two_android_kotlin", "three_android_kotlin"]
					resolution_strategy = ""
				}]
				android = [{
					provider = "ketch"
					consents = ["one_android", "two_android", "three_android"]
					resolution_strategy = ""
				}]
				ios_swift = [{
					provider = "custom"
					resolution_strategy = "and"
					consents = ["one_ios_swift", "two_ios_swift", "three_ios_swift"]
				}]
				ios = [{
					provider = "custom"
					resolution_strategy = "and"
					consents = ["one_ios", "two_ios", "three_ios"]
				}]
				unity = [{
					provider = "custom"
					resolution_strategy = "or"
					consents = ["one_unity", "two_unity", "three_unity"]
				}]
				reactnative = [{
					provider = "custom"
					resolution_strategy = "and"
					consents = ["one_reactnative", "two_reactnative", "three_reactnative"]
				}]
				flutter = [{
					provider = "custom"
					resolution_strategy = "and"
					consents = ["one_flutter", "two_flutter", "three_flutter"]
				}]
				cordova = [{
					provider = "custom"
					resolution_strategy = "and"
					consents = ["one_cordova", "two_cordova", "three_cordova"]
				}]
				cloud = [{
					provider = "custom"
					resolution_strategy = "and"
					consents = ["one_cloud", "two_cloud", "three_cloud"]
				}]
				warehouse = [{
					provider = "custom"
					resolution_strategy = "and"
					consents = ["one_warehouse", "two_warehouse", "three_warehouse"]
				}]
				shopify = [{
					provider = "custom"
					resolution_strategy = "and"
					consents = ["one_shopify", "two_shopify", "three_shopify"]
				}]
			}
			`,
		APIUpdate: `{
				"appToken": "updated_app_token",
				"delay": "5",
				"customMappings": [
					{
						"from": "Product Purchased",
						"to": "abc123"
					},
					{
						"from": "Signup",
						"to": "def456"
					}
				],
				"partnerParamKeys": [
					{
						"from": "userId",
						"to": "user_id"
					}
				],
				"enableInstallAttributionTracking": {
					"android": true,
					"ios": true
				},
				"eventFilteringOption": "whitelistedEvents",
				"whitelistedEvents": [{
					"eventName": "one"
				}, {
					"eventName": "two"
				}, {
					"eventName": "three"
				}],
				"connectionMode": {
					"androidKotlin": "cloud",
					"android": "cloud",
					"iosSwift": "cloud",
					"ios": "cloud",
					"unity": "cloud",
					"reactnative": "cloud",
					"flutter": "cloud",
					"cordova": "cloud",
					"shopify": "cloud",
					"cloud": "cloud",
					"warehouse": "cloud"
				},
				"consentManagement": {
					"androidKotlin": [
						{
							"provider": "ketch",
							"resolutionStrategy": "",
							"consents": [
								{
									"consent": "one_android_kotlin"
								},
								{
									"consent": "two_android_kotlin"
								},
								{
									"consent": "three_android_kotlin"
								}
							]
						}
					],
					"android": [
						{
							"provider": "ketch",
							"resolutionStrategy": "",
							"consents": [
								{
									"consent": "one_android"
								},
								{
									"consent": "two_android"
								},
								{
									"consent": "three_android"
								}
							]
						}
					],
					"iosSwift": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{
									"consent": "one_ios_swift"
								},
								{
									"consent": "two_ios_swift"
								},
								{
									"consent": "three_ios_swift"
								}
							]
						}
					],
					"ios": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{
									"consent": "one_ios"
								},
								{
									"consent": "two_ios"
								},
								{
									"consent": "three_ios"
								}
							]
						}
					],
					"unity": [
						{
							"provider": "custom",
							"resolutionStrategy": "or",
							"consents": [
								{
									"consent": "one_unity"
								},
								{
									"consent": "two_unity"
								},
								{
									"consent": "three_unity"
								}
							]
						}
					],
					"reactnative": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{
									"consent": "one_reactnative"
								},
								{
									"consent": "two_reactnative"
								},
								{
									"consent": "three_reactnative"
								}
							]
						}
					],
					"flutter": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{
									"consent": "one_flutter"
								},
								{
									"consent": "two_flutter"
								},
								{
									"consent": "three_flutter"
								}
							]
						}
					],
					"cordova": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{
									"consent": "one_cordova"
								},
								{
									"consent": "two_cordova"
								},
								{
									"consent": "three_cordova"
								}
							]
						}
					],
					"cloud": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{
									"consent": "one_cloud"
								},
								{
									"consent": "two_cloud"
								},
								{
									"consent": "three_cloud"
								}
							]
						}
					],
					"warehouse": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{
									"consent": "one_warehouse"
								},
								{
									"consent": "two_warehouse"
								},
								{
									"consent": "three_warehouse"
								}
							]
						}
					],
					"shopify": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{
									"consent": "one_shopify"
								},
								{
									"consent": "two_shopify"
								},
								{
									"consent": "three_shopify"
								}
							]
						}
					]
				}
			}`,
	},
}

func TestDestinationResourceAdjust(t *testing.T) {
	cmt.AssertDestination(t, "adjust", adjustTestConfigs)
}

func TestAccDestinationAdjust(t *testing.T) {
	acc.AccAssertDestination(t, "adjust", adjustTestConfigs)
}
