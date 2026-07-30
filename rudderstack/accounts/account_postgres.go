package accounts

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	properties := []c.ConfigProperty{
		c.Simple("options.host", "host"),
		c.Simple("options.dbname", "dbname"),
		c.Simple("options.user", "user"),
		c.Simple("options.port", "port"),
		c.Simple("options.sslMode", "ssl_mode", c.SkipZeroValue),
		c.Simple("secret.password", "password"),
	}
	cfgSchema := map[string]*schema.Schema{
		"host": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Postgres host.",
		},
		"dbname": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Database name.",
		},
		"user": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Database username. Non-secret option field, matching the Postgres source form.",
		},
		"port": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Database port (e.g. \"5432\").",
		},
		"ssl_mode": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "SSL mode: \"disable\" or \"require\". Defaults to \"disable\" when omitted.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(disable|require)$"),
		},
		"password": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "Database password.",
		},
	}
	c.Accounts.Register("postgres", c.ConfigMeta{
		APIType:      "SOURCE_POSTGRES",
		Properties:   properties,
		ConfigSchema: cfgSchema,
	})
}
