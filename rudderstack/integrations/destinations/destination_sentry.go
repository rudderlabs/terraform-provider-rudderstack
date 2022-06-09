package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	c.Destinations.Register("sentry", c.ConfigMeta{
		APIType: "SENTRY",
		Properties: []c.ConfigProperty{
			c.Simple("dsn", "dsn"),
			c.Simple("environment", "environment", c.SkipZeroValue),
			c.Simple("customVersionProperty", "custom_version_property", c.SkipZeroValue),
			c.Simple("release", "release", c.SkipZeroValue),
			c.Simple("serverName", "server_name", c.SkipZeroValue),
			c.Simple("logger", "logger", c.SkipZeroValue),
			c.Simple("debugMode", "debug_mode"),
			c.ArrayWithStrings("ignoreErrors", "ignoreErrors", "ignore_errors"),
			c.ArrayWithStrings("includePaths", "includePaths", "include_paths"),
			c.ArrayWithStrings("allowUrls", "allowUrls", "allow_urls"),
			c.ArrayWithStrings("denyUrls", "denyUrls", "deny_urls"),
			c.Simple("useNativeSDK.web", "use_native_sdk.0.web"),
			c.ArrayWithStrings("whitelistedEvents", "eventName", "event_filtering.0.whitelist"),
			c.ArrayWithStrings("blacklistedEvents", "eventName", "event_filtering.0.blacklist"),
			c.ArrayWithStrings("oneTrustCookieCategories.web", "oneTrustCookieCategory", "onetrust_cookie_categories.0.web"),
		},
		ConfigSchema: map[string]*schema.Schema{
			"dsn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Enter the public DSN of your Sentry project. This is a mandatory field.",
				// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"environment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enter the value you want RudderStack to set as the environment configuration in your Sentry dashboard.",
				// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"custom_version_property": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This field helps you dynamically track the application version in Sentry.",
				// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"release": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This field is used for tracking your application's version in Sentry.",
				// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"server_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This option is used to track the host on which the client is running.",
				// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"logger": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Set the name you want Sentry to use as logger.",
				// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"debug_mode": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If enabled, no events are sent to your Sentry instance.",
			},
			"ignore_errors": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "This option refers to a list of error messages that you do not want Sentry to notify you.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"include_paths": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "This field should contain the regex patterns of URLs that are part of the app in the stack trace.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"allow_urls": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"deny_urls": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "This is the list of the regex patterns or exact URL strings - from which the errors need to be exclusively sent to Sentry.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
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
			"event_filtering": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "With this option, you can determine which events are blocked or allowed to flow through to Sentry.",
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
				MaxItems:    1,
				Optional:    true,
				Description: "Specify the OneTrust category name for mapping the OneTrust consent settings to RudderStack's consent purposes.",
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
