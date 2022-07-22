package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	c.Destinations.Register("snowflake", c.ConfigMeta{
		APIType: "SNOWFLAKE",
		Properties: []c.ConfigProperty{
			c.Simple("account", "account"),
			c.Simple("database", "database"),
			c.Simple("warehouse", "warehouse"),
			c.Simple("user", "user"),
			c.Simple("password", "password"),
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
		},
		ConfigSchema: map[string]*schema.Schema{
			"account": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Account ID of your Snowflake warehouse. This account ID is part of the Snowflake URL.",
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
							Description:      "Enter the name of your S3 bucket.",
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
						},
						"access_key_id": {
							Type:             schema.TypeString,
							Optional:         true,
							Sensitive:        true,
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
						"enable_sse": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "This setting enables server-side encryption.",
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
							Description:      "Staging GCS Object Storage Bucket Name",
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
						},
						"credentials": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "GCP Service Account credentials JSON for RudderStack to use in loading data into your Google Cloud Storage",
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|.+"),
						},
						"storage_integration": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Storage Integration",
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
							Description:      "Staging Azure Blob Storage Container Name",
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
						},
						"account_name": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Azure Blob Storage Account Name",
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
						},
						"account_key": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Azure Blob Storage Account Key",
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
						},
						"storage_integration": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Storage Integration",
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
						},
					},
				},
			},
			"prefix": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Prefix",
				ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].*)|^(.{0,100})$"),
			},
		},
	})
}
