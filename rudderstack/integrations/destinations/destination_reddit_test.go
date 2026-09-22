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

var redditTestConfigs = []c.TestConfig{
	{
		TerraformCreate: `
					rudder_account_id = "__ACCOUNT_ID__"
					account_id        = "reddit-pixel-id"
				`,
		APICreate: `{
					"rudderAccountId": "__ACCOUNT_ID__",
					"accountId": "reddit-pixel-id",
					"version": "v3",
					"hashData": true
				}`,
		TerraformUpdate: `
					rudder_account_id = "__ACCOUNT_ID__"
					account_id        = "reddit-pixel-id-updated"
					version           = "v2"
					hash_data         = false
					events_mapping = [
						{ from = "Order Completed", to = "Purchase" },
						{ from = "Order Completed", to = "Lead" }
					]
				connection_mode {
					web           = "cloud"
					android       = "cloud"
					android_kotlin = "cloud"
					ios           = "cloud"
					ios_swift     = "cloud"
					unity         = "cloud"
					amp           = "cloud"
					cloud         = "cloud"
					warehouse     = "cloud"
					reactnative   = "cloud"
					flutter       = "cloud"
					cordova       = "cloud"
					shopify       = "cloud"
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
					android = [{
						provider = "ketch"
						consents = ["one_android", "two_android", "three_android"]
						resolution_strategy = ""
					}]
					android_kotlin = [{
						provider = "ketch"
						consents = ["one_android_kotlin", "two_android_kotlin", "three_android_kotlin"]
						resolution_strategy = ""
					}]
					ios = [{
						provider = "custom"
						resolution_strategy = "and"
						consents = ["one_ios", "two_ios", "three_ios"]
					}]
					ios_swift = [{
						provider = "custom"
						resolution_strategy = "and"
						consents = ["one_ios_swift", "two_ios_swift", "three_ios_swift"]
					}]
					unity = [{
						provider = "custom"
						resolution_strategy = "or"
						consents = ["one_unity", "two_unity", "three_unity"]
					}]
					amp = [{
						provider = "custom"
						resolution_strategy = "and"
						consents = ["one_amp", "two_amp", "three_amp"]
					}]
					cloud = [{
						provider = "custom"
						resolution_strategy = "and"
						consents = ["one_cloud", "two_cloud", "three_cloud"]
					}]
					warehouse = [{
						provider = "custom"
						resolution_strategy = "and"
						consents = ["one_warehouse", "two_warehouse", "three_warehouse"]
					}]
					reactnative = [{
						provider = "custom"
						resolution_strategy = "and"
						consents = ["one_reactnative", "two_reactnative", "three_reactnative"]
					}]
					flutter = [{
						provider = "custom"
						resolution_strategy = "and"
						consents = ["one_flutter", "two_flutter", "three_flutter"]
					}]
					cordova = [{
						provider = "custom"
						resolution_strategy = "and"
						consents = ["one_cordova", "two_cordova", "three_cordova"]
					}]
					shopify = [{
						provider = "custom"
						resolution_strategy = "and"
						consents = ["one_shopify", "two_shopify", "three_shopify"]
					}]
				}
			`,
		APIUpdate: `{
					"rudderAccountId": "__ACCOUNT_ID__",
					"accountId": "reddit-pixel-id-updated",
					"version": "v2",
					"hashData": false,
					"eventsMapping": [
						{ "from": "Order Completed", "to": "Purchase" },
						{ "from": "Order Completed", "to": "Lead" }
					],
				"connectionMode": {
					"web": "cloud",
					"android": "cloud",
					"androidKotlin": "cloud",
					"ios": "cloud",
					"iosSwift": "cloud",
					"unity": "cloud",
					"amp": "cloud",
					"cloud": "cloud",
					"warehouse": "cloud",
					"reactnative": "cloud",
					"flutter": "cloud",
					"cordova": "cloud",
					"shopify": "cloud"
				},
				"consentManagement": {
					"web": [
						{
							"provider": "oneTrust",
							"resolutionStrategy": "",
							"consents": [
								{"consent": "one_web"},
								{"consent": "two_web"},
								{"consent": "three_web"}
							]
						},
						{
							"provider": "ketch",
							"resolutionStrategy": "",
							"consents": [
								{"consent": "one_web"},
								{"consent": "two_web"},
								{"consent": "three_web"}
							]
						},
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{"consent": "one_web"},
								{"consent": "two_web"},
								{"consent": "three_web"}
							]
						}
					],
					"android": [
						{
							"provider": "ketch",
							"resolutionStrategy": "",
							"consents": [
								{"consent": "one_android"},
								{"consent": "two_android"},
								{"consent": "three_android"}
							]
						}
					],
					"androidKotlin": [
						{
							"provider": "ketch",
							"resolutionStrategy": "",
							"consents": [
								{"consent": "one_android_kotlin"},
								{"consent": "two_android_kotlin"},
								{"consent": "three_android_kotlin"}
							]
						}
					],
					"ios": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{"consent": "one_ios"},
								{"consent": "two_ios"},
								{"consent": "three_ios"}
							]
						}
					],
					"iosSwift": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{"consent": "one_ios_swift"},
								{"consent": "two_ios_swift"},
								{"consent": "three_ios_swift"}
							]
						}
					],
					"unity": [
						{
							"provider": "custom",
							"resolutionStrategy": "or",
							"consents": [
								{"consent": "one_unity"},
								{"consent": "two_unity"},
								{"consent": "three_unity"}
							]
						}
					],
					"amp": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{"consent": "one_amp"},
								{"consent": "two_amp"},
								{"consent": "three_amp"}
							]
						}
					],
					"cloud": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{"consent": "one_cloud"},
								{"consent": "two_cloud"},
								{"consent": "three_cloud"}
							]
						}
					],
					"warehouse": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{"consent": "one_warehouse"},
								{"consent": "two_warehouse"},
								{"consent": "three_warehouse"}
							]
						}
					],
					"reactnative": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{"consent": "one_reactnative"},
								{"consent": "two_reactnative"},
								{"consent": "three_reactnative"}
							]
						}
					],
					"flutter": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{"consent": "one_flutter"},
								{"consent": "two_flutter"},
								{"consent": "three_flutter"}
							]
						}
					],
					"cordova": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{"consent": "one_cordova"},
								{"consent": "two_cordova"},
								{"consent": "three_cordova"}
							]
						}
					],
					"shopify": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{"consent": "one_shopify"},
								{"consent": "two_shopify"},
								{"consent": "three_shopify"}
							]
						}
					]
				}
			}`,
	},
}

