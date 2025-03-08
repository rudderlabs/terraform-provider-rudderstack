package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"web", "cloud", "flutter"}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)

	properties := []c.ConfigProperty{
		c.Simple("yourInstance", "endpoint"),
		c.Simple("teamApiKey", "api_key"),
		c.Simple("useV2Group", "use_v2_group"),
		c.Simple("connectionMode.web", "connection_mode.0.web"),
		c.Simple("connectionMode.cloud", "connection_mode.0.cloud"),
		c.Simple("connectionMode.flutter", "connection_mode.0.flutter"),
		c.Simple("disableSessionRecording.web", "disable_session_recording.0.web"),
		c.Simple("capturePageView.web", "capture_page_view.0.web"),
		c.Simple("useNativeSDK.web", "use_native_sdk.0.web"),
		c.Simple("autocapture.web", "autocapture.0.web"),
		c.Simple("enableLocalStoragePersistence.web", "enable_local_storage_persistence.0.web"),
		c.ArrayWithObjects("xhrHeaders.web", "xhr_headers", map[string]interface{}{
			"key":   "key",
			"value": "value",
		}),
		c.ArrayWithObjects("propertyBlacklist.web", "property_blacklist", map[string]interface{}{
			"property": "property",
		}),
		c.Simple("personProfiles.web", "person_profiles.0.web"),

		c.ArrayWithStrings("whitelistedEvents", "eventName", "event_filtering.0.whitelist"),
		c.ArrayWithStrings("blacklistedEvents", "eventName", "event_filtering.0.blacklist"),
		c.Discriminator("eventFilteringOption", c.DiscriminatorValues{
			"event_filtering.0.whitelist": "whitelistedEvents",
			"event_filtering.0.blacklist": "blacklistedEvents",
		}),
	}

	properties = append(properties, commonProperties...)

	schema := map[string]*schema.Schema{
		"endpoint": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The endpoint of your PostHog service.",
		},
		"api_key": {
			Type:             schema.TypeString,
			Required:         true,
			Sensitive:        true,
			Description:      "Enter your PostHog API key.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"use_v2_group": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Whether to use v2 grouping for your PostHog service.",
			Default:     true,
		},

		"event_filtering": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "This option allows you filter the events you want to send to PostHog.",
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
						ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].*)|^(cloud|device)$"),
					},
					"flutter": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].*)|^(cloud)$"),
					},
					"cloud": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].*)|^(cloud)$"),
					},
				},
			},
		},

		"use_native_sdk": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Use this setting to enable native SDK.",
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

		"autocapture": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Use this setting to enable autocapture.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"web": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  true,
					},
				},
			},
		},
		"capture_page_view": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Use this setting to enable capture page views.",
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
		"disable_session_recording": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Use this setting to disable session recording.",
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
		"enable_local_storage_persistence": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Use this setting to enable local storage persistence.",
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
		"xhr_headers": {
			Type:        schema.TypeList,
			ConfigMode:  schema.SchemaConfigModeAttr,
			Optional:    true,
			Description: "Use this setting to enable property blacklist.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:             schema.TypeString,
						Required:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
					},
					"value": {
						Type:             schema.TypeString,
						Required:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,500})$"),
					},
				},
			},
		},
		"property_blacklist": {
			Type:        schema.TypeList,
			ConfigMode:  schema.SchemaConfigModeAttr,
			Optional:    true,
			Description: "Use this setting to enable property blacklist.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"property": {
						Type:             schema.TypeString,
						Required:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
					},
				},
			},
		},

		"person_profiles": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Use this setting to enable person profiles.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"web": {
						Type:             schema.TypeString,
						Optional:         true,
						Default:          "always",
						ValidateDiagFunc: c.StringMatchesRegexp("(^always$)|(^identified_only$)"),
					},
				},
			},
		},
	}

	for key, value := range commonSchema {
		schema[key] = value
	}

	c.Destinations.Register("posthog", c.ConfigMeta{
		APIType:      "POSTHOG",
		Properties:   properties,
		ConfigSchema: schema,
	})
}
