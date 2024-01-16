package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	c.Destinations.Register("intercom", c.ConfigMeta{
		APIType: "INTERCOM",
		Properties: []c.ConfigProperty{
			c.Simple("appId", "app_id"),
			c.Simple("apiKey", "api_key"),
			c.Simple("useNativeSDK.web", "use_native_sdk.0.web", c.SkipZeroValue),
			c.Simple("useNativeSDK.ios", "use_native_sdk.0.ios", c.SkipZeroValue),
			c.Simple("useNativeSDK.android", "use_native_sdk.0.android", c.SkipZeroValue),
			c.Simple("collectContext", "collect_context", c.SkipZeroValue),
			c.Simple("sendAnonymousId", "send_anonymous_id", c.SkipZeroValue),
			c.Simple("updateLastRequestAt", "update_last_request_at", c.SkipZeroValue),
			c.Simple("mobileApiKeyAndroid.android", "mobile_api_key_android", c.SkipZeroValue),
			c.Simple("mobileApiKeyIOS.ios", "mobile_api_key_ios", c.SkipZeroValue),
			c.ArrayWithStrings("oneTrustCookieCategories", "oneTrustCookieCategory", "onetrust_cookie_categories"),
			c.ArrayWithStrings("whitelistedEvents", "eventName", "event_filtering.0.whitelist"),
			c.ArrayWithStrings("blacklistedEvents", "eventName", "event_filtering.0.blacklist"),
			c.Discriminator("eventFilteringOption", c.DiscriminatorValues{
				"event_filtering.0.whitelist": "whitelistedEvents",
				"event_filtering.0.blacklist": "blacklistedEvents",
			}),
		},
		ConfigSchema: map[string]*schema.Schema{
			"api_key": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Enter Access Token.",
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
			},
			"app_id": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Enter App Id.",
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
			},
			"mobile_api_key_android": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Enter Android API Key.",
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
			},
			"mobile_api_key_ios": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Enable iOS API Key.",
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
			},
			"collect_context": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "This setting enables including Context with Identify Calls.",
			},
			"send_anonymous_id": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "This setting enables sending AnonymousId as Secondary UserId.",
			},
			"update_last_request_at": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "This setting enables the last seen with the current time.",
			},
			"use_native_sdk": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Enable this setting to send the events through SDK.",
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
