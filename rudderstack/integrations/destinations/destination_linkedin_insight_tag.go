package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"web"}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)

	properties := []c.ConfigProperty{
		c.Simple("partnerId", "partner_id"),
		c.ArrayWithStrings("whitelistedEvents", "eventName", "event_filtering.0.whitelist"),
		c.ArrayWithStrings("blacklistedEvents", "eventName", "event_filtering.0.blacklist"),
		c.Discriminator("eventFilteringOption", c.DiscriminatorValues{
			"event_filtering.0.whitelist": "whitelistedEvents",
			"event_filtering.0.blacklist": "blacklistedEvents",
		}),
		c.ArrayWithObjects("eventToConversionIdMap", "event_to_conversion_id_map", map[string]string{
			"from": "from",
			"to":   "to",
		}),
		c.Simple("useNativeSDK.web", "use_native_sdk.0.web"),
	}

	properties = append(properties, commonProperties...)

	schema := map[string]*schema.Schema{
		"partner_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Enter your LinkedIn Partner ID.",
		},
		"event_to_conversion_id_map": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Event Conversion IDs.",
			ConfigMode:  schema.SchemaConfigModeAttr,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"from": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Provide the event name.",
					},
					"to": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Enter the conversion ID.",
					},
				},
			},
		},
		"event_filtering": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "With this option, you can determine which events are blocked or allowed to flow through to LinkedIn.",
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

	c.Destinations.Register("LINKEDIN_INSIGHT_TAG", c.ConfigMeta{
		APIType:      "LINKEDIN_INSIGHT_TAG",
		Properties:   properties,
		ConfigSchema: schema,
	})
}
