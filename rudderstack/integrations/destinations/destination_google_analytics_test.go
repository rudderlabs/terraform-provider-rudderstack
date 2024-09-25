package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceGoogleAnalytics(t *testing.T) {
	cmt.AssertDestination(t, "google_analytics", []c.TestConfig{
		{
			TerraformCreate: `
				tracking_id = "UA-00-0000"
			`,
			APICreate: `{
				"trackingID": "UA-00-0000"
			}`,
			TerraformUpdate: `
				tracking_id = "UA-00-0000"

				double_click              = true
				enhanced_link_attribution = true
				include_search            = true
				disable_md5               = true
				anonymize_ip              = true
				enhanced_ecommerce        = true
				non_interaction           = true
	
				server_side_identify {
					event_category = "..."
					event_action   = "..."
				}
	
				track_categorized_pages {
					web = true
				}
	
				track_named_pages {
					web = true
				}
	
				sample_rate {
					web = "1000"
				}
	
				site_speed_sample_rate {
					web = "1000"
				}
	
				set_all_mapped_props {
					web = true
				}
	
				domain {
					web = "..."
				}
	
				optimize {
					web = "..."
				}
	
				use_google_amp_client_id {
					web = true
				}
	
				named_tracker {
					web = true
				}
	
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

				reset_custom_dimensions_on_page {
				   web = ["one", "two", "three"]
				}				

				content_groupings = [{
				  from = "from"
				  to   = "to"
				}]
			
				dimensions = [{
				  from = "from"
				  to   = "to"
				}]
			`,
			APIUpdate: `{
				"trackingID": "UA-00-0000",
				"doubleClick": true,
				"enhancedLinkAttribution": true,
				"includeSearch": true,
				"disableMd5": true,
				"anonymizeIp": true,
				"enhancedEcommerce": true,
				"nonInteraction": true,
				"blacklistedEvents": [
				  { "eventName": "one" },
				  { "eventName": "two" },
				  { "eventName": "three" }
				],
				"serverSideIdentifyEventCategory": "...",
				"serverSideIdentifyEventAction": "...",
				"enableServerSideIdentify": true,
				"eventFilteringOption": "blacklistedEvents",
				"useNativeSDK": { "web": true },
				"trackCategorizedPages": { "web": true },
				"trackNamedPages": { "web": true },
				"sampleRate": { "web": "1000" },
				"siteSpeedSampleRate": { "web": "1000" },
				"setAllMappedProps": { "web": true },
				"domain": { "web": "..." },
				"optimize": { "web": "..." },
				"useGoogleAmpClientId": { "web": true },
				"namedTracker": { "web": true },
				"resetCustomDimensionsOnPage": {
					"web": [
						{ "resetCustomDimensionsOnPage": "one" },
						{ "resetCustomDimensionsOnPage": "two" },
						{ "resetCustomDimensionsOnPage": "three" }
					]
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
				"contentGroupings": [{ "from": "from", "to": "to" }],
				"dimensions": [{ "from": "from", "to": "to" }]
			}`,
		},
	})
}
