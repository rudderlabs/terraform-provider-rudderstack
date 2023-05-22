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
				onetrust_cookie_categories = ["one", "two", "three"]

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
				"oneTrustCookieCategories": [
					{ "oneTrustCookieCategory": "one" },
					{ "oneTrustCookieCategory": "two" },
					{ "oneTrustCookieCategory": "three" }
				]
			}`,
		},
	})
}
