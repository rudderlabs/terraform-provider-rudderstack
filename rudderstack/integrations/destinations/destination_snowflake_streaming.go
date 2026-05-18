package destinations

import (
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
		privateKeyProperty(),
		c.Simple("privateKeyPassphrase", "private_key_passphrase", c.SkipZeroValue),
		c.Simple("skipTracksTable", "skip_tracks_table"),
		c.Simple("jsonPaths", "json_paths", c.SkipZeroValue),
		c.Simple("enableIceberg", "enable_iceberg"),
		c.Simple("externalVolume", "external_volume", c.SkipZeroValue),
		c.Simple("underscoreDivideNumbers", "underscore_divide_numbers"),
		c.Simple("allowUsersContextTraits", "allow_users_context_traits"),
		c.Simple("connectionMode.web", "connection_mode.0.web", c.SkipZeroValue),
		c.Simple("connectionMode.android", "connection_mode.0.android", c.SkipZeroValue),
		c.Simple("connectionMode.androidKotlin", "connection_mode.0.android_kotlin", c.SkipZeroValue),
		c.Simple("connectionMode.ios", "connection_mode.0.ios", c.SkipZeroValue),
		c.Simple("connectionMode.iosSwift", "connection_mode.0.ios_swift", c.SkipZeroValue),
		c.Simple("connectionMode.unity", "connection_mode.0.unity", c.SkipZeroValue),
		c.Simple("connectionMode.amp", "connection_mode.0.amp", c.SkipZeroValue),
		c.Simple("connectionMode.cloud", "connection_mode.0.cloud", c.SkipZeroValue),
		c.Simple("connectionMode.cloudSource", "connection_mode.0.cloud_source", c.SkipZeroValue),
		c.Simple("connectionMode.reactnative", "connection_mode.0.reactnative", c.SkipZeroValue),
		c.Simple("connectionMode.flutter", "connection_mode.0.flutter", c.SkipZeroValue),
		c.Simple("connectionMode.cordova", "connection_mode.0.cordova", c.SkipZeroValue),
		c.Simple("connectionMode.shopify", "connection_mode.0.shopify", c.SkipZeroValue),
	}

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
				c.StringNotMatchesRegexp("^(pg_|PG_|pG_|Pg_)"),
			),
		},
		"private_key": {
			Type:             schema.TypeString,
			Required:         true,
			Sensitive:        true,
			Description:      "Private key for Snowpipe Streaming auth. Accepts both PEM-formatted keys (with BEGIN/END headers) and raw base64-encoded key bodies. Raw keys are automatically wrapped with PEM headers before being sent to the API.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|.+"),
			DiffSuppressFunc: suppressPEMKeyDiff,
		},
		"private_key_passphrase": {
			Type:             schema.TypeString,
			Optional:         true,
			Sensitive:        true,
			Description:      "Passphrase for encrypted private key.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{0,100})$"),
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
			Description:      "External volume name. Required when `enable_iceberg` is true.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
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
		"connection_mode": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Use this setting to set how you want to route events from your source to destination.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"web": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$"),
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
					"cloud_source": {
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
		schema[key] = value
	}

	c.Destinations.Register("snowflake_streaming", c.ConfigMeta{
		APIType:      "SNOWPIPE_STREAMING",
		Properties:   properties,
		ConfigSchema: schema,
	})
}
