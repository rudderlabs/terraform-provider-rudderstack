package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	c.Destinations.Register("salesforce", c.ConfigMeta{
		APIType: "SALESFORCE",
		Properties: []c.ConfigProperty{
			c.Simple("userName", "user_name"),
			c.Simple("password", "password"),
			c.Simple("initialAccessToken", "initial_access_token"),
			c.Simple("mapProperties", "map_properties", c.SkipZeroValue),
			c.Simple("sandbox", "sandbox", c.SkipZeroValue),
			c.Simple("useContactId", "use_contact_id", c.SkipZeroValue),
			c.ArrayWithStrings("oneTrustCookieCategories", "oneTrustCookieCategory", "onetrust_cookie_categories"),
		},
		ConfigSchema: map[string]*schema.Schema{
			"user_name": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Enter the Salesforce username.",
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"password": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Enter the password for the above user.",
				Sensitive:        true,
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"initial_access_token": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Enter your Salesforce security token.",
				Sensitive:        true,
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"map_properties": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Use this setting to map RudderStack event properties to specific Salesforce fields.",
			},
			"sandbox": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Use this setting to enable Salesforce sandbox mode.",
			},
			"use_contact_id": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "When enabled, RudderStack uses contactId for the converted leads.",
			},
			"onetrust_cookie_categories": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specify the OneTrust category name for mapping the OneTrust consent settings to RudderStack's consent purposes.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	})
}
