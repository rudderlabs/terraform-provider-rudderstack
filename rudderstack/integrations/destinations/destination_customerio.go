package destinations

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

var customerIOSupportedSourceTypes = []string{"web", "android", "androidKotlin", "ios", "iosSwift", "unity", "reactnative", "flutter", "cordova", "amp", "cloud", "warehouse", "shopify"}

func init() {
	commonProperties, commonSchema := GetCommonConfigMeta(customerIOSupportedSourceTypes)

	properties := []c.ConfigProperty{
		c.Simple("siteID", "site_id", c.SkipZeroValue),
		c.Simple("apiKey", "api_key", c.SkipZeroValue),
		c.Simple("apiVersion", "api_version"),
		c.Simple("userIdIdentifierType", "user_id_identifier_type", c.SkipZeroValue),
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
		c.SimpleWithDefault("sdkVersion.web", "sdk_version.0.web", "v2"),
		c.Simple("writeKey.web", "write_key.0.web", c.SkipZeroValue),
		c.Simple("anonymousInApp.web", "anonymous_in_app.0.web"),
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
			Optional:         true,
			Description:      "Enter your Customer.io site ID. Required unless only web device mode with SDK v2 is configured.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"api_key": {
			Type:             schema.TypeString,
			Optional:         true,
			Sensitive:        true,
			Description:      "Enter your Customer.io API key. Required unless only web device mode is configured.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"api_version": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "v2",
			Description:      "Customer.io API version for cloud-mode delivery. Defaults to `v2` for the unified /v2/batch API; set to `v1` for legacy per-endpoint behavior. This setting does not affect device-mode SDK delivery.",
			ValidateDiagFunc: c.StringMatchesRegexp("^(v1|v2)$"),
		},
		"user_id_identifier_type": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Customer.io identifier that receives the RudderStack `userId`. Required when `api_version` is `v2`, regardless of connection mode.",
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
			Description: "Enable this setting to send in-app messages to your website in web device mode with SDK v1.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"web": {
						Type:     schema.TypeBool,
						Optional: true,
					},
				},
			},
		},
		"sdk_version": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Description: "Choose the Customer.io SDK version for web device mode. Defaults to `v2` when this block is omitted.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"web": {
						Type:             schema.TypeString,
						Optional:         true,
						Default:          "v2",
						ValidateDiagFunc: c.StringMatchesRegexp("^(v1|v2)$"),
					},
				},
			},
		},
		"write_key": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Enter the Customer.io Data Pipelines write key for web device mode with SDK v2.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"web": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(.{0,100})$"),
					},
				},
			},
		},
		"anonymous_in_app": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Enable in-app messages for anonymous users in web device mode with SDK v2.",
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
		APIType:             "CUSTOMERIO",
		Version:             1,
		Properties:          properties,
		ConfigSchema:        schema,
		CustomizeConfigDiff: validateCustomerIODestinationConfig,
	})
}

func validateCustomerIODestinationConfig(d *schema.ResourceDiff) error {
	return errors.Join(
		validateCustomerIOWriteKey(d),
		validateCustomerIOAPIKey(d),
		validateCustomerIOSiteID(d),
		validateCustomerIOUserIDIdentifier(d),
	)
}

func validateCustomerIOWriteKey(d *schema.ResourceDiff) error {
	webDeviceMode, webDeviceModeKnown := customerIOOptionalListStringConfigValue(d, "config.0.connection_mode.#", "config.0.connection_mode.0.web", "")
	sdkVersion, sdkVersionKnown := customerIOSDKVersionConfigValue(d)

	if webDeviceModeKnown && sdkVersionKnown && webDeviceMode == "device" && sdkVersion == "v2" {
		writeKey, writeKeyKnown := customerIOOptionalListStringConfigValue(d, "config.0.write_key.#", "config.0.write_key.0.web", "")
		if writeKeyKnown && writeKey == "" {
			return fmt.Errorf("config.0.write_key.0.web must be non-empty when connection_mode.0.web is %q and sdk_version.0.web is %q", "device", "v2")
		}
	}
	return nil
}

func validateCustomerIOAPIKey(d *schema.ResourceDiff) error {
	exactWebDeviceMode, exactWebDeviceModeKnown := customerIOConnectionModeIsWebOnly(d)
	if !exactWebDeviceModeKnown {
		return nil
	}
	if !exactWebDeviceMode {
		apiKey, apiKeyKnown := customerIOStringConfigValue(d, "config.0.api_key")
		if apiKeyKnown && apiKey == "" {
			return fmt.Errorf("config.0.api_key must be non-empty unless connection_mode sets only web to %q", "device")
		}
	}
	return nil
}

