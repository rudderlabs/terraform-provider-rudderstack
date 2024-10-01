package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"web", "android", "ios", "unity", "reactnative", "flutter", "cordova", "amp", "cloud", "warehouse", "shopify"}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)

	properties := []c.ConfigProperty{
		c.Simple("sheetName", "sheet_name"),
		c.Simple("sheetId", "sheet_id"),
		c.Simple("credentials", "credentials"),
		c.ArrayWithObjects("eventKeyMap", "event_key_map", map[string]interface{}{
			"from": "from",
			"to":   "to",
		}),
	}

	properties = append(properties, commonProperties...)

	schema := map[string]*schema.Schema{
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
	}

	for key, value := range commonSchema {
		schema[key] = value
	}

	c.Destinations.Register("google_sheets", c.ConfigMeta{
		APIType:      "GOOGLESHEETS",
		Properties:   properties,
		ConfigSchema: schema,
	})
}
