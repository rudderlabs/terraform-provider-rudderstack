package destinations_test

import (
	"testing"

	acc "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/acc"
	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

const singularConsentManagementTF = `
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
					flutter = [{
						provider            = "oneTrust"
						consents            = ["one_flutter"]
						resolution_strategy = ""
					}]
					reactnative = [{
						provider            = "oneTrust"
						consents            = ["one_reactnative"]
						resolution_strategy = ""
					}]
					cordova = [{
						provider            = "oneTrust"
						consents            = ["one_cordova"]
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
					shopify = [{
						provider            = "oneTrust"
						consents            = ["one_shopify"]
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
				}
`

const singularConsentManagementAPI = `
				"consentManagement": {
					"android": [{"provider": "oneTrust", "resolutionStrategy": "", "consents": [{"consent": "one_android"}]}],
					"androidKotlin": [{"provider": "oneTrust", "resolutionStrategy": "", "consents": [{"consent": "one_android_kotlin"}]}],
					"ios": [{"provider": "oneTrust", "resolutionStrategy": "", "consents": [{"consent": "one_ios"}]}],
					"iosSwift": [{"provider": "oneTrust", "resolutionStrategy": "", "consents": [{"consent": "one_ios_swift"}]}],
					"flutter": [{"provider": "oneTrust", "resolutionStrategy": "", "consents": [{"consent": "one_flutter"}]}],
					"reactnative": [{"provider": "oneTrust", "resolutionStrategy": "", "consents": [{"consent": "one_reactnative"}]}],
					"cordova": [{"provider": "oneTrust", "resolutionStrategy": "", "consents": [{"consent": "one_cordova"}]}],
					"amp": [{"provider": "oneTrust", "resolutionStrategy": "", "consents": [{"consent": "one_amp"}]}],
					"cloud": [{"provider": "oneTrust", "resolutionStrategy": "", "consents": [{"consent": "one_cloud"}]}],
					"warehouse": [{"provider": "oneTrust", "resolutionStrategy": "", "consents": [{"consent": "one_warehouse"}]}],
					"shopify": [{"provider": "oneTrust", "resolutionStrategy": "", "consents": [{"consent": "one_shopify"}]}],
					"web": [{"provider": "oneTrust", "resolutionStrategy": "", "consents": [{"consent": "one_web"}]}],
					"unity": [{"provider": "oneTrust", "resolutionStrategy": "", "consents": [{"consent": "one_unity"}]}]
				}
`

var singularTestConfigs = []c.TestConfig{
	{
		TerraformCreate: `
				api_key = "test_api_key"
				connection_mode {
					android = "device"
					ios     = "device"
				}
			`,
		APICreate: `{
				"apiKey": "test_api_key",
				"connectionMode": {
					"android": "device",
					"ios": "device"
				}
			}`,
		TerraformUpdate: `
				api_key    = "updated_api_key"
				api_secret = "updatedapisecret"
				session_event_list = ["Application Session", "New Campaign Session"]
				match_id {
					unity = "hash"
				}
				use_native_sdk {
					android     = true
					ios         = true
					reactnative = true
					cordova     = true
				}
				event_filtering {
					whitelist = ["one", "two", "three"]
				}
				connection_mode {
					android        = "device"
					ios            = "device"
					reactnative    = "device"
					cordova        = "device"
					android_kotlin = "cloud"
					ios_swift      = "cloud"
					flutter        = "cloud"
					web            = "cloud"
					unity          = "cloud"
					amp            = "cloud"
					shopify        = "cloud"
					cloud          = "cloud"
					warehouse      = "cloud"
				}
` + singularConsentManagementTF,
		APIUpdate: `{
				"apiKey": "updated_api_key",
				"apiSecret": "updatedapisecret",
				"sessionEventList": [
					{"sessionEventName": "Application Session"},
					{"sessionEventName": "New Campaign Session"}
				],
				"match_id": {"unity": "hash"},
				"useNativeSDK": {
					"android": true,
					"ios": true,
					"reactnative": true,
					"cordova": true
				},
				"eventFilteringOption": "whitelistedEvents",
				"whitelistedEvents": [
					{"eventName": "one"},
					{"eventName": "two"},
					{"eventName": "three"}
				],
				"connectionMode": {
					"android": "device",
					"ios": "device",
					"reactnative": "device",
					"cordova": "device",
					"androidKotlin": "cloud",
					"iosSwift": "cloud",
					"flutter": "cloud",
					"web": "cloud",
					"unity": "cloud",
					"amp": "cloud",
					"shopify": "cloud",
					"cloud": "cloud",
					"warehouse": "cloud"
				},` + singularConsentManagementAPI + `
			}`,
	},
}

func TestDestinationResourceSingular(t *testing.T) {
	cmt.AssertDestination(t, "singular", singularTestConfigs)
}

func TestAccDestinationSingular(t *testing.T) {
	acc.AccAssertDestination(t, "singular", singularTestConfigs)
}
