package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	c.Destinations.Register("statsig", c.ConfigMeta{
		APIType: "STATSIG",
		Properties: []c.ConfigProperty{
			c.Simple("secretKey", "secret_key"),
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
			c.ArrayWithStrings("oneTrustCookieCategories", "oneTrustCookieCategory", "onetrust_cookie_categories"),
		},
		ConfigSchema: map[string]*schema.Schema{
			"secret_key": {
				Type:             schema.TypeString,
				Required:         true,
				Sensitive:        true,
				Description:      "Enter the Secret Key.",
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,200})$"),
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
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].*)|^(cloud)$"),
						},
						"ios": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].*)|^(cloud)$"),
						},
						"android": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].*)|^(cloud)$"),
						},
						"react_native": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].*)|^(cloud)$"),
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
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].*)|^(cloud)$"),
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
