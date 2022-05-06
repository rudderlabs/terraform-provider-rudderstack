package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	c.Destinations.Register("amplitude", c.ConfigMeta{
		APIType: "AM",
		Properties: []c.ConfigProperty{
			c.Simple("apiKey", "api_key"),
			c.Simple("apiSecret", "api_secret"),
			c.Simple("groupTypeTrait", "group_type_trait", c.SkipZeroValue),
			c.Simple("groupValueTrait", "group_value_trait", c.SkipZeroValue),
			c.Simple("trackAllPages", "track_all_pages", c.SkipZeroValue),
			c.Simple("trackCategorizedPages", "track_categorized_pages", c.SkipZeroValue),
			c.Simple("trackNamedPages", "track_named_pages", c.SkipZeroValue),
			c.Simple("trackProductsOnce", "track_products_once", c.SkipZeroValue),
			c.Simple("trackRevenuePerProduct", "track_revenue_per_product", c.SkipZeroValue),
			c.Simple("versionName", "version_name", c.SkipZeroValue),
			c.ArrayWithStrings("traitsToIncrement", "traits", "traits_to_increment"),
			c.ArrayWithStrings("traitsToSetOnce", "traits", "traits_to_set_once"),
			c.ArrayWithStrings("traitsToAppend", "traits", "traits_to_append"),
			c.ArrayWithStrings("traitsToPrepend", "traits", "traits_to_prepend"),
			c.Simple("useNativeSDK.web", "use_native_sdk.0.web"),
			c.Simple("useNativeSDK.ios", "use_native_sdk.0.ios"),
			c.Simple("useNativeSDK.android", "use_native_sdk.0.android"),
			c.Simple("useNativeSDK.reactnative", "use_native_sdk.0.react_native"),
			c.Simple("preferAnonymousIdForDeviceId.web", "prefer_anonymous_id_for_device_id.0.web"),
			c.Simple("deviceIdFromUrlParam.web", "device_id_from_url_param.0.web"),
			c.Simple("forceHttps.web", "force_https.0.web"),
			c.Simple("trackGclid.web", "track_gclid.0.web"),
			c.Simple("trackReferrer.web", "track_referrer.0.web"),
			c.Simple("saveParamsReferrerOncePerSession.web", "save_params_referrer_once_per_session.0.web"),
			c.Simple("trackUtmProperties.web", "track_utm_properties.0.web"),
			c.Simple("unsetParamsReferrerOnNewSession.web", "unset_params_referrer_on_new_session.0.web"),
			c.Simple("batchEvents.web", "batch_events.0.web"),
			c.ArrayWithStrings("whitelistedEvents", "eventName", "event_filtering.0.whitelist"),
			c.ArrayWithStrings("blacklistedEvents", "eventName", "event_filtering.0.blacklist"),
			c.Simple("eventUploadPeriodMillis.web", "event_upload_period_millis.0.web"),
			c.Simple("eventUploadPeriodMillis.android", "event_upload_period_millis.0.android"),
			c.Simple("eventUploadPeriodMillis.ios", "event_upload_period_millis.0.ios"),
			c.Simple("eventUploadPeriodMillis.reactnative", "event_upload_period_millis.0.react_native"),
			c.Simple("eventUploadThreshold.web", "event_upload_threshold.0.web"),
			c.Simple("eventUploadThreshold.android", "event_upload_threshold.0.android"),
			c.Simple("eventUploadThreshold.ios", "event_upload_threshold.0.ios"),
			c.Simple("eventUploadThreshold.reactnative", "event_upload_threshold.0.react_native"),
			c.Simple("mapDeviceBrand", "map_device_brand", c.SkipZeroValue),
			c.Simple("enableLocationListening.android", "enable_location_listening.0.android"),
			c.Simple("enableLocationListening.reactnative", "enable_location_listening.0.react_native"),
			c.Simple("trackSessionEvents.android", "track_session_events.0.android"),
			c.Simple("trackSessionEvents.ios", "track_session_events.0.ios"),
			c.Simple("trackSessionEvents.reactnative", "track_session_events.0.react_native"),
			c.Simple("useAdvertisingIdForDeviceId.android", "use_advertising_id_for_device_id.0.android"),
			c.Simple("useAdvertisingIdForDeviceId.reactnative", "use_advertising_id_for_device_id.0.react_native"),
			c.Simple("useIdfaAsDeviceId.ios", "use_idfa_as_device_id.0.ios"),
			c.Simple("useIdfaAsDeviceId.reactnative", "use_idfa_as_device_id.0.react_native"),
			c.ArrayWithStrings("oneTrustCookieCategories.web", "oneTrustCookieCategory", "onetrust_cookie_categories.0.web"),
		},
		ConfigSchema: map[string]*schema.Schema{
			"api_key": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"api_secret": {
				Type:             schema.TypeString,
				Optional:         true,
				Sensitive:        true,
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"group_type_trait": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"group_value_trait": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"track_all_pages": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"track_categorized_pages": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"track_named_pages": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"track_products_once": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"track_revenue_per_product": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"version_name": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"traits_to_increment": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"traits_to_set_once": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"traits_to_append": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"traits_to_prepend": {
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
			"prefer_anonymous_id_for_device_id": {
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
			"device_id_from_url_param": {
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
			"force_https": {
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
			"track_gclid": {
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
			"track_referrer": {
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
			"save_params_referrer_once_per_session": {
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
			"track_utm_properties": {
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
			"unset_params_referrer_on_new_session": {
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
			"batch_events": {
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
			"event_upload_period_millis": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ios": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"android": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"react_native": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"event_upload_threshold": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ios": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"android": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"react_native": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"map_device_brand": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable_location_listening": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
			"track_session_events": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
			"use_advertising_id_for_device_id": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
			"use_idfa_as_device_id": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ios": {
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
