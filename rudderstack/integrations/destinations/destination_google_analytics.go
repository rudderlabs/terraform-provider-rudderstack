package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	c.Destinations.Register("google_analytics", c.ConfigMeta{
		APIType: "GA",
		Properties: []c.ConfigProperty{
			c.Simple("trackingID", "tracking_id"),
			c.Simple("doubleClick", "double_click", c.SkipZeroValue),
			c.Simple("enhancedLinkAttribution", "enhanced_link_attribution", c.SkipZeroValue),
			c.Simple("includeSearch", "include_search", c.SkipZeroValue),
			c.Simple("serverSideIdentifyEventCategory", "server_side_identify.0.event_category"),
			c.Simple("serverSideIdentifyEventAction", "server_side_identify.0.event_action"),
			c.Discriminator("enableServerSideIdentify", c.DiscriminatorValues{
				"server_side_identify.0.event_category": true,
			}),
			c.Simple("disableMd5", "disable_md5", c.SkipZeroValue),
			c.Simple("anonymizeIp", "anonymize_ip", c.SkipZeroValue),
			c.Simple("enhancedEcommerce", "enhanced_ecommerce", c.SkipZeroValue),
			c.Simple("nonInteraction", "non_interaction", c.SkipZeroValue),
			c.Simple("sendUserId", "send_user_id", c.SkipZeroValue),
			c.ArrayWithStrings("whitelistedEvents", "eventName", "event_filtering.0.whitelist"),
			c.ArrayWithStrings("blacklistedEvents", "eventName", "event_filtering.0.blacklist"),
			c.Discriminator("eventFilteringOption", c.DiscriminatorValues{
				"event_filtering.0.whitelist": "whitelistedEvents",
				"event_filtering.0.blacklist": "blacklistedEvents",
			}),
			c.Simple("useNativeSDK.web", "use_native_sdk.0.web"),
			c.Simple("trackCategorizedPages.web", "track_categorized_pages.0.web"),
			c.Simple("trackNamedPages.web", "track_named_pages.0.web"),
			c.Simple("sampleRate.web", "sample_rate.0.web"),
			c.Simple("siteSpeedSampleRate.web", "site_speed_sample_rate.0.web"),
			c.ArrayWithStrings("resetCustomDimensionsOnPage.web", "resetCustomDimensionsOnPage", "reset_custom_dimensions_on_page.0.web"),
			c.Simple("setAllMappedProps.web", "set_all_mapped_props.0.web"),
			c.Simple("domain.web", "domain.0.web"),
			c.Simple("optimize.web", "optimize.0.web"),
			c.Simple("useGoogleAmpClientId.web", "use_google_amp_client_id.0.web"),
			c.Simple("namedTracker.web", "named_tracker.0.web"),
			c.ArrayWithStrings("oneTrustCookieCategories.web", "oneTrustCookieCategory", "onetrust_cookie_categories.0.web"),
			c.Simple("dimensions", "dimensions", c.SkipZeroValue),
			c.Simple("contentGroupings", "content_groupings", c.SkipZeroValue),
		},
		ConfigSchema: map[string]*schema.Schema{
			"tracking_id": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Enter your Google Analytics Tracking ID.",
				ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^(UA|YT|MO)-\\d+-\\d{0,100}$)"),
			},
			"double_click": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enhanced_link_attribution": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable this setting to activate the Google Analytics enhanced link attribution feature.",
			},
			"include_search": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable this setting to include the querystring in `page` views.",
			},
			"server_side_identify": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_category": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Enter the server-side `identify` event category.",
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^(.{0,100})$)"),
						},
						"event_action": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Enter the server-side `identify` event action.",
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^(.{0,100})$)"),
						},
					},
				},
			},
			"disable_md5": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable this setting to disable client ID MD5 encryption.",
			},
			"anonymize_ip": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enabling this setting anonymizes your IP address information.",
			},
			"enhanced_ecommerce": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable this setting to activate the enhanced e-commerce feature.",
			},
			"non_interaction": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable this setting to add the non-interaction flag to all the events.",
			},
			"send_user_id": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable this setting to send the `userId` to Google Analytics.",
			},
			"event_filtering": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "This setting allows you to specify which events should be blocked or allowed to flow through to Google Analytics.",
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
							Description:  "Enter the event names to be blacklisted..",
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
				Description: "Enable this setting to send the events via the web device mode.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"track_categorized_pages": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Enable this setting to track categorized pages.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"track_named_pages": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Enable this setting to track named pages.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"sample_rate": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Enter the sample rate.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^(.{0,100})$)"),
						},
					},
				},
			},
			"site_speed_sample_rate": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Enter the site speed sample rate.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^(.{0,100})$)"),
						},
					},
				},
			},
			"reset_custom_dimensions_on_page": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Use this field to reset the dimensions for the `page` calls.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"set_all_mapped_props": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Use this field to set all the mapped properties.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"domain": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Enter your cookie domain name.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^(.{0,100})$)"),
						},
					},
				},
			},
			"optimize": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Enter your Google Optimize Container ID.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^(.{0,100})$)"),
						},
					},
				},
			},
			"use_google_amp_client_id": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Enable this setting to use the Google AMP Client ID",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},
			"named_tracker": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Enable this setting to send events with the `track` name `rudderGATracker`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web": {
							Type:     schema.TypeBool,
							Required: true,
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
			"content_groupings": {
				Type:       schema.TypeList,
				MaxItems:   10,
				Optional:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
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
							ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
						},
					},
				},
			},
			"dimensions": {
				Type:       schema.TypeList,
				MaxItems:   10,
				Optional:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
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
							ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
						},
					},
				},
			},
		},
	})
}
