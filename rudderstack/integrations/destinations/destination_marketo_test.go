package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceMarketo(t *testing.T) {
	cmt.AssertDestination(t, "marketo", []c.TestConfig{
		{
			TerraformCreate: `
				account_id = "..."
				client_id = "cid"
				client_secret = "cs"
				track_anonymous_events = true
				create_if_not_exist = true
				connection_mode {
					web = "cloud"
					ios = "cloud"
				}
			`,
			APICreate: `{
				"accountId": "...",
				"clientId": "cid",
				"clientSecret": "cs",
				"trackAnonymousEvents": true,
				"createIfNotExist": true,
				"connectionMode": {
					"web": "cloud",
					"ios": "cloud"
				}
			}`,
			TerraformUpdate: `
				account_id = "..."
				client_id = "cid2"
				client_secret = "cs"
				track_anonymous_events = true
				create_if_not_exist = false
				lead_trait_mapping = [
					{
						from = "property0"
						to = "value0"
					}
				]
				rudder_events_mapping = [
					{
						event = "event0"
						marketo_primarykey = "marketoPrimarykey0"
						marketo_activity_id = "marketoActivityId0"
					}
				]
				custom_activity_property_map = [
					{
						from = "property1"
						to = "value1"
					}
				]
				connection_mode {
					web = "cloud"
					ios = "cloud"
				}
				onetrust_cookie_categories {
					web = ["one", "two", "three"]
					android = ["one", "two", "three"]
					ios = ["one", "two", "three"]
					unity = ["one", "two", "three"]
					reactnative = ["one", "two", "three"]
					flutter = ["one", "two", "three"]
					cordova = ["one", "two", "three"]
					amp = ["one", "two", "three"]
					cloud = ["one", "two", "three"]
					warehouse = ["one", "two", "three"]
					shopify = ["one", "two", "three"]
				}
			`,
			APIUpdate: `{
				"accountId": "...",
				"clientId": "cid2",
				"clientSecret": "cs",
				"trackAnonymousEvents": true,
				"createIfNotExist": false,
				"leadTraitMapping": [
					{
						"from": "property0",
						"to": "value0"
					}
				],
				"rudderEventsMapping": [
					{
						"event": "event0",
						"marketoPrimarykey": "marketoPrimarykey0",
						"marketoActivityId": "marketoActivityId0"
					}
				],
				"customActivityPropertyMap": [
					{
						"from": "property1",
						"to": "value1"
					}
				],
				"connectionMode": {
					"web": "cloud",
					"ios": "cloud"
				},
				"oneTrustCookieCategories": {
					"web": [
						{ "oneTrustCookieCategory": "one" },
						{ "oneTrustCookieCategory": "two" },
						{ "oneTrustCookieCategory": "three" }
					],
					"android": [
						{ "oneTrustCookieCategory": "one" },
						{ "oneTrustCookieCategory": "two" },
						{ "oneTrustCookieCategory": "three" }
					],
					"ios": [
						{ "oneTrustCookieCategory": "one" },
						{ "oneTrustCookieCategory": "two" },
						{ "oneTrustCookieCategory": "three" }
					],
					"unity": [
						{ "oneTrustCookieCategory": "one" },
						{ "oneTrustCookieCategory": "two" },
						{ "oneTrustCookieCategory": "three" }
					],
					"reactnative": [
						{ "oneTrustCookieCategory": "one" },
						{ "oneTrustCookieCategory": "two" },
						{ "oneTrustCookieCategory": "three" }
					],
					"flutter": [
						{ "oneTrustCookieCategory": "one" },
						{ "oneTrustCookieCategory": "two" },
						{ "oneTrustCookieCategory": "three" }
					],
					"cordova": [
						{ "oneTrustCookieCategory": "one" },
						{ "oneTrustCookieCategory": "two" },
						{ "oneTrustCookieCategory": "three" }
					],
					"amp": [
						{ "oneTrustCookieCategory": "one" },
						{ "oneTrustCookieCategory": "two" },
						{ "oneTrustCookieCategory": "three" }
					],
					"cloud": [
						{ "oneTrustCookieCategory": "one" },
						{ "oneTrustCookieCategory": "two" },
						{ "oneTrustCookieCategory": "three" }
					],
					"warehouse": [
						{ "oneTrustCookieCategory": "one" },
						{ "oneTrustCookieCategory": "two" },
						{ "oneTrustCookieCategory": "three" }
					],
					"shopify": [
						{ "oneTrustCookieCategory": "one" },
						{ "oneTrustCookieCategory": "two" },
						{ "oneTrustCookieCategory": "three" }
					]
				}
			}`,
		},
	})
}
