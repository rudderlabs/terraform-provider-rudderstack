package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	c.Destinations.Register("s3_datalake", c.ConfigMeta{
		APIType: "S3_DATALAKE",
		Properties: []c.ConfigProperty{
			c.Simple("bucketName", "bucket_name"),
			c.Simple("namespace", "namespace", c.SkipZeroValue),
			c.Simple("prefix", "prefix", c.SkipZeroValue),
			c.Simple("accessKeyID", "access_key_id", c.SkipZeroValue),
			c.Simple("accessKey", "access_key", c.SkipZeroValue),
			c.Simple("enableSSE", "enable_sse", c.SkipZeroValue),
			c.Simple("useGlue", "use_glue"),
			c.Simple("region", "region", c.SkipZeroValue),
			c.Simple("syncFrequency", "sync.0.frequency"),
			c.Simple("syncStartAt", "sync.0.start_at", c.SkipZeroValue),
		},
		ConfigSchema: map[string]*schema.Schema{
			"bucket_name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateDiagFunc: c.ValidateAll(
					c.StringMatchesRegexp("(^env[.].*)|^(.{0,64})$"),
					c.StringNotMatchesRegexp("^(pg_|PG_|pG_|Pg_)"),
				),
			},
			"prefix": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"access_key_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"access_key": {
				Type:             schema.TypeString,
				Optional:         true,
				Sensitive:        true,
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"enable_sse": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"use_glue": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"region": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|.+"),
			},
			"sync": {
				Type:     schema.TypeList,
				MinItems: 1, MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"frequency": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: c.StringMatchesRegexp("^(30|60|180|360|720|1440)$"),
						},
						"start_at": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: c.StringMatchesRegexp("^([01][0-9]|2[0-3]):[0-5][0-9]$"),
						},
					},
				},
			},
		},
	})
}
