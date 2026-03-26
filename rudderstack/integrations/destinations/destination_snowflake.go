package destinations

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"web", "android", "ios", "unity", "reactnative", "flutter", "cordova", "amp", "cloud", "cloudSource", "shopify"}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)

	properties := []c.ConfigProperty{
		c.Simple("account", "account"),
		c.Simple("database", "database"),
		c.Simple("warehouse", "warehouse"),
		c.Simple("user", "user"),
		c.Simple("useKeyPairAuth", "use_key_pair_auth", c.SkipZeroValue),
		c.Simple("password", "password", c.SkipZeroValue),
		privateKeyProperty(),
		c.Simple("privateKeyPassphrase", "private_key_passphrase", c.SkipZeroValue),
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
		c.Simple("preferAppend", "prefer_append"),
		c.Simple("skipUsersTable", "skip_users_table"),
		c.Simple("skipTracksTable", "skip_tracks_table"),
		c.Simple("cleanupObjectStorageFiles", "cleanup_object_storage_files", c.SkipZeroValue),
		c.Conditional("bucketName", "s3.0.bucket_name", c.Equals("cloudProvider", "AWS")),
		c.Simple("accessKeyID", "s3.0.access_key_id", c.SkipZeroValue),
		c.Simple("accessKey", "s3.0.access_key", c.SkipZeroValue),
		c.Simple("enableSSE", "s3.0.enable_sse", c.SkipZeroValue),
		c.Simple("iamRoleARN", "s3.0.role_based_authentication.0.i_am_role_arn", c.SkipZeroValue),
		c.Discriminator("roleBasedAuth", c.DiscriminatorValues{
			"s3.0.access_key":                false,
			"s3.0.access_key_id":             false,
			"s3.0.role_based_authentication": true,
		}),
		c.Conditional("storageIntegration", "s3.0.storage_integration", c.Equals("cloudProvider", "AWS")),
		c.Conditional("bucketName", "gcp.0.bucket_name", c.Equals("cloudProvider", "GCP")),
		c.Simple("credentials", "gcp.0.credentials", c.SkipZeroValue),
		c.Conditional("storageIntegration", "gcp.0.storage_integration", c.Equals("cloudProvider", "GCP")),
		c.Simple("containerName", "azure.0.container_name", c.SkipZeroValue),
		c.Simple("accountName", "azure.0.account_name", c.SkipZeroValue),
		c.Simple("accountKey", "azure.0.account_key", c.SkipZeroValue),
		c.Simple("sasToken", "azure.0.sas_token", c.SkipZeroValue),
		c.Simple("useSASTokens", "azure.0.use_sas_tokens", c.SkipZeroValue),
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
		"use_key_pair_auth": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable this setting to use key pair authentication instead of password-based authentication.",
		},
		"password": {
			Type:             schema.TypeString,
			Optional:         true,
			Sensitive:        true,
			Description:      "Password for the user. Required when use_key_pair_auth is false.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|.+"),
			AtLeastOneOf:     []string{"config.0.password", "config.0.private_key"},
			ConflictsWith:    []string{"config.0.private_key"},
		},
		"private_key": {
			Type:             schema.TypeString,
			Optional:         true,
			Sensitive:        true,
			Description:      "Private key for key pair authentication. Required when use_key_pair_auth is true. Accepts both PEM-formatted keys (with BEGIN/END headers) and raw base64-encoded key bodies. Raw keys are automatically wrapped with PEM headers before being sent to the API.",
			ValidateDiagFunc: c.StringMatchesRegexp(".+"),
			DiffSuppressFunc: suppressPEMKeyDiff,
			AtLeastOneOf:     []string{"config.0.password", "config.0.private_key"},
			ConflictsWith:    []string{"config.0.password"},
		},
		"private_key_passphrase": {
			Type:             schema.TypeString,
			Optional:         true,
			Sensitive:        true,
			Description:      "Passphrase for the private key, if the private key is encrypted.",
			ValidateDiagFunc: c.StringMatchesRegexp("^(.{0,100})$"),
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
						ValidateDiagFunc: c.StringMatchesRegexp("^(5|15|30|60|180|360|720|1440)$"),
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
		"prefer_append": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Disable to move from Append to Merge operation. Switching from Append to Merge ensures 100% non-duplicate data, but would increase warehouse operations time significantly.",
		},
		"skip_users_table": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Disable the creation of a Users table. The table stores all unique users, but note that due to merge operations, it can significantly increase warehouse operation time.",
		},
		"skip_tracks_table": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable this toggle to skip sending the event data to the tracks table.",
		},
		"cleanup_object_storage_files": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable for cleanup of object storage files (deletion) after successful sync.",
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
						AtLeastOneOf:     []string{"config.0.s3.0.access_key_id", "config.0.s3.0.role_based_authentication"},
						ConflictsWith:    []string{"config.0.s3.0.role_based_authentication"},
						RequiredWith:     []string{"config.0.s3.0.access_key"},
					},
					"access_key": {
						Type:             schema.TypeString,
						Optional:         true,
						Sensitive:        true,
						Description:      "Enter your AWS secret access key.",
						ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
						ConflictsWith:    []string{"config.0.s3.0.role_based_authentication"},
						RequiredWith:     []string{"config.0.s3.0.access_key_id"},
					},
					"enable_sse": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "Toggle on this setting to enable server-side encryption for your S3 bucket.",
					},
					"role_based_authentication": {
						Type:          schema.TypeList,
						MaxItems:      1,
						Optional:      true,
						Description:   "Use IAM role-based authentication for S3 access.",
						AtLeastOneOf:  []string{"config.0.s3.0.access_key_id", "config.0.s3.0.role_based_authentication"},
						ConflictsWith: []string{"config.0.s3.0.access_key_id", "config.0.s3.0.access_key"},
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"i_am_role_arn": {
									Type:             schema.TypeString,
									Required:         true,
									Description:      "The IAM role ARN to use for authentication.",
									ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
								},
							},
						},
					},
					"storage_integration": {
						Type:             schema.TypeString,
						Optional:         true,
						Description:      "Create the cloud storage integration in Snowflake and enter the name of integration.",
						ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{0,100})$"),
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
						Optional:         true,
						Sensitive:        true,
						Description:      "Enter the account key for your Azure container.",
						ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
						AtLeastOneOf:     []string{"config.0.azure.0.account_key", "config.0.azure.0.sas_token"},
						ConflictsWith:    []string{"config.0.azure.0.sas_token"},
					},
					"sas_token": {
						Type:             schema.TypeString,
						Optional:         true,
						Sensitive:        true,
						Description:      "Enter the SAS token for your Azure Blob Storage.",
						ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.+)$"),
						AtLeastOneOf:     []string{"config.0.azure.0.account_key", "config.0.azure.0.sas_token"},
						ConflictsWith:    []string{"config.0.azure.0.account_key"},
					},
					"use_sas_tokens": {
						Type:        schema.TypeBool,
						Optional:    true,
						Default:     false,
						Description: "Use shared access signature (SAS) tokens to grant limited access to Azure Storage resources.",
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

// privateKeyProperty creates a ConfigProperty for the privateKey field that
// auto-wraps raw base64-encoded key bodies with PEM headers.
func privateKeyProperty() c.ConfigProperty {
	return c.ConfigProperty{
		FromStateFunc: func(config, state string) (string, error) {
			v := gjson.Get(state, "private_key")
			if !v.Exists() || v.String() == "" {
				return config, nil
			}
			return sjson.Set(config, "privateKey", wrapPEMKey(v.String()))
		},
		ToStateFunc: func(state, config string) (string, error) {
			r := gjson.Get(config, "privateKey")
			if r.Exists() {
				return sjson.Set(state, "private_key", r.Value())
			}
			return state, nil
		},
	}
}

// suppressPEMKeyDiff suppresses diffs between PEM-wrapped and raw key representations.
func suppressPEMKeyDiff(_, old, new string, _ *schema.ResourceData) bool {
	return stripPEMHeaders(old) == stripPEMHeaders(new)
}

// wrapPEMKey wraps a raw key body with PEM headers if not already PEM-formatted.
func wrapPEMKey(key string) string {
	if strings.HasPrefix(key, "-----BEGIN") {
		return key
	}
	return "-----BEGIN PRIVATE KEY-----\n" + key + "\n-----END PRIVATE KEY-----"
}

// stripPEMHeaders removes PEM header/footer lines and whitespace, returning just the key body.
func stripPEMHeaders(key string) string {
	s := key
	for _, marker := range []string{
		"-----BEGIN ENCRYPTED PRIVATE KEY-----",
		"-----END ENCRYPTED PRIVATE KEY-----",
		"-----BEGIN PRIVATE KEY-----",
		"-----END PRIVATE KEY-----",
	} {
		s = strings.ReplaceAll(s, marker, "")
	}
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\r", "")
	return strings.TrimSpace(s)
}
