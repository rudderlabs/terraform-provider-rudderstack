package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"android", "androidKotlin", "ios", "iosSwift", "web", "unity", "amp", "cloud", "reactnative", "flutter", "cordova", "warehouse", "shopify"}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)

	properties := []c.ConfigProperty{
		c.Simple("rudderAccountId", "rudder_account_id"),
		c.Simple("customerId", "customer_id"),
		c.Simple("subAccount", "sub_account"),
		c.Simple("loginCustomerId", "login_customer_id", c.SkipZeroValue),
		c.ArrayWithObjects("eventsToOfflineConversionsTypeMapping", "events_to_offline_conversions_type_mapping", map[string]interface{}{
			"from": "from",
			"to":   "to",
		}),
		c.ArrayWithObjects("eventsToConversionsNamesMapping", "events_to_conversions_names_mapping", map[string]interface{}{
			"from": "from",
			"to":   "to",
		}),
		c.ArrayWithObjects("customVariables", "custom_variables", map[string]interface{}{
			"from": "from",
			"to":   "to",
		}),
		c.Simple("UserIdentifierSource", "user_identifier_source"),
		c.Simple("conversionEnvironment", "conversion_environment"),
		c.Simple("defaultUserIdentifier", "default_user_identifier"),
		c.Simple("hashUserIdentifier", "hash_user_identifier"),
		c.Simple("validateOnly", "validate_only"),
		c.Simple("connectionMode.android", "connection_mode.0.android", c.SkipZeroValue),
		c.Simple("connectionMode.androidKotlin", "connection_mode.0.android_kotlin", c.SkipZeroValue),
		c.Simple("connectionMode.ios", "connection_mode.0.ios", c.SkipZeroValue),
		c.Simple("connectionMode.iosSwift", "connection_mode.0.ios_swift", c.SkipZeroValue),
		c.Simple("connectionMode.web", "connection_mode.0.web", c.SkipZeroValue),
		c.Simple("connectionMode.unity", "connection_mode.0.unity", c.SkipZeroValue),
		c.Simple("connectionMode.amp", "connection_mode.0.amp", c.SkipZeroValue),
		c.Simple("connectionMode.cloud", "connection_mode.0.cloud", c.SkipZeroValue),
		c.Simple("connectionMode.reactnative", "connection_mode.0.reactnative", c.SkipZeroValue),
		c.Simple("connectionMode.flutter", "connection_mode.0.flutter", c.SkipZeroValue),
		c.Simple("connectionMode.cordova", "connection_mode.0.cordova", c.SkipZeroValue),
		c.Simple("connectionMode.warehouse", "connection_mode.0.warehouse", c.SkipZeroValue),
		c.Simple("connectionMode.shopify", "connection_mode.0.shopify", c.SkipZeroValue),
	}

	properties = append(properties, commonProperties...)

	destinationSchema := map[string]*schema.Schema{
		"rudder_account_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The RudderStack account ID for OAuth-based event delivery.",
		},
		"customer_id": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Enter your Google Ads customer ID.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"sub_account": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable this setting if the customer ID belongs to a sub-account.",
		},
		"login_customer_id": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "If the customer ID is from a sub-account, provide the customer ID of the manager account.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"events_to_offline_conversions_type_mapping": {
			Type:        schema.TypeList,
			Optional:    true,
			ConfigMode:  schema.SchemaConfigModeAttr,
			Description: "Map RudderStack event names to Google Ads offline conversion types.",
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
						ValidateDiagFunc: c.StringMatchesRegexp("^(click|call|store|)$"),
					},
				},
			},
		},
		"events_to_conversions_names_mapping": {
			Type:        schema.TypeList,
			Optional:    true,
			ConfigMode:  schema.SchemaConfigModeAttr,
			Description: "Map RudderStack event names to Google Ads conversion names.",
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
		"custom_variables": {
			Type:        schema.TypeList,
			Optional:    true,
			ConfigMode:  schema.SchemaConfigModeAttr,
			Description: "Map RudderStack variable names to custom Google Ads variables.",
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
		"user_identifier_source": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "none",
			Description:      "Source of the user identifier.",
			ValidateDiagFunc: c.StringMatchesRegexp("^(none|UNSPECIFIED|UNKNOWN|FIRST_PARTY|THIRD_PARTY)$"),
		},
		"conversion_environment": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "none",
			Description:      "The environment this conversion was recorded on, such as app or web.",
			ValidateDiagFunc: c.StringMatchesRegexp("^(none|UNSPECIFIED|UNKNOWN|APP|WEB)$"),
		},
		"default_user_identifier": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "email",
			Description:      "The user identifier to use for store and click conversions.",
			ValidateDiagFunc: c.StringMatchesRegexp("^(email|phone)$"),
		},
		"hash_user_identifier": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Hash user identifying information using SHA-256.",
		},
		"validate_only": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable this option to only validate the request.",
		},
		"connection_mode": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Configure the connection mode for Google Ads Offline Conversions.",
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
					"warehouse": {
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

	c.Destinations.Register("google_adwords_offline_conversions", c.ConfigMeta{
		APIType:      "GOOGLE_ADWORDS_OFFLINE_CONVERSIONS",
		Version:      1,
		Properties:   properties,
		ConfigSchema: destinationSchema,
	})
}
