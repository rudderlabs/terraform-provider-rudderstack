package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"web", "android", "ios", "unity", "reactnative", "flutter", "cordova", "amp", "cloud", "warehouse", "shopify"}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)

	properties := []c.ConfigProperty{
		c.Simple("accountId", "account_id"),
		c.Simple("clientId", "client_id"),
		c.Simple("clientSecret", "client_secret"),
		c.Simple("trackAnonymousEvents", "track_anonymous_events"),
		c.Simple("createIfNotExist", "create_if_not_exist"),
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
		c.ArrayWithObjects("rudderEventsMapping", "rudder_events_mapping", map[string]interface{}{
			"event":             "event",
			"marketoPrimarykey": "marketo_primarykey",
			"marketoActivityId": "marketo_activity_id",
		}),
		c.ArrayWithObjects("leadTraitMapping", "lead_trait_mapping", map[string]interface{}{
			"from": "from",
			"to":   "to",
		}),
		c.ArrayWithObjects("customActivityPropertyMap", "custom_activity_property_map", map[string]interface{}{
			"from": "from",
			"to":   "to",
		}),
	}

	properties = append(properties, commonProperties...)

	schema := map[string]*schema.Schema{
		"account_id": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Marketo Account ID",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"client_id": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Marketo Client ID",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"client_secret": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Marketo Client Secret",
			Sensitive:        true,
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"track_anonymous_events": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Boolean flag to track anonymous events",
		},
		"create_if_not_exist": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Boolean flag to create lead if not exist",
		},
		"lead_trait_mapping": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Lead Trait Mapping",
			ConfigMode:  schema.SchemaConfigModeAttr,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"from": {
						Type:     schema.TypeString,
						Required: true,
					},
					"to": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"rudder_events_mapping": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Rudder Events Mapping",
			ConfigMode:  schema.SchemaConfigModeAttr,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"event": {
						Type:     schema.TypeString,
						Required: true,
					},
					"marketo_primarykey": {
						Type:     schema.TypeString,
						Required: true,
					},
					"marketo_activity_id": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"custom_activity_property_map": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Custom Activity Property Map",
			ConfigMode:  schema.SchemaConfigModeAttr,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"from": {
						Type:     schema.TypeString,
						Required: true,
					},
					"to": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
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
	}

	for key, value := range commonSchema {
		schema[key] = value
	}

	c.Destinations.Register("marketo", c.ConfigMeta{
		APIType:      "MARKETO",
		Properties:   properties,
		ConfigSchema: schema,
	})
}
