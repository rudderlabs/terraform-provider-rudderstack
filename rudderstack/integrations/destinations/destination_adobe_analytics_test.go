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
				onetrust_cookie_categories = ["one", "two", "three"]
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
				"oneTrustCookieCategories": [
					{ "oneTrustCookieCategory": "one" },
					{ "oneTrustCookieCategory": "two" },
					{ "oneTrustCookieCategory": "three" }
				],
				"eventFilteringOption": "blacklistedEvents"
			  }`,
		},
	})
}
