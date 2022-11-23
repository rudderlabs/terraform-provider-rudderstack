package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	c.Destinations.Register("adobe_analytics", c.ConfigMeta{
		APIType: "ADOBE_ANALYTICS",
		Properties: []c.ConfigProperty{
			c.Simple("trackingServerUrl", "tracking_server_url", c.SkipZeroValue),
			c.Simple("trackingServerSecureUrl", "tracking_server_secure_url", c.SkipZeroValue),
			c.Simple("reportSuiteIds", "report_suite_ids"),
			c.Simple("sslHeartbeat", "ssl_heartbeat", c.SkipZeroValue),
			c.Simple("heartbeatTrackingServerUrl", "heartbeat_tracking_server_url", c.SkipZeroValue),
			c.Simple("useUtf8Charset", "use_utf8_charset", c.SkipZeroValue),
			c.Simple("useSecureServerSide", "use_secure_server_side", c.SkipZeroValue),
			c.Simple("proxyNormalUrl", "proxy_normal_url", c.SkipZeroValue),
			c.Simple("proxyHeartbeatUrl", "proxy_heartbeat_url", c.SkipZeroValue),
			c.ArrayWithObjects("eventsToTypes", "events_to_types", map[string]string{
				"from": "from",
				"to":   "to",
			}),
			c.Simple("marketingCloudOrgId", "marketing_cloud_org_id", c.SkipZeroValue),
			c.Simple("dropVisitorId", "drop_visitor_id", c.SkipZeroValue),
			c.Simple("timestampOption", "timestamp_option", c.SkipZeroValue),
			c.Simple("timestampOptionalReporting", "timestamp_optional_reporting", c.SkipZeroValue),
			c.Simple("noFallbackVisitorId", "no_fallback_visitor_id", c.SkipZeroValue),
			c.Simple("preferVisitorId", "prefer_visitor_id", c.SkipZeroValue),
			c.ArrayWithObjects("rudderEventsToAdobeEvents", "rudder_events_to_adobe_events", map[string]string{
				"from": "from",
				"to":   "to",
			}),
			c.Simple("trackPageName", "track_page_name", c.SkipZeroValue),
			c.ArrayWithObjects("contextDataMapping", "context_data_mapping", map[string]string{
				"from": "from",
				"to":   "to",
			}),
			c.Simple("contextDataPrefix", "context_data_prefix", c.SkipZeroValue),
			c.Simple("useLegacyLinkName", "use_legacy_link_name", c.SkipZeroValue),
			c.Simple("pageNameFallbackTostring", "page_name_fallback_tostring", c.SkipZeroValue),
			c.ArrayWithObjects("mobileEventMapping", "mobile_event_mapping", map[string]string{
				"from": "from",
				"to":   "to",
			}),
			c.Simple("sendFalseValues", "send_false_values", c.SkipZeroValue),
			c.ArrayWithObjects("eVarMapping", "e_var_mapping", map[string]string{
				"from": "from",
				"to":   "to",
			}),
			c.ArrayWithObjects("hierMapping", "hier_mapping", map[string]string{
				"from": "from",
				"to":   "to",
			}),
			c.ArrayWithObjects("listMapping", "list_mapping", map[string]string{
				"from": "from",
				"to":   "to",
			}),
			c.ArrayWithObjects("listDelimiter", "list_delimiter", map[string]string{
				"from": "from",
				"to":   "to",
			}),

			c.ArrayWithObjects("customPropsMapping", "custom_props_mapping", map[string]string{
				"from": "from",
				"to":   "to",
			}),
			c.ArrayWithObjects("propsDelimiter", "props_delimiter", map[string]string{
				"from": "from",
				"to":   "to",
			}),
			c.ArrayWithObjects("eventMerchEventToAdobeEvent", "event_merch_event_to_adobe_event", map[string]string{
				"from": "from",
				"to":   "to",
			}),
			c.ArrayWithObjects("eventMerchProperties", "event_merch_properties", map[string]string{
				"eventMerchProperties": "property",
			}),

			c.ArrayWithObjects("productMerchEventToAdobeEvent", "product_merch_event_to_adobe_event", map[string]string{
				"from": "from",
				"to":   "to",
			}),
			c.ArrayWithObjects("productMerchProperties", "product_merch_properties", map[string]string{
				"productMerchProperties": "property",
			}),

			c.ArrayWithObjects("productMerchEvarsMap", "product_merch_evars_map", map[string]string{
				"from": "from",
				"to":   "to",
			}),
			c.Simple("productIdentifier", "product_identifier", c.SkipZeroValue),
			c.Simple("useNativeSDK.web", "use_native_sdk.0.web"),
			c.Simple("useNativeSDK.ios", "use_native_sdk.0.ios"),
			c.Simple("useNativeSDK.android", "use_native_sdk.0.android"),
			c.Simple("useNativeSDK.reactnative", "use_native_sdk.0.react_native"),
			c.ArrayWithStrings("whitelistedEvents", "eventName", "event_filtering.0.whitelist"),
			c.ArrayWithStrings("blacklistedEvents", "eventName", "event_filtering.0.blacklist"),
			c.Discriminator("eventFilteringOption", c.DiscriminatorValues{
				"event_filtering.0.whitelist": "whitelistedEvents",
				"event_filtering.0.blacklist": "blacklistedEvents",
			}),
			c.ArrayWithStrings("oneTrustCookieCategories.web", "oneTrustCookieCategory", "onetrust_cookie_categories.0.web"),
		},
		ConfigSchema: map[string]*schema.Schema{
			"tracking_server_url": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Enter your Tracking Server URL",
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
			},
			"tracking_server_secure_url": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Enter your Tracking Server Secure URL",
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
			},
			"report_suite_ids": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Enter your Report Suite ID(s). You can add multiple report suite ids by separated by commas.",
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,300})$"),
			},
			"ssl_heartbeat": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Check for Heartbeat calls to be made over https",
			},
			"heartbeat_tracking_server_url": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Enter your Heartbeat Tracking Server URL",
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
			},
			"use_utf8_charset": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Use UTF-8 charset",
			},
			"use_secure_server_side": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Use Secure URL for Server-side",
			},
			"proxy_normal_url": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Enter your Adobe Analytics Javascript SDK URL",
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
			},
			"proxy_heartbeat_url": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Enter your Adobe Analytics Hearbeat SDK URL",
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
			},

			"events_to_types": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "You can map your Rudder video events with types of Video Events",
				ConfigMode:  schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Provide the Video Event Name.",
							ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
						},
						"to": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Enter the Type of Video Event",
							ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
						},
					},
				},
			},
			"marketing_cloud_org_id": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Enter your Marketing Cloud Organization Id.",
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
			},
			"drop_visitor_id": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Check to Drop Visitor Id.",
			},
			"timestamp_option": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Enter your Timestamp Option.",
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(disabled|hybrid|optional|enabled)$"),
			},
			"timestamp_optional_reporting": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Check to send both Timestamp and VisitorID for Timestamp Optional Reporting Suites",
			},
			"no_fallback_visitor_id": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Check to enable no Fallbacks for Visitor ID",
			},
			"prefer_visitor_id": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Check to prefer Visitor Id",
			},
			"rudder_events_to_adobe_events": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "You can map Rudder Events to Adobe Custom Events.",
				ConfigMode:  schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the Event Name",
							ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
						},
						"to": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the Adobe Custom Event",
							ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
						},
					},
				},
			},
			"track_page_name": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Check to enable pageName for Track Events",
			},
			"context_data_mapping": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "You can map Rudder Context data to Adobe Context Data",
				ConfigMode:  schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the Context Data path.",
							ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
						},
						"to": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the Adobe Context Data property name",
							ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
						},
					},
				},
			},
			"context_data_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enter your prefix to add before all contextData property.",
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
			},
			"use_legacy_link_name": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Check to use Legacy LinkName",
			},
			"page_name_fallback_tostring": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Check to allow Page Name Fallback to Screen",
			},
			"mobile_event_mapping": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "You can map Rudder Mobile events.",
				ConfigMode:  schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the Context Data path.",
							ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
						},
						"to": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the Adobe Context Data property name",
							ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
						},
					},
				},
			},
			"send_false_values": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Check to allow sending false value from properties",
			},
			"e_var_mapping": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "You can map Rudder Properties to Adobe eVars",
				ConfigMode:  schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the Rudder Property",
							ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
						},
						"to": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the eVar Index",
							ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
						},
					},
				},
			},
			"hier_mapping": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "You can map Rudder Properties to Adobe Hierarchy properties",
				ConfigMode:  schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the Rudder Property",
							ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
						},
						"to": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the eVar Index",
							ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
						},
					},
				},
			},
			"list_mapping": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "You can map Rudder Properties to Adobe list properties",
				ConfigMode:  schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the Rudder Property as an array/string seperated by commas",
							ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
						},
						"to": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the eVar Index",
							ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
						},
					},
				},
			},
			"list_delimiter": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "You can map your Rudder Property with Delimiters for list properties",
				ConfigMode:  schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the Rudder Property.",
							ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
						},
						"to": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the List Delimiter.",
							ValidateDiagFunc: c.StringMatchesRegexp("^$|(^env[.].+)|^(\\||:|,|;|\\/)$"),
						},
					},
				},
			},
			"custom_props_mapping": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "You can map Rudder Properties to Adobe Custom properties",
				ConfigMode:  schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the Rudder Property.",
							ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
						},
						"to": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the prop Index.",
							ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
						},
					},
				},
			},
			"props_delimiter": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "You can map your Rudder Property with Delimiters for Adobe Custom properties",
				ConfigMode:  schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the Rudder Property.",
							ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
						},
						"to": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the List Delimiter.",
							ValidateDiagFunc: c.StringMatchesRegexp("^$|(^env[.].+)|^(\\||:|,|;|\\/)$"),
						},
					},
				},
			},
			"event_merch_event_to_adobe_event": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "You can map Rudder Events to Adobe Event Merchandise events.",
				ConfigMode:  schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the Rudder Event.",
							ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
						},
						"to": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the Adobe Event.",
							ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
						},
					},
				},
			},
			"event_merch_properties": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Currency/Incremental properties to add to merchandise events at event level",
				ConfigMode:  schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"property": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the property.",
							ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
						},
					},
				},
			},
			"product_merch_event_to_adobe_event": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "You can map Rudder Events to Adobe Product Merchandise events",
				ConfigMode:  schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the Rudder Event.",
							ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
						},
						"to": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the Adobe Event.",
							ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
						},
					},
				},
			},
			"product_merch_properties": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Currency/Incremental properties to add to merchandise events at product level",
				ConfigMode:  schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"property": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the property.",
							ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
						},
					},
				},
			},
			"product_merch_evars_map": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "You can map Rudder Properties to eVars at product level",
				ConfigMode:  schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the Rudder Event.",
							ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
						},
						"to": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the eVar Index.",
							ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
						},
					},
				},
			},
			"product_identifier": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enter your Product Identifier",
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(name|id|sku)$"),
			},
			"use_native_sdk": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Enable this setting to send events to Adobe Analytics via the device mode.",
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
						"react_native": {
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
				Description: "This option allows you filter the events you want to send to Adobe Analytics.",
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
