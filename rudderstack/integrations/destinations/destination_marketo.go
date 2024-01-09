package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	c.Destinations.Register("marketo", c.ConfigMeta{
		APIType: "MARKETO",
		Properties: []c.ConfigProperty{
			c.Simple("accountId", "account_id"),
			c.Simple("clientId", "client_id"),
			c.Simple("clientSecret", "client_secret"),
			c.Simple("trackAnonymousEvents", "track_anonymous_events"),
			c.Simple("createIfNotExist", "create_if_not_exist"),
			c.ArrayWithStrings("oneTrustCookieCategories", "oneTrustCookieCategory", "onetrust_cookie_categories"),
			c.ArrayWithObjects("rudderEventsMapping", "rudder_events_mapping", map[string]string{
				"event": "event",
				"marketoPrimarykey":   "marketoPrimarykey",
				"marketoActivityId":   "marketoActivityId",
			}),
			c.ArrayWithObjects("leadTraitMapping", "lead_trait_mapping", map[string]string{
				"from": "from",
				"to":   "to"
			}),
			c.ArrayWithObjects("customActivityPropertyMap", "custom_activity_property_map", map[string]string{
				"from": "from",
				"to":   "to"
			}),
		},
		ConfigSchema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Marketo Account ID",
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Marketo Client ID",
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Marketo Client Secret",
				Sensitive:   true,
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"track_anonymous_events": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Boolean flag to track anonymous events",
			},
			"create_if_not_exist": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Boolean flag to create lead if not exist",
			},
			"lead_trait_mapping": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Lead Trait Mapping",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from": {
							Type:     schema.TypeString,
							Required: true,
						},
						"to": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"rudder_events_mapping": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Rudder Events Mapping",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event": {
							Type:     schema.TypeString,
							Required: true,
						},
						"marketoPrimarykey": {
							Type:     schema.TypeString,
							Required: true,
						},
						"marketoActivityId": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"custom_activity_property_map": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Custom Activity Property Map",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from": {
							Type:     schema.TypeString,
							Required: true,
						},
						"to": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
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
