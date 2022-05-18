package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	c.Destinations.Register("gcs", c.ConfigMeta{
		APIType: "GCS",
		Properties: []c.ConfigProperty{
			c.Simple("bucketName", "bucket_name"),
			c.Simple("prefix", "prefix", c.SkipZeroValue),
			c.Simple("credentials", "credentials", c.SkipZeroValue),
		},
		ConfigSchema: map[string]*schema.Schema{
			"bucket_name": {
				Type:     schema.TypeString,
				Required: true,
				// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"prefix": {
				Type:     schema.TypeString,
				Optional: true,
				// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"credentials": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
		},
	})
}
