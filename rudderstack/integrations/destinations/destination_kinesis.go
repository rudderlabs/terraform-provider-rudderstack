package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"web", "android", "ios", "unity", "reactnative", "flutter", "cordova", "amp", "cloud", "warehouse", "shopify"}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)

	properties := []c.ConfigProperty{
		c.Simple("region", "region"),
		c.Simple("stream", "stream"),
		c.Simple("accessKeyID", "key_based_authentication.0.access_key_id", c.SkipZeroValue),
		c.Simple("accessKey", "key_based_authentication.0.access_key", c.SkipZeroValue),
		c.Discriminator("roleBasedAuth", c.DiscriminatorValues{
			"role_based_authentication": true,
			"key_based_authentication":  false,
		}),
		c.Simple("iamRoleARN", "role_based_authentication.0.i_am_role_arn", c.SkipZeroValue),
		c.Simple("useMessageId", "use_message_id", c.SkipZeroValue),
	}

	properties = append(properties, commonProperties...)

	schema := map[string]*schema.Schema{
		"region": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Enter the region.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"stream": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Enter the stream name.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"use_message_id": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Use MessageId as Partition Key",
		},
		"role_based_authentication": {
			Type:         schema.TypeList,
			MaxItems:     1,
			Optional:     true,
			Description:  "This option allows you select the arn based authentication.",
			ExactlyOneOf: []string{"config.0.role_based_authentication", "config.0.key_based_authentication"},
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
		"key_based_authentication": {
			Type:         schema.TypeList,
			MaxItems:     1,
			Optional:     true,
			Description:  "This option allows you select the key based authentication.",
			ExactlyOneOf: []string{"config.0.role_based_authentication", "config.0.key_based_authentication"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"access_key_id": {
						Type:        schema.TypeString,
						Required:    true,
						Sensitive:   true,
						Description: "Enter the AWS Access Key ID.",
					},
					"access_key": {
						Type:        schema.TypeString,
						Required:    true,
						Sensitive:   true,
						Description: "Enter the AWS Secret Access Key.",
					},
				},
			},
		},
	}

	for key, value := range commonSchema {
		schema[key] = value
	}

	c.Destinations.Register("kinesis", c.ConfigMeta{
		APIType:      "KINESIS",
		Properties:   properties,
		ConfigSchema: schema,
	})
}
