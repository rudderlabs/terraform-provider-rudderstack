package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceIterable(t *testing.T) {
	cmt.AssertDestination(t, "iterable", []c.TestConfig{
		{
			TerraformCreate: `
				api_key = "73983282843839749873"	
				map_to_single_event = true
				track_categorized_pages = true
				track_named_pages = true
			`,
			APICreate: `{
				"apiKey": "73983282843839749873",
				"mapToSingleEvent": true,
				"trackCategorisedPages": true,
				"trackNamedPages": true
			}`,
			TerraformUpdate: `
				api_key = "83983282843839749873"
				map_to_single_event = false
				track_all_pages = true
				track_categorized_pages = true
				track_named_pages = true
				use_native_sdk {
					web = true
				}
				initialisation_identifier { 
					web = "email" 
				}
				get_in_app_event_mapping {
					web = ["one", "two", "three"]
				}
				purchase_event_mapping { 
					web = ["one", "two", "three"]
				}
				send_track_for_inapp { 
					web = true 
				}
				animation_duration { 
					web = "200" 
				}
				display_interval { 
					web = "2500" 
				}
				on_open_screen_reader_message { 
					web =  "..." 
				}
				on_open_node_to_take_focus { 
					web =  "..." 
				}
				package_name { 
					web = "my-package-test" 
				}
				right_offset { 
					web = "15" 
				}
				top_offset { 
					web = "11" 
				}
				bottom_offset { 
					web = "24%" 
				}
				handle_links { 
					web = "open-all-new-tab" 
				}
				close_button_color { 
					web = "blue" 
				}
				close_button_size { 
					web = "..."
				}
				close_button_color_top_offset { 
					web = "3%"
				}
				close_button_color_side_offset { 
					web = "2%" 
				}
				icon_path { 
					web = "..." 
				}
				is_required_to_dismiss_message { 
					web = true 
				}
				close_button_position { 
					web = "..." 
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

			`,
			APIUpdate: `{
				"apiKey": "83983282843839749873",
				"mapToSingleEvent": false,
				"trackAllPages": true,
				"trackCategorisedPages": true,
				"trackNamedPages": true,
				"useNativeSDK": {
					"web": true
				},
				"initialisationIdentifier": { "web": "email" },
				"getInAppEventMapping": {
					"web": [
						{ "eventName": "one" }, 
						{ "eventName": "two" }, 
						{ "eventName": "three" }
					]
				},
				"purchaseEventMapping": {
					 "web": [
						{ "eventName": "one" },
						{ "eventName": "two" },
						{ "eventName": "three" }
					] 
				},
				"sendTrackForInapp": { "web": true },
				"animationDuration": { "web": "200" },
				"displayInterval": { "web": "2500" },
				"onOpenScreenReaderMessage": { "web": "..." },
				"onOpenNodeToTakeFocus": { "web": "..." },
				"packageName": { "web": "my-package-test" },
				"rightOffset": { "web": "15" },
				"topOffset": { "web": "11" },
				"bottomOffset": { "web": "24%" },
				"handleLinks": { "web": "open-all-new-tab" },
				"closeButtonColor": { "web": "blue" },
				"closeButtonSize": { "web": "..." },
				"closeButtonColorTopOffset": { "web": "3%" },
				"closeButtonColorSideOffset": { "web": "2%" },
				"iconPath": { "web": "..." },
				"isRequiredToDismissMessage": { "web": true },
				"closeButtonPosition": { "web": "..." },
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
				}
			}`,
		},
	})
}
