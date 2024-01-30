package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	c.Destinations.Register("google_sheets", c.ConfigMeta{
		APIType: "GOOGLESHEETS",
		Properties: []c.ConfigProperty{
			c.Simple("sheetName", "sheet_name"),
			c.Simple("sheetId", "sheet_id"),
			c.Simple("credentials", "credentials"),
			c.ArrayWithObjects("eventKeyMap", "event_key_map", map[string]string{
				"from": "from",
				"to":   "to",
			}),
			c.ArrayWithStrings("oneTrustCookieCategories", "oneTrustCookieCategory", "onetrust_cookie_categories"),
		},
		ConfigSchema: map[string]*schema.Schema{
			"sheet_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specify the name of the Google spreadsheet to which you want to send the data.",
			},
			"sheet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Enter your Google sheet ID. You can find it in the spreadsheet URL.",
			},
			"credentials": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Enter the credentials JSON used by the client library to access the Google Sheets API.",
			},
			"event_key_map": {
				Type:        schema.TypeList,
				ConfigMode:  schema.SchemaConfigModeAttr,
				Optional:    true,
				Description: "Add Event Properties to map to Google-Sheets Column.",
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
