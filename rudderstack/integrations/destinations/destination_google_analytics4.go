package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"web", "android", "ios", "unity", "reactnative", "flutter", "cordova", "amp", "cloud", "warehouse", "shopify"}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)

	properties := []c.ConfigProperty{
		c.Simple("apiSecret", "api_secret", c.SkipZeroValue),
		c.Simple("typesOfClient", "types_of_client", c.SkipZeroValue),
		c.Simple("measurementId", "measurement_id", c.SkipZeroValue),
		c.Simple("firebaseAppId", "firebase_app_id", c.SkipZeroValue),
		c.Simple("blockPageViewEvent", "block_page_view_event", c.SkipZeroValue),
		c.Simple("extendPageViewParams", "extend_page_view_params", c.SkipZeroValue),
		c.Simple("sendUserId", "send_user_id", c.SkipZeroValue),
		c.ArrayWithStrings("whitelistedEvents", "eventName", "event_filtering.0.whitelist"),
		c.ArrayWithStrings("blacklistedEvents", "eventName", "event_filtering.0.blacklist"),
		c.Discriminator("eventFilteringOption", c.DiscriminatorValues{
			"event_filtering.0.whitelist": "whitelistedEvents",
			"event_filtering.0.blacklist": "blacklistedEvents",
		}),
		c.Simple("useNativeSDK.web", "use_native_sdk.0.web"),
	}

	properties = append(properties, commonProperties...)

	schema := map[string]*schema.Schema{
		"api_secret": {
			Type:             schema.TypeString,
			Optional:         true,
			Sensitive:        true,
			Description:      "This field is required only for the cloud mode setup where you can enter the API Secret generated through the Google Analytics dashboard.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
		},
		"types_of_client": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Select the client type as gtag or Firebase.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(gtag|firebase)$"),
		},
		"measurement_id": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Enter the Measurement Id which is the identifier for a data stream.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(G-.{1,100})$|^$"),
		},
		"firebase_app_id": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Enter the Firebase App ID which is the identifier for Firebase app.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
		},
		"block_page_view_event": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Enable this setting to disable sending `page_view` events on load. This setting is applicable only for device mode.",
		},
		"extend_page_view_params": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Enable this setting to send `url` and `search` along with any other custom property to the `page` call of the RudderStack SDK. This setting is applicable only for device mode.",
		},
		"send_user_id": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "If enabled, the user ID is set to the identified visitors and sent to Google Analytics 4.",
		},
		"event_filtering": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "With this option, you can determine which events are blocked or allowed to flow through to Google Analytics 4.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"whitelist": {
						Type:         schema.TypeList,
						Optional:     true,
						Description:  "Enter the event names to be allowlisted.",
						ExactlyOneOf: []string{"config.0.event_filtering.0.whitelist", "config.0.event_filtering.0.blacklist"},
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"blacklist": {
						Type:         schema.TypeList,
						Optional:     true,
						Description:  "Enter the event names to be denylisted.",
						ExactlyOneOf: []string{"config.0.event_filtering.0.whitelist", "config.0.event_filtering.0.blacklist"},
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
		"use_native_sdk": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Enable this setting to send the events via the device mode.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"web": {
						Type:     schema.TypeBool,
						Optional: true,
					},
				},
			},
		},
	}

	for key, value := range commonSchema {
		schema[key] = value
	}

	c.Destinations.Register("google_analytics4", c.ConfigMeta{
		APIType:      "GA4",
		Properties:   properties,
		ConfigSchema: schema,
	})
}
