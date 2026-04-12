package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"web", "android", "ios", "androidKotlin", "iosSwift", "unity", "reactnative", "flutter", "cordova", "amp", "cloud", "warehouse", "shopify"}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)
	ga4URLPattern := "(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(?:http(s)?:\\/\\/)?[\\w.-]+(?:\\.[\\w\\.-]+)+[\\w\\-\\._~:\\/?#[\\]@!\\$&'\\(\\)\\*\\+,;=.]*|^$"
	ga4NgrokPattern := ".*\\.ngrok\\.io.*"

	properties := []c.ConfigProperty{
		c.Simple("apiSecret", "api_secret"),
		c.Simple("typesOfClient", "client_type"),
		c.Simple("measurementId", "measurement_id", c.SkipZeroValue),
		c.Simple("firebaseAppId", "firebase_app_id", c.SkipZeroValue),
		c.Simple("debugMode", "debug_mode"),
		c.Simple("blockPageViewEvent", "block_page_view_event", c.SkipZeroValue),
		c.Simple("extendPageViewParams", "extend_page_view_params", c.SkipZeroValue),
		c.Simple("sendUserId", "send_user_id", c.SkipZeroValue),
		c.Simple("sdkBaseUrl", "sdk_base_url", c.SkipZeroValue),
		c.Simple("serverContainerUrl", "server_container_url", c.SkipZeroValue),
		c.ArrayWithObjects("piiPropertiesToIgnore", "pii_properties_to_ignore", map[string]interface{}{
			"piiProperty": "pii_property",
		}),
		c.ArrayWithStrings("whitelistedEvents", "eventName", "event_filtering.0.whitelist"),
		c.ArrayWithStrings("blacklistedEvents", "eventName", "event_filtering.0.blacklist"),
		c.Discriminator("eventFilteringOption", c.DiscriminatorValues{
			"event_filtering.0.whitelist": "whitelistedEvents",
			"event_filtering.0.blacklist": "blacklistedEvents",
		}),
		c.Simple("useNativeSDK.web", "use_native_sdk.0.web"),
		c.Simple("useNativeSDK.android", "use_native_sdk.0.android"),
		c.Simple("useNativeSDK.ios", "use_native_sdk.0.ios"),
		c.Simple("capturePageView.web", "capture_page_view.0.web"),
		c.Simple("debugView.web", "debug_view.0.web"),
		c.Simple("overrideClientAndSessionId.web", "override_client_and_session_ids.0.web"),
	}

	properties = append(properties, commonProperties...)

	schema := map[string]*schema.Schema{
		"api_secret": {
			Type:             schema.TypeString,
			Required:         true,
			Sensitive:        true,
			Description:      "Enter the API Secret generated through the Google Analytics dashboard.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"client_type": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "gtag",
			Description:      "Select the client type as gtag or Firebase.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(gtag|firebase)$"),
		},
		"measurement_id": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Enter the Measurement Id which is the identifier for a data stream.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(G-.{1,100})$|^$"),
		},
		"firebase_app_id": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Enter the Firebase App ID which is the identifier for Firebase app.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
		},
		"debug_mode": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable this setting to send events to GA4's validation server instead of reporting them. This allows you to check validation responses in Live Events, but these events will not show up in reports.",
		},
		"sdk_base_url": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Enter your GA4 Custom Domain URL. By default, it is https://www.googletagmanager.com.",
			ValidateDiagFunc: c.ValidateAll(c.StringMatchesRegexp(ga4URLPattern), c.StringNotMatchesRegexp(ga4NgrokPattern)),
		},
		"server_container_url": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Enter your GA4 Server Side Container URL.",
			ValidateDiagFunc: c.ValidateAll(c.StringMatchesRegexp(ga4URLPattern), c.StringNotMatchesRegexp(ga4NgrokPattern)),
		},
		"pii_properties_to_ignore": {
			Type:        schema.TypeList,
			ConfigMode:  schema.SchemaConfigModeAttr,
			Optional:    true,
			Description: "Use this field to filter sensitive PII fields from your events before sending them to GA4.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"pii_property": {
						Type:             schema.TypeString,
						Required:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
					},
				},
			},
		},
		"block_page_view_event": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Enable this setting to disable sending `page_view` events on load. This setting is applicable only for device mode.",
		},
		"extend_page_view_params": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Enable this setting to send `url` and `search` along with any other custom property to the `page` call of the RudderStack SDK. This setting is applicable only for device mode.",
		},
		"send_user_id": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "If enabled, the user ID is set to the identified visitors and sent to Google Analytics 4.",
		},
		"event_filtering": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "With this option, you can determine which events are blocked or allowed to flow through to Google Analytics 4.",
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
					"android": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"ios": {
						Type:     schema.TypeBool,
						Optional: true,
					},
				},
			},
		},
		"capture_page_view": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Choose whether to send page view events through the RudderStack JS SDK or through automatic collection using GA4 Enhanced Measurement.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"web": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(rs|gtag)$"),
					},
				},
			},
		},
		"debug_view": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Enable this setting to monitor your device mode events in GA4 DebugView.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"web": {
						Type:     schema.TypeBool,
						Optional: true,
					},
				},
			},
		},
		"override_client_and_session_ids": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Override the gtag client ID and session ID with RudderStack's to ensure attribution is properly unified across page and track events.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"web": {
						Type:     schema.TypeBool,
						Optional: true,
					},
				},
			},
		},
	}

	for key, value := range commonSchema {
		schema[key] = value
	}

	c.Destinations.Register("google_analytics4", c.ConfigMeta{
		APIType:      "GA4",
		Properties:   properties,
		ConfigSchema: schema,
	})
}
