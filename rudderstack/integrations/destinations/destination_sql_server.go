package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"web", "android", "ios", "unity", "reactnative", "flutter", "cordova", "amp", "cloud", "cloudSource", "shopify"}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)

	properties := []c.ConfigProperty{
		c.Simple("host", "host"),
		c.Simple("database", "database"),
		c.Simple("user", "user"),
		c.Simple("password", "password"),
		c.Simple("port", "port"),
		c.Simple("namespace", "namespace", c.SkipZeroValue),
		c.Simple("sslMode", "ssl_mode"),
		c.Simple("syncFrequency", "sync_frequency"),
		c.Simple("syncStartAt", "sync_start_at", c.SkipZeroValue),
		c.Simple("excludeWindow.excludeWindowStartTime", "exclude_window.0.exclude_window_start_time", c.SkipZeroValue),
		c.Simple("excludeWindow.excludeWindowEndTime", "exclude_window.0.exclude_window_end_time", c.SkipZeroValue),
		c.Simple("useRudderStorage", "use_rudder_storage", c.SkipZeroValue),
		c.Simple("bucketProvider", "bucket_provider", c.SkipZeroValue),
		c.Simple("bucketName", "bucket_name", c.SkipZeroValue),
		c.Simple("accessKeyID", "access_key_id", c.SkipZeroValue),
		c.Simple("accessKey", "access_key", c.SkipZeroValue),
		c.Simple("accountName", "account_name", c.SkipZeroValue),
		c.Simple("accountKey", "account_key", c.SkipZeroValue),
		c.Simple("credentials", "credentials", c.SkipZeroValue),
	}

	properties = append(properties, commonProperties...)

	schema := map[string]*schema.Schema{
		"host": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Enter the host name of your SQL Server database.",
		},
		"database": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Enter the name of your SQL Server database.",
		},
		"user": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Enter the username of your SQL Server database.",
		},
		"password": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "Enter the password for your SQL Server database user.",
		},
		"port": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "1433",
			Description: "Enter the port number of your SQL Server database instance.",
		},
		"namespace": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Enter the schema name where RudderStack will create all the tables. If not specified, RudderStack will set this to the source name by default.",
		},
		"ssl_mode": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "disable",
			Description:      "Choose the SSL mode through which RudderStack will connect to your SQL Server instance.",
			ValidateDiagFunc: c.StringMatchesRegexp("^(disable|true|false)$"),
		},
		"sync_frequency": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "30",
			Description:      "Specify how often RudderStack should sync the data to your SQL Server database (in minutes).",
			ValidateDiagFunc: c.StringMatchesRegexp("^(30|60|180|360|720|1440)$"),
		},
		"sync_start_at": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Specify the time (UTC) when RudderStack should start syncing data to SQL Server.",
			ValidateDiagFunc: c.StringMatchesRegexp("^([01][0-9]|2[0-3]):[0-5][0-9]$"),
		},
		"exclude_window": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Set a time window (UTC) during which RudderStack will not sync data to SQL Server.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"exclude_window_start_time": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^([01][0-9]|2[0-3]):[0-5][0-9]$"),
					},
					"exclude_window_end_time": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^([01][0-9]|2[0-3]):[0-5][0-9]$"),
					},
				},
			},
		},
		"use_rudder_storage": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "Enable this to use RudderStack-managed object storage for staging files.",
		},
		"bucket_provider": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The cloud object storage provider to use when use_rudder_storage is disabled.",
			ValidateDiagFunc: c.StringMatchesRegexp("^(S3|GCS|AZURE_BLOB|MINIO)$"),
		},
		"bucket_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The name of the object storage bucket.",
		},
		"access_key_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The AWS Access Key ID (for S3 bucket provider).",
		},
		"access_key": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "The AWS Secret Access Key (for S3 bucket provider).",
		},
		"account_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The Azure Blob Storage account name.",
		},
		"account_key": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "The Azure Blob Storage account key.",
		},
		"credentials": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "The GCS service account credentials JSON.",
		},
	}

	for key, value := range commonSchema {
		schema[key] = value
	}

	c.Destinations.Register("sql_server", c.ConfigMeta{
		APIType:      "MSSQL",
		Properties:   properties,
		ConfigSchema: schema,
	})
}
