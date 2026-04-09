package destinations_test

import (
	"testing"

	acc "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/acc"
	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

var facebookConversionsTestConfigs = []c.TestConfig{
	{
			TerraformCreate: `
				dataset_id   = "1234567898765"
				access_token = "my-access-token"
			`,
			APICreate: `{
				"datasetId": "1234567898765",
				"accessToken": "my-access-token",
				"actionSource": "website"
			}`,
			TerraformUpdate: `
				dataset_id   = "1234567898765"
				access_token = "my-access-token"

				action_source      = "app"
				limited_data_usage = true
				test_destination   = true
				test_event_code    = "TEST80569"
				remove_external_id = true

				events_to_events = [{
					from = "Product Searched"
					to   = "Search"
				}, {
					from = "Order Completed"
					to   = "Purchase"
				}]

				blacklist_pii_properties = [
					{
						property = "phone"
						hash     = false
					},
					{
						property = "email"
						hash     = true
					}
				]

				whitelist_pii_properties = [
					{
						property = "name"
					},
					{
						property = "address"
					}
				]

				consent_management {
					web = [
						{
							provider = "oneTrust"
							consents = ["one_web", "two_web", "three_web"]
							resolution_strategy = ""
						},
						{
							provider = "ketch"
							consents = ["one_web", "two_web", "three_web"]
							resolution_strategy = ""
						},
						{
							provider = "custom"
							resolution_strategy = "and"
							consents = ["one_web", "two_web", "three_web"]
						}
					]
					android = [{
						provider = "ketch"
						consents = ["one_android", "two_android", "three_android"]
						resolution_strategy = ""
					}]
					android_kotlin = [{
						provider = "ketch"
						consents = ["one_android_kotlin", "two_android_kotlin", "three_android_kotlin"]
						resolution_strategy = ""
					}]
					ios = [{
						provider = "custom"
						resolution_strategy = "and"
						consents = ["one_ios", "two_ios", "three_ios"]
					}]
					ios_swift = [{
						provider = "custom"
						resolution_strategy = "and"
						consents = ["one_ios_swift", "two_ios_swift", "three_ios_swift"]
					}]
					unity = [{
						provider = "custom"
						resolution_strategy = "or"
						consents = ["one_unity", "two_unity", "three_unity"]
					}]
					amp = [{
						provider = "custom"
						resolution_strategy = "and"
						consents = ["one_amp", "two_amp", "three_amp"]
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
					shopify = [{
						provider = "custom"
						resolution_strategy = "and"
						consents = ["one_shopify", "two_shopify", "three_shopify"]
					}]
				}
			`,
			APIUpdate: `{
				"datasetId": "1234567898765",
				"accessToken": "my-access-token",
				"actionSource": "app",
				"limitedDataUSage": true,
				"testDestination": true,
				"testEventCode": "TEST80569",
				"removeExternalId": true,
				"eventsToEvents": [
					{ "from": "Product Searched", "to": "Search" },
					{ "from": "Order Completed", "to": "Purchase" }
				],
				"blacklistPiiProperties": [
					{ "blacklistPiiProperties": "phone", "blacklistPiiHash": false },
					{ "blacklistPiiProperties": "email", "blacklistPiiHash": true }
				],
				"whitelistPiiProperties": [
					{ "whitelistPiiProperties": "name" },
					{ "whitelistPiiProperties": "address" }
				],
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

func TestDestinationResourceFacebookConversions(t *testing.T) {
	cmt.AssertDestination(t, "facebook_conversions", facebookConversionsTestConfigs)
}

func TestAccDestinationFacebookConversions(t *testing.T) {
	acc.AccAssertDestination(t, "facebook_conversions", facebookConversionsTestConfigs)
}
