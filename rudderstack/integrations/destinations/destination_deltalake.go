package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"web", "android", "androidKotlin", "ios", "iosSwift", "unity", "reactnative", "flutter", "cordova", "amp", "cloud", "cloudSource", "shopify"}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)

	properties := []c.ConfigProperty{
		c.Simple("host", "host"),
		c.Simple("port", "port"),
		c.Simple("path", "path"),
		c.Simple("useOauth", "use_oauth"),
		c.Simple("token", "token", c.SkipZeroValue),
		c.Simple("oauthClientID", "oauth_client_id", c.SkipZeroValue),
		c.Simple("oauthClientSecret", "oauth_client_secret", c.SkipZeroValue),
		c.Simple("catalog", "catalog", c.SkipZeroValue),
		c.Simple("namespace", "namespace", c.SkipZeroValue),
		c.Simple("syncFrequency", "sync.0.frequency"),
		c.Simple("syncStartAt", "sync.0.start_at", c.SkipZeroValue),
		c.Simple("excludeWindow.excludeWindowStartTime", "sync.0.exclude_window_start_time", c.SkipZeroValue),
		c.Simple("excludeWindow.excludeWindowEndTime", "sync.0.exclude_window_end_time", c.SkipZeroValue),
		c.Simple("skipTracksTable", "skip_tracks_table"),
		c.Simple("skipUsersTable", "skip_users_table"),
		c.Simple("preferAppend", "prefer_append"),
		c.Simple("underscoreDivideNumbers", "underscore_divide_numbers"),
		c.Simple("allowUsersContextTraits", "allow_users_context_traits"),
		c.Simple("enableExternalLocation", "enable_external_location"),
		c.Simple("externalLocation", "external_location", c.SkipZeroValue),
		c.Simple("useRudderStorage", "use_rudder_storage"),
		c.Simple("cleanupObjectStorageFiles", "cleanup_object_storage_files"),
		c.Simple("bucketProvider", "bucket_provider"),
		c.Simple("prefix", "prefix", c.SkipZeroValue),
		// bucketName is shared by the S3 and GCS staging blocks, so it maps
		// conditionally on bucketProvider when read back from the API.
		c.Conditional("bucketName", "s3.0.bucket_name", c.Equals("bucketProvider", "S3")),
		c.Simple("enableSSE", "s3.0.enable_sse"),
		c.Simple("useSTSTokens", "s3.0.use_sts_tokens"),
		c.Simple("accessKeyID", "s3.0.access_key_id", c.SkipZeroValue),
		c.Simple("accessKey", "s3.0.access_key", c.SkipZeroValue),
		c.Simple("iamRoleARN", "s3.0.role_based_authentication.0.i_am_role_arn", c.SkipZeroValue),
		c.Discriminator("roleBasedAuth", c.DiscriminatorValues{
			"s3.0.access_key":                false,
			"s3.0.access_key_id":             false,
			"s3.0.role_based_authentication": true,
		}),
		c.Conditional("bucketName", "gcs.0.bucket_name", c.Equals("bucketProvider", "GCS")),
		c.Simple("credentials", "gcs.0.credentials", c.SkipZeroValue),
		c.Simple("containerName", "azure_blob.0.container_name", c.SkipZeroValue),
		c.Simple("accountName", "azure_blob.0.account_name", c.SkipZeroValue),
		c.Simple("accountKey", "azure_blob.0.account_key", c.SkipZeroValue),
		c.Simple("useSASTokens", "azure_blob.0.use_sas_tokens"),
		c.Simple("sasToken", "azure_blob.0.sas_token", c.SkipZeroValue),
		c.Simple("enableHierarchicalNamespace", "azure_blob.0.enable_hierarchical_namespace"),
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
		"host": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Server hostname of your Databricks SQL warehouse, without the protocol. For example, `abc-12345de-6789.cloud.databricks.com`.",
			// The destination definition spells the ngrok exclusion as a
			// negative lookahead, which Go's RE2 engine does not support, so it
			// is a separate check here.
			ValidateDiagFunc: c.ValidateAll(
				c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
				c.StringNotMatchesRegexp("\\.ngrok\\.io"),
			),
		},
		"port": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Port of your Databricks SQL warehouse. Usually `443`.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
		},
		"path": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "HTTP path of your Databricks SQL warehouse. For example, `/sql/1.0/endpoints/<endpointId>`.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
		},
		"use_oauth": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Authenticate with a Databricks machine-to-machine OAuth service principal instead of a personal access token. Set `oauth_client_id` and `oauth_client_secret` when enabled.",
		},
		"token": {
			Type:             schema.TypeString,
			Optional:         true,
			Sensitive:        true,
			Description:      "Databricks personal access token. Required when `use_oauth` is false.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
			ConflictsWith:    []string{"config.0.oauth_client_id", "config.0.oauth_client_secret"},
		},
		"oauth_client_id": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Client ID of the Databricks OAuth service principal. Required when `use_oauth` is true.",
			ValidateDiagFunc: c.StringMatchesRegexp("^(.{1,100})$"),
			RequiredWith:     []string{"config.0.oauth_client_secret"},
		},
		"oauth_client_secret": {
			Type:             schema.TypeString,
			Optional:         true,
			Sensitive:        true,
			Description:      "Client secret of the Databricks OAuth service principal. Required when `use_oauth` is true.",
			ValidateDiagFunc: c.StringMatchesRegexp("^(.{1,100})$"),
			RequiredWith:     []string{"config.0.oauth_client_id"},
		},
		"catalog": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Unity Catalog name where the data assets are organised. Leave empty to use the `hive_metastore` catalog.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].*)|^(.*)$"),
		},
		"namespace": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Schema name within the catalog where RudderStack creates the tables.",
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
						Description:      "Specify how often RudderStack should sync the data to your Delta Lake, in minutes.",
						ValidateDiagFunc: c.StringMatchesRegexp("^(5|10|15|30|60|180|360|720|1440)$"),
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
						Description:      "This optional setting lets you set a time window when RudderStack will not sync the data to your warehouse.",
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
		"skip_tracks_table": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable this toggle to skip sending the event data to the tracks table.",
		},
		"skip_users_table": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Disable the creation of a users table.",
		},
		"prefer_append": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Enable append operation for warehouse syncs. Disable to use merge, which rules out duplicate rows but makes syncs considerably slower.",
		},
		"underscore_divide_numbers": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
			// underscore_divide_numbers and allow_users_context_traits are
			// declared "rs-immutable" in rudder-integrations-config. They are
			// one-time flags distinguishing old destination generations from
			// new: flipping one rewrites the column set emitted into a
			// warehouse that already holds data written under the old
			// convention. The control plane refuses the update outright, so an
			// in-place change can never apply. ForceNew makes Terraform replace
			// the destination instead of planning something the API rejects.
			ForceNew:    true,
			Description: "Use underscores to split numeric segments in column names. Changing this forces a new destination to be created, because the control plane treats the field as immutable.",
		},
		"allow_users_context_traits": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
			// Immutable upstream — see the note on underscore_divide_numbers.
			ForceNew:    true,
			Description: "Promote `context.traits.*` fields to top-level traits. Changing this forces a new destination to be created, because the control plane treats the field as immutable.",
		},
		"enable_external_location": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
			// Boolean, not a string: rudder-server reads this key with
			// GetBoolDestinationConfig and the destination definition declares
			// it as `"type": "boolean"`. The definition's conditional branch
			// that requires externalLocation compares against the string
			// "true", so that branch never fires for a boolean and the
			// requirement is not enforced server-side. Set external_location
			// yourself whenever this is enabled.
			Description: "Create the Delta tables in an external location instead of the Databricks-managed one. Set `external_location` when enabled.",
		},
		"external_location": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Path of the external location for the Delta tables. Applies only to tables RudderStack creates after the change; existing tables are left as they are.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].*)|.*"),
			RequiredWith:     []string{"config.0.enable_external_location"},
		},
		"use_rudder_storage": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable this setting to use RudderStack-managed buckets for object storage. Only available on RudderStack-hosted data planes.",
		},
		"cleanup_object_storage_files": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Delete the staged files from object storage after a successful sync. Applies only when `use_rudder_storage` is false.",
		},
		"bucket_provider": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Object storage provider used to stage the data before loading it into Delta Lake. One of `S3`, `GCS` or `AZURE_BLOB`. Must match the staging block you declare below.",
			ValidateDiagFunc: c.StringMatchesRegexp("^(S3|GCS|AZURE_BLOB)$"),
		},
		"prefix": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "If specified, RudderStack will create a folder in the bucket with this prefix and push all the data within that folder.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].*)|^(.{0,100})$"),
		},
		"s3": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			Description:   "S3 staging settings. Declare this block when `bucket_provider` is `S3` and `use_rudder_storage` is false.",
			ConflictsWith: []string{"config.0.gcs", "config.0.azure_blob"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"bucket_name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Name of the S3 bucket where RudderStack stores the data before loading it into Delta Lake.",
						// The S3 bucket-name rule in the destination definition
						// uses negative lookaheads (no `xn--` prefix, no `..`,
						// not an IP address), unsupported by Go's RE2 engine, so
						// they are separate checks here.
						ValidateDiagFunc: c.ValidateAll(
							c.StringMatchesRegexp("(^env[.].*)|^[a-z0-9][a-z0-9-.]{1,61}[a-z0-9]$"),
							c.StringNotMatchesRegexp("^xn--"),
							c.StringNotMatchesRegexp("\\.\\."),
							c.StringNotMatchesRegexp("^[0-9]+\\.[0-9]+\\.[0-9]+\\.[0-9]+$"),
						),
					},
					"enable_sse": {
						Type:        schema.TypeBool,
						Optional:    true,
						Default:     false,
						Description: "Enable server-side encryption for your S3 bucket.",
					},
					"use_sts_tokens": {
						Type:        schema.TypeBool,
						Optional:    true,
						Default:     false,
						Description: "Use STS tokens to copy the staging files. Supported only with the Databricks S3A client. Set the access keys or the role-based authentication block when enabled.",
					},
					"access_key_id": {
						Type:             schema.TypeString,
						Optional:         true,
						Sensitive:        true,
						Description:      "Enter your AWS access key ID obtained from the AWS console.",
						ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{0,100})$"),
						ConflictsWith:    []string{"config.0.s3.0.role_based_authentication"},
						RequiredWith:     []string{"config.0.s3.0.access_key"},
					},
					"access_key": {
						Type:             schema.TypeString,
						Optional:         true,
						Sensitive:        true,
						Description:      "Enter your AWS secret access key.",
						ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{0,100})$"),
						ConflictsWith:    []string{"config.0.s3.0.role_based_authentication"},
						RequiredWith:     []string{"config.0.s3.0.access_key_id"},
					},
					"role_based_authentication": {
						Type:          schema.TypeList,
						MaxItems:      1,
						Optional:      true,
						Description:   "Use IAM role-based authentication for S3 access instead of access keys.",
						ConflictsWith: []string{"config.0.s3.0.access_key_id", "config.0.s3.0.access_key"},
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"i_am_role_arn": {
									Type:             schema.TypeString,
									Required:         true,
									Description:      "The IAM role ARN to use for authentication.",
									ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{0,100})$"),
								},
							},
						},
					},
				},
			},
		},
		"gcs": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			Description:   "Google Cloud Storage staging settings. Declare this block when `bucket_provider` is `GCS` and `use_rudder_storage` is false.",
			ConflictsWith: []string{"config.0.s3", "config.0.azure_blob"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"bucket_name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Name of the GCS bucket where RudderStack stores the data before loading it into Delta Lake.",
						// Same lookahead story as the S3 bucket name: GCS
						// forbids a `goog` prefix, the substring `google` and
						// `..` anywhere in the name.
						ValidateDiagFunc: c.ValidateAll(
							c.StringMatchesRegexp("(^env[.].*)|^[a-z0-9][a-z0-9-._]{1,61}[a-z0-9]$"),
							c.StringNotMatchesRegexp("^goog"),
							c.StringNotMatchesRegexp("google"),
							c.StringNotMatchesRegexp("\\.\\."),
						),
					},
					"credentials": {
						Type:             schema.TypeString,
						Required:         true,
						Sensitive:        true,
						Description:      "GCP service account credentials JSON for RudderStack to use when loading data into your Google Cloud Storage bucket.",
						ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|.+"),
					},
				},
			},
		},
		"azure_blob": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			Description:   "Azure Blob Storage staging settings. Declare this block when `bucket_provider` is `AZURE_BLOB` and `use_rudder_storage` is false.",
			ConflictsWith: []string{"config.0.s3", "config.0.gcs"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"container_name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Name of the Azure container where RudderStack stores the data before loading it into Delta Lake.",
						// The definition expresses the 3-63 character bound as a
						// lookahead, which RE2 rejects, so the length is checked
						// on its own here.
						ValidateDiagFunc: c.ValidateAll(
							c.StringMatchesRegexp("(^env[.].*)|^[a-z0-9]+(-[a-z0-9]+)*$"),
							c.StringMatchesRegexp("(^env[.].*)|^.{3,63}$"),
						),
					},
					"account_name": {
						Type:             schema.TypeString,
						Required:         true,
						Sensitive:        true,
						Description:      "Enter the account name for the Azure container.",
						ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].*)|^(.{1,100})$"),
					},
					"account_key": {
						Type:             schema.TypeString,
						Optional:         true,
						Sensitive:        true,
						Description:      "Enter the account key for your Azure container. Required unless `use_sas_tokens` is enabled.",
						ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
						ConflictsWith:    []string{"config.0.azure_blob.0.sas_token"},
					},
					"use_sas_tokens": {
						Type:        schema.TypeBool,
						Optional:    true,
						Default:     false,
						Description: "Grant limited access to the Azure storage resources with a shared access signature instead of the account key.",
					},
					"sas_token": {
						Type:             schema.TypeString,
						Optional:         true,
						Sensitive:        true,
						Description:      "Shared access signature token for your Azure container. Required when `use_sas_tokens` is enabled.",
						ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.+)$"),
						ConflictsWith:    []string{"config.0.azure_blob.0.account_key"},
						RequiredWith:     []string{"config.0.azure_blob.0.use_sas_tokens"},
					},
					"enable_hierarchical_namespace": {
						Type:        schema.TypeBool,
						Optional:    true,
						Default:     false,
						Description: "Enable this setting if the storage account has a hierarchical namespace (Azure Data Lake Storage Gen2).",
					},
				},
			},
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

	c.Destinations.Register("deltalake", c.ConfigMeta{
		APIType:      "DELTALAKE",
		Version:      1,
		Properties:   properties,
		ConfigSchema: schema,
	})
}
