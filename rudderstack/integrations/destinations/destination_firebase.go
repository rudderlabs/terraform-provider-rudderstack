package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"android", "ios", "unity", "reactnative", "flutter"}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)

	properties := []c.ConfigProperty{
		c.ArrayWithStrings("whitelistedEvents", "eventName", "event_filtering.0.whitelist"),
		c.ArrayWithStrings("blacklistedEvents", "eventName", "event_filtering.0.blacklist"),
		c.Discriminator("eventFilteringOption", c.DiscriminatorValues{
			"event_filtering.0.whitelist": "whitelistedEvents",
			"event_filtering.0.blacklist": "blacklistedEvents",
		}),

		c.Simple("connectionMode.android", "connection_mode.0.android", c.SkipZeroValue),
		c.Simple("connectionMode.ios", "connection_mode.0.ios", c.SkipZeroValue),
		c.Simple("connectionMode.unity", "connection_mode.0.unity", c.SkipZeroValue),
		c.Simple("connectionMode.reactnative", "connection_mode.0.reactnative", c.SkipZeroValue),
		c.Simple("connectionMode.flutter", "connection_mode.0.flutter", c.SkipZeroValue),
	}

	properties = append(properties, commonProperties...)

	schema := map[string]*schema.Schema{
		"event_filtering": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "This option allows you filter the events you want to send to Firebase.",
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
			Description: "Use this setting to set how you want to route events from your source to destination..",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"ios": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].*)|^(device)$"),
					},
					"reactnative": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].*)|^(device)$"),
					},
					"unity": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].*)|^(device)$"),
					},
					"flutter": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].*)|^(device)$"),
					},
					"android": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].*)|^(device)$"),
					},
				},
			},
		},
	}

	for key, value := range commonSchema {
		schema[key] = value
	}

	c.Destinations.Register("firebase", c.ConfigMeta{
		APIType:      "FIREBASE",
		Properties:   properties,
		ConfigSchema: schema,
	})
}
