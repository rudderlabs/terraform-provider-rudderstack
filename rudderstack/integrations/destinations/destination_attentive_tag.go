package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"web", "android", "ios", "unity", "amp", "cloud", "warehouse", "reactnative", "flutter", "cordova", "shopify"}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)

	properties := []c.ConfigProperty{
		c.Simple("apiKey", "api_key"),
		c.Simple("signUpSourceId", "sign_up_source_id", c.SkipZeroValue),
	}

	properties = append(properties, commonProperties...)

	schema := map[string]*schema.Schema{
		"api_key": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "Enter your API Key.",
		},
		"sign_up_source_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Enter your Sign Up Source ID.",
		},
	}

	for key, value := range commonSchema {
		schema[key] = value
	}

	c.Destinations.Register("attentive_tag", c.ConfigMeta{
		APIType:      "ATTENTIVE_TAG",
		Properties:   properties,
		ConfigSchema: schema,
	})
}