func TestDestinationResourceReddit(t *testing.T) {
	cmt.AssertDestination(t, "reddit", redditTestConfigs)
}

func TestDestinationResourceRedditDefaultsAndZeroValues(t *testing.T) {
	cm := c.Destinations.Entries()["reddit"]
	got, err := cm.StateToAPI(`{"rudder_account_id":"oauth-account-id","account_id":"reddit-pixel-id","version":"v3","hash_data":false}`)
	if err != nil {
		t.Fatalf("StateToAPI failed: %v", err)
	}
	want := `{"rudderAccountId":"oauth-account-id","accountId":"reddit-pixel-id","version":"v3","hashData":false}`
	if !testutil.JSONEq(got, want) {
		t.Fatalf("API config mismatch\ngot: %s\nwant: %s", got, want)
	}
	if strings.Contains(got, "connectionMode") {
		t.Fatalf("expected connectionMode omitted when unset, got: %s", got)
	}
}

func TestDestinationResourceRedditValidation(t *testing.T) {
	configSchema := c.Destinations.Entries()["reddit"].ConfigSchema
	if configSchema["rudder_account_id"].Sensitive || configSchema["account_id"].Sensitive {
		t.Fatal("expected no sensitive Reddit fields because secretKeys is empty")
	}
	if _, ok := configSchema["use_native_sdk"]; ok {
		t.Fatal("reddit must not expose schema-only useNativeSDK")
	}
	accountIDSchema := configSchema["account_id"]
	if diags := accountIDSchema.ValidateDiagFunc("", cty.Path{}); !diags.HasError() {
		t.Fatal("expected empty account_id to fail")
	}
	if diags := accountIDSchema.ValidateDiagFunc("{{ pixel || fallback }}", cty.Path{}); diags.HasError() {
		t.Fatalf("expected dynamic account_id to pass: %v", diags)
	}
	mappingSchema := configSchema["events_mapping"].Elem.(*schema.Resource)
	if diags := mappingSchema.Schema["to"].ValidateDiagFunc("PageVisit", cty.Path{}); diags.HasError() {
		t.Fatalf("expected PageVisit to pass: %v", diags)
	}
	if diags := mappingSchema.Schema["to"].ValidateDiagFunc("Custom", cty.Path{}); !diags.HasError() {
		t.Fatal("expected unsupported target to fail")
	}
}

func TestAccDestinationReddit(t *testing.T) {
	acc.AccAssertOAuthDestination(t, "reddit", redditTestConfigs)
}
