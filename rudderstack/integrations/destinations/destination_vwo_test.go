package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceVWO(t *testing.T) {
	cmt.AssertDestination(t, "vwo", []c.TestConfig{
		{
			TerraformCreate: `
				account_id = "..."
			`,
			APICreate: `{
				"accountId": "..."
			}`,
			TerraformUpdate: `
				account_id = "..."

				is_spa                   = true
				send_experiment_track    = true
				send_experiment_identify = true
			
				library_tolerance  = "2000"
				settings_tolerance = "2000"
			
				use_existing_jquery = false
			
				use_native_sdk {
					web = true
				}
			
				event_filtering {
					whitelist = ["one", "two", "three"]
				}
				onetrust_cookie_categories {
					web = ["one", "two", "three"]
				}
			`,
			APIUpdate: `{
				"accountId": "...",
				"isSPA": true,
				"sendExperimentTrack": true,
				"sendExperimentIdentify": true,
				"libraryTolerance": "2000",
				"settingsTolerance": "2000",
				"whitelistedEvents": [
				  {
					"eventName": "one"
				  },
				  {
					"eventName": "two"
				  },
				  {
					"eventName": "three"
				  }
				],
				"eventFilteringOption": "whitelistedEvents",
				"useNativeSDK": {
				  "web": true
				},
				"oneTrustCookieCategories": {
					"web": [
						{ "oneTrustCookieCategory": "one" },
						{ "oneTrustCookieCategory": "two" },
						{ "oneTrustCookieCategory": "three" }
					]
				}
			}`,
		},
	})
}
