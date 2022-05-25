package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	c.Destinations.Register("redshift", c.ConfigMeta{
		APIType: "RS",
		Properties: []c.ConfigProperty{
			c.Simple("host", "host"),
			c.Simple("port", "port"),
			c.Simple("database", "database"),
			c.Simple("user", "user"),
			c.Simple("password", "password"),
			c.Simple("namespace", "namespace", c.SkipZeroValue),
			c.Simple("enableSSE", "enable_sse", c.SkipZeroValue),
			c.Simple("useRudderStorage", "use_rudder_storage"),
			c.Simple("syncFrequency", "sync.0.frequency"),
			c.Simple("syncStartAt", "sync.0.start_at", c.SkipZeroValue),
			c.Simple("excludeWindow.excludeWindowStartTime", "sync.0.exclude_window_start_time", c.SkipZeroValue),
			c.Simple("excludeWindow.excludeWindowEndTime", "sync.0.exclude_window_end_time", c.SkipZeroValue),
			c.Simple("bucketName", "s3.0.bucket_name"),
			c.Simple("accessKeyID", "s3.0.access_key_id"),
			c.Simple("accessKey", "s3.0.access_key"),
		},
		ConfigSchema: map[string]*schema.Schema{
			"host": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,255})$"),
			},
			"port": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
			},
			"user": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
			},
			"password": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
				Sensitive:        true,
			},
			"database": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: c.ValidateAll(
					c.StringMatchesRegexp("(^env[.].*)|^(.{0,64})$"),
					c.StringNotMatchesRegexp("^(pg_|PG_|pG_|Pg_)"),
				),
			},
			"enable_sse": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"use_rudder_storage": {
				Type:     schema.TypeBool,
				Required: true,
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
						"exclude_window_start_time": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: c.StringMatchesRegexp("^([01][0-9]|2[0-3]):[0-5][0-9]$"),
						},
						"exclude_window_end_time": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: c.StringMatchesRegexp("^([01][0-9]|2[0-3]):[0-5][0-9]$"),
						},
					},
				},
			},
			"s3": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket_name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
						},
						"access_key_id": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
						},
						"access_key": {
							Type:             schema.TypeString,
							Optional:         true,
							Sensitive:        true,
							ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
						},
					},
				},
			},
		},
	})
}
