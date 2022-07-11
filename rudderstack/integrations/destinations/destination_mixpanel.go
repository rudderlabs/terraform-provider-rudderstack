package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	c.Destinations.Register("mixpanel", c.ConfigMeta{
		APIType: "MP",
		Properties: []c.ConfigProperty{
			c.Simple("token", "token"),
			c.Simple("apiSecret", "api_secret"),
			c.Simple("dataResidency", "data_residency"),
			c.Simple("people", "people"),
			c.Simple("setAllTraitsByDefault", "set_all_traits_by_default"),
			c.Simple("consolidatedPageCalls", "consolidated_page_calls"),
			c.Simple("trackCategorizedPages", "track_categorized_pages"),
			c.Simple("trackNamedPages", "track_named_pages"),
			c.Simple("sourceName", "source_name"),
			c.Simple("crossSubdomainCookie", "cross_subdomain_cookie"),
			c.Simple("persistence", "persistence"),
			c.Simple("secureCookie", "secure_cookie"),
			c.ArrayWithStrings("superProperties", "property", "super_properties"),
			c.ArrayWithStrings("peopleProperties", "property", "people_properties"),
			c.ArrayWithStrings("eventIncrements", "property", "event_increments"),
			c.ArrayWithStrings("propIncrements", "property", "prop_increments"),
			c.ArrayWithStrings("groupKeySettings", "groupKey", "group_key_settings"),
			c.Simple("useNativeSDK.web", "use_native_sdk.0.web"),
			c.ArrayWithStrings("whitelistedEvents", "eventName", "event_filtering.0.whitelist"),
			c.ArrayWithStrings("blacklistedEvents", "eventName", "event_filtering.0.blacklist"),
			c.ArrayWithStrings("oneTrustCookieCategories.web", "oneTrustCookieCategory", "onetrust_cookie_categories.0.web"),
		},
		ConfigSchema: map[string]*schema.Schema{
			"token": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Mixpanel API Token",
				Sensitive:        true,
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"api_secret": {
				Type:             schema.TypeString,
				Optional:         true,
				Sensitive:        true,
				Description:      "Mixpanel API secret",
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
			},
			"data_residency": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Mixpanel Server region either us/eu",
				ValidateDiagFunc: c.StringMatchesRegexp("^(us|eu)$"),
			},
			"people": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Boolean flag to send all of your identify calls to Mixpanel's People feature ",
			},
			"set_all_traits_by_default": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "While this is checked, our integration automatically sets all traits on identify calls as super properties and people properties if Mixpanel People is checked as well.",
			},
			"consolidated_page_calls": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "This will track Loaded a Page events to Mixpanel for all page method calls. We enable this by default as it's how Mixpanel suggests sending these calls.",
			},
			"track_categorized_pages": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "This will track events to Mixpanel for page method calls that have a category associated with them. For example page('Docs', 'Index') would translate to Viewed Docs Index Page.",
			},
			"track_named_pages": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "This will track events to Mixpanel for page method calls that have a name associated with them. For example page('Signup') would translate to Viewed Signup Page.",
			},
			"source_name": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "This value, if it's not blank, will be sent as rudderstack_source_name to Mixpanel for every event/page/screen call.",
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
			},
			"cross_subdomain_cookie": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "This will allow the Mixpanel cookie to persist between different pages of your application.",
			},
			"persistence": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Choose persistence for Mixpanel SDK. One of none|cookie|localStorage",
				ValidateDiagFunc: c.StringMatchesRegexp("^(none|cookie|localStorage)$"),
			},
			"secure_cookie": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "This will mark the Mixpanel cookie as secure, meaning it will only be transmitted over https",
			},
			"super_properties": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Property to send as super Properties",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"people_properties": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Traits to set as People Properties",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"event_increments": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Events to increment in People",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"prop_increments": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Properties to increment in People",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"group_key_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Group Key",
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
