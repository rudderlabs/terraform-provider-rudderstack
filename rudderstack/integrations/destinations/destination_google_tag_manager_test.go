package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceGoogleTagManager(t *testing.T) {
	cmt.AssertDestination(t, "google_tag_manager", []c.TestConfig{
		{
			TerraformCreate: `
				container_id = "GTM-000000"
			`,
			APICreate: `{
				"containerID": "GTM-000000"
			}`,
			TerraformUpdate: `
				container_id = "GTM-000000"

				server_url = "https://example.com"
			
				use_native_sdk {
					web = true
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
				}
			`,
			APIUpdate: `{
				"containerID": "GTM-000000",
				"serverUrl": "https://example.com",
				"blacklistedEvents": [
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
				"eventFilteringOption": "blacklistedEvents",
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
