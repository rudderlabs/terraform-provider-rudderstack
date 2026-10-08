package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"android", "androidKotlin", "ios", "iosSwift", "web", "unity", "amp", "cloud", "warehouse", "reactnative", "flutter", "cordova", "shopify"}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)

	properties := []c.ConfigProperty{
		c.Simple("rudderAccountId", "rudder_account_id"),
		c.Simple("version", "version"),
		c.Simple("accountId", "account_id"),
		c.Simple("hashData", "hash_data"),
		c.ArrayWithObjects("eventsMapping", "events_mapping", map[string]interface{}{
			"from": "from",
			"to":   "to",
		}),
		c.Simple("connectionMode.web", "connection_mode.0.web", c.SkipZeroValue),
		c.Simple("connectionMode.android", "connection_mode.0.android", c.SkipZeroValue),
		c.Simple("connectionMode.androidKotlin", "connection_mode.0.android_kotlin", c.SkipZeroValue),
		c.Simple("connectionMode.ios", "connection_mode.0.ios", c.SkipZeroValue),
		c.Simple("connectionMode.iosSwift", "connection_mode.0.ios_swift", c.SkipZeroValue),
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
		"rudder_account_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The RudderStack account ID for OAuth-based event delivery.",
		},
		"version": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "v3",
			Description:      "The Reddit Conversions API version.",
			ValidateDiagFunc: c.StringMatchesRegexp("^(v3|v2)$"),
		},
		"account_id": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "The Pixel ID of the Reddit Ads account associated with the conversion events.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|^(.{1,100})$"),
		},
		"hash_data": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Hash data before delivery. Disable this when sending pre-hashed email, user ID, IP, and advertiser ID values.",
		},
		"events_mapping": {
			Type:        schema.TypeList,
			Optional:    true,
			ConfigMode:  schema.SchemaConfigModeAttr,
			Description: "Map RudderStack events to Reddit events.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"from": {
						Type:             schema.TypeString,
						Required:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|^(.{1,100})$"),
					},
					"to": {
						Type:     schema.TypeString,
						Required: true,
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
							"ViewContent",
							"Search",
							"AddToCart",
							"AddToWishlist",
							"Purchase",
							"SignUp",
							"Lead",
							"PageVisit",
						}, false)),
					},
				},
			},
		},
		"connection_mode": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Configure the connection mode for Reddit.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"web":            {Type: schema.TypeString, Optional: true, ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$")},
					"android":        {Type: schema.TypeString, Optional: true, ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$")},
					"android_kotlin": {Type: schema.TypeString, Optional: true, ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$")},
					"ios":            {Type: schema.TypeString, Optional: true, ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$")},
					"ios_swift":      {Type: schema.TypeString, Optional: true, ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$")},
					"unity":          {Type: schema.TypeString, Optional: true, ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$")},
					"amp":            {Type: schema.TypeString, Optional: true, ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$")},
					"cloud":          {Type: schema.TypeString, Optional: true, ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$")},
					"warehouse":      {Type: schema.TypeString, Optional: true, ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$")},
					"reactnative":    {Type: schema.TypeString, Optional: true, ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$")},
					"flutter":        {Type: schema.TypeString, Optional: true, ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$")},
					"cordova":        {Type: schema.TypeString, Optional: true, ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$")},
					"shopify":        {Type: schema.TypeString, Optional: true, ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$")},
				},
			},
		},
	}

	for key, value := range commonSchema {
		schema[key] = value
	}

	c.Destinations.Register("reddit", c.ConfigMeta{
		APIType:      "REDDIT",
		Version:      1,
		Properties:   properties,
		ConfigSchema: schema,
	})
}
