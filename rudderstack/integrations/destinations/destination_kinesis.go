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
			c.Simple("accessKeyID", "key_based_authentication.0.access_key_id", c.SkipZeroValue),
			c.Simple("accessKey", "key_based_authentication.0.access_key", c.SkipZeroValue),
			c.Discriminator("roleBasedAuth", c.DiscriminatorValues{
				"role_based_authentication": true,
				"key_based_authentication":  false,
			}),
			c.Simple("iamRoleARN", "role_based_authentication.0.i_am_role_arn", c.SkipZeroValue),
			c.Simple("useMessageId", "use_message_id", c.SkipZeroValue),
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
			"use_message_id": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Use MessageId as Partition Key",
			},
			"role_based_authentication": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "This option allows you filter the events you want to send to Amplitude.",
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
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "This option allows you filter the events you want to send to Amplitude.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_key_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Enter the event names to be blacklisted.",
						},
						"access_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Enter the event names to be blacklisted.",
						},
					},
				},
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
