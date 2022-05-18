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
		},
		ConfigSchema: map[string]*schema.Schema{
			"dsn": {
				Type:     schema.TypeString,
				Required: true,
				// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"environment": {
				Type:     schema.TypeString,
				Optional: true,
				// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"custom_version_property": {
				Type:     schema.TypeString,
				Optional: true,
				// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"release": {
				Type:     schema.TypeString,
				Optional: true,
				// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"server_name": {
				Type:     schema.TypeString,
				Optional: true,
				// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"logger": {
				Type:     schema.TypeString,
				Optional: true,
				// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"debug_mode": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ignore_errors": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"include_paths": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"allow_urls": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"deny_urls": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
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
		},
	})
}
