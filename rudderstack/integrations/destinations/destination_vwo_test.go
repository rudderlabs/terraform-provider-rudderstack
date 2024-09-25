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
					]
				}
			}`,
		},
	})
}
