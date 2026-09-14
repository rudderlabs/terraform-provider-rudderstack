package destinations

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"web", "android", "androidKotlin", "ios", "iosSwift", "unity", "reactnative", "flutter", "cordova", "amp", "cloud", "cloudSource", "shopify"}
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
		c.Simple("syncStartAt", "sync_start_at"),
		c.Simple("excludeWindow.excludeWindowStartTime", "exclude_window.0.exclude_window_start_time"),
		c.Simple("excludeWindow.excludeWindowEndTime", "exclude_window.0.exclude_window_end_time"),
		c.Simple("jsonPaths", "json_paths"),
		c.Simple("useRudderStorage", "use_rudder_storage"),
		c.Discriminator("bucketProvider", c.DiscriminatorValues{
			"s3":    "S3",
			"gcp":   "GCS",
			"azure": "AZURE_BLOB",
			"minio": "MINIO",
		}),
		c.Conditional("bucketName", "s3.0.bucket_name", c.Equals("bucketProvider", "S3"), c.SkipZeroValue),
		c.Conditional("bucketName", "gcp.0.bucket_name", c.Equals("bucketProvider", "GCS"), c.SkipZeroValue),
		c.Conditional("bucketName", "minio.0.bucket_name", c.Equals("bucketProvider", "MINIO"), c.SkipZeroValue),
		c.Simple("clientKey", "client_key", c.SkipZeroValue),
		c.Simple("clientCert", "client_cert", c.SkipZeroValue),
		c.Simple("serverCA", "server_ca", c.SkipZeroValue),
		c.Conditional("iamRoleARN", "s3.0.role_based_authentication.0.i_am_role_arn", c.Equals("bucketProvider", "S3"), c.SkipZeroValue),
		c.Discriminator("roleBasedAuth", c.DiscriminatorValues{
			"s3.0.access_key":                false,
			"s3.0.access_key_id":             false,
			"s3.0.role_based_authentication": true,
		}),
		c.Conditional("accessKeyID", "s3.0.access_key_id", c.Equals("bucketProvider", "S3"), c.SkipZeroValue),
		c.Conditional("accessKey", "s3.0.access_key", c.Equals("bucketProvider", "S3"), c.SkipZeroValue),
		c.Conditional("accountName", "azure.0.account_name", c.Equals("bucketProvider", "AZURE_BLOB"), c.SkipZeroValue),
		c.Conditional("accountKey", "azure.0.account_key", c.Equals("bucketProvider", "AZURE_BLOB"), c.SkipZeroValue),
		c.Conditional("sasToken", "azure.0.sas_token", c.Equals("bucketProvider", "AZURE_BLOB"), c.SkipZeroValue),
		c.Conditional("useSASTokens", "azure.0.use_sas_tokens", c.Equals("bucketProvider", "AZURE_BLOB")),
		c.Conditional("containerName", "azure.0.container_name", c.Equals("bucketProvider", "AZURE_BLOB"), c.SkipZeroValue),
		c.Conditional("credentials", "gcp.0.credentials", c.Equals("bucketProvider", "GCS"), c.SkipZeroValue),
		c.Conditional("endPoint", "minio.0.end_point", c.Equals("bucketProvider", "MINIO"), c.SkipZeroValue),
		c.Conditional("accessKeyID", "minio.0.access_key_id", c.Equals("bucketProvider", "MINIO"), c.SkipZeroValue),
		c.Conditional("secretAccessKey", "minio.0.secret_access_key", c.Equals("bucketProvider", "MINIO"), c.SkipZeroValue),
		c.Conditional("useSSL", "minio.0.use_ssl", c.Equals("bucketProvider", "MINIO")),
	}

	properties = append(properties, commonProperties...)

	schema := map[string]*schema.Schema{
		"host": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Enter the host name of your PostgreSQL database.",
		},
		"database": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Enter the name of your PostgreSQL database.",
		},
		"user": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Enter the username of your PostgreSQL database.",
		},
		"password": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "Enter the password of your PostgreSQL database.",
		},
		"port": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "5432",
			Description: "Enter the port number of your PostgreSQL database.",
		},
		"namespace": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Enter the namespace of your PostgreSQL database.",
		},
		"ssl_mode": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "disable",
			Description:      "Enter the SSL mode of your PostgreSQL database.",
			ValidateDiagFunc: c.StringMatchesRegexp("^(disable|require|verify-ca)$"),
		},
		"sync_frequency": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "30",
			Description:      "Enter the frequency at which the data should be synced from your PostgreSQL database.",
			ValidateDiagFunc: c.StringMatchesRegexp("^(30|60|180|360|720|1440)$"),
		},
		"client_key": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Enter your Client Key Pem File",
		},
		"client_cert": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Enter your Client Cert Pem File",
		},
		"server_ca": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Enter your Server CA Pem File",
		},
		"use_rudder_storage": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "Enable this setting to use RudderStack's data warehouse to store the data from your PostgreSQL database.",
		},
		"s3": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			Description:   "Configure S3 object storage for staging PostgreSQL warehouse data.",
			ConflictsWith: []string{"config.0.gcp", "config.0.azure", "config.0.minio"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"bucket_name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Specify the name of your S3 bucket where RudderStack will store the data before loading it into PostgreSQL.",
					},
					"access_key_id": {
						Type:          schema.TypeString,
						Optional:      true,
						Sensitive:     true,
						Description:   "Enter your AWS access key ID obtained from the AWS console.",
						AtLeastOneOf:  []string{"config.0.s3.0.access_key_id", "config.0.s3.0.role_based_authentication"},
						ConflictsWith: []string{"config.0.s3.0.role_based_authentication"},
						RequiredWith:  []string{"config.0.s3.0.access_key"},
					},
					"access_key": {
						Type:          schema.TypeString,
						Optional:      true,
						Sensitive:     true,
						Description:   "Enter your AWS secret access key.",
						ConflictsWith: []string{"config.0.s3.0.role_based_authentication"},
						RequiredWith:  []string{"config.0.s3.0.access_key_id"},
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
									Type:        schema.TypeString,
									Required:    true,
									Description: "The IAM role ARN to use for authentication.",
								},
							},
						},
					},
				},
			},
		},
		"gcp": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			Description:   "Configure Google Cloud Storage for staging PostgreSQL warehouse data.",
			ConflictsWith: []string{"config.0.s3", "config.0.azure", "config.0.minio"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"bucket_name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Specify the name of your GCS bucket where RudderStack will store the data before loading it into PostgreSQL.",
					},
					"credentials": {
						Type:        schema.TypeString,
						Required:    true,
						Sensitive:   true,
						Description: "GCP Service Account credentials JSON for RudderStack to use in loading data into your Google Cloud Storage.",
					},
				},
			},
		},
		"azure": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			Description:   "Configure Azure Blob Storage for staging PostgreSQL warehouse data.",
			ConflictsWith: []string{"config.0.s3", "config.0.gcp", "config.0.minio"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"container_name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Specify the name of your Azure container where RudderStack will store the data before loading it into PostgreSQL.",
					},
					"account_name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Enter the account name for the Azure container.",
					},
					"account_key": {
						Type:          schema.TypeString,
						Optional:      true,
						Sensitive:     true,
						Description:   "Enter the account key for your Azure container. Required when `use_sas_tokens` is `false`.",
						ConflictsWith: []string{"config.0.azure.0.sas_token"},
					},
					"sas_token": {
						Type:          schema.TypeString,
						Optional:      true,
						Sensitive:     true,
						Description:   "Enter the SAS token for your Azure container. Required when `use_sas_tokens` is `true`.",
						ConflictsWith: []string{"config.0.azure.0.account_key"},
					},
					"use_sas_tokens": {
						Type:        schema.TypeBool,
						Optional:    true,
						Default:     false,
						Description: "Use a shared access signature (SAS) token instead of an account key.",
					},
				},
			},
		},
		"minio": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			Description:   "Configure MinIO object storage for staging PostgreSQL warehouse data.",
			ConflictsWith: []string{"config.0.s3", "config.0.gcp", "config.0.azure"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"bucket_name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Specify the name of your MinIO bucket where RudderStack will store the data before loading it into PostgreSQL.",
					},
					"end_point": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Enter the MinIO endpoint.",
					},
					"access_key_id": {
						Type:        schema.TypeString,
						Required:    true,
						Sensitive:   true,
						Description: "Enter your MinIO access key ID.",
					},
					"secret_access_key": {
						Type:        schema.TypeString,
						Required:    true,
						Sensitive:   true,
						Description: "Enter your MinIO secret access key.",
					},
					"use_ssl": {
						Type:        schema.TypeBool,
						Optional:    true,
						Default:     true,
						Description: "Use SSL for the MinIO connection.",
					},
				},
			},
		},
	}

	for key, value := range commonSchema {
		schema[key] = value
	}

	c.Destinations.Register("postgres", c.ConfigMeta{
		APIType:       "POSTGRES",
		Version:       1,
		Properties:    properties,
		ConfigSchema:  schema,
		CustomizeDiff: customizePostgresDiff,
	})
}

func customizePostgresDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	useRudderStorage, _ := d.Get("config.0.use_rudder_storage").(bool)
	storageBlocks := []string{"s3", "gcp", "azure", "minio"}

	configuredBlocks := []string{}
	for _, block := range storageBlocks {
		if postgresStorageBlockConfigured(d, block) {
			configuredBlocks = append(configuredBlocks, block)
		}
	}

	if useRudderStorage {
		if len(configuredBlocks) > 0 {
			return fmt.Errorf("config.0.%s cannot be configured when use_rudder_storage is true", configuredBlocks[0])
		}
		return nil
	}

	if len(configuredBlocks) != 1 {
		return fmt.Errorf("exactly one of config.0.s3, config.0.gcp, config.0.azure, or config.0.minio must be configured when use_rudder_storage is false")
	}

	if configuredBlocks[0] == "azure" {
		return validatePostgresAzureStorage(d)
	}

	return nil
}

func postgresStorageBlockConfigured(d *schema.ResourceDiff, block string) bool {
	v, ok := d.Get("config.0." + block).([]interface{})
	return ok && len(v) > 0
}

func validatePostgresAzureStorage(d *schema.ResourceDiff) error {
	useSASTokens, _ := d.Get("config.0.azure.0.use_sas_tokens").(bool)
	accountKey := d.Get("config.0.azure.0.account_key").(string)
	sasToken := d.Get("config.0.azure.0.sas_token").(string)

	if useSASTokens {
		if accountKey != "" {
			return fmt.Errorf("config.0.azure.0.account_key cannot be configured when use_sas_tokens is true")
		}
		if sasToken == "" {
			return fmt.Errorf("config.0.azure.0.sas_token is required when use_sas_tokens is true")
		}
		return nil
	}

	if sasToken != "" {
		return fmt.Errorf("config.0.azure.0.sas_token cannot be configured when use_sas_tokens is false")
	}
	if accountKey == "" {
		return fmt.Errorf("config.0.azure.0.account_key is required when use_sas_tokens is false")
	}

	return nil
}
