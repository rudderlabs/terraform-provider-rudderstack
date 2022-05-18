package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	c.Destinations.Register("google_analytics4", c.ConfigMeta{
		APIType: "GA4",
		Properties: []c.ConfigProperty{
			c.Simple("apiSecret", "api_secret", c.SkipZeroValue),
			c.Simple("typesOfClient", "types_of_client", c.SkipZeroValue),
			c.Simple("measurementId", "measurement_id", c.SkipZeroValue),
			c.Simple("firebaseAppId", "firebase_app_id", c.SkipZeroValue),
			c.Simple("blockPageViewEvent", "block_page_view_event", c.SkipZeroValue),
			c.Simple("extendPageViewParams", "extend_page_view_params", c.SkipZeroValue),
			c.Simple("sendUserId", "send_user_id", c.SkipZeroValue),
			c.ArrayWithStrings("whitelistedEvents", "eventName", "event_filtering.0.whitelist"),
			c.ArrayWithStrings("blacklistedEvents", "eventName", "event_filtering.0.blacklist"),
			c.Discriminator("eventFilteringOption", c.DiscriminatorValues{
				"event_filtering.0.whitelist": "whitelistedEvents",
				"event_filtering.0.blacklist": "blacklistedEvents",
			}),
			c.Simple("useNativeSDK.web", "use_native_sdk.0.web"),
			c.ArrayWithStrings("oneTrustCookieCategories.web", "oneTrustCookieCategory", "onetrust_cookie_categories.0.web"),
		},
		ConfigSchema: map[string]*schema.Schema{
			"api_secret": {
				Type:             schema.TypeString,
				Optional:         true,
				Sensitive:        true,
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
			},
			"types_of_client": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(gtag|firebase)$"),
			},
			"measurement_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(G-.{1,100})$|^$"),
			},
			"firebase_app_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
			},
			"block_page_view_event": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"extend_page_view_params": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"send_user_id": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"event_filtering": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"whitelist": {
							Type:         schema.TypeList,
							Optional:     true,
							ExactlyOneOf: []string{"config.0.event_filtering.0.whitelist", "config.0.event_filtering.0.blacklist"},
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"blacklist": {
							Type:         schema.TypeList,
							Optional:     true,
							ExactlyOneOf: []string{"config.0.event_filtering.0.whitelist", "config.0.event_filtering.0.blacklist"},
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"use_native_sdk": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"onetrust_cookie_categories": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	})
}
