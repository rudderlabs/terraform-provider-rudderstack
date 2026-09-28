package destinations_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	acc "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/acc"
	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

var googleAdsTestConfigs = []c.TestConfig{
	{
		TerraformCreate: `
				conversion_id = "AW-00000000"
			`,
		APICreate: `{
				"conversionID": "AW-00000000",
				"conversionLinker": true,
				"sendPageView": true
			}`,
		TerraformUpdate: `
				conversion_id = "AW-00000000"

				default_page_conversion = "..."
			
				page_load_conversions = [
					{
						"label" = "..."
						"name"  = "..."
					}
				]
			
				click_event_conversions = [
					{
						"label" = "..."
						"name"  = "..."
					}
				]
			
				dynamic_remarketing {
					web = true
				}
			
				conversion_linker          = true
				send_page_view             = true
				disable_ad_personalization = true
			
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
				"conversionID": "AW-00000000",
				"pageLoadConversions": [
				  {
					"conversionLabel": "...",
					"name": "..."
				  }
				],
				"clickEventConversions": [
				  {
					"conversionLabel": "...",
					"name": "..."
				  }
				],
				"defaultPageConversion": "...",
				"dynamicRemarketing": {
				  "web": true
				},
				"conversionLinker": true,
				"sendPageView": true,
				"disableAdPersonalization": true,
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
}

func TestDestinationResourceGoogleAds(t *testing.T) {
	cmt.AssertDestination(t, "google_ads", googleAdsTestConfigs)
}

func TestDestinationResourceGoogleAdsConnectionModeOmittedWhenUnset(t *testing.T) {
	cm := c.Destinations.Entries()["google_ads"]

	got, err := cm.StateToAPI(`{"conversion_id": "AW-00000000"}`)
	if err != nil {
		t.Fatalf("StateToAPI failed: %v", err)
	}
	if strings.Contains(got, "connectionMode") {
		t.Fatalf("expected connectionMode to stay out of API config when unset, got: %s", got)
	}
}

func TestDestinationResourceGoogleAdsConnectionModeValidation(t *testing.T) {
	configSchema := c.Destinations.Entries()["google_ads"].ConfigSchema
	connectionModeSchema := configSchema["connection_mode"].Elem.(*schema.Resource)
	webModeSchema := connectionModeSchema.Schema["web"]

	if diags := webModeSchema.ValidateDiagFunc("cloud", cty.Path{}); !diags.HasError() {
		t.Fatal("expected cloud connection mode to fail validation because Google Ads only supports device mode")
	}
	if diags := webModeSchema.ValidateDiagFunc("device", cty.Path{}); diags.HasError() {
		t.Fatalf("expected device connection mode to pass validation: %v", diags)
	}
}

func TestAccDestinationGoogleAds(t *testing.T) {
	acc.AccAssertDestination(t, "google_ads", googleAdsTestConfigs)
}
