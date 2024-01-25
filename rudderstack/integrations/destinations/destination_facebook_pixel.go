package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"web", "android", "ios", "unity", "reactnative", "flutter", "cordova", "amp", "cloud", "warehouse", "shopify"}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)

	properties := []c.ConfigProperty{
		c.Simple("pixelId", "pixel_id"),
		c.Simple("accessToken", "access_token", c.SkipZeroValue),
		c.Simple("standardPageCall", "standard_page_call", c.SkipZeroValue),
		c.Simple("valueFieldIdentifier", "value_field_identifier", c.SkipZeroValue),
		c.Simple("advancedMapping", "advanced_mapping", c.SkipZeroValue),
		c.Simple("testDestination", "test_destination", c.SkipZeroValue),
		c.Simple("testEventCode", "test_event_code", c.SkipZeroValue),
		c.Simple("eventsToEvents", "events_to_events", c.SkipZeroValue),
		c.ArrayWithStrings("eventCustomProperties", "eventCustomProperties", "event_custom_properties"),
		c.ArrayWithObjects("blacklistPiiProperties", "blacklist_pii_properties", map[string]interface{}{
			"blacklistPiiProperties": "property",
			"blacklistPiiHash":       "hash",
		}),
		c.ArrayWithObjects("whitelistPiiProperties", "whitelist_pii_properties", map[string]interface{}{
			"whitelistPiiProperties": "property",
		}),
		c.Simple("categoryToContent", "category_to_content", c.SkipZeroValue),
		c.Simple("legacyConversionPixelId.from", "legacy_conversion_pixel_id.0.from"),
		c.Simple("legacyConversionPixelId.to", "legacy_conversion_pixel_id.0.to"),
		c.Simple("useNativeSDK.web", "use_native_sdk.0.web"),
		c.ArrayWithStrings("whitelistedEvents", "eventName", "event_filtering.0.whitelist"),
		c.ArrayWithStrings("blacklistedEvents", "eventName", "event_filtering.0.blacklist"),
		c.Discriminator("eventFilteringOption", c.DiscriminatorValues{
			"event_filtering.0.whitelist": "whitelistedEvents",
			"event_filtering.0.blacklist": "blacklistedEvents",
		}),
	}

	properties = append(properties, commonProperties...)

	schema := map[string]*schema.Schema{
		"pixel_id": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Enter your Facebook Pixel ID.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"access_token": {
			Type:             schema.TypeString,
			Optional:         true,
			Sensitive:        true,
			Description:      "Enter your Facebook business access token required to send the events via the cloud mode.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,208})$"),
		},
		"standard_page_call": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "If this setting is enabled, RudderStack sets `pageview` as a standard event for all the `page` and `screen` calls.",
		},
		"value_field_identifier": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "You can set this field to `properties.price` or `properties.value`. RudderStack will then assign this to the value field of the Facebook payload.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(properties.value|properties.price)$"),
		},
		"advanced_mapping": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "With this setting, you can enable the advanced mapping feature.",
		},
		"test_destination": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Enable this setting if you are using this destination for testing purposes.",
		},
		"test_event_code": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "If the above setting is enabled, enter the relevant test event code.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
		},
		"events_to_events": {
			Type:        schema.TypeList,
			MaxItems:    10,
			Optional:    true,
			Description: "You can map your events to standard Facebook events using this setting.",
			ConfigMode:  schema.SchemaConfigModeAttr,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"from": {
						Type:             schema.TypeString,
						Required:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
					},
					"to": {
						Type:             schema.TypeString,
						Required:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
					},
				},
			},
		},
		"event_custom_properties": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "For the standard events, some predefined properties are taken by Facebook. If you want to send more properties for your events, mention those properties in this field.",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"blacklist_pii_properties": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Enter the PII properties to be blacklisted.",
			ConfigMode:  schema.SchemaConfigModeAttr,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"property": {
						Type:             schema.TypeString,
						Required:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
					},
					"hash": {
						Type:     schema.TypeBool,
						Required: true,
					},
				},
			},
		},
		"whitelist_pii_properties": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Enter the PII properties to be whitelisted.",
			ConfigMode:  schema.SchemaConfigModeAttr,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"property": {
						Type:             schema.TypeString,
						Required:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
					},
				},
			},
		},
		"category_to_content": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "This option lets you specify the category fields to specific Facebook content type.",
			ConfigMode:  schema.SchemaConfigModeAttr,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"from": {
						Type:             schema.TypeString,
						Required:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
					},
					"to": {
						Type:             schema.TypeString,
						Required:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
					},
				},
			},
		},
		"legacy_conversion_pixel_id": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "With this setting, you can send specific events to a legacy conversion Pixel by specifying the event-Pixel ID mapping. Note that this option is available only when sending events via the device mode.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"from": {
						Type:             schema.TypeString,
						Required:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
					},
					"to": {
						Type:             schema.TypeString,
						Required:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
					},
				},
			},
		},
		"use_native_sdk": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Enable this setting to send events from the web SDK to Facebook Pixel via the device mode.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"web": {
						Type:     schema.TypeBool,
						Optional: true,
					},
				},
			},
		},
		"event_filtering": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "This setting lets you determine which events are blocked or allowed to flowed through to Facebook Pixel.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"whitelist": {
						Type:         schema.TypeList,
						Optional:     true,
						ExactlyOneOf: []string{"config.0.event_filtering.0.whitelist", "config.0.event_filtering.0.blacklist"},
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"blacklist": {
						Type:         schema.TypeList,
						Optional:     true,
						ExactlyOneOf: []string{"config.0.event_filtering.0.whitelist", "config.0.event_filtering.0.blacklist"},
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
	}

	for key, value := range commonSchema {
		schema[key] = value
	}

	c.Destinations.Register("facebook_pixel", c.ConfigMeta{
		APIType: "FACEBOOK_PIXEL",
		Properties: properties,
		ConfigSchema: schema,
	})
}
