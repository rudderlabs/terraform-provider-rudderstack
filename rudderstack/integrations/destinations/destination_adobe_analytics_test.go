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
			}`,
			TerraformUpdate: `{
			}`,
			APIUpdate: `{
				"trackingServerUrl": "https://www.rudderstack.com/docs/destinations/streaming-destinations/adobe-analytics",
				"trackingServerSecureUrl": "https://www.rudderstack.com/docs/destinations/streaming-destinations/adobe-analytics/",
				"reportSuiteIds": "id001, id002",
				"sslHeartbeat": true,
				"heartbeatTrackingServerUrl": "https://heartbeat.com",
				"useUtf8Charset": true,
				"useSecureServerSide": true,
				"proxyNormalUrl": "",
				"proxyHeartbeatUrl": "",
				"marketingCloudOrgId": "org10121",
				"dropVisitorId": true,
				"timestampOption": "enabled",
				"timestampOptionalReporting": true,
				"noFallbackVisitorId": false,
				"preferVisitorId": true,
				"trackPageName": true,
				"contextDataPrefix": "rudderstack-context",
				"useLegacyLinkName": true,
				"pageNameFallbackTostring": false,
				"sendFalseValues": true,
				"productIdentifier": "sku",
				"eventFilteringOption": "blacklistedEvents",
				"eventsToTypes": [
				  {
					"from": "video start",
					"to": "heartbeatPlaybackStarted"
				  },
				  {
					"from": "video end",
					"to": "heartbeatPlaybackCompleted"
				  }
				],
				"listDelimiter": [
				  {
					"from": "listPhone",
					"to": ","
				  },
				  {
					"from": "listMobile",
					"to": "|"
				  }
				],
				"propsDelimiter": [
				  {
					"from": "customPhone",
					"to": ","
				  },
				  {
					"from": "customMobile",
					"to": "|"
				  }
				],
				"eventMerchProperties": [
				  {
					"eventMerchProperties": "currency"
				  },
				  {
					"eventMerchProperties": "location"
				  }
				],
				"productMerchProperties": [
				  {
					"productMerchProperties": "productName"
				  },
				  {
					"productMerchProperties": "productInventory"
				  }
				],
				"whitelistedEvents": [
				  {
					"eventName": ""
				  }
				],
				"blacklistedEvents": [
				  {
					"eventName": "Restart Checkout"
				  },
				  {
					"eventName": "Restart Initiated"
				  }
				],
				"rudderEventsToAdobeEvents": [
				  {
					"from": "product searched",
					"to": "ps1,ps2"
				  }
				],
				"contextDataMapping": [
				  {
					"from": "page.name",
					"to": "pName"
				  },
				  {
					"from": "page.url",
					"to": "pUrl"
				  }
				],
				"mobileEventMapping": [
				  {
					"from": "page.name",
					"to": "pName"
				  },
				  {
					"from": "page.url",
					"to": "pUrl"
				  }
				],
				"eVarMapping": [
				  {
					"from": "phone",
					"to": "1"
				  },
				  {
					"from": "mobile",
					"to": "2"
				  }
				],
				"hierMapping": [
				  {
					"from": "phone",
					"to": "1"
				  },
				  {
					"from": "mobile",
					"to": "2"
				  }
				],
				"listMapping": [
				  {
					"from": "listPhone",
					"to": "1"
				  },
				  {
					"from": "listMobile",
					"to": "2"
				  }
				],
				"customPropsMapping": [
				  {
					"from": "phone",
					"to": "1"
				  },
				  {
					"from": "mobile",
					"to": "2"
				  }
				],
				"eventMerchEventToAdobeEvent": [
				  {
					"from": "Order Completed",
					"to": "merchEvent1"
				  }
				],
				"productMerchEventToAdobeEvent": [
				  {
					"from": "Product Ordered",
					"to": "MerchProduct1"
				  }
				],
				"productMerchEvarsMap": [
				  {
					"from": "phone",
					"to": "1"
				  },
				  {
					"from": "mobile",
					"to": "2"
				  }
				],
				"useNativeSDK": {
				  "web": true
				},
				"oneTrustCookieCategories": {
				  "web": [
					{
					  "oneTrustCookieCategory": "Marketing"
					}
				  ]
				}
			  }`,
		},
	})
}
