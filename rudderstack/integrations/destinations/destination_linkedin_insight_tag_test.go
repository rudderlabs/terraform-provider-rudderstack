package destinations_test

import (
	"testing"

	acc "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/acc"
	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

var linkedinInsightTagTestConfigs = []c.TestConfig{
	{
		TerraformCreate: `
				partner_id = "p-id"
			`,
		APICreate: `{
				"partnerId": "p-id"
			}`,
		TerraformUpdate: `
				partner_id = "p-id"
				event_to_conversion_id_map = [
				{
					from = "a1"
					to   = "b1"
				}, 
				{
					from = "a2"
					to   = "b2"
				}]
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
		APIUpdate: `
			{
				"partnerId": "p-id",
				"eventToConversionIdMap": [
				  { "from": "a1", "to": "b1" },
				  { "from": "a2", "to": "b2" }
				],
				"useNativeSDK": {
					"web": true
				},
				"eventFilteringOption": "whitelistedEvents",
				"whitelistedEvents": [{
						"eventName": "one"
					},
					{
						"eventName": "two"
					},
					{
						"eventName": "three"
					}
				],
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
			}			
			`,
	},
}

func TestDestinationResourceLinkedinInsightTag(t *testing.T) {
	cmt.AssertDestination(t, "LINKEDIN_INSIGHT_TAG", linkedinInsightTagTestConfigs)
}

func TestAccDestinationLinkedinInsightTag(t *testing.T) {
	acc.AccAssertDestination(t, "LINKEDIN_INSIGHT_TAG", linkedinInsightTagTestConfigs)
}
