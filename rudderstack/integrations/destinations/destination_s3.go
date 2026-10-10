package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"web", "android", "androidKotlin", "ios", "iosSwift", "unity", "reactnative", "flutter", "cordova", "amp", "cloud", "warehouse", "shopify"}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)

	properties := []c.ConfigProperty{
		c.Simple("bucketName", "bucket_name"),
		c.Simple("prefix", "prefix", c.SkipZeroValue),
		c.Simple("accessKeyID", "access_key_id", c.SkipZeroValue),
		c.Simple("accessKey", "access_key", c.SkipZeroValue),
		c.Simple("iamRoleARN", "role_based_authentication.0.i_am_role_arn", c.SkipZeroValue),
		c.Discriminator("roleBasedAuth", c.DiscriminatorValues{
			"access_key_id":             false,
			"access_key":                false,
			"role_based_authentication": true,
		}),
		c.Simple("enableSSE", "enable_sse", c.SkipZeroValue),
	}

	properties = append(properties, commonProperties...)

	schema := map[string]*schema.Schema{
		"bucket_name": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Enter the name of your S3 bucket.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"prefix": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Enter a prefix which RudderStack associates as the path prefix to all the files stored in your S3 bucket.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"access_key_id": {
			Type:             schema.TypeString,
			Optional:         true,
			Sensitive:        true,
			Description:      "Enter your AWS access key ID.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			ConflictsWith:    []string{"config.0.role_based_authentication"},
		},
		"access_key": {
			Type:             schema.TypeString,
			Optional:         true,
			Sensitive:        true,
			Description:      "Enter your AWS secret access key.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			ConflictsWith:    []string{"config.0.role_based_authentication"},
		},
		"role_based_authentication": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			Description:   "Use IAM role-based authentication instead of access keys.",
			ConflictsWith: []string{"config.0.access_key_id", "config.0.access_key"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"i_am_role_arn": {
						Type:             schema.TypeString,
						Required:         true,
						Description:      "The IAM role ARN to use for authentication.",
						ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
					},
				},
			},
		},
		"enable_sse": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "This setting enables server-side encryption.",
		},
	}

	for key, value := range commonSchema {
		schema[key] = value
	}

	c.Destinations.Register("s3", c.ConfigMeta{
		APIType:      "S3",
		Version:      1,
		Properties:   properties,
		ConfigSchema: schema,
	})
}
