package source_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/accounts"
)

func TestAccountResourceSourceDatabricks(t *testing.T) {
	cmt.AssertAccount(t, accounts.CategorySource, "databricks", []cmt.AccountsTestConfig{
		{
			TerraformCreate: `
				host    = "some-host"
				path	  = "some-path"
				token   = "some-token"
			`,
			APICreate: cmt.TestConfigAPIPayload{
				Options: `{
					"host": "some-host",
					"path": "some-path",
					"port": 443
				}`,
				Secret: `{
					"token": "some-token"
				}`,
			},
			TerraformUpdate: `
				host    = "some-host-updated"
				token   = "some-token-updated"
				path	  = "some-path"
				port    = 444
				catalog = "some-catalog"
			`,
			APIUpdate: cmt.TestConfigAPIPayload{
				Options: `{
					"host": "some-host-updated",
					"path": "some-path",
					"port": 444,
					"catalog": "some-catalog"
				}`,
				Secret: `{
					"token": "some-token-updated"
				}`,
			},
		},
	})
}
