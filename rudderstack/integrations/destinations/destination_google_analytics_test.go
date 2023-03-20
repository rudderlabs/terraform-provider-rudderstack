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
		
			    onetrust_cookie_categories = ["one", "two", "three"]

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
				"oneTrustCookieCategories": [
					{ "oneTrustCookieCategory": "one" },
					{ "oneTrustCookieCategory": "two" },
					{ "oneTrustCookieCategory": "three" }
				],
				"contentGroupings": [{ "from": "from", "to": "to" }],
				"dimensions": [{ "from": "from", "to": "to" }]
			}`,
		},
	})
}
