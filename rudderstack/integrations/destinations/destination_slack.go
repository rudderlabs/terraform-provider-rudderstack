package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"web", "android", "ios", "unity", "reactnative", "flutter", "cordova", "amp", "cloud", "warehouse", "shopify"}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)

	properties := []c.ConfigProperty{
		c.Simple("webhookUrl", "webhook_url"),
		c.Simple("identifyTemplate", "identify_template", c.SkipZeroValue),
		c.ArrayWithObjects("eventChannelSettings", "event_channel_settings", map[string]interface{}{
			"eventName":    "name",
			"eventChannel": "channel",
			"eventRegex":   "regex",
		}),
		c.ArrayWithObjects("eventTemplateSettings", "event_template_settings", map[string]interface{}{
			"eventName":     "name",
			"eventTemplate": "template",
			"eventRegex":    "regex",
		}),
		c.ArrayWithStrings("whitelistedTraitsSettings", "trait", "whitelisted_trait_settings"),
	}

	properties = append(properties, commonProperties...)

	schema := map[string]*schema.Schema{
		"webhook_url": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Enter your Slack's incoming webhook URL.",
			// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"identify_template": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Specify the template that you want the `identify` event to be transformed to before it is sent to Slack.",
			// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"event_channel_settings": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Specify your event channel settings.",
			ConfigMode:  schema.SchemaConfigModeAttr,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Enter the event name or the regex to match the RudderStack event name.",
					},
					"channel": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Enter the name of the Slack channel where the event will be sent.",
					},
					"regex": {
						Type:        schema.TypeBool,
						Required:    true,
						Description: "Enable this setting if the event name in the first parameter is a regular expression.",
					},
				},
			},
		},
		"event_template_settings": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Specify your event template settings.",
			ConfigMode:  schema.SchemaConfigModeAttr,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Enter the event name or the regex to match the RudderStack event name.",
					},
					"template": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Specify the template for the above event names matching the regex.",
					},
					"regex": {
						Type:        schema.TypeBool,
						Required:    true,
						Description: "Enable this setting if the event name is a regex in the first parameter.",
					},
				},
			},
		},
		"whitelisted_trait_settings": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Only the traits listed in this section are considered to be a part of the identify template. The rest are sent to Slack.",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}

	for key, value := range commonSchema {
		schema[key] = value
	}

	c.Destinations.Register("slack", c.ConfigMeta{
		APIType:      "SLACK",
		Properties:   properties,
		ConfigSchema: schema,
	})
}
