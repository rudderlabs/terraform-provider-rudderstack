package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{
		"android", "androidKotlin", "ios", "iosSwift", "web",
		"unity", "amp", "cloud", "warehouse", "reactnative",
		"flutter", "cordova", "shopify",
	}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)

	properties := []c.ConfigProperty{
		c.Simple("accountSID", "account_sid"),
		c.Simple("apiKey", "api_key"),
		c.Simple("campaignId", "campaign_id"),
		c.Simple("eventTypeId", "event_type_id", c.SkipZeroValue),
		c.Simple("impactAppId", "impact_app_id", c.SkipZeroValue),
		c.Simple("enableEmailHashing", "enable_email_hashing", c.SkipZeroValue),
		c.Simple("enableIdentifyEvents", "enable_identify_events", c.SkipZeroValue),
		c.Simple("enablePageEvents", "enable_page_events", c.SkipZeroValue),
		c.Simple("enableScreenEvents", "enable_screen_events", c.SkipZeroValue),
		c.ArrayWithObjects("rudderToImpactProperty", "rudder_to_impact_property", map[string]interface{}{
			"from": "from",
			"to":   "to",
		}),
		c.ArrayWithObjects("productsMapping", "products_mapping", map[string]interface{}{
			"from": "from",
			"to":   "to",
		}),
		c.ArrayWithStrings("actionEventNames", "eventName", "action_event_names"),
		c.ArrayWithStrings("installEventNames", "eventName", "install_event_names"),
		c.Simple("connectionMode.android", "connection_mode.0.android", c.SkipZeroValue),
		c.Simple("connectionMode.androidKotlin", "connection_mode.0.android_kotlin", c.SkipZeroValue),
		c.Simple("connectionMode.ios", "connection_mode.0.ios", c.SkipZeroValue),
		c.Simple("connectionMode.iosSwift", "connection_mode.0.ios_swift", c.SkipZeroValue),
		c.Simple("connectionMode.web", "connection_mode.0.web", c.SkipZeroValue),
		c.Simple("connectionMode.unity", "connection_mode.0.unity", c.SkipZeroValue),
		c.Simple("connectionMode.amp", "connection_mode.0.amp", c.SkipZeroValue),
		c.Simple("connectionMode.cloud", "connection_mode.0.cloud", c.SkipZeroValue),
		c.Simple("connectionMode.warehouse", "connection_mode.0.warehouse", c.SkipZeroValue),
		c.Simple("connectionMode.reactnative", "connection_mode.0.reactnative", c.SkipZeroValue),
		c.Simple("connectionMode.flutter", "connection_mode.0.flutter", c.SkipZeroValue),
		c.Simple("connectionMode.cordova", "connection_mode.0.cordova", c.SkipZeroValue),
		c.Simple("connectionMode.shopify", "connection_mode.0.shopify", c.SkipZeroValue),
	}

	properties = append(properties, commonProperties...)

	schema := map[string]*schema.Schema{
		"account_sid": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Your impact.com Account SID.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|^(.{1,100})$"),
		},
		"api_key": {
			Type:             schema.TypeString,
			Required:         true,
			Sensitive:        true,
			Description:      "Your impact.com API Key.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|^(.{1,100})$"),
		},
		"campaign_id": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Your impact.com Campaign ID.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|^[0-9]+$"),
		},
		"event_type_id": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Your impact.com Event Type ID.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|^[0-9]+$|^$"),
		},
		"impact_app_id": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Your impact.com App ID.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|^[0-9]+$|^$"),
		},
		"enable_email_hashing": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable email hashing.",
		},
		"enable_identify_events": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable identify events.",
		},
		"enable_page_events": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable page events.",
		},
		"enable_screen_events": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable screen events.",
		},
		"rudder_to_impact_property": {
			Type:        schema.TypeList,
			Optional:    true,
			ConfigMode:  schema.SchemaConfigModeAttr,
			Description: "Map RudderStack properties to impact.com properties.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"from": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|^(.{0,100})$"),
					},
					"to": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|^(.{0,100})$"),
					},
				},
			},
		},
		"products_mapping": {
			Type:        schema.TypeList,
			Optional:    true,
			ConfigMode:  schema.SchemaConfigModeAttr,
			Description: "Map RudderStack product properties to impact.com product properties.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"from": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|^(.{0,100})$"),
					},
					"to": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|^(.{0,100})$"),
					},
				},
			},
		},
		"action_event_names": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of action event names.",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"install_event_names": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of install event names.",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"connection_mode": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Connection mode per source type.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"android": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$"),
					},
					"android_kotlin": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$"),
					},
					"ios": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$"),
					},
					"ios_swift": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$"),
					},
					"web": {
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
				},
			},
		},
	}

	for key, value := range commonSchema {
		schema[key] = value
	}

	c.Destinations.Register("impact", c.ConfigMeta{
		APIType:      "IMPACT",
		Version:      1,
		Properties:   properties,
		ConfigSchema: schema,
	})
}
