package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	c.Destinations.Register("bigquery", c.ConfigMeta{
		APIType: "BQ",
		Properties: []c.ConfigProperty{
			c.Simple("project", "project"),
			c.Simple("location", "location", c.SkipZeroValue),
			c.Simple("bucketName", "bucket_name"),
			c.Simple("prefix", "prefix", c.SkipZeroValue),
			c.Simple("namespace", "namespace", c.SkipZeroValue),
			c.Simple("credentials", "credentials"),
			c.Simple("syncFrequency", "sync.0.frequency"),
			c.Simple("syncStartAt", "sync.0.start_at", c.SkipZeroValue),
			c.Simple("excludeWindow.excludeWindowStartTime", "sync.0.exclude_window_start_time", c.SkipZeroValue),
			c.Simple("excludeWindow.excludeWindowEndTime", "sync.0.exclude_window_end_time", c.SkipZeroValue),
		},
		ConfigSchema: map[string]*schema.Schema{
			"project": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
			},
			"location": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
			},
			"bucket_name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
			},
			"prefix": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
			},
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
				// TODO: panic: error parsing regexp: invalid or unsupported Perl syntax: `(?!`
				// ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].*)|^((?!pg_|PG_|pG_|Pg_).{0,64})$"),
			},
			"credentials": {
				Type:             schema.TypeString,
				Required:         true,
				Sensitive:        true,
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
		},
	})
}
