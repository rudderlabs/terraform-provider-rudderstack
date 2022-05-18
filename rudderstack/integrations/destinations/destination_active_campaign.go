package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	c.Destinations.Register("active_campaign", c.ConfigMeta{
		APIType: "ACTIVE_CAMPAIGN",
		Properties: []c.ConfigProperty{
			c.Simple("apiUrl", "api_url"),
			c.Simple("apiKey", "api_key", c.SkipZeroValue),
			c.Simple("actid", "actid", c.SkipZeroValue),
			c.Simple("eventKey", "event_key", c.SkipZeroValue),
		},
		ConfigSchema: map[string]*schema.Schema{
			"api_url": {
				Type:     schema.TypeString,
				Required: true,
				// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"api_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"actid": {
				Type:     schema.TypeString,
				Optional: true,
				// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"event_key": {
				Type:     schema.TypeString,
				Optional: true,
				// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
		},
	})
}
