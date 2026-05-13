package destinations

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"android", "androidKotlin", "ios", "iosSwift", "web", "unity", "amp", "cloud", "reactnative", "cloudSource", "flutter", "cordova", "shopify"}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)

	properties := []c.ConfigProperty{
		c.Simple("account", "account"),
		c.Simple("database", "database"),
		c.Simple("warehouse", "warehouse"),
		c.Simple("user", "user"),
		c.Simple("role", "role", c.SkipZeroValue),
		c.Simple("namespace", "namespace"),
		c.Simple("privateKey", "private_key"),
		c.Simple("privateKeyPassphrase", "private_key_passphrase", c.SkipZeroValue),
		c.Simple("skipTracksTable", "skip_tracks_table"),
		c.Simple("jsonPaths", "json_paths", c.SkipZeroValue),
		c.Simple("enableIceberg", "enable_iceberg"),
		c.Simple("externalVolume", "external_volume", c.SkipZeroValue),
		c.Simple("underscoreDivideNumbers", "underscore_divide_numbers"),
		c.Simple("allowUsersContextTraits", "allow_users_context_traits"),
		c.Simple("connectionMode.android", "connection_mode.0.android", c.SkipZeroValue),
		c.Simple("connectionMode.androidKotlin", "connection_mode.0.android_kotlin", c.SkipZeroValue),
		c.Simple("connectionMode.ios", "connection_mode.0.ios", c.SkipZeroValue),
		c.Simple("connectionMode.iosSwift", "connection_mode.0.ios_swift", c.SkipZeroValue),
		c.Simple("connectionMode.web", "connection_mode.0.web", c.SkipZeroValue),
		c.Simple("connectionMode.unity", "connection_mode.0.unity", c.SkipZeroValue),
		c.Simple("connectionMode.amp", "connection_mode.0.amp", c.SkipZeroValue),
		c.Simple("connectionMode.reactnative", "connection_mode.0.reactnative", c.SkipZeroValue),
		c.Simple("connectionMode.flutter", "connection_mode.0.flutter", c.SkipZeroValue),
		c.Simple("connectionMode.cordova", "connection_mode.0.cordova", c.SkipZeroValue),
		c.Simple("connectionMode.shopify", "connection_mode.0.shopify", c.SkipZeroValue),
		c.Simple("connectionMode.cloud", "connection_mode.0.cloud", c.SkipZeroValue),
		c.Simple("connectionMode.cloudSource", "connection_mode.0.cloud_source", c.SkipZeroValue),
	}

	properties = append(properties, oneTrustCookieCategoriesProperties(supportedSourceTypes)...)
	properties = append(properties, ketchConsentPurposesProperties(supportedSourceTypes)...)
	properties = append(properties, commonProperties...)

	schema := map[string]*schema.Schema{
		"account": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Account ID of your Snowflake warehouse.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"database": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Name of the Snowflake database.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"warehouse": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Name of the Snowflake warehouse.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"user": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Snowflake user name.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"role": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Optional Snowflake role.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
		},
		"namespace": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Schema name where tables are created.",
			ValidateDiagFunc: c.ValidateAll(
				c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^.{1,64}$"),
				c.StringNotMatchesRegexp("(?i)^pg_"),
			),
		},
		"private_key": {
			Type:             schema.TypeString,
			Required:         true,
			Sensitive:        true,
			Description:      "Private key for Snowpipe Streaming auth.",
			ValidateDiagFunc: c.StringMatchesRegexp("-----BEGIN (?:ENCRYPTED )?PRIVATE KEY-----[\\s\\S]+?-----END (?:ENCRYPTED )?PRIVATE KEY-----"),
		},
		"private_key_passphrase": {
			Type:             schema.TypeString,
			Optional:         true,
			Sensitive:        true,
			Description:      "Passphrase for encrypted private key.",
			ValidateDiagFunc: c.StringMatchesRegexp("^(.{0,100})$"),
		},
		"skip_tracks_table": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Skip writing events to tracks table.",
		},
		"json_paths": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "JSON columns in dot notation separated by commas.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.*)$"),
		},
		"enable_iceberg": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Create Iceberg tables.",
		},
		"external_volume": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "External volume name, required when Iceberg is enabled.",
			ValidateDiagFunc: c.StringMatchesRegexp("^(.{1,100})$"),
		},
		"underscore_divide_numbers": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Use underscores to split numeric segments.",
		},
		"allow_users_context_traits": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Allow users context traits.",
		},
		"connection_mode": connectionModeSchema(),
		"one_trust_cookie_categories": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Legacy OneTrust cookie category mapping by source type.",
			Elem:        &schema.Resource{Schema: oneTrustSourceSchema(supportedSourceTypes)},
		},
		"ketch_consent_purposes": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Legacy Ketch purpose mapping by source type.",
			Elem:        &schema.Resource{Schema: ketchSourceSchema(supportedSourceTypes)},
		},
	}

	for key, value := range commonSchema {
		schema[key] = value
	}

	c.Destinations.Register("snowflake_streaming", c.ConfigMeta{
		APIType:      "SNOWPIPE_STREAMING",
		Properties:   properties,
		ConfigSchema: schema,
	})
}

func connectionModeSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		MaxItems:    1,
		Optional:    true,
		Description: "Connection mode configuration per source type.",
		Elem: &schema.Resource{Schema: map[string]*schema.Schema{
			"android":        connectionModeEntry(),
			"android_kotlin": connectionModeEntry(),
			"ios":            connectionModeEntry(),
			"ios_swift":      connectionModeEntry(),
			"web":            connectionModeEntry(),
			"unity":          connectionModeEntry(),
			"amp":            connectionModeEntry(),
			"reactnative":    connectionModeEntry(),
			"cloud":          connectionModeEntry(),
			"cloud_source":   connectionModeEntry(),
			"flutter":        connectionModeEntry(),
			"cordova":        connectionModeEntry(),
			"shopify":        connectionModeEntry(),
		}},
	}
}

func connectionModeEntry() *schema.Schema {
	return &schema.Schema{
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$"),
	}
}

func oneTrustCookieCategoriesProperties(sourceTypes []string) []c.ConfigProperty {
	properties := make([]c.ConfigProperty, 0, len(sourceTypes))
	for _, sourceType := range sourceTypes {
		properties = append(properties, c.ArrayWithObjects(
			fmt.Sprintf("oneTrustCookieCategories.%s", sourceType),
			fmt.Sprintf("one_trust_cookie_categories.0.%s", camelToSnake(sourceType)),
			map[string]interface{}{
				"oneTrustCookieCategory": "one_trust_cookie_category",
			},
		))
	}
	return properties
}

func ketchConsentPurposesProperties(sourceTypes []string) []c.ConfigProperty {
	properties := make([]c.ConfigProperty, 0, len(sourceTypes))
	for _, sourceType := range sourceTypes {
		properties = append(properties, c.ArrayWithObjects(
			fmt.Sprintf("ketchConsentPurposes.%s", sourceType),
			fmt.Sprintf("ketch_consent_purposes.0.%s", camelToSnake(sourceType)),
			map[string]interface{}{
				"purpose": "purpose",
			},
		))
	}
	return properties
}

func oneTrustSourceSchema(sourceTypes []string) map[string]*schema.Schema {
	elements := map[string]*schema.Schema{}
	for _, sourceType := range sourceTypes {
		elements[camelToSnake(sourceType)] = &schema.Schema{
			Type:       schema.TypeList,
			Optional:   true,
			ConfigMode: schema.SchemaConfigModeAttr,
			Elem: &schema.Resource{Schema: map[string]*schema.Schema{
				"one_trust_cookie_category": {
					Type:             schema.TypeString,
					Optional:         true,
					ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
				},
			}},
		}
	}
	return elements
}

func ketchSourceSchema(sourceTypes []string) map[string]*schema.Schema {
	elements := map[string]*schema.Schema{}
	for _, sourceType := range sourceTypes {
		elements[camelToSnake(sourceType)] = &schema.Schema{
			Type:       schema.TypeList,
			Optional:   true,
			ConfigMode: schema.SchemaConfigModeAttr,
			Elem: &schema.Resource{Schema: map[string]*schema.Schema{
				"purpose": {
					Type:             schema.TypeString,
					Optional:         true,
					ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
				},
			}},
		}
	}
	return elements
}
