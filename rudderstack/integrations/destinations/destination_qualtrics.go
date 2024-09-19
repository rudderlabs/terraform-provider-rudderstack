package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"web", "android", "ios"}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)

	properties := []c.ConfigProperty{
		c.Simple("projectId", "project_id"),
		c.Simple("brandId", "brand_id"),
		c.Simple("enableGenericPageTitle.web", "enable_generic_page_title", c.SkipZeroValue),
		c.Simple("useNativeSDK.web", "use_native_sdk.0.web", c.SkipZeroValue),
		c.Simple("useNativeSDK.ios", "use_native_sdk.0.ios", c.SkipZeroValue),
		c.Simple("useNativeSDK.android", "use_native_sdk.0.android", c.SkipZeroValue),
		c.ArrayWithStrings("whitelistedEvents", "eventName", "event_filtering.0.whitelist"),
		c.ArrayWithStrings("blacklistedEvents", "eventName", "event_filtering.0.blacklist"),
		c.Discriminator("eventFilteringOption", c.DiscriminatorValues{
			"event_filtering.0.whitelist": "whitelistedEvents",
			"event_filtering.0.blacklist": "blacklistedEvents",
		}),
	}

	properties = append(properties, commonProperties...)

	schema := map[string]*schema.Schema{
		"project_id": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "Enter your Project ID.",
		},
		"brand_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Enter your Brand ID.",
		},
		"enable_generic_page_title": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "This setting enables Generic Page Title.",
		},
		"use_native_sdk": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Enable this setting to send the events through SDKs.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"web": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"ios": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"android": {
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
	}

	for key, value := range commonSchema {
		schema[key] = value
	}

	c.Destinations.Register("qualtrics", c.ConfigMeta{
		APIType:      "QUALTRICS",
		Properties:   properties,
		ConfigSchema: schema,
	})
}
