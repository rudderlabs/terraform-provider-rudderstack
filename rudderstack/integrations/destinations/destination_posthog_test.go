package destinations_test

import (
	"testing"

	acc "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/acc"
	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

var posthogTestConfigs = []c.TestConfig{
	{
		TerraformCreate: `
				endpoint = "https://app.posthog.com"
				api_key = "api_key"
				use_v2_group = true
				connection_mode {
					ios_swift = "cloud"
				}
			`,
		APICreate: `{
				"yourInstance": "https://app.posthog.com",
				"teamApiKey": "api_key",
				"useV2Group": true,
				"connectionMode": {
					"iosSwift": "cloud"
				}
			}`,
		TerraformUpdate: `
				endpoint = "https://app.posthog.com"
				api_key = "api_key"
				use_v2_group = true
				connection_mode {
					web = "device"
					flutter = "cloud"
					cloud = "cloud"
					android = "cloud"
					android_kotlin = "cloud"
					ios = "cloud"
					ios_swift = "cloud"
					unity = "cloud"
					amp = "cloud"
					reactnative = "cloud"
					cordova = "cloud"
					shopify = "cloud"
					warehouse = "cloud"
				}
				event_filtering {
					blacklist = ["event3", "event4"]
				}
				autocapture {
					web = true
				}
				use_native_sdk {
					web = true
				}
				capture_page_view {
					web = true
				}
				disable_session_recording {
					web = true
				}
				enable_local_storage_persistence {
					web = true
				}
			`,
		APIUpdate: `{
				"yourInstance": "https://app.posthog.com",
				"teamApiKey": "api_key",
				"useV2Group": true,
				"connectionMode": {
					"web": "device",
					"flutter": "cloud",
					"cloud": "cloud",
					"android": "cloud",
					"androidKotlin": "cloud",
					"ios": "cloud",
					"iosSwift": "cloud",
					"unity": "cloud",
					"amp": "cloud",
					"reactnative": "cloud",
					"cordova": "cloud",
					"shopify": "cloud",
					"warehouse": "cloud"
				},
				"eventFilteringOption": "blacklistedEvents",
				"blacklistedEvents": [
					{
						"eventName": "event3"
					},
					{
						"eventName": "event4"
					}
				],
				"autocapture": {
					"web": true
				},
				"useNativeSDK": {
					"web": true
				},
				"capturePageView": {
					"web": true
				},
				"disableSessionRecording": {
					"web": true
				},
				"enableLocalStoragePersistence": {
					"web": true
				}
			}`,
	},
}

func TestDestinationResourcePosthog(t *testing.T) {
	cmt.AssertDestination(t, "posthog", posthogTestConfigs)
}

func TestAccDestinationPosthog(t *testing.T) {
	acc.AccAssertDestination(t, "posthog", posthogTestConfigs)
}
