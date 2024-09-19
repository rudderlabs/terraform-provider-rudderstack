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
		c.Simple("accessKeyID", "access_key_id", c.SkipZeroValue),
		c.Simple("accessKey", "access_key", c.SkipZeroValue),
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
