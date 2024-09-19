package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"web", "android", "ios", "unity", "reactnative", "flutter", "cordova", "amp", "cloud", "warehouse", "shopify"}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)

	properties := []c.ConfigProperty{
		c.Simple("appId", "app_id"),
		c.Simple("apiKey", "api_key"),
		c.Simple("useNativeSDK.web", "use_native_sdk.0.web", c.SkipZeroValue),
		c.Simple("useNativeSDK.ios", "use_native_sdk.0.ios", c.SkipZeroValue),
		c.Simple("useNativeSDK.android", "use_native_sdk.0.android", c.SkipZeroValue),
		c.Simple("collectContext", "collect_context", c.SkipZeroValue),
		c.Simple("sendAnonymousId", "send_anonymous_id", c.SkipZeroValue),
		c.Simple("updateLastRequestAt", "update_last_request_at", c.SkipZeroValue),
		c.Simple("mobileApiKeyAndroid.android", "mobile_api_key_android", c.SkipZeroValue),
		c.Simple("mobileApiKeyIOS.ios", "mobile_api_key_ios", c.SkipZeroValue),
		c.ArrayWithStrings("whitelistedEvents", "eventName", "event_filtering.0.whitelist"),
		c.ArrayWithStrings("blacklistedEvents", "eventName", "event_filtering.0.blacklist"),
		c.Discriminator("eventFilteringOption", c.DiscriminatorValues{
			"event_filtering.0.whitelist": "whitelistedEvents",
			"event_filtering.0.blacklist": "blacklistedEvents",
		}),
	}

	properties = append(properties, commonProperties...)

	schema := map[string]*schema.Schema{
		"api_key": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Enter your Intercom access token.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
		},
		"app_id": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Enter your app ID.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
		},
		"mobile_api_key_android": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Enter the Android API Key.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
		},
		"mobile_api_key_ios": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Enter the iOS API Key.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
		},
		"collect_context": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Enable this setting to include the user context along with your identify calls.",
		},
		"send_anonymous_id": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Enable this setting to send anonymousId as the secondary userId.",
		},
		"update_last_request_at": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Enable this setting to send the last seen information with the current time.",
		},
		"use_native_sdk": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Enable this setting to send the events through device mode, that is, using the native SDK.",
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
				},
			},
		},
		"event_filtering": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Use this setting to determine which events should be blocked or allowed to flow through.",
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
		"onetrust_cookie_categories": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Specify the OneTrust category name for mapping the OneTrust consent settings to RudderStack's consent purposes.",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}

	for key, value := range commonSchema {
		schema[key] = value
	}

	c.Destinations.Register("intercom", c.ConfigMeta{
		APIType:      "INTERCOM",
		Properties:   properties,
		ConfigSchema: schema,
	})
}
