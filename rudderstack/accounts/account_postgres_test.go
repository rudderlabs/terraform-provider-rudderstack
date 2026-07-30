package accounts_test

import (
	"testing"

	acc "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/acc"
	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	_ "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/accounts"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

var postgresAccountTestConfigs = []c.TestConfig{
	{
		// Create with the required fields only; the optional `ssl_mode` is exercised
		// in the update step below.
		TerraformCreate: `
			host     = "db.example.com"
			dbname   = "analytics"
			user     = "rudder"
			port     = "5432"
			password = "s3cr3t"
		`,
		APICreate: `{
			"name": "example",
			"accountDefinitionName": "SOURCE_POSTGRES",
			"options": { "host": "db.example.com", "dbname": "analytics", "user": "rudder", "port": "5432" },
			"secret":  { "password": "s3cr3t" }
		}`,
		TerraformUpdate: `
			host     = "db.example.com"
			dbname   = "analytics"
			user     = "rudder"
			port     = "5432"
			ssl_mode = "require"
			password = "s3cr3t"
		`,
		APIUpdate: `{
			"name": "example-updated",
			"accountDefinitionName": "SOURCE_POSTGRES",
			"options": { "host": "db.example.com", "dbname": "analytics", "user": "rudder", "port": "5432", "sslMode": "require" },
			"secret":  { "password": "s3cr3t" }
		}`,
	},
}

func TestAccountResourcePostgres(t *testing.T) {
	cmt.AssertAccount(t, "postgres", postgresAccountTestConfigs)
}

func TestAccAccountPostgres(t *testing.T) {
	acc.AccAssertAccount(t, "postgres", postgresAccountTestConfigs)
}
