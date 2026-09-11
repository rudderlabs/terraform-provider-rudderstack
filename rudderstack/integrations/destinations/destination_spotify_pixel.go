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
		c.Simple("pixelId", "pixel_id"),
		c.Simple("enableAliasCall", "enable_alias_call", c.SkipZeroValue),
		c.Simple("eventsToSpotifyPixelEvents", "events_to_spotify_pixel_events", c.SkipZeroValue),
		c.ArrayWithStrings("whitelistedEvents", "eventName", "event_filtering.0.whitelist"),
		c.ArrayWithStrings("blacklistedEvents", "eventName", "event_filtering.0.blacklist"),
		c.Discriminator("eventFilteringOption", c.DiscriminatorValues{
			"event_filtering.0.whitelist": "whitelistedEvents",
			"event_filtering.0.blacklist": "blacklistedEvents",
		}),
		c.Simple("connectionMode.web", "connection_mode.0.web", c.SkipZeroValue),
	}

	properties = append(properties, commonProperties...)

	schema := map[string]*schema.Schema{
		"pixel_id": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Enter your Spotify Pixel ID.",
			ValidateDiagFunc: c.StringMatchesRegexp("^(.{1,100})$"),
		},
		"enable_alias_call": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Enable this setting to send identify calls to Spotify Pixel as alias events.",
		},
		"events_to_spotify_pixel_events": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "You can map your events to standard Spotify Pixel events using this setting.",
			ConfigMode:  schema.SchemaConfigModeAttr,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"from": {
						Type:             schema.TypeString,
						Required:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(.{0,100})$"),
					},
					"to": {
						Type:     schema.TypeString,
						Required: true,
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
							"lead",
							"product",
							"addtocart",
							"checkout",
							"purchase",
						}, false)),
					},
				},
			},
		},
		"event_filtering": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "This setting lets you determine which events are blocked or allowed to flow through to Spotify Pixel.",
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
			Optional:    true,
			Description: "Set the connection mode used to send events to Spotify Pixel for each source type.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"web": {
						Type:             schema.TypeString,
						Optional:         true,
						Description:      "Connection mode for web sources. Spotify Pixel only supports device mode.",
						ValidateDiagFunc: c.StringMatchesRegexp("^(device)$"),
					},
				},
			},
		},
	}

	for key, value := range commonSchema {
		schema[key] = value
	}

	c.Destinations.Register("spotify_pixel", c.ConfigMeta{
		APIType:      "SPOTIFYPIXEL",
		Version:      1,
		Properties:   properties,
		ConfigSchema: schema,
	})
}
