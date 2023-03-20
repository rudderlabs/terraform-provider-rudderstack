package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	c.Destinations.Register("s3", c.ConfigMeta{
		APIType: "S3",
		Properties: []c.ConfigProperty{
			c.Simple("bucketName", "bucket_name"),
			c.Simple("prefix", "prefix", c.SkipZeroValue),
			c.Simple("accessKeyID", "access_key_id", c.SkipZeroValue),
			c.Simple("accessKey", "access_key", c.SkipZeroValue),
			c.Simple("enableSSE", "enable_sse", c.SkipZeroValue),
			c.ArrayWithStrings("oneTrustCookieCategories", "oneTrustCookieCategory", "onetrust_cookie_categories"),
		},
		ConfigSchema: map[string]*schema.Schema{
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
