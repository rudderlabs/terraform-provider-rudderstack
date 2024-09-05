package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"web", "android", "ios", "unity", "reactnative", "flutter", "cordova", "amp", "cloud", "warehouse", "shopify"}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)

	properties := []c.ConfigProperty{
		c.Simple("email", "email"),
		c.Simple("apiToken", "api_token"),
		c.Simple("domain", "domain"),
		c.Simple("createUsersAsVerified", "create_users_as_verified", c.SkipZeroValue),
		c.Simple("sendGroupCallsWithoutUserId", "send_group_calls_without_user_id", c.SkipZeroValue),
		c.Simple("removeUsersFromOrganization", "remove_users_from_organization", c.SkipZeroValue),
		c.Simple("searchByExternalId", "search_by_external_id", c.SkipZeroValue),
	}

	properties = append(properties, commonProperties...)

	schema := map[string]*schema.Schema{
		"email": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Enter the email used to log into your Zendesk dashboard.",
			// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"api_token": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "Enter the Zendesk API token used to authenticate the request. To create an API token, refer to this [Zendesk support page](https://support.zendesk.com/hc/en-us/articles/226022787-Generating-a-new-API-token-).",
			// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"domain": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Enter your Zendesk subdomain without `.zendesk.com`",
			// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"create_users_as_verified": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Enabling this setting creates verified users in Zendesk, that is, the email verification is skipped.",
		},
		"send_group_calls_without_user_id": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Enable this setting if you don't want to associate the user with a group. ",
		},
		"remove_users_from_organization": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Enable this setting to remove users from an organization.",
		},
		"search_by_external_id": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Update user's primary email.",
		},
	}

	for key, value := range commonSchema {
		schema[key] = value
	}

	c.Destinations.Register("zendesk", c.ConfigMeta{
		APIType:      "ZENDESK",
		Properties:   properties,
		ConfigSchema: schema,
	})
}
