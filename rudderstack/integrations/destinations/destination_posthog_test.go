package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourcePosthog(t *testing.T) {
	cmt.AssertDestination(t, "posthog", []c.TestConfig{
		{
			TerraformCreate: `
				endpoint = "https://app.posthog.com"
				api_key = "api_key"
				use_v2_group = true
				connection_mode {
					web = "cloud"
					flutter = "cloud"
					cloud = "cloud"
				}
			`,
			APICreate: `{
				"yourInstance": "https://app.posthog.com",
				"teamApiKey": "api_key",
				"useV2Group": true,
				"connectionMode": {
					"web": "cloud",
					"flutter": "cloud",
					"cloud": "cloud"
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
				property_blacklist = [
					{
						property = "property1"
					},
					{
						property = "property2"
					}
				]
			`,
			APIUpdate: `{
				"yourInstance": "https://app.posthog.com",
				"teamApiKey": "api_key",
				"useV2Group": true,
				"connectionMode": {
					"web": "device",
					"flutter": "cloud",
					"cloud": "cloud"
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
				},
				"propertyBlacklist": {
					"web": [
						{
							"property": "property1"
						},
						{
							"property": "property2"
						}
					]
				}
			}`,
		},
	})
}
