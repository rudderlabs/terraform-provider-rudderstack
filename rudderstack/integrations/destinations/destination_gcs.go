package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"web", "android", "ios", "unity", "reactnative", "flutter", "cordova", "amp", "cloud", "warehouse", "shopify"}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)

	properties := []c.ConfigProperty{
		c.Simple("bucketName", "bucket_name"),
		c.Simple("prefix", "prefix", c.SkipZeroValue),
		c.Simple("credentials", "credentials", c.SkipZeroValue),
	}

	properties = append(properties, commonProperties...)

	schema := map[string]*schema.Schema{
		"bucket_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Enter your Google Cloud Storage bucket name.",
			// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"prefix": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Enter your prefix which RudderStack associates with your GCS bucket before loading all the data into it.",
			// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"credentials": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "Enter the contents of your Google Cloud connection credentials JSON.",
			// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
	}

	for key, value := range commonSchema {
		schema[key] = value
	}

	c.Destinations.Register("gcs", c.ConfigMeta{
		APIType:      "GCS",
		Properties:   properties,
		ConfigSchema: schema,
	})
}
