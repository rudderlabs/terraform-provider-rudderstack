package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"web"}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)

	properties := []c.ConfigProperty{
		c.Simple("conversionID", "conversion_id"),
		c.ArrayWithObjects("pageLoadConversions", "page_load_conversions", map[string]string{
			"conversionLabel": "label",
			"name":            "name",
		}),
		c.ArrayWithObjects("clickEventConversions", "click_event_conversions", map[string]string{
			"conversionLabel": "label",
			"name":            "name",
		}),
		c.Simple("defaultPageConversion", "default_page_conversion", c.SkipZeroValue),
		c.Simple("dynamicRemarketing.web", "dynamic_remarketing.0.web"),
		c.Simple("conversionLinker", "conversion_linker", c.SkipZeroValue),
		c.Simple("sendPageView", "send_page_view", c.SkipZeroValue),
		c.Simple("disableAdPersonalization", "disable_ad_personalization", c.SkipZeroValue),
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
		"conversion_id": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Enter your Google Ads Conversion ID.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^AW-(.{0,100})$"),
		},
		"page_load_conversions": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "You can configure the page load conversions for multiple instances.",
			ConfigMode:  schema.SchemaConfigModeAttr,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"label": {
						Type:             schema.TypeString,
						Required:         true,
						Description:      "Provide the conversion label from Google Ads.",
						ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
					},
					"name": {
						Type:             schema.TypeString,
						Required:         true,
						Description:      "Enter the name of the `page` event to be sent to Google Ads.",
						ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
					},
				},
			},
		},
		"click_event_conversions": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "For `track` calls, you can configure these fields.",
			ConfigMode:  schema.SchemaConfigModeAttr,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"label": {
						Type:             schema.TypeString,
						Required:         true,
						Description:      "Enter your Google Ads conversion label.",
						ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
					},
					"name": {
						Type:             schema.TypeString,
						Required:         true,
						Description:      "Enter the name of the `track` event to be sent to Google Ads.",
						ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
					},
				},
			},
		},
		"default_page_conversion": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Enter the default conversion label.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
		},
		"dynamic_remarketing": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Enabling this tracking mode allows RudderStack to leverage Google Ads' Dynamic Remarketing feature for event tracking.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"web": {
						Type:     schema.TypeBool,
						Optional: true,
					},
				},
			},
		},
		"conversion_linker": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "This setting is enabled by default. If you don't want the global site tag (gtag.js) to set first-party cookies on your website domain, you should disable this setting.",
		},
		"send_page_view": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Enabling this setting configures Google Ads to automatically send your `page` events.",
		},
		"disable_ad_personalization": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Enable this setting to programmatically disable ad personalization.",
		},
		"event_filtering": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "With this option, you can determine which events are blocked or allowed to flow through to Google Ads.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"whitelist": {
						Type:         schema.TypeList,
						Optional:     true,
						Description:  "Enter the event names to be whitelisted.",
						ExactlyOneOf: []string{"config.0.event_filtering.0.whitelist", "config.0.event_filtering.0.blacklist"},
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"blacklist": {
						Type:         schema.TypeList,
						Optional:     true,
						Description:  "Enter the event names to be blacklisted.",
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
			Description: "As this is a device mode destination, this setting will always be enabled.",
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

	c.Destinations.Register("google_ads", c.ConfigMeta{
		APIType: "GOOGLEADS",
		Properties: properties,
		ConfigSchema: schema,
	})
}
