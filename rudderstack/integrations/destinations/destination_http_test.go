package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceHTTP(t *testing.T) {
	cmt.AssertDestination(t, "http", []c.TestConfig{
		{
			TerraformCreate: `
				api_url  = "https://example.com/base"
				auth     = "basicAuth"
				username = "myuser"
				password = "mypass"
			`,
			APICreate: `{
				"apiUrl": "https://example.com/base",
				"auth": "basicAuth",
				"username": "myuser",
				"password": "mypass",
				"method": "POST",
				"format": "JSON",
				"isBatchingEnabled": false,
				"isDefaultMapping": true
			}`,
			TerraformUpdate: `
				api_url  = "https://example.com/base"
				auth     = "basicAuth"
				username = "myuser"
				password = "mypass"
			`,
			APIUpdate: `{
				"apiUrl": "https://example.com/base",
				"auth": "basicAuth",
				"username": "myuser",
				"password": "mypass",
				"method": "POST",
				"format": "JSON",
				"isBatchingEnabled": false,
				"isDefaultMapping": true
			}`,
		},
		{
			TerraformCreate: `
				api_url = "https://example.com/base"
			`,
			APICreate: `{
				"apiUrl": "https://example.com/base",
				"auth": "noAuth",
				"method": "POST",
				"format": "JSON",
				"isBatchingEnabled": false,
				"isDefaultMapping": true
			}`,
			TerraformUpdate: `
				api_url       = "https://example.com/base"
				auth          = "apiKeyAuth"
				api_key_name  = "x-api-key"
				api_key_value = "secret-api-key"
				method        = "PATCH"
				format        = "XML"
				xml_root_key  = "events"
				is_default_mapping = false
				is_batching_enabled = true
				max_batch_size      = "10"

				path_params = [
					{
						path = "users"
					},
					{
						path = "$.userId"
					}
				]

				query_params = [
					{
						to   = "customerId"
						from = "$.userId"
					},
					{
						to   = "env"
						from = "prod"
					}
				]

				headers = [
					{
						to   = "content-type"
						from = "application/xml"
					},
					{
						to   = "x-source"
						from = "$.type"
					}
				]

				properties_mapping = [
					{
						to   = "$.events.body"
						from = "$.properties.payload"
					},
					{
						to   = "$.events.kind"
						from = "track"
					}
				]

				event_filtering {
					blacklist = ["one", "two", "three"]
				}

				connection_mode {
					android        = "cloud"
					android_kotlin = "cloud"
					ios            = "cloud"
					ios_swift      = "cloud"
					web            = "cloud"
					unity          = "cloud"
					amp            = "cloud"
					cloud          = "cloud"
					warehouse      = "cloud"
					reactnative    = "cloud"
					flutter        = "cloud"
					cordova        = "cloud"
					shopify        = "cloud"
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
					shopify = [{
						provider = "custom"
						resolution_strategy = "and"
						consents = ["one_shopify", "two_shopify", "three_shopify"]
					}]
				}
			`,
			APIUpdate: `{
				"apiUrl": "https://example.com/base",
				"auth": "apiKeyAuth",
				"apiKeyName": "x-api-key",
				"apiKeyValue": "secret-api-key",
				"method": "PATCH",
				"format": "XML",
				"xmlRootKey": "events",
				"pathParams": [
					{ "path": "users" },
					{ "path": "$.userId" }
				],
				"queryParams": [
					{ "to": "customerId", "from": "$.userId" },
					{ "to": "env", "from": "prod" }
				],
				"headers": [
					{ "to": "content-type", "from": "application/xml" },
					{ "to": "x-source", "from": "$.type" }
				],
				"propertiesMapping": [
					{ "to": "$.events.body", "from": "$.properties.payload" },
					{ "to": "$.events.kind", "from": "track" }
				],
				"isBatchingEnabled": true,
				"maxBatchSize": "10",
				"blacklistedEvents": [
					{ "eventName": "one" },
					{ "eventName": "two" },
					{ "eventName": "three" }
				],
				"eventFilteringOption": "blacklistedEvents",
				"isDefaultMapping": false,
				"connectionMode": {
					"android": "cloud",
					"androidKotlin": "cloud",
					"ios": "cloud",
					"iosSwift": "cloud",
					"web": "cloud",
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
					],
					"android": [
						{
							"provider": "ketch",
							"resolutionStrategy": "",
							"consents": [
								{ "consent": "one_android" },
								{ "consent": "two_android" },
								{ "consent": "three_android" }
							]
						}
					],
					"androidKotlin": [
						{
							"provider": "ketch",
							"resolutionStrategy": "",
							"consents": [
								{ "consent": "one_android_kotlin" },
								{ "consent": "two_android_kotlin" },
								{ "consent": "three_android_kotlin" }
							]
						}
					],
					"ios": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{ "consent": "one_ios" },
								{ "consent": "two_ios" },
								{ "consent": "three_ios" }
							]
						}
					],
					"iosSwift": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{ "consent": "one_ios_swift" },
								{ "consent": "two_ios_swift" },
								{ "consent": "three_ios_swift" }
							]
						}
					],
					"unity": [
						{
							"provider": "custom",
							"resolutionStrategy": "or",
							"consents": [
								{ "consent": "one_unity" },
								{ "consent": "two_unity" },
								{ "consent": "three_unity" }
							]
						}
					],
					"reactnative": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{ "consent": "one_reactnative" },
								{ "consent": "two_reactnative" },
								{ "consent": "three_reactnative" }
							]
						}
					],
					"flutter": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{ "consent": "one_flutter" },
								{ "consent": "two_flutter" },
								{ "consent": "three_flutter" }
							]
						}
					],
					"cordova": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{ "consent": "one_cordova" },
								{ "consent": "two_cordova" },
								{ "consent": "three_cordova" }
							]
						}
					],
					"amp": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{ "consent": "one_amp" },
								{ "consent": "two_amp" },
								{ "consent": "three_amp" }
							]
						}
					],
					"cloud": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{ "consent": "one_cloud" },
								{ "consent": "two_cloud" },
								{ "consent": "three_cloud" }
							]
						}
					],
					"warehouse": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{ "consent": "one_warehouse" },
								{ "consent": "two_warehouse" },
								{ "consent": "three_warehouse" }
							]
						}
					],
					"shopify": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{ "consent": "one_shopify" },
								{ "consent": "two_shopify" },
								{ "consent": "three_shopify" }
							]
						}
					]
				}
			}`,
		},
	})
}
