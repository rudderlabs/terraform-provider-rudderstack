package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"android", "androidKotlin", "ios", "iosSwift", "web", "unity", "amp", "cloud", "warehouse", "reactnative", "flutter", "cordova", "shopify"}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)

	properties := []c.ConfigProperty{
		c.Simple("bootstrapServer", "bootstrap_server"),
		c.Simple("topic", "topic"),
		c.Simple("apiKey", "api_key"),
		c.Simple("apiSecret", "api_secret"),
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

	destinationSchema := map[string]*schema.Schema{
		"bootstrap_server": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Bootstrap server for your Confluent Cloud cluster.",
			ValidateDiagFunc: c.StringMatchesRegexp("^(.{0,100})$"),
		},
		"topic": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Topic name in your Confluent Cloud cluster where the data will be sent.",
			ValidateDiagFunc: c.StringMatchesRegexp("^(.{0,100})$"),
		},
		"api_key": {
			Type:             schema.TypeString,
			Required:         true,
			Sensitive:        true,
			Description:      "API key for your Confluent Cloud cluster.",
			ValidateDiagFunc: c.StringMatchesRegexp("^(.{0,100})$"),
		},
		"api_secret": {
			Type:             schema.TypeString,
			Required:         true,
			Sensitive:        true,
			Description:      "API secret for your Confluent Cloud cluster.",
			ValidateDiagFunc: c.StringMatchesRegexp("^(.{0,100})$"),
		},
		"connection_mode": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Use this setting to set how you want to route events from your source to destination.",
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
		destinationSchema[key] = value
	}

	c.Destinations.Register("confluent_cloud", c.ConfigMeta{
		APIType:      "CONFLUENT_CLOUD",
		Properties:   properties,
		ConfigSchema: destinationSchema,
	})
}
