package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"android", "androidKotlin", "ios", "iosSwift", "web", "unity", "amp", "cloud", "warehouse", "reactnative", "flutter", "cordova", "shopify"}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)

	properties := []c.ConfigProperty{
		c.Simple("authorizationType", "authorization_type"),
		c.Simple("apiVersion", "api_version"),
		c.Simple("apiKey", "api_key", c.SkipZeroValue),
		c.Simple("accessToken", "access_token", c.SkipZeroValue),
		c.Simple("hubID", "hub_id", c.SkipZeroValue),
		c.Simple("lookupField", "lookup_field", c.SkipZeroValue),
		c.Simple("doAssociation", "do_association", c.SkipZeroValue),
		c.ArrayWithObjects("hubspotEvents", "hubspot_events", map[string]interface{}{
			"rsEventName":      "rs_event_name",
			"hubspotEventName": "hubspot_event_name",
			"eventProperties":  "event_properties",
		}),
		c.ArrayWithStrings("whitelistedEvents", "eventName", "event_filtering.0.whitelist"),
		c.ArrayWithStrings("blacklistedEvents", "eventName", "event_filtering.0.blacklist"),
		c.Discriminator("eventFilteringOption", c.DiscriminatorValues{
			"event_filtering.0.whitelist": "whitelistedEvents",
			"event_filtering.0.blacklist": "blacklistedEvents",
		}),
		c.Simple("useNativeSDK.web", "use_native_sdk.0.web", c.SkipZeroValue),
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

	destinationSchema := map[string]*schema.Schema{
		"authorization_type": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Authorization type: API Key (legacy) or Private Apps.",
			ValidateDiagFunc: c.StringMatchesRegexp("^(legacyApiKey|newPrivateAppApi)$"),
		},
		"api_version": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "HubSpot API version to use.",
			ValidateDiagFunc: c.StringMatchesRegexp("^(newApi|legacyApi)$"),
		},
		"api_key": {
			Type:             schema.TypeString,
			Optional:         true,
			Sensitive:        true,
			Description:      "Your API Key (Settings -> Integrations -> API Key). Required when Authorization Type is API Key.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"access_token": {
			Type:             schema.TypeString,
			Optional:         true,
			Sensitive:        true,
			Description:      "Your private app access token. Required when Authorization Type is Private Apps.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"hub_id": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Your Hub ID (under account name).",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
		},
		"lookup_field": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "HubSpot property name to be used for upsert. Required when API Version is New API (v3).",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"do_association": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Create associations between object records.",
		},
		"hubspot_events": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Map RudderStack event names to HubSpot Custom Behavioral event names.",
			ConfigMode:  schema.SchemaConfigModeAttr,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"rs_event_name": {
						Type:             schema.TypeString,
						Required:         true,
						Description:      "RudderStack event name.",
						ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
					},
					"hubspot_event_name": {
						Type:             schema.TypeString,
						Required:         true,
						Description:      "HubSpot Custom Behavioral event name.",
						ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
					},
					"event_properties": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "Map RudderStack event properties to HubSpot event properties.",
						ConfigMode:  schema.SchemaConfigModeAttr,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"from": {
									Type:             schema.TypeString,
									Required:         true,
									Description:      "RudderStack property name.",
									ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
								},
								"to": {
									Type:             schema.TypeString,
									Required:         true,
									Description:      "HubSpot property name.",
									ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
								},
							},
						},
					},
				},
			},
		},
		"event_filtering": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Client-side event filtering: allowlist or denylist events.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"whitelist": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "Allowlisted event names.",
						Elem:        &schema.Schema{Type: schema.TypeString},
					},
					"blacklist": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "Denylisted event names.",
						Elem:        &schema.Schema{Type: schema.TypeString},
					},
				},
			},
		},

		"use_native_sdk": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Use native SDK for specific source types.",
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
			Description: "Configure the connection mode for HubSpot.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"web": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud|device)$"),
					},
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

	c.Destinations.Register("hs", c.ConfigMeta{
		APIType:      "HS",
		Properties:   properties,
		ConfigSchema: destinationSchema,
	})
}
