package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceSentry(t *testing.T) {
	cmt.AssertDestination(t, "sentry", []c.TestConfig{
		{
			TerraformCreate: `
				dsn = "https://some-url"
			`,
			APICreate: `{
				"dsn": "https://some-url",
				"debugMode": false
			}`,
			TerraformUpdate: `
				dsn = "https://some-url"
				server_name             = "..."
				release                 = "..."
				environment             = "..."
				custom_version_property = "..."
				logger                  = "..."
				debug_mode              = true
			
				ignore_errors = ["one", "two", "three"]
				include_paths = ["one", "two", "three"]
				allow_urls    = ["one", "two", "three"]
				deny_urls     = ["one", "two", "three"]
			
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
				"dsn": "https://some-url",
				"serverName": "...",
				"release": "...",
				"environment": "...",
				"customVersionProperty": "...",
				"logger": "...",
				"debugMode": true,
				"ignoreErrors": [
					{ "ignoreErrors": "one" },
					{ "ignoreErrors": "two" },
					{ "ignoreErrors": "three" }
				],
				"includePaths": [
					{ "includePaths": "one" },
					{ "includePaths": "two" },
					{ "includePaths": "three" }
				],
				"allowUrls": [
					{ "allowUrls": "one" },
					{ "allowUrls": "two" },
					{ "allowUrls": "three" }
				],
				"denyUrls": [
					{ "denyUrls": "one" },
					{ "denyUrls": "two" },
					{ "denyUrls": "three" }
				],
				"useNativeSDK": {
					"web": true
				},
				"whitelistedEvents": [
					{ "eventName": "one" },
					{ "eventName": "two" },
					{ "eventName": "three" }
				],
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
