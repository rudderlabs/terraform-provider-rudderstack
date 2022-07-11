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
			c.ArrayWithObjects("superProperties", "super_properties", map[string]string{
				"property": "property",
			}),
			c.ArrayWithObjects("peopleProperties", "people_properties", map[string]string{
				"property": "property",
			}),
			c.ArrayWithObjects("eventIncrements", "event_increments", map[string]string{
				"property": "property",
			}),
			c.ArrayWithObjects("propIncrements", "prop_increments", map[string]string{
				"property": "property",
			}),
			c.ArrayWithObjects("groupKeySettings", "group_key_settings", map[string]string{
				"groupKey": "group_key",
			}),
			c.Simple("useNativeSDK.web", "use_native_sdk.0.web"),
			c.Discriminator("eventFilteringOption", c.DiscriminatorValues{
				"event_filtering.0.whitelist": "whitelistedEvents",
				"event_filtering.0.blacklist": "blacklistedEvents",
			}),
			c.ArrayWithStrings("oneTrustCookieCategories.web", "oneTrustCookieCategory", "onetrust_cookie_categories.0.web"),
		},
		ConfigSchema: map[string]*schema.Schema{
			"token": {
				Type: schema.TypeString,
				Required: true,
				Description: "Mixpanel API Token",
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"api_secret": {
				Type: schema.TypeString,
				Optional: true,
				Description: " Mixpanel API secret",
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
			},
			"data_residency": {
				Type: schema.TypeString,
				Optional: true,
				Description: "Mixpanel Server region either us/eu",	
				ValidateDiagFunc: c.StringMatchesRegexp("^(us|eu)$"),			
			},
			"people": {
				Type: schema.TypeBool,
				Optional: true,
				Description: "Boolean flag to send all of your identify calls to Mixpanel's People feature ",
			},
			"set_all_traits_by_default":{
				Type: schema.TypeBool,
				Optional: true,
				Description: "While this is checked, our integration automatically sets all traits on identify calls as super properties and people properties if Mixpanel People is checked as well.",
			},
			"consolidated_page_calls": {
				Type: schema.TypeBool,
				Optional: true,
				Description: "This will track Loaded a Page events to Mixpanel for all page method calls. We enable this by default as it's how Mixpanel suggests sending these calls.",
			},
			"track_categorized_pages": {
				Type: schema.TypeBool,
				Optional: true,
				Description: "This will track events to Mixpanel for page method calls that have a category associated with them. For example page('Docs', 'Index') would translate to Viewed Docs Index Page.",
			},
			"track_named_pages": {
				Type: schema.TypeBool,
				Optional: true,
				Description: "This will track events to Mixpanel for page method calls that have a name associated with them. For example page('Signup') would translate to Viewed Signup Page.",
			},
			"source_name": {
				Type: schema.TypeString,
				Optional: true,
				Description: "This value, if it's not blank, will be sent as rudderstack_source_name to Mixpanel for every event/page/screen call.",
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
			},
			"cross_subdomain_cookie": {
				Type: schema.TypeBool,
				Optional: true,
				Description: "This will allow the Mixpanel cookie to persist between different pages of your application.",
			},
			"persistence": {
				Type: schema.TypeString,
				Optional: true,
				Description: "Choose persistence for Mixpanel SDK. One of none|cookie|localStorage",
				ValidateDiagFunc: c.StringMatchesRegexp("^(none|cookie|localStorage)$"),
				
			},
			"secure_cookie": {
				Type: schema.TypeBool,
				Optional: true,
				Description: "This will mark the Mixpanel cookie as secure, meaning it will only be transmitted over https",
			},
			"super_properties": {
				Type: schema.TypeList,
				Optional: true,
				Description: "Property to send as super Properties",
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"property": {
							Type: schema.TypeString,
							Required: true,
							Description: "",
							ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
						},
					},
				},
			},
			"people_properties": {
				Type: schema.TypeList,
				Optional: true,
				Description: "",
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"property": {
							Type: schema.TypeString,
							Required: true,
							Description: "",
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{0,100})$"),
						},
					},
				},
			},
			"event_increments": {
				Type: schema.TypeList,
				Optional: true,
				Description: "",
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"property": {
							Type: schema.TypeString,
							Required: true,
							Description: "",
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{0,100})$"),
						},
					},
				},
			},
			"prop_increments": {
				Type: schema.TypeList,
				Optional: true,
				Description: "",
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"property": {
							Type: schema.TypeString,
							Required: true,
							Description: "",
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{0,100})$"),
						},
					},
				},
			},
			"group_key_settings": {
				Type: schema.TypeList,
				Optional: true,
				Description: "",
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_key": {
							Type: schema.TypeString,
							Required: true,
							Description: "",
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{0,100})$"),
						},
					},
				},
			},
			"use_native_sdk": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Enable this setting to send events to Mixpanel via the device mode.",
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
				Description: "This option allows you filter the events you want to send to Amplitude.",
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
