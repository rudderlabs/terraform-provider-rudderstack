package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	c.Destinations.Register("kinesis", c.ConfigMeta{
		APIType: "KINESIS",
		Properties: []c.ConfigProperty{
			c.Simple("region", "region"),
			c.Simple("stream", "stream"),
			c.Simple("accessKeyID", "access_key_id", c.SkipZeroValue),
			c.Simple("accessKey", "access_key", c.SkipZeroValue),
			c.Simple("roleBasedAuth", "role_based_auth"),
			c.Simple("iamRoleARN", "i_am_role_arn", c.SkipZeroValue),
			c.Simple("useMessageId", "use_message_id"),
			c.ArrayWithStrings("oneTrustCookieCategories", "oneTrustCookieCategory", "onetrust_cookie_categories"),
		},
		ConfigSchema: map[string]*schema.Schema{
			"region": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Region",
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"stream": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Stream Name",
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"access_key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS Access Key ID",
			},
			"access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "AWS Secret Access Key",
			},
			"i_am_role_arn": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "IAM Role ARN",
			},
			"use_message_id": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Use MessageId as Partition Key",
			},
			"role_based_auth": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Role Based Authentication",
			},
			"onetrust_cookie_categories": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specify the OneTrust category name for mapping the OneTrust consent settings to RudderStack's consent purposes.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	})
}
