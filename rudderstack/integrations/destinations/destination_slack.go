package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	c.Destinations.Register("slack", c.ConfigMeta{
		APIType: "SLACK",
		Properties: []c.ConfigProperty{
			c.Simple("webhookUrl", "webhook_url"),
			c.Simple("identifyTemplate", "identify_template", c.SkipZeroValue),
			c.ArrayWithObjects("eventChannelSettings", "event_channel_settings", map[string]string{
				"eventName":    "name",
				"eventChannel": "channel",
				"eventRegex":   "regex",
			}),
			c.ArrayWithObjects("eventTemplateSettings", "event_template_settings", map[string]string{
				"eventName":     "name",
				"eventTemplate": "template",
				"eventRegex":    "regex",
			}),
			c.ArrayWithStrings("whitelistedTraitsSettings", "trait", "whitelisted_trait_settings"),
		},
		ConfigSchema: map[string]*schema.Schema{
			"webhook_url": {
				Type:     schema.TypeString,
				Required: true,
				// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"identify_template": {
				Type:     schema.TypeString,
				Optional: true,
				// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"event_channel_settings": {
				Type:       schema.TypeList,
				Optional:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"channel": {
							Type:     schema.TypeString,
							Required: true,
						},
						"regex": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},
			"event_template_settings": {
				Type:       schema.TypeList,
				Optional:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"template": {
							Type:     schema.TypeString,
							Required: true,
						},
						"regex": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},
			"whitelisted_trait_settings": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	})
}
