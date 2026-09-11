package destinations_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil"
	acc "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/acc"
	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

var spotifyPixelTestConfigs = []c.TestConfig{
	{
		TerraformCreate: `
					pixel_id = "spotify-pixel-id"
				`,
		APICreate: `{
					"pixelId": "spotify-pixel-id"
				}`,
		TerraformUpdate: `
					pixel_id = "spotify-pixel-id-updated"

					enable_alias_call = true

					events_to_spotify_pixel_events = [{
						from = "Lead Captured"
						to   = "lead"
					}, {
						from = "Product Added"
						to   = "addtocart"
					}, {
						from = "Order Completed"
						to   = "purchase"
					}]

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
					"pixelId": "spotify-pixel-id-updated",
					"enableAliasCall": true,
					"eventsToSpotifyPixelEvents": [
						{ "from": "Lead Captured", "to": "lead" },
						{ "from": "Product Added", "to": "addtocart" },
						{ "from": "Order Completed", "to": "purchase" }
					],
					"whitelistedEvents": [
						{ "eventName": "one" },
						{ "eventName": "two" },
						{ "eventName": "three" }
					],
					"eventFilteringOption": "whitelistedEvents",
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
						]
					}
				}`,
	},
}

func TestDestinationResourceSpotifyPixel(t *testing.T) {
	cmt.AssertDestination(t, "spotify_pixel", spotifyPixelTestConfigs)
}

func TestDestinationResourceSpotifyPixelEventFilteringBranches(t *testing.T) {
	cm := c.Destinations.Entries()["spotify_pixel"]
	tests := []struct {
		name  string
		state string
		want  string
	}{
		{
			name: "whitelist",
			state: `{
				"pixel_id": "spotify-pixel-id",
				"event_filtering": [{"whitelist": ["one", "two"]}]
			}`,
			want: `{
				"pixelId": "spotify-pixel-id",
				"whitelistedEvents": [{"eventName": "one"}, {"eventName": "two"}],
				"eventFilteringOption": "whitelistedEvents"
			}`,
		},
		{
			name: "blacklist",
			state: `{
				"pixel_id": "spotify-pixel-id",
				"event_filtering": [{"blacklist": ["one", "two"]}]
			}`,
			want: `{
				"pixelId": "spotify-pixel-id",
				"blacklistedEvents": [{"eventName": "one"}, {"eventName": "two"}],
				"eventFilteringOption": "blacklistedEvents"
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cm.StateToAPI(tt.state)
			if err != nil {
				t.Fatalf("StateToAPI failed: %v", err)
			}

			if !testutil.JSONEq(got, tt.want) {
				t.Fatalf("API config mismatch\ngot:  %s\nwant: %s", got, tt.want)
			}
		})
	}
}

func TestDestinationResourceSpotifyPixelDoesNotExposeNativeSDKToggle(t *testing.T) {
	cm := c.Destinations.Entries()["spotify_pixel"]
	if _, ok := cm.ConfigSchema["use_native_sdk"]; ok {
		t.Fatal("spotify_pixel must not expose use_native_sdk because Spotify Pixel is web device-mode only")
	}

	states := map[string]string{
		"omitted": `{
			"pixel_id": "spotify-pixel-id"
		}`,
		"false": `{
			"pixel_id": "spotify-pixel-id",
			"use_native_sdk": [{"web": false}]
		}`,
	}

	for name, state := range states {
		t.Run(name, func(t *testing.T) {
			got, err := cm.StateToAPI(state)
			if err != nil {
				t.Fatalf("StateToAPI failed: %v", err)
			}
			if strings.Contains(got, "useNativeSDK") {
				t.Fatalf("expected use_native_sdk to stay out of API config, got: %s", got)
			}
		})
	}
}

func TestDestinationResourceSpotifyPixelValidation(t *testing.T) {
	configSchema := c.Destinations.Entries()["spotify_pixel"].ConfigSchema

	pixelIDSchema := configSchema["pixel_id"]
	if pixelIDSchema.Sensitive {
		t.Fatal("expected pixel_id not to be sensitive because pixelId is not listed in Spotify Pixel secretKeys")
	}
	if diags := pixelIDSchema.ValidateDiagFunc("", cty.Path{}); !diags.HasError() {
		t.Fatal("expected empty pixel_id to fail validation")
	}
	if diags := pixelIDSchema.ValidateDiagFunc(strings.Repeat("a", 101), cty.Path{}); !diags.HasError() {
		t.Fatal("expected pixel_id over 100 characters to fail validation")
	}
	if diags := pixelIDSchema.ValidateDiagFunc(strings.Repeat("a", 100), cty.Path{}); diags.HasError() {
		t.Fatalf("expected 100-character pixel_id to pass validation: %v", diags)
	}

	mappingSchema := configSchema["events_to_spotify_pixel_events"].Elem.(*schema.Resource)
	toSchema := mappingSchema.Schema["to"]
	if diags := toSchema.ValidateDiagFunc("custom", cty.Path{}); !diags.HasError() {
		t.Fatal("expected custom Spotify event mapping target to fail validation")
	}
	if diags := toSchema.ValidateDiagFunc("purchase", cty.Path{}); diags.HasError() {
		t.Fatalf("expected valid Spotify event mapping target to pass validation: %v", diags)
	}
}

func TestAccDestinationSpotifyPixel(t *testing.T) {
	acc.AccAssertDestination(t, "spotify_pixel", spotifyPixelTestConfigs)
}
