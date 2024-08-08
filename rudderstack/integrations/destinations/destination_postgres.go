package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	c.Destinations.Register("postgres", c.ConfigMeta{
		APIType: "POSTGRES",
		Properties: []c.ConfigProperty{
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
			c.Simple("useRudderStorage", "use_rudder_storage", c.SkipZeroValue), // boolean
			c.Simple("bucketProvider", "bucket_provider"),
			c.ArrayWithStrings("oneTrustCookieCategories", "oneTrustCookieCategory", "onetrust_cookie_categories"),
			c.Simple("bucketName", "bucket_name"),
			c.Simple("clientKey", "client_key", c.SkipZeroValue),
			c.Simple("clientCert", "client_cert", c.SkipZeroValue),
			c.Simple("serverCA", "server_ca", c.SkipZeroValue),
			c.Simple("roleBasedAuth", "role_based_auth", c.SkipZeroValue), // boolean
			c.Simple("iamRoleARN", "iam_role_arn", c.SkipZeroValue),
			c.Simple("accessKeyID", "access_key_id", c.SkipZeroValue),
			c.Simple("accessKey", "access_key", c.SkipZeroValue),
			c.Simple("accountName", "account_name", c.SkipZeroValue),
			c.Simple("accountKey", "account_key", c.SkipZeroValue),
			c.Simple("sasToken", "sas_token", c.SkipZeroValue),
			c.Simple("useSASTokens", "use_sas_tokens", c.SkipZeroValue), // boolean
			c.Simple("credentials", "credentials", c.SkipZeroValue),
			c.Simple("endPoint", "end_point", c.SkipZeroValue),
			c.Simple("accessKeyID", "access_key_id", c.SkipZeroValue),
			c.Simple("secretAccessKey", "secret_access_key", c.SkipZeroValue),
			c.Simple("useSSL", "use_ssl", c.SkipZeroValue), // boolean
		},
		ConfigSchema: map[string]*schema.Schema{
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
		}})
}
