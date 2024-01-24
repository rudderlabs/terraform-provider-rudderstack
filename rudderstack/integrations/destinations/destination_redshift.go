package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"web", "android", "ios", "unity", "reactnative", "flutter", "cordova", "amp", "cloud", "warehouse", "shopify"}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)

	properties := []c.ConfigProperty{
		c.Simple("host", "host"),
		c.Simple("port", "port"),
		c.Simple("database", "database"),
		c.Simple("user", "user"),
		c.Simple("password", "password"),
		c.Simple("namespace", "namespace", c.SkipZeroValue),
		c.Simple("enableSSE", "enable_sse", c.SkipZeroValue),
		c.Simple("useRudderStorage", "use_rudder_storage"),
		c.Simple("syncFrequency", "sync.0.frequency"),
		c.Simple("syncStartAt", "sync.0.start_at", c.SkipZeroValue),
		c.Simple("excludeWindow.excludeWindowStartTime", "sync.0.exclude_window_start_time", c.SkipZeroValue),
		c.Simple("excludeWindow.excludeWindowEndTime", "sync.0.exclude_window_end_time", c.SkipZeroValue),
		c.Simple("bucketName", "s3.0.bucket_name"),
		c.Simple("accessKeyID", "s3.0.access_key_id"),
		c.Simple("accessKey", "s3.0.access_key"),
	}

	properties = append(properties, commonProperties...)

	schema := map[string]*schema.Schema{
		"host": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "The host name of your Redshift service.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,255})$"),
		},
		"port": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "The port number associated with the Redshift database instance.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
		},
		"user": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "The name of the user with the required read/write access to the above database.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
		},
		"password": {
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
			Sensitive:        true,
			Description:      "The password for the above user.",
		},
		"database": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "The database name in your Redshift instance where the data will be sent.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
		},
		"namespace": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Enter the schema name where RudderStack will create all the tables. If you don't specify any namespace, RudderStack will set this to the source name, by default.",
			ValidateDiagFunc: c.ValidateAll(
				c.StringMatchesRegexp("(^env[.].*)|^(.{0,64})$"),
				c.StringNotMatchesRegexp("^(pg_|PG_|pG_|Pg_)"),
			),
		},
		"enable_sse": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "This setting enables server-side encryption.",
		},
		"use_rudder_storage": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "Enable this setting to use the RudderStack-hosted object storage.",
		},
		"sync": {
			Type:     schema.TypeList,
			MinItems: 1, MaxItems: 1,
			Required:    true,
			Description: "Specify your sync settings.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"frequency": {
						Type:             schema.TypeString,
						Required:         true,
						Description:      "Specify how often RudderStack should sync the data to your Redshift database.",
						ValidateDiagFunc: c.StringMatchesRegexp("^(30|60|180|360|720|1440)$"),
					},
					"start_at": {
						Type:             schema.TypeString,
						Optional:         true,
						Description:      "Specify the particular time of the day (in UTC) when you want RudderStack to sync the data to the warehouse.",
						ValidateDiagFunc: c.StringMatchesRegexp("^([01][0-9]|2[0-3]):[0-5][0-9]$"),
					},
					"exclude_window_start_time": {
						Type:             schema.TypeString,
						Optional:         true,
						Description:      "This optional setting lets you set a time window when RudderStack will not sync the data to your database.",
						ValidateDiagFunc: c.StringMatchesRegexp("^([01][0-9]|2[0-3]):[0-5][0-9]$"),
					},
					"exclude_window_end_time": {
						Type:             schema.TypeString,
						Optional:         true,
						Description:      "Set the end time of the exclusion window.",
						ValidateDiagFunc: c.StringMatchesRegexp("^([01][0-9]|2[0-3]):[0-5][0-9]$"),
					},
				},
			},
		},
		"s3": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"bucket_name": {
						Type:             schema.TypeString,
						Required:         true,
						Description:      "Enter the name of your S3 bucket.",
						ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
					},
					"access_key_id": {
						Type:             schema.TypeString,
						Optional:         true,
						Description:      "Enter your AWS access key ID.",
						ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
					},
					"access_key": {
						Type:             schema.TypeString,
						Optional:         true,
						Sensitive:        true,
						Description:      "Enter your AWS secret access key.",
						ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
					},
				},
			},
		},
	}

	for key, value := range commonSchema {
		schema[key] = value
	}

	c.Destinations.Register("redshift", c.ConfigMeta{
		APIType: "RS",
		Properties: properties,
		ConfigSchema: schema,
	})
}
