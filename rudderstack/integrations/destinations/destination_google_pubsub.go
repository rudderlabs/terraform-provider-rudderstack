package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"android", "ios", "web", "unity", "amp", "cloud", "warehouse", "reactnative", "flutter", "cordova", "shopify"}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)

	properties := []c.ConfigProperty{
		c.Simple("projectId", "project_id"),
		c.Simple("credentials", "credentials"),
		c.ArrayWithObjects("eventToTopicMap", "event_to_topic_map", map[string]interface{}{
			"from": "from",
			"to":   "to",
		}),
		c.ArrayWithObjects("eventToAttributesMap", "event_to_attribute_map", map[string]interface{}{
			"from": "from",
			"to":   "to",
		}),
	}

	properties = append(properties, commonProperties...)

	schema := map[string]*schema.Schema{
		"project_id": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Specify the Google Cloud Project ID.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"credentials": {
			Type:             schema.TypeString,
			Required:         true,
			Sensitive:        true,
			Description:      "The credentials JSON is used by the client library to access the Pub/Sub API.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|.+"),
		},
		"event_to_topic_map": {
			Type:        schema.TypeList,
			ConfigMode:  schema.SchemaConfigModeAttr,
			Optional:    true,
			Description: "Map RudderStack events to Google pub/sub topics",
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
		"event_to_attribute_map": {
			Type:        schema.TypeList,
			ConfigMode:  schema.SchemaConfigModeAttr,
			Optional:    true,
			Description: "Map message properties to Google pub/sub message Attribute Key",
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

	c.Destinations.Register("google_pubsub", c.ConfigMeta{
		APIType:      "GOOGLEPUBSUB",
		Properties:   properties,
		ConfigSchema: schema,
	})
}
