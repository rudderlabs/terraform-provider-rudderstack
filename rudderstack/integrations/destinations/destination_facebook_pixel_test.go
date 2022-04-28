package destinations_test

import (
	"testing"

	"github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceFacebookPixel(t *testing.T) {
	testutil.AssertDestination(t, "facebook_pixel", []c.TestConfig{
		{
			TerraformCreate: `
				pixel_id = "abc123"
			`,
			APICreate: `{
				"pixelId": "abc123"
			}`,
			TerraformUpdate: `
				pixel_id = "facebook pixel id"

				access_token = "facebook access token"
			
				standard_page_call     = true
				value_field_identifier = "properties.price"
				advanced_mapping       = true
				test_destination       = true
				test_event_code        = "..."
			
				events_to_events = [{
					from = "a1"
					to   = "b1"
				}, {
					from = "a2"
					to   = "b2"
				}]
			
				event_custom_properties = ["one", "two", "three"]
			
				# blacklist_pii_properties {
				# 
				# }
			
				# whitelist_pii_properties {
				# 
				# }
			
				category_to_content = [{
				  from = "from"
				  to   = "to"
				}]
			
				legacy_conversion_pixel_id {
				  from = "from"
				  to   = "to"
				}
			
				use_native_sdk {
				  web          = true
				}
			
				event_filtering {
				  blacklist = ["one", "two", "three"]
				}
			
				onetrust_cookie_categories {
				  web = ["one", "two", "three"]
				}
			`,
			APIUpdate: `{
				"pixelId": "facebook pixel id",
				"accessToken": "facebook access token",
				"standardPageCall": true,
				"valueFieldIdentifier": "properties.price",
				"advancedMapping": true,
				"testDestination": true,
				"testEventCode": "...",
				"eventsToEvents": [
				  { "from": "a1", "to": "b1" },
				  { "from": "a2", "to": "b2" }
				],
				"eventCustomProperties": [
				  { "eventCustomProperties": "one" },
				  { "eventCustomProperties": "two" },
				  { "eventCustomProperties": "three" }
				],
				"categoryToContent": [{ "from": "from", "to": "to" }],
				"legacyConversionPixelId": { "from": "from", "to": "to" },
				"useNativeSDK": { "web": true },
				"oneTrustCookieCategories": {
				  "web": [
					{ "oneTrustCookieCategory": "one" },
					{ "oneTrustCookieCategory": "two" },
					{ "oneTrustCookieCategory": "three" }
				  ]
				},
				"blacklistedEvents": [
				  { "eventName": "one" },
				  { "eventName": "two" },
				  { "eventName": "three" }
				],
				"eventFilteringOption": "blacklistedEvents"
			}`,
		},
	})
}
