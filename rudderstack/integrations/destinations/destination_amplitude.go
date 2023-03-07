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
			c.Discriminator("eventFilteringOption", c.DiscriminatorValues{
				"event_filtering.0.whitelist": "whitelistedEvents",
				"event_filtering.0.blacklist": "blacklistedEvents",
			}),
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
			c.ArrayWithStrings("oneTrustCookieCategories", "oneTrustCookieCategory", "onetrust_cookie_categories"),
			c.Simple("residencyServer", "residency_server", c.SkipZeroValue),
		},
		ConfigSchema: map[string]*schema.Schema{
			"api_key": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Enter your Amplitude API key.",
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"api_secret": {
				Type:             schema.TypeString,
				Optional:         true,
				Sensitive:        true,
				Description:      "Enter the Amplitude API Secret key required for user deletion.",
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"group_type_trait": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "RudderStack will use this value as `groupType` in the `group` calls.",
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"group_value_trait": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "RudderStack will use this value as `groupValue` in the `group` calls.",
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"residency_server": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].*)|^(standard|EU)$"),
			},
			"track_all_pages": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If this setting is enabled, RudderStack sends an event named `Loaded a page` / `Loaded a Screen` to Amplitude.",
			},
			"track_categorized_pages": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If this setting is enabled and if `category` is present in a `page` / `screen` call, then an event named `Viewed {category} page` / `Viewed {category} Screen` will be sent to Amplitude.",
			},
			"track_named_pages": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If this setting is enabled and `name` is present in a `page` call, then an event named `Viewed {name} page` will be sent to Amplitude.",
			},
			"track_products_once": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If this setting is enabled and if the event payload contains an array of products, then the event is tracked with the original event name and all the products as its property. Otherwise, each product is tracked with event as `Product purchased`.",
			},
			"track_revenue_per_product": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If this setting is enabled and if the event payload contains multiple products, each product's revenue is tracked individually.",
			},
			"version_name": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "The value of this field is set as the `versionName` of the Amplitude SDK.",
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"traits_to_increment": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "If this setting is enabled, the value of the corresponding trait will be incremented at Amplitude, with the value provided against the trait in an `identify` call.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"traits_to_set_once": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "If this setting is enabled, the value of the corresponding trait will be set once at Amplitude with the value provided against the trait in an `identify` call.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"traits_to_append": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "If this setting is enabled, the value of the corresponding trait will be appended to the corresponding trait array at Amplitude.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"traits_to_prepend": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "If this setting is enabled, the value of the corresponding trait will be prepended to the corresponding trait array at Amplitude.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"use_native_sdk": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Enable this setting to send events to Amplitude via the device mode.",
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
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "If this setting is enabled, the device ID will be set as the `anonymousId` generated by RudderStack SDK or by the `anonymousId` set via RudderStack's `setAnonymousId()` method.",
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
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "If this setting is enabled, the Amplitude SDK will parse the URL parameter and set the device ID from `amp_device_id`.",
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
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "If this setting is enabled, the events will always be uploaded by the Amplitude SDK to the HTTPS endpoint, otherwise it will use the embedding site's protocol.",
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
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "If this setting is enabled, the Amplitude SDK will capture the `gclid` URL parameters along with the user's `initial_gclid` parameters.",
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
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "If this setting is enabled, the Amplitude SDK will capture the `referrer` and `referring_domain` for each session along with the user's `initial_referrer` and `initial_referring_domain`.",
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
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "If this setting is enabled, the corresponding tracking of `gclid`, referrer, UTM parameters will be done once per session.",
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
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "If this setting is enabled, the Amplitude SDK parses the UTM parameters in the query string or `_utmz` cookie and includes them as user properties in all uploaded events.",
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
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "If this setting is disabled, the existing `referrer` and `utm_parameter` values will be passed to each new session. If enabled, `referrer` and `utm_parameter` properties will be set to `null` upon instantiating a new session.",
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
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "If this setting is enabled, the events are batched together and uploaded by the Amplitude SDK.",
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
			"event_upload_period_millis": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "If the batch events settings is enabled, this is the amount of time that the SDK waits to upload the events.",
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
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "If the batch events settings is enabled, this is the minimum number of events to batch together by the Amplitude SDK.",
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
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable this setting for RudderStack to send the device brand information (`context.device.brand`) to Amplitude.",
			},
			"enable_location_listening": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Enable this setting to activate location listening.",
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
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Enable this setting to track the session events.",
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
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Enable this setting to set the advertising ID as the device ID.",
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
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Enable this setting to set the IDFA as the device ID.",
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
