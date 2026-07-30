package accounts_test

import (
	"testing"

	acc "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/acc"
	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	_ "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/accounts"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

// keyPair auth: privateKey (+ optional passphrase on update). The unset `password`
// is omitted from the API request by SkipZeroValue.
var snowflakeKeyPairAccountTestConfigs = []c.TestConfig{
	{
		TerraformCreate: `
			account             = "xy12345.eu-west-1"
			dbname              = "ANALYTICS"
			warehouse           = "COMPUTE_WH"
			user                = "RUDDER"
			authentication_type = "keyPair"
			private_key         = "dummy-snowflake-private-key"
		`,
		APICreate: `{
			"name": "example",
			"accountDefinitionName": "SOURCE_SNOWFLAKE",
			"options": { "account": "xy12345.eu-west-1", "dbname": "ANALYTICS", "warehouse": "COMPUTE_WH", "user": "RUDDER", "authenticationType": "keyPair" },
			"secret":  { "privateKey": "dummy-snowflake-private-key" }
		}`,
		TerraformUpdate: `
			account                = "xy12345.eu-west-1"
			dbname                 = "ANALYTICS"
			warehouse              = "COMPUTE_WH"
			user                   = "RUDDER"
			role                   = "ANALYST"
			authentication_type    = "keyPair"
			private_key            = "dummy-snowflake-private-key"
			private_key_passphrase = "pp"
		`,
		APIUpdate: `{
			"name": "example-updated",
			"accountDefinitionName": "SOURCE_SNOWFLAKE",
			"options": { "account": "xy12345.eu-west-1", "dbname": "ANALYTICS", "warehouse": "COMPUTE_WH", "user": "RUDDER", "role": "ANALYST", "authenticationType": "keyPair" },
			"secret":  { "privateKey": "dummy-snowflake-private-key", "privateKeyPassphrase": "pp" }
		}`,
	},
}

// password auth: only `password`. privateKey / passphrase are omitted from the request.
var snowflakePasswordAccountTestConfigs = []c.TestConfig{
	{
		TerraformCreate: `
			account             = "xy12345.eu-west-1"
			dbname              = "ANALYTICS"
			warehouse           = "COMPUTE_WH"
			user                = "RUDDER"
			authentication_type = "password"
			password            = "s3cr3t"
		`,
		APICreate: `{
			"name": "example",
			"accountDefinitionName": "SOURCE_SNOWFLAKE",
			"options": { "account": "xy12345.eu-west-1", "dbname": "ANALYTICS", "warehouse": "COMPUTE_WH", "user": "RUDDER", "authenticationType": "password" },
			"secret":  { "password": "s3cr3t" }
		}`,
		TerraformUpdate: `
			account             = "xy12345.eu-west-1"
			dbname              = "ANALYTICS"
			warehouse           = "COMPUTE_WH"
			user                = "RUDDER"
			role                = "ANALYST"
			authentication_type = "password"
			password            = "s3cr3t"
		`,
		APIUpdate: `{
			"name": "example-updated",
			"accountDefinitionName": "SOURCE_SNOWFLAKE",
			"options": { "account": "xy12345.eu-west-1", "dbname": "ANALYTICS", "warehouse": "COMPUTE_WH", "user": "RUDDER", "role": "ANALYST", "authenticationType": "password" },
			"secret":  { "password": "s3cr3t" }
		}`,
	},
}

func TestAccountResourceSnowflakeKeyPair(t *testing.T) {
	cmt.AssertAccount(t, "snowflake", snowflakeKeyPairAccountTestConfigs)
}

func TestAccountResourceSnowflakePassword(t *testing.T) {
	cmt.AssertAccount(t, "snowflake", snowflakePasswordAccountTestConfigs)
}

func TestAccAccountSnowflakeKeyPair(t *testing.T) {
	acc.AccAssertAccount(t, "snowflake", snowflakeKeyPairAccountTestConfigs)
}

func TestAccAccountSnowflakePassword(t *testing.T) {
	acc.AccAssertAccount(t, "snowflake", snowflakePasswordAccountTestConfigs)
}
