package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"web"}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)

	properties := []c.ConfigProperty{
		c.Simple("advertiserId", "advertiser_id"),
		c.ArrayWithObjects("eventMappingFromConfig", "event_mapping_from_config", map[string]interface{}{
			"from": "from",
			"to":   "to",
		}),
		c.ArrayWithStrings("whitelistedEvents", "eventName", "event_filtering.0.whitelist"),
		c.ArrayWithStrings("blacklistedEvents", "eventName", "event_filtering.0.blacklist"),
		c.Discriminator("eventFilteringOption", c.DiscriminatorValues{
			"event_filtering.0.whitelist": "whitelistedEvents",
			"event_filtering.0.blacklist": "blacklistedEvents",
		}),
		c.Simple("useNativeSDK.web", "use_native_sdk.0.web"),
		c.Simple("connectionMode.web", "connection_mode.0.web", c.SkipZeroValue),
	}

	properties = append(properties, commonProperties...)

	schema := map[string]*schema.Schema{
		"advertiser_id": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "The Reddit Pixel ID.",
			ValidateDiagFunc: c.StringMatchesRegexp("^(.{0,100})$"),
		},
		"event_mapping_from_config": {
			Type:        schema.TypeList,
			Optional:    true,
			ConfigMode:  schema.SchemaConfigModeAttr,
			Description: "Map RudderStack events to Reddit Pixel events.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"from": {
						Type:             schema.TypeString,
						Required:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|^(.{0,100})$"),
					},
					"to": {
						Type:     schema.TypeString,
						Required: true,
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
							"AddToCart",
							"AddToWishlist",
							"Purchase",
							"Lead",
							"ViewContent",
							"Search",
							"SignUp",
							"",
						}, false)),
					},
				},
			},
		},
		"event_filtering": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Determine which events are allowed or blocked from flowing to Reddit Pixel.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"whitelist": {
						Type:         schema.TypeList,
						Optional:     true,
						Description:  "Event names to allow.",
						ExactlyOneOf: []string{"config.0.event_filtering.0.whitelist", "config.0.event_filtering.0.blacklist"},
						Elem:         &schema.Schema{Type: schema.TypeString},
					},
					"blacklist": {
						Type:         schema.TypeList,
						Optional:     true,
						Description:  "Event names to block.",
						ExactlyOneOf: []string{"config.0.event_filtering.0.whitelist", "config.0.event_filtering.0.blacklist"},
						Elem:         &schema.Schema{Type: schema.TypeString},
					},
				},
			},
		},
		"use_native_sdk": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Configure whether web events are sent using the native Reddit Pixel SDK.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"web": {
						Type:     schema.TypeBool,
						Optional: true,
					},
				},
			},
		},
		"connection_mode": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Set the connection mode used to send web events to Reddit Pixel.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"web": {
						Type:             schema.TypeString,
						Optional:         true,
						Description:      "Connection mode for web sources. Reddit Pixel only supports device mode.",
						ValidateDiagFunc: c.StringMatchesRegexp("^(device)$"),
					},
				},
			},
		},
	}

	for key, value := range commonSchema {
		schema[key] = value
	}

	c.Destinations.Register("reddit_pixel", c.ConfigMeta{
		APIType:      "REDDIT_PIXEL",
		Version:      1,
		Properties:   properties,
		ConfigSchema: schema,
	})
}
