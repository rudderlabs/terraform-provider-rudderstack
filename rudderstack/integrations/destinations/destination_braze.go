package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	c.Destinations.Register("braze", c.ConfigMeta{
		APIType: "BRAZE",
		Properties: []c.ConfigProperty{
			c.Simple("restApiKey", "rest_api_key"),
			c.Simple("appKey", "app_key"),
			c.Simple("dataCenter", "data_center"),
			c.Simple("enableSubscriptionGroupInGroupCall", "enable_subscription_group_in_group_call", c.SkipZeroValue),
			c.Simple("enableNestedArrayOperations", "enable_nested_array_operations", c.SkipZeroValue),
			c.Simple("sendPurchaseEventWithExtraProperties", "send_purchase_event_with_extra_properties", c.SkipZeroValue),
			c.Simple("trackAnonymousUser", "track_anonymous_user", c.SkipZeroValue),
			c.Simple("supportDedup", "support_dedup", c.SkipZeroValue),
			c.Simple("enableBrazeLogging", "enable_braze_logging", c.SkipZeroValue),
			c.Simple("enablePushNotification", "enable_push_notification", c.SkipZeroValue),
			c.Simple("allowUserSuppliedJavascript", "allow_user_supplied_javascript", c.SkipZeroValue),
			c.ArrayWithStrings("whitelistedEvents", "eventName", "event_filtering.0.whitelist"),
			c.ArrayWithStrings("blacklistedEvents", "eventName", "event_filtering.0.blacklist"),
			c.Discriminator("eventFilteringOption", c.DiscriminatorValues{
				"event_filtering.0.whitelist": "whitelistedEvents",
				"event_filtering.0.blacklist": "blacklistedEvents",
			}),
			c.ArrayWithStrings("oneTrustCookieCategories", "oneTrustCookieCategory", "onetrust_cookie_categories"),
			c.Simple("connectionMode.web", "connection_mode.0.web", c.SkipZeroValue),
			c.Simple("connectionMode.ios", "connection_mode.0.ios", c.SkipZeroValue),
			c.Simple("connectionMode.android", "connection_mode.0.android", c.SkipZeroValue),
			c.Simple("connectionMode.reactnative", "connection_mode.0.react_native", c.SkipZeroValue),
			c.Simple("connectionMode.unity", "connection_mode.0.unity", c.SkipZeroValue),
			c.Simple("connectionMode.amp", "connection_mode.0.amp", c.SkipZeroValue),
			c.Simple("connectionMode.flutter", "connection_mode.0.flutter", c.SkipZeroValue),
			c.Simple("connectionMode.cordova", "connection_mode.0.cordova", c.SkipZeroValue),
			c.Simple("connectionMode.shopify", "connection_mode.0.shopify", c.SkipZeroValue),
			c.Simple("connectionMode.cloud", "connection_mode.0.cloud", c.SkipZeroValue),
			c.Simple("connectionMode.warehouse", "connection_mode.0.warehouse", c.SkipZeroValue),
		},
		ConfigSchema: map[string]*schema.Schema{
			"rest_api_key": {
				Type:             schema.TypeString,
				Optional:         true,
				Sensitive:        true,
				Description:      "Enter your Braze Rest Api Key. Required for cloud mode.",
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"app_key": {
				Type:             schema.TypeString,
				Optional:         true,
				Sensitive:        true,
				Description:      "Enter your Braze APP Key. Required for Device mode.",
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"data_center": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Enter your Braze Data Center.",
			},
			"enable_subscription_group_in_group_call": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Use this setting to Enable subscription groups in group call",
			},
			"enable_nested_array_operations": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Use this setting to use Custom Attributes Operation.",
			},
			"send_purchase_event_with_extra_properties": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Use this setting to Enable to send purchase events with custom properties.",
			},
			"track_anonymous_user": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Use this setting to Track events for anonymous users.",
			},
			"support_dedup": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Use this setting to enable Deduplicate Traits on identify and track.",
			},
			"enable_braze_logging": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Use this setting to show braze logs to customer.",
			},
			"enable_push_notification": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Use this setting to use push notification. It requires service worker setup by client.",
			},
			"allow_user_supplied_javascript": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Use this setting to enable HTML in-app messages.",
			},
			"event_filtering": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "This option allows you filter the events you want to send to Braze.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"whitelist": {
							Type:         schema.TypeList,
							Optional:     true,
							Description:  "Enter the event names to be allowlisted.",
							ExactlyOneOf: []string{"config.0.event_filtering.0.whitelist", "config.0.event_filtering.0.blacklist"},
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"blacklist": {
							Type:         schema.TypeList,
							Optional:     true,
							Description:  "Enter the event names to be denylisted.",
							ExactlyOneOf: []string{"config.0.event_filtering.0.whitelist", "config.0.event_filtering.0.blacklist"},
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"connection_mode": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "Use this setting to set how you want to route events from your source to destination..",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].*)|^(cloud|device|hybrid)$"),
						},
						"ios": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].*)|^(cloud|device|hybrid)$"),
						},
						"android": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].*)|^(cloud|device|hybrid)$"),
						},
						"react_native": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].*)|^(cloud|device)$"),
						},
						"unity": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].*)|^(cloud)$"),
						},
						"amp": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].*)|^(cloud)$"),
						},
						"flutter": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].*)|^(cloud|device)$"),
						},
						"cordova": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].*)|^(cloud)$"),
						},
						"shopify": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].*)|^(cloud)$"),
						},
						"cloud": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].*)|^(cloud)$"),
						},
						"warehouse": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].*)|^(cloud)$"),
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
