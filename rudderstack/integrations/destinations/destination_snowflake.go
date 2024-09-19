package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"web", "android", "ios", "unity", "reactnative", "flutter", "cordova", "amp", "cloud", "cloud_source", "shopify"}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)

	properties := []c.ConfigProperty{
		c.Simple("account", "account"),
		c.Simple("database", "database"),
		c.Simple("warehouse", "warehouse"),
		c.Simple("user", "user"),
		c.Simple("password", "password"),
		c.Simple("role", "role", c.SkipZeroValue),
		c.Simple("namespace", "namespace", c.SkipZeroValue),
		c.Simple("syncFrequency", "sync.0.frequency"),
		c.Simple("syncStartAt", "sync.0.start_at", c.SkipZeroValue),
		c.Simple("excludeWindow.excludeWindowStartTime", "sync.0.exclude_window_start_time", c.SkipZeroValue),
		c.Simple("excludeWindow.excludeWindowEndTime", "sync.0.exclude_window_end_time", c.SkipZeroValue),
		c.Simple("jsonPaths", "json_paths", c.SkipZeroValue),
		c.Simple("useRudderStorage", "use_rudder_storage"),
		c.Discriminator("cloudProvider", c.DiscriminatorValues{
			"s3":    "AWS",
			"gcp":   "GCP",
			"azure": "AZURE",
		}),
		c.Simple("additionalProperties", "additional_properties"),
		c.Conditional("bucketName", "s3.0.bucket_name", c.Equals("cloudProvider", "AWS")),
		c.Simple("accessKeyID", "s3.0.access_key_id", c.SkipZeroValue),
		c.Simple("accessKey", "s3.0.access_key", c.SkipZeroValue),
		c.Simple("enableSSE", "s3.0.enable_sse", c.SkipZeroValue),
		c.Conditional("bucketName", "gcp.0.bucket_name", c.Equals("cloudProvider", "GCP")),
		c.Simple("credentials", "gcp.0.credentials", c.SkipZeroValue),
		c.Conditional("storageIntegration", "gcp.0.storage_integration", c.Equals("cloudProvider", "GCP")),
		c.Simple("containerName", "azure.0.container_name", c.SkipZeroValue),
		c.Simple("accountName", "azure.0.account_name", c.SkipZeroValue),
		c.Simple("accountKey", "azure.0.account_key", c.SkipZeroValue),
		c.Conditional("storageIntegration", "azure.0.storage_integration", c.Equals("cloudProvider", "AZURE")),
		c.Simple("prefix", "prefix", c.SkipZeroValue),
	}

	properties = append(properties, commonProperties...)

	schema := map[string]*schema.Schema{
		"account": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Account ID of your Snowflake warehouse. This account ID is part of the Snowflake URL. Example : https://www.rudderstack.com/docs/destinations/warehouse-destinations/faq/#while-configuring-the-snowflake-destination-what-should-i-enter-in-the-account-field",
			ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
		},
		"database": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Name of the database.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
		},
		"warehouse": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Name of the warehouse.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
		},
		"user": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Name of the user.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
		},
		"password": {
			Type:             schema.TypeString,
			Required:         true,
			Sensitive:        true,
			Description:      "Password for the user.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|.+"),
		},
		"role": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Role for the user. If not specified, the default role is used",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
		},
		"namespace": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Schema name for the warehouse where the tables are created by Rudderstack.",
			ValidateDiagFunc: c.ValidateAll(
				c.StringMatchesRegexp("(^env[.].*)|^(.{0,64})$"),
				c.StringNotMatchesRegexp("^(pg_|PG_|pG_|Pg_)"),
			),
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
						Description:      "Specify how often RudderStack should sync the data to your snowflake database.",
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
		"json_paths": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Specify required json properties in dot notation separated by commas.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].*)|.*"),
		},
		"use_rudder_storage": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Enable this setting to use RudderStack-managed buckets for object storage.",
			Default:     false,
		},
		"additional_properties": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"s3": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			Description:   "",
			ConflictsWith: []string{"config.0.gcp", "config.0.azure"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"bucket_name": {
						Type:             schema.TypeString,
						Required:         true,
						Description:      "Specify the name of your S3 bucket where RudderStack will store the data before loading it into Snowflake.",
						ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
					},
					"access_key_id": {
						Type:             schema.TypeString,
						Optional:         true,
						Sensitive:        true,
						Description:      "Enter your AWS access key ID obtained from the AWS console.",
						ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
					},
					"access_key": {
						Type:             schema.TypeString,
						Optional:         true,
						Sensitive:        true,
						Description:      "Enter your AWS secret access key.",
						ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
					},
					"enable_sse": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "Toggle on this setting to enable server-side encryption for your S3 bucket.",
					},
				},
			},
		},
		"gcp": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			Description:   "",
			ConflictsWith: []string{"config.0.s3", "config.0.azure"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"bucket_name": {
						Type:             schema.TypeString,
						Required:         true,
						Description:      "Specify the name of your GCS bucket where RudderStack will store the data before loading it into Snowflake.",
						ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
					},
					"credentials": {
						Type:             schema.TypeString,
						Required:         true,
						Description:      "GCP Service Account credentials JSON for RudderStack to use in loading data into your Google Cloud Storage.",
						ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|.+"),
					},
					"storage_integration": {
						Type:             schema.TypeString,
						Required:         true,
						Description:      "Create the cloud storage integration in Snowflake and enter the name of integration.Please refer to this for more details -> https://www.rudderstack.com/docs/destinations/warehouse-destinations/snowflake/#configuring-cloud-storage-integration-with-snowflake",
						ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
					},
				},
			},
		},
		"azure": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			Description:   "",
			ConflictsWith: []string{"config.0.s3", "config.0.gcp"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"container_name": {
						Type:             schema.TypeString,
						Required:         true,
						Description:      "Specify the name of your Azure container where RudderStack will store the data before loading it into Snowflake.",
						ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
					},
					"account_name": {
						Type:             schema.TypeString,
						Required:         true,
						Description:      "Enter the account name for the Azure container.",
						ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
					},
					"account_key": {
						Type:             schema.TypeString,
						Required:         true,
						Description:      "Enter the account key for your Azure container.",
						ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
					},
					"storage_integration": {
						Type:             schema.TypeString,
						Required:         true,
						Description:      "Create the cloud storage integration in Snowflake and enter the name of integration. Please refer to this for more details -> https://www.rudderstack.com/docs/destinations/warehouse-destinations/snowflake/#configuring-cloud-storage-integration-with-snowflake",
						ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
					},
				},
			},
		},
		"prefix": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "If specified, RudderStack will create a folder in the bucket with this prefix and push all the data within that folder.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].*)|^(.{0,100})$"),
		},
	}

	for key, value := range commonSchema {
		schema[key] = value
	}

	c.Destinations.Register("snowflake", c.ConfigMeta{
		APIType:      "SNOWFLAKE",
		Properties:   properties,
		ConfigSchema: schema,
	})
}
