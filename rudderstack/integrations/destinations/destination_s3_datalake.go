package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"web", "android", "ios", "unity", "reactnative", "flutter", "cordova", "amp", "cloud", "cloudSource", "shopify"}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)

	properties := []c.ConfigProperty{
		c.Simple("bucketName", "bucket_name"),
		c.Simple("namespace", "namespace", c.SkipZeroValue),
		c.Simple("prefix", "prefix", c.SkipZeroValue),
		c.Simple("accessKeyID", "access_key_id", c.SkipZeroValue),
		c.Simple("accessKey", "access_key", c.SkipZeroValue),
		c.Simple("iamRoleARN", "role_based_authentication.0.i_am_role_arn", c.SkipZeroValue),
		c.Discriminator("roleBasedAuth", c.DiscriminatorValues{
			"role_based_authentication": true,
			"access_key":                false,
			"access_key_id":             false,
		}),
		c.Simple("enableSSE", "enable_sse", c.SkipZeroValue),
		c.Simple("useGlue", "use_glue"),
		c.Simple("region", "region", c.SkipZeroValue),
		c.Simple("syncFrequency", "sync.0.frequency"),
		c.Simple("syncStartAt", "sync.0.start_at", c.SkipZeroValue),
	}

	properties = append(properties, commonProperties...)

	schema := map[string]*schema.Schema{
		"bucket_name": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "The name of the S3 bucket that will be used to store the data before loading it into the S3 data lake.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"namespace": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "If specified, all the data for the destination will be pushed to `s3://<bucketName>/<prefix>/rudder-datalake/<namespace>`. ",
			ValidateDiagFunc: c.ValidateAll(
				c.StringMatchesRegexp("(^env[.].*)|^(.{0,64})$"),
				c.StringNotMatchesRegexp("^(pg_|PG_|pG_|Pg_)"),
			),
		},
		"prefix": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "If specified, RudderStack creates a folder in the bucket with this prefix and push all the data within that folder.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"access_key_id": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Enter your AWS access key ID.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"access_key": {
			Type:             schema.TypeString,
			Optional:         true,
			Sensitive:        true,
			Description:      "Enter your AWS secret access key.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"role_based_authentication": {
			Type:         schema.TypeList,
			MaxItems:     1,
			Optional:     true,
			Description:  "This option allows you select the arn based authentication.",
			RequiredWith: []string{"config.0.role_based_authentication"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"i_am_role_arn": {
						Type:        schema.TypeString,
						Optional:    true,
						Default:     true,
						Description: "Role Based Authentication",
					},
				},
			},
		},
		"enable_sse": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "This setting enables server-side encryption.",
		},
		"use_glue": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "This setting enables AWS Glue.",
		},
		"region": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Enter your AWS Glue region. For example, for N.Virginia, it would be `us-east-1`.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|.+"),
		},
		"sync": {
			Type:     schema.TypeList,
			MinItems: 1, MaxItems: 1,
			Required:    true,
			Description: "Specify your data sync settings.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"frequency": {
						Type:             schema.TypeString,
						Required:         true,
						Description:      "Specify how often RudderStack should sync the data to your S3 Datalake.",
						ValidateDiagFunc: c.StringMatchesRegexp("^(30|60|180|360|720|1440)$"),
					},
					"start_at": {
						Type:             schema.TypeString,
						Optional:         true,
						Description:      "This optional setting lets you specify the particular time of the day (in UTC) when you want RudderStack to sync the data to the warehouse.",
						ValidateDiagFunc: c.StringMatchesRegexp("^([01][0-9]|2[0-3]):[0-5][0-9]$"),
					},
				},
			},
		},
	}

	for key, value := range commonSchema {
		schema[key] = value
	}

	c.Destinations.Register("s3_datalake", c.ConfigMeta{
		APIType:      "S3_DATALAKE",
		Properties:   properties,
		ConfigSchema: schema,
	})
}
