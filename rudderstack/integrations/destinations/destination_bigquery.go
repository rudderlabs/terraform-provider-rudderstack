package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"web", "android", "ios", "unity", "reactnative", "flutter", "cordova", "amp", "cloud", "cloudSource", "shopify"}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)

	properties := []c.ConfigProperty{
		c.Simple("project", "project"),
		c.Simple("location", "location", c.SkipZeroValue),
		c.Simple("bucketName", "bucket_name"),
		c.Simple("prefix", "prefix", c.SkipZeroValue),
		c.Simple("namespace", "namespace", c.SkipZeroValue),
		c.Simple("credentials", "credentials"),
		c.Simple("skipTracksTable", "skip_tracks_table"),
		c.Simple("skipViews", "skip_views"),
		c.Simple("skipUsersTable", "skip_users_table"),
		c.Simple("partitionColumn", "partition_column"),
		c.Simple("partitionType", "partition_type"),
		c.Simple("syncFrequency", "sync.0.frequency"),
		c.Simple("syncStartAt", "sync.0.start_at", c.SkipZeroValue),
		c.Simple("excludeWindow.excludeWindowStartTime", "sync.0.exclude_window_start_time", c.SkipZeroValue),
		c.Simple("excludeWindow.excludeWindowEndTime", "sync.0.exclude_window_end_time", c.SkipZeroValue),
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
		c.Simple("connectionMode.cloudSource", "connection_mode.0.cloud_source", c.SkipZeroValue),
		c.Simple("jsonPaths", "json_paths", c.SkipZeroValue),
		c.Simple("cleanupObjectStorageFiles", "cleanup_object_storage_files"),
	}

	properties = append(properties, commonProperties...)

	schema := map[string]*schema.Schema{
		"project": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Enter your GCP project ID where the BigQuery database is located.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
		},
		"location": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Enter the GCP region of your project dataset.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{0,100})$"),
		},
		"bucket_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Enter the name of your staging storage bucket.",
			ValidateDiagFunc: c.ValidateAll(
				// Allow env variables or validate bucket name format
				c.StringMatchesRegexp("(^env[.].+)|^[a-z0-9][a-z0-9-._]{1,61}[a-z0-9]$"),
				// Reject bucket names starting with "goog"
				c.StringNotMatchesRegexp("^goog"),
				// Reject bucket names containing "google"
				c.StringNotMatchesRegexp("google"),
				// Reject bucket names that look like IP addresses
				c.StringNotMatchesRegexp("^\\d+\\.\\d+\\.\\d+\\.\\d+$"),
				// Reject bucket names with consecutive dots
				c.StringNotMatchesRegexp("\\.\\."),
			),
		},
		"prefix": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "If specified, RudderStack creates a folder in the bucket with this prefix and loads all the data in it.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].*)|^(.{0,100})$"),
		},
		"namespace": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Enter the schema name where RudderStack will create all the tables. If not specified, RudderStack will set this to the source name by default.",
			ValidateDiagFunc: c.ValidateAll(
				// Allow env variables or validate max length
				c.StringMatchesRegexp("(^env[.].*)|^.{0,64}$"),
				// Reject names starting with pg_ (case insensitive)
				c.StringNotMatchesRegexp("(?i)^pg_"),
			),
		},
		"credentials": {
			Type:             schema.TypeString,
			Required:         true,
			Sensitive:        true,
			Description:      "Enter your GCP service account credentials JSON.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|.+"),
		},
		"sync": {
			Type:     schema.TypeList,
			MinItems: 1, MaxItems: 1,
			Optional:    true,
			Description: "Enter the sync settings for the following fields:",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"frequency": {
						Type:             schema.TypeString,
						Optional:         true,
						Description:      "Specify how often RudderStack should sync the data to your BigQuery dataset.",
						ValidateDiagFunc: c.StringMatchesRegexp("^(5|15|30|60|180|360|720|1440)$"),
					},
					"start_at": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Specify the particular time of the day (in UTC) when you want RudderStack to sync the data to BigQuery.",
					},
					"exclude_window_start_time": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Set a time window when RudderStack will not sync the data to your database.",
					},
					"exclude_window_end_time": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Specify the end time for the exclusion window.",
					},
				},
			},
		},
		"skip_tracks_table": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "If enabled, RudderStack will skip sending the event data to the tracks table.",
		},
		"skip_views": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "If enabled, RudderStack will skip creating views.",
		},
		"skip_users_table": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "If enabled, RudderStack will skip sending the event data to the users table.",
		},
		"partition_column": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "_PARTITIONTIME",
			Description:      "Column to use for partitioning",
			ValidateDiagFunc: c.StringMatchesRegexp("^(\\_PARTITIONTIME|loaded_at|received_at|timestamp|sent_at|original_timestamp)$"),
		},
		"partition_type": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "day",
			Description:      "Partition type",
			ValidateDiagFunc: c.StringMatchesRegexp("^(hour|day)$"),
		},
		"connection_mode": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Use this setting to set how you want to route events from your source to destination..",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"web": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$"),
					},
					"ios": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$"),
					},
					"android": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$"),
					},
					"react_native": {
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
				},
			},
		},
		"cleanup_object_storage_files": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable for cleanup of object storage files (deletion) after successful sync",
		},
		"json_paths": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Specify required JSON properties in dot notation separated by commas.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].*)|.*$"),
		},
	}

	for key, value := range commonSchema {
		schema[key] = value
	}

	c.Destinations.Register("bigquery", c.ConfigMeta{
		APIType:      "BQ",
		Properties:   properties,
		ConfigSchema: schema,
	})
}
