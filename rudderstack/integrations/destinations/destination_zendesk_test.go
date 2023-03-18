package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceZendesk(t *testing.T) {
	cmt.AssertDestination(t, "zendesk", []c.TestConfig{
		{
			TerraformCreate: `
				email     = "test@example.com"
				api_token = "..."
				domain    = "..."
			`,
			APICreate: `{
				"email": "test@example.com",
				"apiToken": "...",
				"domain": "..."
			}`,
			TerraformUpdate: `
				email     = "test@example.com"
				api_token = "..."
				domain    = "..."
			
				create_users_as_verified         = true
				send_group_calls_without_user_id = true
				remove_users_from_organization   = true
				search_by_external_id = false
			`,
			APIUpdate: `{
				"email": "test@example.com",
				"apiToken": "...",
				"domain": "...",
				"createUsersAsVerified": true,
				"sendGroupCallsWithoutUserId": true,
				"removeUsersFromOrganization": true,
				"searchByExternalId": false
			}`,
		},
	})
}
