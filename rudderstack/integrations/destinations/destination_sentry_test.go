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
