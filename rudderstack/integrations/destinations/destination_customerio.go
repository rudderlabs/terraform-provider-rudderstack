package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"web", "android", "ios", "unity", "reactnative", "flutter", "cordova", "amp", "cloud", "warehouse", "shopify"}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)

	properties := []c.ConfigProperty{
		c.Simple("siteID", "site_id"),
		c.Simple("apiKey", "api_key"),
		c.Simple("deviceTokenEventName", "device_token_event_name", c.SkipZeroValue),
		c.Simple("datacenterEU", "datacenter_eu", c.SkipZeroValue),
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
		"site_id": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Enter your Customer.io site ID.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"api_key": {
			Type:             schema.TypeString,
			Required:         true,
			Sensitive:        true,
			Description:      "Enter your Customer.io API key.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"device_token_event_name": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Enter the name of the event that is fired immediately after setting the device token.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"datacenter_eu": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Enable this option in case your account is based in the EU region.",
		},
		"use_native_sdk": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Enable this setting to send the events through Customer.io's native JavaScript SDK.",
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
			Description: "RudderStack lets you determine which events should be allowed to flow through or blocked.",
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
	}

	for key, value := range commonSchema {
		schema[key] = value
	}

	c.Destinations.Register("customerio", c.ConfigMeta{
		APIType:      "CUSTOMERIO",
		Properties:   properties,
		ConfigSchema: schema,
	})
}