func validateCustomerIOSiteID(d *schema.ResourceDiff) error {
	exactWebDeviceMode, exactWebDeviceModeKnown := customerIOConnectionModeIsWebOnly(d)
	sdkVersion, sdkVersionKnown := customerIOSDKVersionConfigValue(d)
	siteIDRequired := (exactWebDeviceModeKnown && !exactWebDeviceMode) || (sdkVersionKnown && sdkVersion != "v2")
	if !siteIDRequired {
		return nil
	}

	siteID, siteIDKnown := customerIOStringConfigValue(d, "config.0.site_id")
	if siteIDKnown && siteID == "" {
		return fmt.Errorf("config.0.site_id must be non-empty unless connection_mode sets only web to %q and sdk_version.0.web is %q", "device", "v2")
	}

	return nil
}

func validateCustomerIOUserIDIdentifier(d *schema.ResourceDiff) error {
	apiVersion, apiVersionKnown := customerIOStringConfigValue(d, "config.0.api_version")
	if !apiVersionKnown || apiVersion != "v2" {
		return nil
	}

	identifier, identifierKnown := customerIOStringConfigValue(d, "config.0.user_id_identifier_type")
	if identifierKnown && identifier == "" {
		return errors.New(`config.0.user_id_identifier_type must be set when api_version is "v2"`)
	}
	return nil
}

func customerIOOptionalListStringConfigValue(d *schema.ResourceDiff, countKey, valueKey, defaultValue string) (string, bool) {
	count, countKnown := customerIOListBlockCount(d, countKey)
	if !countKnown {
		return "", false
	}
	if count == 0 {
		return defaultValue, true
	}
	return customerIOStringConfigValue(d, valueKey)
}

func customerIOSDKVersionConfigValue(d *schema.ResourceDiff) (string, bool) {
	sdkVersion, known := customerIOOptionalListStringConfigValue(d, "config.0.sdk_version.#", "config.0.sdk_version.0.web", "v2")
	if known {
		return sdkVersion, true
	}

	// On create, Terraform reports the count of an omitted Optional+Computed
	// block as unknown. The raw config distinguishes that case from a dynamic
	// block whose count is genuinely unknown. Preserve the latter as unknown so
	// validation can defer until apply.
	if customerIORawConfigListBlockOmitted(d.GetRawConfig(), "sdk_version") {
		return "v2", true
	}
	return "", false
}

func customerIORawConfigListBlockOmitted(raw cty.Value, blockName string) bool {
	if !raw.IsKnown() || raw.IsNull() || !raw.Type().IsObjectType() || !raw.Type().HasAttribute("config") {
		return false
	}

	config := raw.GetAttr("config")
	if !config.IsKnown() || config.IsNull() || !config.CanIterateElements() {
		return false
	}
	configIterator := config.ElementIterator()
	if !configIterator.Next() {
		return false
	}
	_, configValue := configIterator.Element()
	if !configValue.IsKnown() || configValue.IsNull() || !configValue.Type().IsObjectType() || !configValue.Type().HasAttribute(blockName) {
		return false
	}

	block := configValue.GetAttr(blockName)
	if !block.IsKnown() {
		return false
	}
	if block.IsNull() {
		return true
	}
	if !block.CanIterateElements() {
		return false
	}
	return block.LengthInt() == 0
}

func customerIOStringConfigValue(d *schema.ResourceDiff, key string) (string, bool) {
	if !d.NewValueKnown(key) {
		return "", false
	}

	stringValue, _ := d.Get(key).(string)
	return stringValue, true
}

func customerIOConnectionModeIsWebOnly(d *schema.ResourceDiff) (bool, bool) {
	count, countKnown := customerIOListBlockCount(d, "config.0.connection_mode.#")
	if !countKnown {
		return false, false
	}
	if count == 0 {
		return false, true
	}

	valueUnknown := false
	for _, sourceType := range customerIOSupportedSourceTypes {
		terraformSourceType := camelToSnake(sourceType)
		mode, known := customerIOStringConfigValue(d, "config.0.connection_mode.0."+terraformSourceType)
		if !known {
			valueUnknown = true
			continue
		}
		if terraformSourceType == "web" {
			if mode != "device" {
				return false, true
			}
			continue
		}
		if mode != "" {
			return false, true
		}
	}
	if valueUnknown {
		return false, false
	}
	return true, true
}

func customerIOListBlockCount(d *schema.ResourceDiff, key string) (int, bool) {
	if !d.NewValueKnown(key) {
		return 0, false
	}

	count, _ := d.Get(key).(int)
	return count, true
}
