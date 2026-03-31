package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceTiktokAds(t *testing.T) {
	cmt.AssertDestination(t, "tiktok_ads", []c.TestConfig{
		{
			TerraformCreate: `
				pixel_code = "A1T8T4XXXXVIQA8ORZMX9"
			`,
			APICreate: `{
				"pixelCode": "A1T8T4XXXXVIQA8ORZMX9",
				"version": "v2",
				"hashUserProperties": true
			}`,
			TerraformUpdate: `
				pixel_code           = "A1T8T4XXXXVIQA8ORZMX9"
				access_token         = "test-access-token"
				version              = "v1"
				hash_user_properties = false
				send_custom_events   = true

				events_to_standard = [{
					from = "Sign up completed"
					to   = "CompleteRegistration"
				}, {
					from = "Product Added"
					to   = "AddToCart"
				}]

				event_filtering {
					blacklist = ["one", "two", "three"]
				}

				use_native_sdk {
					web = true
				}

				connection_mode {
					web            = "cloud"
					cloud          = "cloud"
					ios            = "cloud"
					ios_swift      = "cloud"
					android        = "cloud"
					android_kotlin = "cloud"
					unity          = "cloud"
					amp            = "cloud"
					warehouse      = "cloud"
					react_native   = "cloud"
					flutter        = "cloud"
					cordova        = "cloud"
					shopify        = "cloud"
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
				"pixelCode": "A1T8T4XXXXVIQA8ORZMX9",
				"accessToken": "test-access-token",
				"version": "v1",
				"hashUserProperties": false,
				"sendCustomEvents": true,
				"eventsToStandard": [
					{ "from": "Sign up completed", "to": "CompleteRegistration" },
					{ "from": "Product Added", "to": "AddToCart" }
				],
				"blacklistedEvents": [
					{ "eventName": "one" },
					{ "eventName": "two" },
					{ "eventName": "three" }
				],
				"eventFilteringOption": "blacklistedEvents",
				"useNativeSDK": { "web": true },
				"connectionMode": {
					"web": "cloud",
					"cloud": "cloud",
					"ios": "cloud",
					"iosSwift": "cloud",
					"android": "cloud",
					"androidKotlin": "cloud",
					"unity": "cloud",
					"amp": "cloud",
					"warehouse": "cloud",
					"reactnative": "cloud",
					"flutter": "cloud",
					"cordova": "cloud",
					"shopify": "cloud"
				},
				"consentManagement": {
					"web": [
						{
							"provider": "oneTrust",
							"resolutionStrategy": "",
							"consents": [
								{ "consent": "one_web" },
								{ "consent": "two_web" },
								{ "consent": "three_web" }
							]
						},
						{
							"provider": "ketch",
							"resolutionStrategy": "",
							"consents": [
								{ "consent": "one_web" },
								{ "consent": "two_web" },
								{ "consent": "three_web" }
							]
						},
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{ "consent": "one_web" },
								{ "consent": "two_web" },
								{ "consent": "three_web" }
							]
						}
					],
					"android": [
						{
							"provider": "ketch",
							"resolutionStrategy": "",
							"consents": [
								{ "consent": "one_android" },
								{ "consent": "two_android" },
								{ "consent": "three_android" }
							]
						}
					],
					"androidKotlin": [
						{
							"provider": "ketch",
							"resolutionStrategy": "",
							"consents": [
								{ "consent": "one_android_kotlin" },
								{ "consent": "two_android_kotlin" },
								{ "consent": "three_android_kotlin" }
							]
						}
					],
					"ios": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{ "consent": "one_ios" },
								{ "consent": "two_ios" },
								{ "consent": "three_ios" }
							]
						}
					],
					"iosSwift": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{ "consent": "one_ios_swift" },
								{ "consent": "two_ios_swift" },
								{ "consent": "three_ios_swift" }
							]
						}
					],
					"unity": [
						{
							"provider": "custom",
							"resolutionStrategy": "or",
							"consents": [
								{ "consent": "one_unity" },
								{ "consent": "two_unity" },
								{ "consent": "three_unity" }
							]
						}
					],
					"reactnative": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{ "consent": "one_reactnative" },
								{ "consent": "two_reactnative" },
								{ "consent": "three_reactnative" }
							]
						}
					],
					"flutter": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{ "consent": "one_flutter" },
								{ "consent": "two_flutter" },
								{ "consent": "three_flutter" }
							]
						}
					],
					"cordova": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{ "consent": "one_cordova" },
								{ "consent": "two_cordova" },
								{ "consent": "three_cordova" }
							]
						}
					],
					"amp": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{ "consent": "one_amp" },
								{ "consent": "two_amp" },
								{ "consent": "three_amp" }
							]
						}
					],
					"cloud": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{ "consent": "one_cloud" },
								{ "consent": "two_cloud" },
								{ "consent": "three_cloud" }
							]
						}
					],
					"warehouse": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{ "consent": "one_warehouse" },
								{ "consent": "two_warehouse" },
								{ "consent": "three_warehouse" }
							]
						}
					],
					"shopify": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{ "consent": "one_shopify" },
								{ "consent": "two_shopify" },
								{ "consent": "three_shopify" }
							]
						}
					]
				}
			}`,
		},
	})
}
