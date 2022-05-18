package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	c.Destinations.Register("facebook_pixel", c.ConfigMeta{
		APIType: "FACEBOOK_PIXEL",
		Properties: []c.ConfigProperty{
			c.Simple("pixelId", "pixel_id"),
			c.Simple("accessToken", "access_token", c.SkipZeroValue),
			c.Simple("standardPageCall", "standard_page_call", c.SkipZeroValue),
			c.Simple("valueFieldIdentifier", "value_field_identifier", c.SkipZeroValue),
			c.Simple("advancedMapping", "advanced_mapping", c.SkipZeroValue),
			c.Simple("testDestination", "test_destination", c.SkipZeroValue),
			c.Simple("testEventCode", "test_event_code", c.SkipZeroValue),
			c.Simple("eventsToEvents", "events_to_events", c.SkipZeroValue),
			c.ArrayWithStrings("eventCustomProperties", "eventCustomProperties", "event_custom_properties"),
			// TODO: figure out why blacklistPiiProperties is different than whitelistPiiProperties and what blacklist hash is about
			c.ArrayWithStrings("blacklistPiiProperties", "blacklistPiiProperties", "blacklist_pii_properties"),
			c.ArrayWithStrings("whitelistPiiProperties", "whitelistPiiProperties", "whitelist_pii_properties"),

			c.Simple("categoryToContent", "category_to_content", c.SkipZeroValue),
			c.Simple("legacyConversionPixelId.from", "legacy_conversion_pixel_id.0.from"),
			c.Simple("legacyConversionPixelId.to", "legacy_conversion_pixel_id.0.to"),
			c.Simple("useNativeSDK.web", "use_native_sdk.0.web"),
			c.ArrayWithStrings("oneTrustCookieCategories.web", "oneTrustCookieCategory", "onetrust_cookie_categories.0.web"),
			c.ArrayWithStrings("whitelistedEvents", "eventName", "event_filtering.0.whitelist"),
			c.ArrayWithStrings("blacklistedEvents", "eventName", "event_filtering.0.blacklist"),
			c.Discriminator("eventFilteringOption", c.DiscriminatorValues{
				"event_filtering.0.whitelist": "whitelistedEvents",
				"event_filtering.0.blacklist": "blacklistedEvents",
			}),
		},
		ConfigSchema: map[string]*schema.Schema{
			"pixel_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"access_token": {
				Type:             schema.TypeString,
				Optional:         true,
				Sensitive:        true,
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,208})$"),
			},
			"standard_page_call": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"value_field_identifier": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(properties.value|properties.price)$"),
			},
			"advanced_mapping": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"test_destination": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"test_event_code": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
			},
			"events_to_events": {
				Type:       schema.TypeList,
				MaxItems:   10,
				Optional:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
						},
						"to": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
						},
					},
				},
			},
			"event_custom_properties": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"blacklist_pii_properties": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"whitelist_pii_properties": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"category_to_content": {
				Type:       schema.TypeList,
				Optional:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
						},
						"to": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
						},
					},
				},
			},
			"legacy_conversion_pixel_id": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
						},
						"to": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
						},
					},
				},
			},
			"use_native_sdk": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"event_filtering": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"whitelist": {
							Type:         schema.TypeList,
							Optional:     true,
							ExactlyOneOf: []string{"config.0.event_filtering.0.whitelist", "config.0.event_filtering.0.blacklist"},
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"blacklist": {
							Type:         schema.TypeList,
							Optional:     true,
							ExactlyOneOf: []string{"config.0.event_filtering.0.whitelist", "config.0.event_filtering.0.blacklist"},
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"onetrust_cookie_categories": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web": {
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
