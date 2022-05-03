package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	c.Destinations.Register("webhook", c.ConfigMeta{
		APIType: "WEBHOOK",
		Properties: []c.ConfigProperty{
			c.Simple("webhookUrl", "webhook_url"),
			c.Simple("webhookMethod", "webhook_method", c.SkipZeroValue),
			c.Simple("headers", "headers", c.SkipZeroValue),
		},
		ConfigSchema: map[string]*schema.Schema{
			"webhook_url": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(?:http(s)?:\\/\\/)?[\\w.-]+(?:\\.[\\w\\.-]+)+[\\w\\-\\._~:/?#[\\]@!\\$&'\\(\\)\\*\\+,;=.]+$"),
			},
			"webhook_method": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|(^(POST|PUT|GET)$)"),
			},
			"headers": {
				Type:       schema.TypeList,
				Optional:   true,
				Sensitive:  true,
				ConfigMode: schema.SchemaConfigModeAttr,
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
							ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,500})$"),
						},
					},
				},
			},
		},
	})
}
