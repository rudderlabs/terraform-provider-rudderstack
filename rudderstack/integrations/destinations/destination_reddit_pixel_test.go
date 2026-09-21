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

var redditPixelTestConfigs = []c.TestConfig{
	{
		TerraformCreate: `advertiser_id = "reddit-pixel-id"`,
		APICreate:       `{"advertiserId":"reddit-pixel-id"}`,
		TerraformUpdate: `
			advertiser_id = "reddit-pixel-id-updated"
			event_mapping_from_config = [
				{ from = "Product Viewed", to = "ViewContent" },
				{ from = "Product Added", to = "AddToCart" },
				{ from = "Order Completed", to = "Purchase" }
			]
			event_filtering { whitelist = ["one", "two", "three"] }
			use_native_sdk { web = false }
			connection_mode { web = "device" }
			consent_management {
				web = [
					{ provider = "oneTrust", consents = ["one_web", "two_web", "three_web"], resolution_strategy = "" },
					{ provider = "ketch", consents = ["one_web", "two_web", "three_web"], resolution_strategy = "" },
					{ provider = "custom", consents = ["one_web", "two_web", "three_web"], resolution_strategy = "and" }
				]
			}
		`,
		APIUpdate: `{
			"advertiserId":"reddit-pixel-id-updated",
			"eventMappingFromConfig":[
				{"from":"Product Viewed","to":"ViewContent"},
				{"from":"Product Added","to":"AddToCart"},
				{"from":"Order Completed","to":"Purchase"}
			],
			"whitelistedEvents":[{"eventName":"one"},{"eventName":"two"},{"eventName":"three"}],
			"eventFilteringOption":"whitelistedEvents",
			"useNativeSDK":{"web":false},
			"connectionMode":{"web":"device"},
			"consentManagement":{"web":[
				{"provider":"oneTrust","resolutionStrategy":"","consents":[{"consent":"one_web"},{"consent":"two_web"},{"consent":"three_web"}]},
				{"provider":"ketch","resolutionStrategy":"","consents":[{"consent":"one_web"},{"consent":"two_web"},{"consent":"three_web"}]},
				{"provider":"custom","resolutionStrategy":"and","consents":[{"consent":"one_web"},{"consent":"two_web"},{"consent":"three_web"}]}
			]}
		}`,
	},
}

func TestDestinationResourceRedditPixel(t *testing.T) {
	cmt.AssertDestination(t, "reddit_pixel", redditPixelTestConfigs)
}

func TestDestinationResourceRedditPixelEventFilteringBranches(t *testing.T) {
	cm := c.Destinations.Entries()["reddit_pixel"]
	tests := []struct {
		name  string
		state string
		want  string
	}{
		{
			name: "whitelist",
			state: `{
				"advertiser_id": "reddit-pixel-id",
				"event_filtering": [{"whitelist": ["one", "two"]}]
			}`,
			want: `{
				"advertiserId": "reddit-pixel-id",
				"whitelistedEvents": [{"eventName": "one"}, {"eventName": "two"}],
				"eventFilteringOption": "whitelistedEvents"
			}`,
		},
		{
			name: "blacklist",
			state: `{
				"advertiser_id": "reddit-pixel-id",
				"event_filtering": [{"blacklist": ["one", "two"]}]
			}`,
			want: `{
				"advertiserId": "reddit-pixel-id",
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

func TestDestinationResourceRedditPixelNativeSDKFalseIsPreserved(t *testing.T) {
	cm := c.Destinations.Entries()["reddit_pixel"]
	got, err := cm.StateToAPI(`{
		"advertiser_id": "reddit-pixel-id",
		"use_native_sdk": [{"web": false}]
	}`)
	if err != nil {
		t.Fatalf("StateToAPI failed: %v", err)
	}
	if !testutil.JSONEq(got, `{
		"advertiserId": "reddit-pixel-id",
		"useNativeSDK": {"web": false}
	}`) {
		t.Fatalf("expected explicit use_native_sdk.web=false in API config, got: %s", got)
	}
}

func TestDestinationResourceRedditPixelConnectionModeOmittedWhenUnset(t *testing.T) {
	// connection_mode is Optional with no provider-injected default, matching the
	// c.Simple + c.SkipZeroValue convention every other destination in this
	// provider uses for connectionMode.{sourceType}. If the user never sets it,
	// it is left out of the API payload entirely rather than the provider
	// assuming "device" on their behalf.
	cm := c.Destinations.Entries()["reddit_pixel"]

	got, err := cm.StateToAPI(`{"advertiser_id": "reddit-pixel-id"}`)
	if err != nil {
		t.Fatalf("StateToAPI failed: %v", err)
	}
	if strings.Contains(got, "connectionMode") {
		t.Fatalf("expected connectionMode to stay out of API config when unset, got: %s", got)
	}
}

func TestDestinationResourceRedditPixelValidation(t *testing.T) {
	configSchema := c.Destinations.Entries()["reddit_pixel"].ConfigSchema

	advertiserIDSchema := configSchema["advertiser_id"]
	if advertiserIDSchema.Sensitive {
		t.Fatal("expected advertiser_id not to be sensitive because advertiserId is not listed in Reddit Pixel secretKeys")
	}
	if diags := advertiserIDSchema.ValidateDiagFunc("", cty.Path{}); diags.HasError() {
		t.Fatalf("expected empty advertiser_id to pass the upstream length validator: %v", diags)
	}
	if diags := advertiserIDSchema.ValidateDiagFunc(strings.Repeat("a", 101), cty.Path{}); !diags.HasError() {
		t.Fatal("expected advertiser_id over 100 characters to fail validation")
	}
	if diags := advertiserIDSchema.ValidateDiagFunc(strings.Repeat("a", 100), cty.Path{}); diags.HasError() {
		t.Fatalf("expected 100-character advertiser_id to pass validation: %v", diags)
	}

	mappingSchema := configSchema["event_mapping_from_config"].Elem.(*schema.Resource)
	toSchema := mappingSchema.Schema["to"]
	if diags := toSchema.ValidateDiagFunc("PageVisit", cty.Path{}); !diags.HasError() {
		t.Fatal("expected PageVisit Reddit Pixel event mapping target to fail validation")
	}
	if diags := toSchema.ValidateDiagFunc("Purchase", cty.Path{}); diags.HasError() {
		t.Fatalf("expected valid Reddit Pixel event mapping target to pass validation: %v", diags)
	}
	if diags := toSchema.ValidateDiagFunc("", cty.Path{}); diags.HasError() {
		t.Fatalf("expected empty Reddit Pixel event mapping target to pass validation: %v", diags)
	}
	fromSchema := mappingSchema.Schema["from"]
	if diags := fromSchema.ValidateDiagFunc("{{ event || fallback }}", cty.Path{}); diags.HasError() {
		t.Fatalf("expected dynamic event mapping source to pass validation: %v", diags)
	}

	connectionModeSchema := configSchema["connection_mode"].Elem.(*schema.Resource)
	webModeSchema := connectionModeSchema.Schema["web"]
	if diags := webModeSchema.ValidateDiagFunc("cloud", cty.Path{}); !diags.HasError() {
		t.Fatal("expected cloud connection mode to fail validation because Reddit Pixel only supports device mode")
	}
	if diags := webModeSchema.ValidateDiagFunc("device", cty.Path{}); diags.HasError() {
		t.Fatalf("expected device connection mode to pass validation: %v", diags)
	}
}

func TestAccDestinationRedditPixel(t *testing.T) {
	acc.AccAssertDestination(t, "reddit_pixel", redditPixelTestConfigs)
}
