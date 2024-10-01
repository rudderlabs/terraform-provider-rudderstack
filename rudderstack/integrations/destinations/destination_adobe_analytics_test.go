package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceAdobeAnalytics(t *testing.T) {
	cmt.AssertDestination(t, "adobe_analytics", []c.TestConfig{
		{
			TerraformCreate: `
				report_suite_ids = "id001, id002"
							`,
			APICreate: `{
				"reportSuiteIds": "id001, id002",
  				"sslHeartbeat": true,
  				"useUtf8Charset": true,
  				"useSecureServerSide": true,
  				"dropVisitorId": true,
  				"timestampOptionalReporting": false,
  				"noFallbackVisitorId": false,
  				"preferVisitorId": false,
  				"trackPageName": true,
  				"useLegacyLinkName": true,
  				"pageNameFallbackTostring": true,
  				"sendFalseValues": true
			}`,
			TerraformUpdate: `
				report_suite_ids = "id003, id004"
				events_to_types = [{
					from = "video start"
					to = "heartbeatPlaybackStarted"
					}]
				list_delimiter = [{
					from = "listPhone"
					to = ","
					}]
				props_delimiter = [{
					from = "customPhone"
					to = ","
					}]
				event_merch_properties = [
					"currency"
					]
				product_merch_properties = [
					"currency"
					]
				event_filtering {
					blacklist = ["one", "two", "three"]
				}
				rudder_events_to_adobe_events = [{
					from = "product searched"
					to = "ps1,ps2"
					}]
				context_data_mapping = [{
					from = "page.name"
					to = "pName"
					}]
				mobile_event_mapping = [{
					from = "page.name"
					to = "pName"
					}]
				e_var_mapping = [{
					from = "phone"
					to = "1"
					}]
				hier_mapping = [{
					from = "phone"
					to = "1"
					}]
				list_mapping = [{
					from = "listPhone"
					to = "1"
					}]
				custom_props_mapping = [{
					from = "phone"
					to = "1"
					}]
				event_merch_event_to_adobe_event = [{
					from = "Order Completed"
					to = "merchEvent1"
					}]
				product_merch_event_to_adobe_event = [{
					from = "Product Ordered"
					to = "MerchProduct1"
					}]
				product_merch_evars_map = [{
					from = "phone"
					to = "1"
					}]
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
			  `,
			APIUpdate: `{
				"reportSuiteIds": "id003, id004",
				"sslHeartbeat": true,
				"useUtf8Charset": true,
				"useSecureServerSide": true,
				"eventsToTypes": [
				  {
					"from": "video start",
					"to": "heartbeatPlaybackStarted"
				  }
				],
				"dropVisitorId": true,
				"timestampOptionalReporting": false,
				"noFallbackVisitorId": false,
				"preferVisitorId": false,
				"rudderEventsToAdobeEvents": [
				  {
					"from": "product searched",
					"to": "ps1,ps2"
				  }
				],
				"trackPageName": true,
				"contextDataMapping": [
				  {
					"from": "page.name",
					"to": "pName"
				  }
				],
				"useLegacyLinkName": true,
				"pageNameFallbackTostring": true,
				"mobileEventMapping": [
				  {
					"from": "page.name",
					"to": "pName"
				  }
				],
				"sendFalseValues": true,
				"eVarMapping": [
				  {
					"from": "phone",
					"to": "1"
				  }
				],
				"hierMapping": [
				  {
					"from": "phone",
					"to": "1"
				  }
				],
				"listMapping": [
				  {
					"from": "listPhone",
					"to": "1"
				  }
				],
				"listDelimiter": [
				  {
					"from": "listPhone",
					"to": ","
				  }
				],
				"customPropsMapping": [
				  {
					"from": "phone",
					"to": "1"
				  }
				],
				"propsDelimiter": [
				  {
					"from": "customPhone",
					"to": ","
				  }
				],
				"eventMerchEventToAdobeEvent": [
				  {
					"from": "Order Completed",
					"to": "merchEvent1"
				  }
				],
				"eventMerchProperties": [
				  {
					"eventMerchProperties": "currency"
				  }
				],
				"productMerchEventToAdobeEvent": [
				  {
					"from": "Product Ordered",
					"to": "MerchProduct1"
				  }
				],
				"productMerchProperties": [
				  {
					"productMerchProperties": "currency"
				  }
				],
				"productMerchEvarsMap": [
				  {
					"from": "phone",
					"to": "1"
				  }
				],
				"blacklistedEvents": [
				  {
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
				"eventFilteringOption": "blacklistedEvents"
			  }`,
		},
	})
}
