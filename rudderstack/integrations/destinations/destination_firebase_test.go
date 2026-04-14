package destinations_test

import (
	"testing"

	acc "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/acc"
	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

var firebaseTestConfigs = []c.TestConfig{
	{
		TerraformCreate: `
			connection_mode {
				android = "device"
				ios = "device"
			}
			event_filtering {
				whitelist = ["event1", "event2"]
			}
			`,
		APICreate: `{
				"connectionMode": {
					"android": "device",
					"ios": "device"
				},
				"whitelistedEvents": [
					{
						"eventName": "event1"
					},
					{
						"eventName": "event2"
					}
				],
				"eventFilteringOption": "whitelistedEvents"
			}`,
		TerraformUpdate: `
			connection_mode {
				android = "device"
				ios = "device"
				unity = "device"
				reactnative = "device"
				flutter = "device"
			}
			event_filtering {
				blacklist = ["event3", "event4", "event5"]
			}
			`,
		APIUpdate: `{
				"connectionMode": {
					"android": "device",
					"ios": "device",
					"unity": "device",
					"reactnative": "device",
					"flutter": "device"
				},
				"blacklistedEvents": [
					{
						"eventName": "event3"
					},
					{
						"eventName": "event4"
					},
					{
						"eventName": "event5"
					}
				],
				"eventFilteringOption": "blacklistedEvents"
			}`,
	},
}

func TestDestinationResourceFirebase(t *testing.T) {
	cmt.AssertDestination(t, "firebase", firebaseTestConfigs)
}

func TestAccDestinationFirebase(t *testing.T) {
	acc.AccAssertDestination(t, "firebase", firebaseTestConfigs)
}
