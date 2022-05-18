package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	c.Destinations.Register("zendesk", c.ConfigMeta{
		APIType: "ZENDESK",
		Properties: []c.ConfigProperty{
			c.Simple("email", "email"),
			c.Simple("apiToken", "api_token"),
			c.Simple("domain", "domain"),
			c.Simple("createUsersAsVerified", "create_users_as_verified", c.SkipZeroValue),
			c.Simple("sendGroupCallsWithoutUserId", "send_group_calls_without_user_id", c.SkipZeroValue),
			c.Simple("removeUsersFromOrganization", "remove_users_from_organization", c.SkipZeroValue),
		},
		ConfigSchema: map[string]*schema.Schema{
			"email": {
				Type:     schema.TypeString,
				Required: true,
				// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"api_token": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"create_users_as_verified": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"send_group_calls_without_user_id": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"remove_users_from_organization": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	})
}
