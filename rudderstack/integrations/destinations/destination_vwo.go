package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	c.Destinations.Register("vwo", c.ConfigMeta{
		APIType: "VWO",
		Properties: []c.ConfigProperty{
			c.Simple("accountId", "account_id"),
			c.Simple("isSPA", "is_spa", c.SkipZeroValue),
			c.Simple("sendExperimentTrack", "send_experiment_track", c.SkipZeroValue),
			c.Simple("sendExperimentIdentify", "send_experiment_identify", c.SkipZeroValue),
			c.Simple("libraryTolerance", "library_tolerance", c.SkipZeroValue),
			c.Simple("settingsTolerance", "settings_tolerance", c.SkipZeroValue),
			c.Simple("useExistingJquery", "use_existing_jquery", c.SkipZeroValue),
			c.ArrayWithStrings("whitelistedEvents", "eventName", "event_filtering.0.whitelist"),
			c.ArrayWithStrings("blacklistedEvents", "eventName", "event_filtering.0.blacklist"),
			c.Discriminator("eventFilteringOption", c.DiscriminatorValues{
				"event_filtering.0.whitelist": "whitelistedEvents",
				"event_filtering.0.blacklist": "blacklistedEvents",
			}),
			c.Simple("useNativeSDK.web", "use_native_sdk.0.web"),
			c.ArrayWithStrings("oneTrustCookieCategories", "oneTrustCookieCategory", "onetrust_cookie_categories"),
		},
		ConfigSchema: map[string]*schema.Schema{
			"account_id": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Enter your VWO account ID.",
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"is_spa": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable this setting if the page is a single page application (SPA).",
			},
			"send_experiment_track": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable this setting to send the experiment data as `track` events.",
			},
			"send_experiment_identify": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable this setting to send the experiments viewed as `identify` traits.",
			},
			"library_tolerance": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Enter the value for the library tolerance setting.",
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"settings_tolerance": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Enter the value for the setting tolerance.",
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"use_existing_jquery": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable this setting to use the existing jQuery.",
			},
			"event_filtering": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Specify which events should be blocked or allowed to flow through to VWO.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"whitelist": {
							Type:         schema.TypeList,
							Optional:     true,
							Description:  "Enter the event nams to be allowlisted.",
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
				Description: "Enable this setting to send the events via the device mode.",
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
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specify the OneTrust category name for mapping the OneTrust consent settings to RudderStack's consent purposes.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	})
}
