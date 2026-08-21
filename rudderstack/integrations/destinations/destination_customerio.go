package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"web", "android", "androidKotlin", "ios", "iosSwift", "unity", "reactnative", "flutter", "cordova", "amp", "cloud", "warehouse", "shopify"}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)

	properties := []c.ConfigProperty{
		c.Simple("siteID", "site_id"),
		c.Simple("apiKey", "api_key"),
		customerIOAPIVersionProperty(),
		c.Simple("userIdMapping", "user_id_mapping", c.SkipZeroValue),
		c.Simple("deviceTokenEventName", "device_token_event_name", c.SkipZeroValue),
		c.Simple("datacenter", "datacenter"),
		c.Simple("connectionMode.web", "connection_mode.0.web", c.SkipZeroValue),
		c.Simple("connectionMode.android", "connection_mode.0.android", c.SkipZeroValue),
		c.Simple("connectionMode.androidKotlin", "connection_mode.0.android_kotlin", c.SkipZeroValue),
		c.Simple("connectionMode.ios", "connection_mode.0.ios", c.SkipZeroValue),
		c.Simple("connectionMode.iosSwift", "connection_mode.0.ios_swift", c.SkipZeroValue),
		c.Simple("connectionMode.unity", "connection_mode.0.unity", c.SkipZeroValue),
		c.Simple("connectionMode.amp", "connection_mode.0.amp", c.SkipZeroValue),
		c.Simple("connectionMode.reactnative", "connection_mode.0.reactnative", c.SkipZeroValue),
		c.Simple("connectionMode.flutter", "connection_mode.0.flutter", c.SkipZeroValue),
		c.Simple("connectionMode.cordova", "connection_mode.0.cordova", c.SkipZeroValue),
		c.Simple("connectionMode.shopify", "connection_mode.0.shopify", c.SkipZeroValue),
		c.Simple("connectionMode.cloud", "connection_mode.0.cloud", c.SkipZeroValue),
		c.Simple("connectionMode.warehouse", "connection_mode.0.warehouse", c.SkipZeroValue),
		c.Simple("useNativeSDK.web", "use_native_sdk.0.web"),
		c.Simple("useNativeSDK.android", "use_native_sdk.0.android"),
		c.Simple("useNativeSDK.ios", "use_native_sdk.0.ios"),
		c.Simple("sendPageNameInSDK.web", "send_page_name_in_sdk.0.web"),
		c.Simple("dataUseInApp.web", "data_use_in_app.0.web"),
		c.Simple("autoTrackDeviceAttributes.android", "auto_track_device_attributes.0.android"),
		c.Simple("autoTrackDeviceAttributes.ios", "auto_track_device_attributes.0.ios"),
		c.Simple("backgroundQueueMinNumberOfTasks.android", "background_queue_min_number_of_tasks.0.android", c.SkipZeroValue),
		c.Simple("backgroundQueueSecondsDelay.android", "background_queue_seconds_delay.0.android", c.SkipZeroValue),
		c.ArrayWithStrings("whitelistedEvents", "eventName", "event_filtering.0.whitelist"),
		c.ArrayWithStrings("blacklistedEvents", "eventName", "event_filtering.0.blacklist"),
		c.Discriminator("eventFilteringOption", c.DiscriminatorValues{
			"event_filtering.0.whitelist": "whitelistedEvents",
			"event_filtering.0.blacklist": "blacklistedEvents",
		}),
	}

	properties = append(properties, commonProperties...)

	schema := map[string]*schema.Schema{
		"site_id": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Enter your Customer.io site ID.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"api_key": {
			Type:             schema.TypeString,
			Required:         true,
			Sensitive:        true,
			Description:      "Enter your Customer.io API key.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"api_version": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "v1",
			Description:      "Customer.io API version for cloud-mode delivery. `v1` uses the existing per-endpoint behavior; `v2` uses the unified /v2/batch API. This setting does not affect device-mode SDK delivery.",
			ValidateDiagFunc: c.StringMatchesRegexp("^(v1|v2)$"),
		},
		"user_id_mapping": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Customer.io identifier that receives the RudderStack `userId` for cloud-mode delivery when `api_version` is `v2`. This setting does not affect device-mode SDK delivery.",
			ValidateDiagFunc: c.StringMatchesRegexp("^(id|email|phone|cio_id)$"),
		},
		"device_token_event_name": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Enter the name of the event that is fired immediately after setting the device token.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"datacenter": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "US",
			Description: "Input your Customer.io Data Center. (US or EU)",
		},
		"connection_mode": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Configure the connection mode per source type for Customer.io.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"web": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud|device)$"),
					},
					"android": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud|device)$"),
					},
					"android_kotlin": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$"),
					},
					"ios": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud|device)$"),
					},
					"ios_swift": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$"),
					},
					"unity": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$"),
					},
					"amp": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$"),
					},
					"reactnative": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$"),
					},
					"flutter": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$"),
					},
					"cordova": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$"),
					},
					"shopify": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$"),
					},
					"cloud": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$"),
					},
					"warehouse": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$"),
					},
				},
			},
		},
		"use_native_sdk": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Enable this setting to send the events through Customer.io's native SDK.",
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
		"send_page_name_in_sdk": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Configure whether to send the page name in SDK mode.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"web": {
						Type:     schema.TypeBool,
						Optional: true,
					},
				},
			},
		},
		"data_use_in_app": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Enable this setting to send in-app messages to your website.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"web": {
						Type:     schema.TypeBool,
						Optional: true,
					},
				},
			},
		},
		"auto_track_device_attributes": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Enable this setting to automatically track device attributes in SDK mode.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
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
		"background_queue_min_number_of_tasks": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Configure the minimum number of tasks in the background queue.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"android": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
					},
				},
			},
		},
		"background_queue_seconds_delay": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Configure the delay in seconds for the background queue.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"android": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
					},
				},
			},
		},
		"event_filtering": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "RudderStack lets you determine which events should be allowed to flow through or blocked.",
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
	}

	for key, value := range commonSchema {
		schema[key] = value
	}

	c.Destinations.Register("customerio", c.ConfigMeta{
		APIType:      "CUSTOMERIO",
		Version:      1,
		Properties:   properties,
		ConfigSchema: schema,
	})
}

func customerIOAPIVersionProperty() c.ConfigProperty {
	property := c.Simple("apiVersion", "api_version", skipDefaultCustomerIOAPIVersion)
	toState := property.ToStateFunc
	property.ToStateFunc = func(state, config string) (string, error) {
		result, err := toState(state, config)
		if err != nil {
			return result, err
		}
		if gjson.Get(result, "api_version").Exists() {
			return result, nil
		}
		return sjson.Set(result, "api_version", "v1")
	}

	return property
}

func skipDefaultCustomerIOAPIVersion(value interface{}) bool {
	return value == "v1"
}
