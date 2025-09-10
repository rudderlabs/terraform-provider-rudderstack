package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceMixpanel(t *testing.T) {
	cmt.AssertDestination(t, "mixpanel", []c.TestConfig{
		{
			TerraformCreate: `
				token = "..."
				identity_merge_api = "simplified"
				data_residency = "us"
				persistence = "none"
				consolidated_page_calls = false
				connection_modes {
					web = "cloud"
				}
			`,
			APICreate: `{
				"token": "...",
				"dataResidency": "us",
				"connectionModes": {
					"web": "cloud"
				},
				"identityMergeApi": "simplified",
				"consolidatedPageCalls": false,
				"persistence": "none",
				"useNewMapping": true,
				"useUserDefinedPageEventName": false,
				"userDefinedPageEventTemplate":  "Viewed {{ category }} {{ name }} page",
				"useUserDefinedScreenEventName": false,
				"userDefinedScreenEventTemplate":  "Viewed {{ category }} {{ name }} screen",
				"dropTraitsInTrackEvent": false,
				"strictMode": false,
				"userDeletionApi": "engage"
			}`,
			TerraformUpdate: `
				token = "..."
				data_residency = "eu"
				connection_modes {
					web = "cloud"
				}
				identity_merge_api = "simplified"
				persistence = "localStorage"
				api_secret = "..."
				people = true
				set_all_traits_by_default = true
				consolidated_page_calls = true
				track_categorized_pages = true
				track_named_pages = true
				source_name = "my-mixpanel"
				cross_subdomain_cookie = true
				secure_cookie = true
				super_properties = ["one","two","three"]
				people_properties = ["one","two","three"]
				event_increments = ["one","two","three"]
				prop_increments = ["one","two","three"]
				group_key_settings = ["one","two","three"]
				
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
					android = [{
						provider = "ketch"
						consents = ["one_android", "two_android", "three_android"]
						resolution_strategy = ""
					}]
					ios = [{
						provider = "custom"
						resolution_strategy = "and"
						consents = ["one_ios", "two_ios", "three_ios"]
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
				use_new_mapping = true
				use_user_defined_page_event_name = true
				user_defined_page_event_template = "Viewed {{ category }} {{ name }} page"
				use_user_defined_screen_event_name = true
				user_defined_screen_event_template = "Viewed {{ category }} {{ name }} screen"
				drop_traits_in_track_event = true
				strict_mode = true
				set_once_properties = ["one", "two", "three"]
				union_properties = ["one", "two", "three"]
				append_properties = ["one", "two", "three"]
				user_deletion_api = "task"
				gdpr_api_token = "..."
			`,
			APIUpdate: `
			{
				"token": "...",
				"dataResidency": "eu",
				"connectionModes": {
					"web": "cloud"
				},
				"identityMergeApi": "simplified",
				"persistence": "localStorage",
				"apiSecret": "...",
				"people": true,
				"setAllTraitsByDefault": true,
				"consolidatedPageCalls": true,
				"trackCategorizedPages": true,
				"trackNamedPages": true,
				"sourceName": "my-mixpanel",
				"crossSubdomainCookie": true,
				"secureCookie": true,
				"superProperties": [{
						"property": "one"
					},
					{
						"property": "two"
					},
					{
						"property": "three"
					}
				],
				"peopleProperties": [{
						"property": "one"
					},
					{
						"property": "two"
					},
					{
						"property": "three"
					}
				],
				"eventIncrements": [{
						"property": "one"
					},
					{
						"property": "two"
					},
					{
						"property": "three"
					}
				],
				"propIncrements": [{
						"property": "one"
					},
					{
						"property": "two"
					},
					{
						"property": "three"
					}
				],
				"groupKeySettings": [{
						"groupKey": "one"
					},
					{
						"groupKey": "two"
					},
					{
						"groupKey": "three"
					}
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
					],
					"android": [
						{
							"provider": "ketch",
							"resolutionStrategy": "",
							"consents": [
								{
									"consent": "one_android"
								},
								{
									"consent": "two_android"
								},
								{
									"consent": "three_android"
								}
							]
						}
					],
					"ios": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{
									"consent": "one_ios"
								},
								{
									"consent": "two_ios"
								},
								{
									"consent": "three_ios"
								}
							]
						}
					],
					"unity": [
						{
							"provider": "custom",
							"resolutionStrategy": "or",
							"consents": [
								{
									"consent": "one_unity"
								},
								{
									"consent": "two_unity"
								},
								{
									"consent": "three_unity"
								}
							]
						}
					],
					"reactnative": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{
									"consent": "one_reactnative"
								},
								{
									"consent": "two_reactnative"
								},
								{
									"consent": "three_reactnative"
								}
							]
						}
					],
					"flutter": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{
									"consent": "one_flutter"
								},
								{
									"consent": "two_flutter"
								},
								{
									"consent": "three_flutter"
								}
							]
						}
					],
					"cordova": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{
									"consent": "one_cordova"
								},
								{
									"consent": "two_cordova"
								},
								{
									"consent": "three_cordova"
								}
							]
						}
					],
					"amp": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{
									"consent": "one_amp"
								},
								{
									"consent": "two_amp"
								},
								{
									"consent": "three_amp"
								}
							]
						}
					],
					"cloud": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{
									"consent": "one_cloud"
								},
								{
									"consent": "two_cloud"
								},
								{
									"consent": "three_cloud"
								}
							]
						}
					],
					"warehouse": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{
									"consent": "one_warehouse"
								},
								{
									"consent": "two_warehouse"
								},
								{
									"consent": "three_warehouse"
								}
							]
						}
					],
					"shopify": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{
									"consent": "one_shopify"
								},
								{
									"consent": "two_shopify"
								},
								{
									"consent": "three_shopify"
								}
							]
						}
					]
				},
				"useNewMapping": true,
				"useUserDefinedPageEventName": true,
				"userDefinedPageEventTemplate":  "Viewed {{ category }} {{ name }} page",
				"useUserDefinedScreenEventName": true,
				"userDefinedScreenEventTemplate":  "Viewed {{ category }} {{ name }} screen",
				"dropTraitsInTrackEvent": true,
				"strictMode": true,
				"setOnceProperties": [
					{
						"property" : "one"
					},
					{
						"property" : "two"
					},
					{
						"property" : "three"
					}
				],
				"unionProperties": [
					{
						"property" : "one"
					},
					{
						"property" : "two"
					},
					{
						"property" : "three"
					}
				],
				"appendProperties": [
					{
						"property" : "one"
					},
					{
						"property" : "two"
					},
					{
						"property" : "three"
					}
				],
				"userDeletionApi": "task",
				"gdprApiToken": "..."
			}
			`,
		},
	})
}
