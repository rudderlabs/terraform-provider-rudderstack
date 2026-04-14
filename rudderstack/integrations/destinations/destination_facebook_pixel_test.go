package destinations_test

import (
	"testing"

	acc "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/acc"
	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

var facebookPixelTestConfigs = []c.TestConfig{
	{
		TerraformCreate: `
				pixel_id     = "abc123"
				access_token = "placeholder-access-token"
			`,
		APICreate: `{
				"pixelId": "abc123",
				"accessToken": "placeholder-access-token",
				"valueFieldIdentifier": "properties.price"
			}`,
		TerraformUpdate: `
				pixel_id = "facebook pixel id"

				access_token = "facebook access token"
			
				standard_page_call     = true
				value_field_identifier = "properties.price"
				advanced_mapping       = true
				test_destination       = true
				test_event_code        = "..."
			
				events_to_events = [{
					from = "Order Completed"
					to   = "Purchase"
				}, {
					from = "Product Added"
					to   = "AddToCart"
				}]

				blacklist_pii_properties = [
					{ 
						property = "one"
						hash     = false
					},
					{ 
						property = "two"
						hash     = true
					}
				]

				whitelist_pii_properties = [
					{
						property = "one"
					},
					{
						property = "two"
					}
				]

				legacy_conversion_pixel_id {
				  from = "from"
				  to   = "to"
				}
			
				use_native_sdk {
				  web          = true
				}
			
				event_filtering {
				  blacklist = ["one", "two", "three"]
				}
			
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
					shopify = [{
						provider = "custom"
						resolution_strategy = "and"
						consents = ["one_shopify", "two_shopify", "three_shopify"]
					}]
				}
			`,
		APIUpdate: `{
				"pixelId": "facebook pixel id",
				"accessToken": "facebook access token",
				"standardPageCall": true,
				"valueFieldIdentifier": "properties.price",
				"advancedMapping": true,
				"testDestination": true,
				"testEventCode": "...",
				"eventsToEvents": [
				  { "from": "Order Completed", "to": "Purchase" },
				  { "from": "Product Added", "to": "AddToCart" }
				],
				"blacklistPiiProperties": [
				  { "blacklistPiiProperties": "one", "blacklistPiiHash": false },
				  { "blacklistPiiProperties": "two", "blacklistPiiHash": true }
				],
				"whitelistPiiProperties": [
				  { "whitelistPiiProperties": "one" },
				  { "whitelistPiiProperties": "two" }
				],
				"legacyConversionPixelId": { "from": "from", "to": "to" },
				"useNativeSDK": { "web": true },
				"consentManagement": {
					"web": [
						{
							"provider": "oneTrust",
							"resolutionStrategy": "",
							"consents": [
								{
									"consent": "one_web"
								},
								{
									"consent": "two_web"
								},
								{
									"consent": "three_web"
								}
							]
						},
						{
							"provider": "ketch",
							"resolutionStrategy": "",
							"consents": [
								{
									"consent": "one_web"
								},
								{
									"consent": "two_web"
								},
								{
									"consent": "three_web"
								}
							]
						},
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{
									"consent": "one_web"
								},
								{
									"consent": "two_web"
								},
								{
									"consent": "three_web"
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
					"amp": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{
									"consent": "one_amp"
								},
								{
									"consent": "two_amp"
								},
								{
									"consent": "three_amp"
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
				},
				"blacklistedEvents": [
				  { "eventName": "one" },
				  { "eventName": "two" },
				  { "eventName": "three" }
				],
				"eventFilteringOption": "blacklistedEvents"
			}`,
	},
}

func TestDestinationResourceFacebookPixel(t *testing.T) {
	cmt.AssertDestination(t, "facebook_pixel", facebookPixelTestConfigs)
}

func TestAccDestinationFacebookPixel(t *testing.T) {
	acc.AccAssertDestination(t, "facebook_pixel", facebookPixelTestConfigs)
}
