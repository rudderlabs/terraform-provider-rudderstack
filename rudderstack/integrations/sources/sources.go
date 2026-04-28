package sources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

var defaultSettingsSchema = map[string]*schema.Schema{
	"geo_enrichment_enabled": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
		Description: "Enable geolocation enrichment on events from this source. " +
			"Requires a paid plan.",
	},
	"temporarily_store_events_for_retries": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Temporarily store events in the gateway staging area to enable retry on failure. Requires a paid plan.",
	},
}

var defaultSettingsProperties = []c.ConfigProperty{
	c.Simple("geoEnrichmentEnabled", "geo_enrichment_enabled"),
	c.Simple("transient", "temporarily_store_events_for_retries"),
}

func init() {
	c.Sources.Register("braze", c.ConfigMeta{
		APIType:            "Braze",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})
	c.Sources.Register("cordova", c.ConfigMeta{
		APIType:            "Cordova",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})
	c.Sources.Register("go", c.ConfigMeta{
		APIType:            "Go",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})
	c.Sources.Register("http", c.ConfigMeta{
		APIType:            "HTTP",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})
	c.Sources.Register("android", c.ConfigMeta{
		APIType:            "Android",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})
	c.Sources.Register("ios", c.ConfigMeta{
		APIType:            "iOS",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})
	c.Sources.Register("java", c.ConfigMeta{
		APIType:            "Java",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})
	c.Sources.Register("javascript", c.ConfigMeta{
		APIType:            "Javascript",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})
	c.Sources.Register("node", c.ConfigMeta{
		APIType:            "Node",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})
	c.Sources.Register("reactnative", c.ConfigMeta{
		APIType:            "ReactNative",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})
	c.Sources.Register("ruby", c.ConfigMeta{
		APIType:            "Ruby",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})
	c.Sources.Register("webhook", c.ConfigMeta{
		APIType:            "webhook",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})
	c.Sources.Register("webhook_shopify", c.ConfigMeta{
		APIType:            "Shopify",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})
	c.Sources.Register("python", c.ConfigMeta{
		APIType:            "Python",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})
	c.Sources.Register("php", c.ConfigMeta{
		APIType:            "PHP",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})

	c.Sources.Register("dotnet", c.ConfigMeta{
		APIType:            "DotNet",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})

	c.Sources.Register("flutter", c.ConfigMeta{
		APIType:            "Flutter",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})

	c.Sources.Register("customerio", c.ConfigMeta{
		APIType:            "Customerio",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})

	c.Sources.Register("facebook_lead_ads", c.ConfigMeta{
		APIType:            "facebook_lead_ads",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})

	c.Sources.Register("adjust", c.ConfigMeta{
		APIType:            "adjust",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})

	c.Sources.Register("rust", c.ConfigMeta{
		APIType:            "Rust",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})

	c.Sources.Register("amp", c.ConfigMeta{
		APIType:            "AMP",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})

	c.Sources.Register("android_kotlin", c.ConfigMeta{
		APIType:            "android_kotlin",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})

	c.Sources.Register("ios_swift", c.ConfigMeta{
		APIType:            "ios_swift",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})

	c.Sources.Register("unity", c.ConfigMeta{
		APIType:            "Unity",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})

	c.Sources.Register("appcenter", c.ConfigMeta{
		APIType:            "Appcenter",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})

	c.Sources.Register("appsflyer", c.ConfigMeta{
		APIType:            "appsflyer",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})

	c.Sources.Register("auth0", c.ConfigMeta{
		APIType:            "Auth0",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})

	c.Sources.Register("canny", c.ConfigMeta{
		APIType:            "canny",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})

	c.Sources.Register("close_crm", c.ConfigMeta{
		APIType:            "close_crm",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})

	c.Sources.Register("cordial", c.ConfigMeta{
		APIType:            "cordial",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})

	c.Sources.Register("extole", c.ConfigMeta{
		APIType:            "Extole",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})

	c.Sources.Register("formsort", c.ConfigMeta{
		APIType:            "formsort",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})

	c.Sources.Register("gainsightpx", c.ConfigMeta{
		APIType:            "GainsightPX",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})

	c.Sources.Register("iterable", c.ConfigMeta{
		APIType:            "iterable",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})

	c.Sources.Register("looker", c.ConfigMeta{
		APIType:            "Looker",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})

	c.Sources.Register("mailjet", c.ConfigMeta{
		APIType:            "mailjet",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})

	c.Sources.Register("mailmodo", c.ConfigMeta{
		APIType:            "mailmodo",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})

	c.Sources.Register("moengage", c.ConfigMeta{
		APIType:            "moengage",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})

	c.Sources.Register("monday", c.ConfigMeta{
		APIType:            "monday",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})

	c.Sources.Register("olark", c.ConfigMeta{
		APIType:            "olark",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})

	c.Sources.Register("ortto", c.ConfigMeta{
		APIType:            "ortto",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})

	c.Sources.Register("pagerduty", c.ConfigMeta{
		APIType:            "pagerduty",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})

	c.Sources.Register("pipedream", c.ConfigMeta{
		APIType:            "pipedream",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})

	c.Sources.Register("refiner", c.ConfigMeta{
		APIType:            "refiner",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})

	c.Sources.Register("revenuecat", c.ConfigMeta{
		APIType:            "revenuecat",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})

	c.Sources.Register("satismeter", c.ConfigMeta{
		APIType:            "satismeter",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})

	c.Sources.Register("segment", c.ConfigMeta{
		APIType:            "Segment",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})

	c.Sources.Register("signl4", c.ConfigMeta{
		APIType:            "signl4",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})

	c.Sources.Register("slack", c.ConfigMeta{
		APIType:            "slack",
		Properties:         []c.ConfigProperty{},
		SkipConfig:         true,
		SettingsSchema:     defaultSettingsSchema,
		SettingsProperties: defaultSettingsProperties,
	})
}
