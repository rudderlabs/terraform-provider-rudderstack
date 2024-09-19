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
		c.Simple("webhookMethod", "webhook_method", c.SkipZeroValue),
		c.Simple("headers", "headers", c.SkipZeroValue),
	}

	properties = append(properties, commonProperties...)

	schema := map[string]*schema.Schema{
		"webhook_url": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Enter the endpoint where RudderStack will send the events.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(?:http(s)?:\\/\\/)?[\\w.-]+(?:\\.[\\w\\.-]+)+[\\w\\-\\._~:/?#[\\]@!\\$&'\\(\\)\\*\\+,;=.]+$"),
		},
		"webhook_method": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "This is the HTTP method of the request sent to the configured endpoint. By default, `POST` is used.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|(^(POST|PUT|GET)$)"),
		},
		"headers": {
			Type:        schema.TypeList,
			Optional:    true,
			Sensitive:   true,
			Description: "Add custom headers for your events via this option. These headers will be added to the request made from RudderStack to your webhook.",
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
						ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,500})$"),
					},
				},
			},
		},
	}

	for key, value := range commonSchema {
		schema[key] = value
	}

	c.Destinations.Register("webhook", c.ConfigMeta{
		APIType:      "WEBHOOK",
		Properties:   properties,
		ConfigSchema: schema,
	})
}
