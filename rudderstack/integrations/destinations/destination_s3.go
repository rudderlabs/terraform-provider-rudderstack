package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"web", "android", "ios", "unity", "reactnative", "flutter", "cordova", "amp", "cloud", "warehouse", "shopify"}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)

	properties := []c.ConfigProperty{
		c.Simple("bucketName", "bucket_name"),
		c.Simple("prefix", "prefix", c.SkipZeroValue),
		c.Simple("accessKeyID", "key_based_authentication.0.access_key_id", c.SkipZeroValue),
		c.Simple("accessKey", "key_based_authentication.0.access_key", c.SkipZeroValue),
		c.Discriminator("roleBasedAuth", c.DiscriminatorValues{
			"role_based_authentication": true,
			"key_based_authentication":  false,
		}),
		c.Simple("iamRoleARN", "role_based_authentication.0.i_am_role_arn", c.SkipZeroValue),
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
		"role_based_authentication": {
			Type:         schema.TypeList,
			MaxItems:     1,
			Optional:     true,
			Description:  "Use an AWS IAM Role ARN for authentication instead of access keys.",
			ExactlyOneOf: []string{"config.0.role_based_authentication", "config.0.key_based_authentication"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"i_am_role_arn": {
						Type:             schema.TypeString,
						Required:         true,
						Description:      "Enter the AWS IAM Role ARN to use for authentication.",
						ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
					},
				},
			},
		},
		"key_based_authentication": {
			Type:         schema.TypeList,
			MaxItems:     1,
			Optional:     true,
			Description:  "Use AWS access key and secret for authentication.",
			ExactlyOneOf: []string{"config.0.role_based_authentication", "config.0.key_based_authentication"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
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
		Properties:   properties,
		ConfigSchema: schema,
	})
}
