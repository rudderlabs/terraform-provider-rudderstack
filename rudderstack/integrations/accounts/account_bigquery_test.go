package accounts_test

import (
	"testing"

	acc "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/acc"
	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
	_ "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/integrations/accounts"
)

var bigqueryAccountTestConfigs = []c.TestConfig{
	{
		TerraformCreate: `
			project     = "my-gcp-project"
			location    = "US"
			credentials = "{\"type\":\"service_account\"}"
		`,
		APICreate: `{
			"name": "example",
			"accountDefinitionName": "SOURCE_BIGQUERY",
			"options": { "projectId": "my-gcp-project", "location": "US" },
			"secret":  { "credentials": "{\"type\":\"service_account\"}" }
		}`,
		TerraformUpdate: `
			project     = "my-gcp-project"
			location    = "EU"
			credentials = "{\"type\":\"service_account\"}"
		`,
		APIUpdate: `{
			"name": "example-updated",
			"accountDefinitionName": "SOURCE_BIGQUERY",
			"options": { "projectId": "my-gcp-project", "location": "EU" },
			"secret":  { "credentials": "{\"type\":\"service_account\"}" }
		}`,
	},
}

func TestAccountResourceBigQuery(t *testing.T) {
	cmt.AssertAccount(t, "bigquery", bigqueryAccountTestConfigs)
}

func TestAccAccountBigQuery(t *testing.T) {
	acc.AccAssertAccount(t, "bigquery", bigqueryAccountTestConfigs)
}
