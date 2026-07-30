package accounts

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	properties := []c.ConfigProperty{
		c.Simple("options.account", "account"),
		c.Simple("options.dbname", "dbname"),
		c.Simple("options.warehouse", "warehouse"),
		c.Simple("options.user", "user"),
		c.Simple("options.role", "role", c.SkipZeroValue),
		c.Simple("options.authenticationType", "authentication_type"),
		// Exactly one auth secret is set at a time; SkipZeroValue omits the unset
		// ones from the API request. The keyPair<->secret pairing is enforced
		// server-side by the SOURCE_SNOWFLAKE combinedSchema, not here.
		c.Simple("secret.password", "password", c.SkipZeroValue),
		c.Simple("secret.privateKey", "private_key", c.SkipZeroValue),
		c.Simple("secret.privateKeyPassphrase", "private_key_passphrase", c.SkipZeroValue),
	}
	cfgSchema := map[string]*schema.Schema{
		"account": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Snowflake account identifier.",
		},
		"dbname": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Database name.",
		},
		"warehouse": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Virtual warehouse.",
		},
		"user": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Login name. Non-secret option field, matching the Snowflake source form.",
		},
		"role": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Snowflake role.",
		},
		"authentication_type": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Auth method: \"keyPair\" or \"password\". Determines which secret is required.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(keyPair|password)$"),
		},
		"password": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "Password. Required when authentication_type is \"password\".",
		},
		"private_key": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "PEM private key. Required when authentication_type is \"keyPair\". Load it with Terraform's file() function.",
		},
		"private_key_passphrase": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "Passphrase protecting the private key (optional, keyPair auth only).",
		},
	}
	c.Accounts.Register("snowflake", c.ConfigMeta{
		APIType:      "SOURCE_SNOWFLAKE",
		Properties:   properties,
		ConfigSchema: cfgSchema,
	})
}
