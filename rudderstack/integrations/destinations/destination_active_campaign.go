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
				Description: "Enter your ActiveCampaign API URL. You can find it in your account in the Settings page under the Developer tab.",
				// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"api_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				Description: "Enter your ActiveCampaign API key.",
				// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"actid": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Enter your ActID here. To obtain the ActID unique to your ActiveCampaign account, go to Settings > Tracking > Event Tracking API.",
				// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"event_key": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Enter the event key unique to your ActiveCampaign account. To obtain the event key, go to your ActiveCampaign account > Settings > Tracking > Event Tracking.",
				// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
		},
	})
}
