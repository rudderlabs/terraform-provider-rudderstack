package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	c.Destinations.Register("active_campaign", c.ConfigMeta{
		APIType: "ACTIVE_CAMPAIGN",
		Properties: []c.ConfigProperty{
			c.Simple("apiUrl", "api_url"),
			c.Simple("apiKey", "api_key", c.SkipZeroValue),
			c.Simple("actid", "actid", c.SkipZeroValue),
			c.Simple("eventKey", "event_key", c.SkipZeroValue),
			c.ArrayWithStrings("oneTrustCookieCategories.web", "oneTrustCookieCategory", "onetrust_cookie_categories.0.web"),
			c.ArrayWithStrings("oneTrustCookieCategories.android", "oneTrustCookieCategory", "onetrust_cookie_categories.0.android"),
			c.ArrayWithStrings("oneTrustCookieCategories.ios", "oneTrustCookieCategory", "onetrust_cookie_categories.0.ios"),
			c.ArrayWithStrings("oneTrustCookieCategories.unity", "oneTrustCookieCategory", "onetrust_cookie_categories.0.unity"),
			c.ArrayWithStrings("oneTrustCookieCategories.reactnative", "oneTrustCookieCategory", "onetrust_cookie_categories.0.reactnative"),
			c.ArrayWithStrings("oneTrustCookieCategories.flutter", "oneTrustCookieCategory", "onetrust_cookie_categories.0.flutter"),
			c.ArrayWithStrings("oneTrustCookieCategories.cordova", "oneTrustCookieCategory", "onetrust_cookie_categories.0.cordova"),
			c.ArrayWithStrings("oneTrustCookieCategories.amp", "oneTrustCookieCategory", "onetrust_cookie_categories.0.amp"),
			c.ArrayWithStrings("oneTrustCookieCategories.cloud", "oneTrustCookieCategory", "onetrust_cookie_categories.0.cloud"),
			c.ArrayWithStrings("oneTrustCookieCategories.warehouse", "oneTrustCookieCategory", "onetrust_cookie_categories.0.warehouse"),
			c.ArrayWithStrings("oneTrustCookieCategories.shopify", "oneTrustCookieCategory", "onetrust_cookie_categories.0.shopify"),
		},
		ConfigSchema: map[string]*schema.Schema{
			"api_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Enter your ActiveCampaign API URL. You can find it in your account in the Settings page under the Developer tab.",
				// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Enter your ActiveCampaign API key.",
				// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"actid": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enter your ActID here. To obtain the ActID unique to your ActiveCampaign account, go to Settings > Tracking > Event Tracking API.",
				// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"event_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enter the event key unique to your ActiveCampaign account. To obtain the event key, go to your ActiveCampaign account > Settings > Tracking > Event Tracking.",
				// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"onetrust_cookie_categories": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Specify OneTrust category IDs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"android": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"ios": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"unity": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"reactnative": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"flutter": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cordova": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"amp": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cloud": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"warehouse": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"shopify": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	})
}
