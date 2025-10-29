package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"web", "android", "ios", "unity", "reactnative", "flutter", "cordova", "amp", "cloud", "warehouse", "shopify"}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)

	properties := []c.ConfigProperty{
		c.Simple("restApiKey", "rest_api_key", c.SkipZeroValue),
		c.Simple("appKey", "app_key", c.SkipZeroValue),
		c.Simple("dataCenter", "data_center"),
		c.Simple("enableSubscriptionGroupInGroupCall", "enable_subscription_group_in_group_call", c.SkipZeroValue),
		c.Simple("enableNestedArrayOperations", "enable_nested_array_operations", c.SkipZeroValue),
		c.Simple("sendPurchaseEventWithExtraProperties", "send_purchase_event_with_extra_properties", c.SkipZeroValue),
		c.Simple("trackAnonymousUser.web", "track_anonymous_user.0.web"),
		c.Simple("supportDedup", "support_dedup", c.SkipZeroValue),
		c.Simple("enableBrazeLogging.web", "enable_braze_logging.0.web"),
		c.Simple("enablePushNotification.web", "enable_push_notification.0.web"),
		c.Simple("allowUserSuppliedJavascript.web", "allow_user_supplied_javascript.0.web"),
		c.ArrayWithStrings("whitelistedEvents", "eventName", "event_filtering.0.whitelist"),
		c.ArrayWithStrings("blacklistedEvents", "eventName", "event_filtering.0.blacklist"),
		c.Discriminator("eventFilteringOption", c.DiscriminatorValues{
			"event_filtering.0.whitelist": "whitelistedEvents",
			"event_filtering.0.blacklist": "blacklistedEvents",
		}),
		c.Simple("connectionMode.web", "connection_mode.0.web", c.SkipZeroValue),
		c.Simple("connectionMode.ios", "connection_mode.0.ios", c.SkipZeroValue),
		c.Simple("connectionMode.android", "connection_mode.0.android", c.SkipZeroValue),
		c.Simple("connectionMode.reactnative", "connection_mode.0.reactnative", c.SkipZeroValue),
		c.Simple("connectionMode.unity", "connection_mode.0.unity", c.SkipZeroValue),
		c.Simple("connectionMode.amp", "connection_mode.0.amp", c.SkipZeroValue),
		c.Simple("connectionMode.flutter", "connection_mode.0.flutter", c.SkipZeroValue),
		c.Simple("connectionMode.cordova", "connection_mode.0.cordova", c.SkipZeroValue),
		c.Simple("connectionMode.shopify", "connection_mode.0.shopify", c.SkipZeroValue),
		c.Simple("connectionMode.cloud", "connection_mode.0.cloud", c.SkipZeroValue),
		c.Simple("connectionMode.warehouse", "connection_mode.0.warehouse", c.SkipZeroValue),
	}

	properties = append(properties, commonProperties...)

	schema := map[string]*schema.Schema{
		"rest_api_key": {
			Type:             schema.TypeString,
			Optional:         true,
			Sensitive:        true,
			Description:      "Enter your Braze Rest Api Key. Required for cloud mode.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
		},
		"app_key": {
			Type:             schema.TypeString,
			Optional:         true,
			Sensitive:        true,
			Description:      "Enter your Braze APP Key. Required for Device mode.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
		},
		"data_center": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Enter your Braze Data Center.",
			ValidateDiagFunc: c.StringMatchesRegexp("^($|US-01|US-02|US-03|US-04|US-05|US-06|US-07|US-08|EU-01|EU-02|EU-03|AU-01)$"),
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
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Use this setting to Track events for anonymous users.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"web": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},
		"support_dedup": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Use this setting to enable Deduplicate Traits on identify and track.",
		},
		"enable_braze_logging": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Use this setting to show braze logs to customer.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"web": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},
		"enable_push_notification": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Use this setting to use push notification. It requires service worker setup by client.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"web": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},
		"allow_user_supplied_javascript": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Use this setting to enable HTML in-app messages.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"web": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
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
			Optional:    true,
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
					"reactnative": {
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
	}

	for key, value := range commonSchema {
		schema[key] = value
	}

	c.Destinations.Register("braze", c.ConfigMeta{
		APIType:      "BRAZE",
		Properties:   properties,
		ConfigSchema: schema,
	})
}
