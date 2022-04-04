package configs

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	c.Destinations.Register("postgres", c.ConfigMeta{
		APIType: "POSTGRES",
		Properties: []c.ConfigProperty{
			c.Simple("host", "host"),
			c.Simple("port", "port"),
			c.Simple("user", "user"),
			c.Simple("password", "password"),
			c.Simple("database", "database"),
			c.Simple("namespace", "namespace"),
			c.Simple("sslMode", "ssl_mode"),
			c.Simple("useRudderStorage", "use_rudder_storage"),
			c.Simple("verifyCAProperties.clientKey", "verify_ca.0.client_key"),
			c.Simple("verifyCAProperties.clientCert", "verify_ca.0.client_cert"),
			c.Simple("verifyCAProperties.serverCA", "verify_ca.0.server_ca"),
			c.Simple("syncFrequency", "sync.0.frequency"),
			c.Simple("syncStartAt", "sync.0.start_at"),
			c.Simple("excludeWindow.excludeWindowStartTime", "sync.0.exclude_window_start_time"),
			c.Simple("excludeWindow.excludeWindowEndTime", "sync.0.exclude_window_end_time"),
			c.Discriminator("bucketProvider", c.DiscriminatorValues{
				"s3":         "S3",
				"gcs":        "GCS",
				"azure_blob": "AZURE_BLOB",
				"minio":      "MINIO",
			}),
			c.Conditional("bucketName", "s3.0.bucket_name", c.Equals("bucketProvider", "S3")),
			c.Conditional("accessKeyID", "s3.0.access_key_id", c.Equals("bucketProvider", "S3")),
			c.Conditional("accessKey", "s3.0.access_key", c.Equals("bucketProvider", "S3")),
			c.Conditional("bucketName", "gcs.0.bucket_name", c.Equals("bucketProvider", "GCS")),
			c.Conditional("credentials", "gcs.0.credentials", c.Equals("bucketProvider", "GCS")),
			c.Conditional("containerName", "azure_blob.0.container_name", c.Equals("bucketProvider", "AZURE_BLOB")),
			c.Conditional("accountName", "azure_blob.0.account_name", c.Equals("bucketProvider", "AZURE_BLOB")),
			c.Conditional("accountKey", "azure_blob.0.account_key", c.Equals("bucketProvider", "AZURE_BLOB")),
			c.Conditional("bucketName", "minio.0.bucket_name", c.Equals("bucketProvider", "MINIO")),
			c.Conditional("endPoint", "minio.0.endpoint", c.Equals("bucketProvider", "MINIO")),
			c.Conditional("accessKeyID", "minio.0.access_key_id", c.Equals("bucketProvider", "MINIO")),
			c.Conditional("secretAccessKey", "minio.0.secret_access_key", c.Equals("bucketProvider", "MINIO")),
			c.Conditional("useSSL", "minio.0.use_ssl", c.Equals("bucketProvider", "MINIO")),
		},
		ConfigSchema: map[string]*schema.Schema{
			"host": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
			},
			"port": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
			},
			"user": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
			},
			"password": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
				Sensitive:        true,
			},
			"database": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
				// ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].*)|^(?!pg_|PG_|pG_|Pg_).*"),
			},
			"ssl_mode": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: c.StringMatchesRegexp("^(disable|require)$"),
			},

			"use_rudder_storage": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"bucket_provider": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sync": {
				Type:     schema.TypeList,
				MinItems: 1, MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"frequency": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: c.StringMatchesRegexp("^(30|60|180|360|720|1440)$"),
						},
						"start_at": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: c.StringMatchesRegexp("^([01][0-9]|2[0-3]):[0-5][0-9]$"),
						},
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
			"verify_ca": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_key": {
							Type:             schema.TypeString,
							Required:         true,
							Sensitive:        true,
							ValidateDiagFunc: c.StringMatchesRegexp("-----BEGIN RSA PRIVATE KEY-----.*-----END CERTIFICATE-----"),
						},
						"client_cert": {
							Type:             schema.TypeString,
							Required:         true,
							Sensitive:        true,
							ValidateDiagFunc: c.StringMatchesRegexp("-----BEGIN RSA PRIVATE KEY-----.*-----END CERTIFICATE-----"),
						},
						"server_ca": {
							Type:             schema.TypeString,
							Required:         true,
							Sensitive:        true,
							ValidateDiagFunc: c.StringMatchesRegexp("-----BEGIN RSA PRIVATE KEY-----.*-----END CERTIFICATE-----"),
						},
					},
				},
			},
			"s3": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ConflictsWith: []string{"config.0.gcs", "config.0.azure_blob", "config.0.minio"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket_name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
						},
						"access_key_id": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
						},
						"access_key": {
							Type:             schema.TypeString,
							Required:         true,
							Sensitive:        true,
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
						},
					},
				},
			},
			"gcs": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ConflictsWith: []string{"config.0.s3", "config.0.azure_blob", "config.0.minio"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket_name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
						},
						"credentials": {
							Type:             schema.TypeString,
							Required:         true,
							Sensitive:        true,
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|.+"),
						},
					},
				},
			},
			"azure_blob": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ConflictsWith: []string{"config.0.s3", "config.0.gcs", "config.0.minio"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"container_name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
						},
						"account_name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
						},
						"account_key": {
							Type:             schema.TypeString,
							Required:         true,
							Sensitive:        true,
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
						},
					},
				},
			},
			"minio": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ConflictsWith: []string{"config.0.s3", "config.0.gcs", "config.0.azure_blob"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket_name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
						},
						"endpoint": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
						},
						"access_key_id": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
						},
						"secret_access_key": {
							Type:             schema.TypeString,
							Required:         true,
							Sensitive:        true,
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
						},
						"use_ssl": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
		},
		TestConfigs: []c.TestConfig{
			{
				TerraformCreate: `
					host = "example.com"
					port = "5432"
					user = "example-user"
					password = "example-password"
					namespace = "example-namespace"
					ssl_mode = "disable"
					database = "example-database"
					use_rudder_storage = true
					s3 {
						bucket_name = "some-bucket-name"
						access_key_id = "some-access-key-id"
						access_key = "some-access-key"
					}
					`,
				APICreate: `{
					"host": "example.com",
					"port": "5432",
					"user": "example-user",
					"password": "example-password",
					"database": "example-database",
					"namespace": "example-namespace",
					"sslMode": "disable",
					"useRudderStorage": true,
					"bucketProvider": "S3",
					"bucketName": "some-bucket-name",
					"accessKeyID": "some-access-key-id",
					"accessKey": "some-access-key"
				}`,
				TerraformUpdate: `
					host = "example.com"
					port = "5432"
					user = "example-user-updated"
					password = "example-password"
					namespace = "example-namespace"
					ssl_mode = "disable"
					database = "example-database"
					use_rudder_storage = true
					gcs {
						bucket_name = "some-bucket-name"
						credentials = "some-credentials"
					}
					`,
				APIUpdate: `{
					"host": "example.com",
					"port": "5432",
					"user": "example-user-updated",
					"password": "example-password",
					"database": "example-database",
					"namespace": "example-namespace",
					"sslMode": "disable",
					"useRudderStorage": true,
					"bucketProvider": "GCS",
					"bucketName": "some-bucket-name",
					"credentials": "some-credentials"
				}`,
			},
		},
	})
}