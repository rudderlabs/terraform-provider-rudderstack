package destinations_test

import (
	"testing"

	acc "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/acc"
	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

var hsHubspotEventsTestConfigs = []c.TestConfig{
	{
		// Create: single event, no event_properties
		TerraformCreate: `
				authorization_type = "legacyApiKey"
				api_version        = "newApi"
				api_key            = "my-api-key"
				lookup_field       = "email"

				hubspot_events = [
					{
						rs_event_name      = "Order Completed"
						hubspot_event_name = "pe12345678_order_completed"
						event_properties   = []
					}
				]
			`,
		APICreate: `{
				"authorizationType": "legacyApiKey",
				"apiVersion": "newApi",
				"apiKey": "my-api-key",
				"lookupField": "email",
				"hubspotEvents": [
					{
						"rsEventName": "Order Completed",
						"hubspotEventName": "pe12345678_order_completed",
						"eventProperties": []
					}
				]
			}`,
		// Update: multiple events, some with event_properties and some without
		TerraformUpdate: `
				authorization_type = "legacyApiKey"
				api_version        = "newApi"
				api_key            = "my-api-key"
				lookup_field       = "email"

				hubspot_events = [
					{
						rs_event_name      = "Product Searched"
						hubspot_event_name = "pe12345678_search"
						event_properties = [
							{
								from = "query"
								to   = "hs_search_query"
							},
							{
								from = "results"
								to   = "hs_search_results"
							}
						]
					},
					{
						rs_event_name      = "Order Completed"
						hubspot_event_name = "pe12345678_order_completed"
						event_properties   = []
					}
				]
			`,
		APIUpdate: `{
				"authorizationType": "legacyApiKey",
				"apiVersion": "newApi",
				"apiKey": "my-api-key",
				"lookupField": "email",
				"hubspotEvents": [
					{
						"rsEventName": "Product Searched",
						"hubspotEventName": "pe12345678_search",
						"eventProperties": [
							{"from": "query", "to": "hs_search_query"},
							{"from": "results", "to": "hs_search_results"}
						]
					},
					{
						"rsEventName": "Order Completed",
						"hubspotEventName": "pe12345678_order_completed",
						"eventProperties": []
					}
				]
			}`,
	},
}

var hsTestConfigs = []c.TestConfig{
		{
			TerraformCreate: `
				authorization_type = "legacyApiKey"
				api_version        = "newApi"
				api_key            = "my-api-key"
				lookup_field       = "email"
			`,
			APICreate: `{
				"authorizationType": "legacyApiKey",
				"apiVersion": "newApi",
				"apiKey": "my-api-key",
				"lookupField": "email"
			}`,
			TerraformUpdate: `
				authorization_type = "newPrivateAppApi"
				api_version        = "newApi"
				access_token       = "my-access-token"
				hub_id             = "74X991"
				lookup_field       = "email"
				do_association     = true

				hubspot_events = [
					{
						rs_event_name      = "Product Searched"
						hubspot_event_name = "pe12345678_search"
						event_properties = [
							{
								from = "query"
								to   = "hs_search_query"
							},
							{
								from = "results"
								to   = "hs_search_results"
							}
						]
					},
					{
						rs_event_name      = "Order Completed"
						hubspot_event_name = "pe12345678_order_completed"
						event_properties   = []
					}
				]

				event_filtering {
					whitelist = ["Product Searched", "Order Completed"]
				}

				use_native_sdk {
					web = true
				}

				connection_mode {
					web = "device"
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
				"authorizationType": "newPrivateAppApi",
				"apiVersion": "newApi",
				"accessToken": "my-access-token",
				"hubID": "74X991",
				"lookupField": "email",
				"doAssociation": true,
				"hubspotEvents": [
					{
						"rsEventName": "Product Searched",
						"hubspotEventName": "pe12345678_search",
						"eventProperties": [
							{"from": "query", "to": "hs_search_query"},
							{"from": "results", "to": "hs_search_results"}
						]
					},
					{
						"rsEventName": "Order Completed",
						"hubspotEventName": "pe12345678_order_completed",
						"eventProperties": []
					}
				],
				"whitelistedEvents": [
					{"eventName": "Product Searched"},
					{"eventName": "Order Completed"}
				],
				"eventFilteringOption": "whitelistedEvents",
				"useNativeSDK": {"web": true},
				"connectionMode": {"web": "device"},
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

func TestDestinationResourceHsHubspotEvents(t *testing.T) {
	cmt.AssertDestination(t, "hs", hsHubspotEventsTestConfigs)
}

func TestDestinationResourceHs(t *testing.T) {
	cmt.AssertDestination(t, "hs", hsTestConfigs)
}

func TestAccDestinationHsHubspotEvents(t *testing.T) {
	acc.AccAssertDestination(t, "hs", hsHubspotEventsTestConfigs)
}

func TestAccDestinationHs(t *testing.T) {
	acc.AccAssertDestination(t, "hs", hsTestConfigs)
}
