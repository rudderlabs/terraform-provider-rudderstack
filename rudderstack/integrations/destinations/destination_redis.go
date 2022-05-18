package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	c.Destinations.Register("redis", c.ConfigMeta{
		APIType: "REDIS",
		Properties: []c.ConfigProperty{
			c.Simple("address", "address"),
			c.Simple("password", "password", c.SkipZeroValue),
			c.Simple("clusterMode", "cluster_mode"),
			c.Simple("secure", "secure"),
			c.Simple("prefix", "prefix", c.SkipZeroValue),
			c.Simple("database", "database", c.SkipZeroValue),
			c.Simple("caCertificate", "ca_certificate", c.SkipZeroValue),
			c.Simple("skipVerify", "skip_verify"),
		},
		ConfigSchema: map[string]*schema.Schema{
			"address": {
				Type:     schema.TypeString,
				Required: true,
				// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"cluster_mode": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"secure": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"prefix": {
				Type:     schema.TypeString,
				Optional: true,
				// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"database": {
				Type:     schema.TypeString,
				Optional: true,
				// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"ca_certificate": {
				Type:     schema.TypeString,
				Optional: true,
				// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"skip_verify": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	})
}
