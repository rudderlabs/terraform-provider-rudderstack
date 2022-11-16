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
				"from": "label",
				"to":   "name",
			}),
			c.Simple("marketingCloudOrgId", "marketing_cloud_org_id", c.SkipZeroValue),
			c.Simple("dropVisitorId", "drop_visitor_id", c.SkipZeroValue),
			c.Simple("timestampOption", "timestamp_option", c.SkipZeroValue),
			c.Simple("timestampOptionalReporting", "timestamp_optional_reporting", c.SkipZeroValue),
			c.Simple("noFallbackVisitorId", "no_fallback_visitor_id", c.SkipZeroValue),
			c.Simple("preferVisitorId", "prefer_visitor_id", c.SkipZeroValue),
			c.ArrayWithObjects("rudderEventsToAdobeEvents", "rudder_events_to_adobe_events", map[string]string{
				"from": "label",
				"to":   "name",
			}),
			c.Simple("trackPageName", "track_page_name", c.SkipZeroValue),
			c.ArrayWithObjects("contextDataMapping", "context_data_mapping", map[string]string{
				"from": "label",
				"to":   "name",
			}),
			c.Simple("contextDataPrefix", "context_data_prefix", c.SkipZeroValue),
			c.Simple("useLegacyLinkName", "use_legacy_link_name", c.SkipZeroValue),
			c.Simple("pageNameFallbackTostring", "page_name_fallback_tostring", c.SkipZeroValue),
			c.ArrayWithObjects("mobileEventMapping", "mobile_event_mapping", map[string]string{
				"from": "label",
				"to":   "name",
			}),
			c.Simple("sendFalseValues", "send_false_values", c.SkipZeroValue),
			c.ArrayWithObjects("eVarMapping", "e_var_mapping", map[string]string{
				"from": "label",
				"to":   "name",
			}),
			c.ArrayWithObjects("hierMapping", "hier_mapping", map[string]string{
				"from": "label",
				"to":   "name",
			}),
			c.ArrayWithObjects("listMapping", "list_mapping", map[string]string{
				"from": "label",
				"to":   "name",
			}),
			c.ArrayWithObjects("listDelimiter", "list_delimiter", map[string]string{
				"from": "label",
				"to":   "name",
			}),

			c.ArrayWithObjects("customPropsMapping", "custom_props_mapping", map[string]string{
				"from": "label",
				"to":   "name",
			}),
			c.ArrayWithObjects("propsDelimiter", "props_delimiter", map[string]string{
				"from": "label",
				"to":   "name",
			}),
			c.ArrayWithObjects("eventMerchEventToAdobeEvent", "event_merch_event_to_adobe_event", map[string]string{
				"from": "label",
				"to":   "name",
			}),
			c.ArrayWithStrings("eventMerchProperties", "event_merch_properties", "event_merch_properties"),

			c.ArrayWithObjects("productMerchEventToAdobeEvent", "product_merch_event_to_adobe_event", map[string]string{
				"from": "label",
				"to":   "name",
			}),
			c.ArrayWithStrings("productMerchProperties", "product_merch_properties", "product_merch_properties"),

			c.ArrayWithObjects("productMerchEvarsMap", "product_merch_evars_map", map[string]string{
				"from": "label",
				"to":   "name",
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
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enter your Tracking Server URL",
			},
			"tracking_server_secure_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enter your Tracking Server Secure URL",
			},
			"report_suite_ids": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Enter your Report Suite ID(s). You can add multiple report suite ids by separated by commas.",
			},
			"ssl_heartbeat": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Check for Heartbeat calls to be made over https",
			},
			"heartbeat_tracking_server_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enter your Heartbeat Tracking Server URL",
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
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enter your Adobe Analytics Javascript SDK URL",
			},
			"proxy_heartbeat_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enter your Adobe Analytics Hearbeat SDK URL",
			},

			"events_to_types": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "You can map your Rudder video events with types of Video Events",
				ConfigMode:  schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"label": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Provide the Video Event Name.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the Type of Video Event",
						},
					},
				},
			},
			"marketing_cloud_org_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enter your Marketing Cloud Organization Id.",
			},
			"drop_visitor_id": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Check to Drop Visitor Id.",
			},
			"timestamp_option": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enter your Timestamp Option.",
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
						"label": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the Event Name",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the Adobe Custom Event",
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
						"label": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the Context Data path.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the Adobe Context Data property name",
						},
					},
				},
			},
			"context_data_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enter your prefix to add before all contextData property.",
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
						"label": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the Context Data path.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the Adobe Context Data property name",
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
						"label": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the Rudder Property",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the eVar Index",
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
						"label": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the Rudder Property",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the eVar Index",
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
						"label": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the Rudder Property as an array/string seperated by commas",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the eVar Index",
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
						"label": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the Rudder Property.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the List Delimiter.",
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
						"label": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the Rudder Property.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the prop Index.",
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
						"label": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the Rudder Property.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the List Delimiter.",
						},
					},
				},
			},
			"event_merch_event_to_adobe_event": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "You can map Rudder Events to Adobe Merchandise events.",
				ConfigMode:  schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"label": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the Rudder Event.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the Adobe Event.",
						},
					},
				},
			},
			"event_merch_properties": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Currency/Incremental properties to add to merchandise events at event level",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"product_merch_event_to_adobe_event": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "You can map Rudder Events to Adobe Merchandise events",
				ConfigMode:  schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"label": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the Rudder Event.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the Adobe Event.",
						},
					},
				},
			},
			"product_merch_properties": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Currency/Incremental properties to add to merchandise events at product level",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"product_merch_evars_map": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "You can map Rudder Properties to eVars at product level",
				ConfigMode:  schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"label": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the Rudder Event.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the eVar Index.",
						},
					},
				},
			},
			"product_identifier": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enter your Product Identifier",
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
