package sources

import c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"

func init() {
	c.Sources.Register("braze", c.ConfigMeta{
		APIType:    "Braze",
		Properties: []c.ConfigProperty{},
		SkipConfig: true,
	})
	c.Sources.Register("cordova", c.ConfigMeta{
		APIType:    "Cordova",
		Properties: []c.ConfigProperty{},
		SkipConfig: true,
	})
	c.Sources.Register("go", c.ConfigMeta{
		APIType:    "Go",
		Properties: []c.ConfigProperty{},
		SkipConfig: true,
	})
	c.Sources.Register("http", c.ConfigMeta{
		APIType:    "HTTP",
		Properties: []c.ConfigProperty{},
		SkipConfig: true,
	})
	c.Sources.Register("android", c.ConfigMeta{
		APIType:    "Android",
		Properties: []c.ConfigProperty{},
		SkipConfig: true,
	})
	c.Sources.Register("ios", c.ConfigMeta{
		APIType:    "iOS",
		Properties: []c.ConfigProperty{},
		SkipConfig: true,
	})
	c.Sources.Register("java", c.ConfigMeta{
		APIType:    "Java",
		Properties: []c.ConfigProperty{},
		SkipConfig: true,
	})
	c.Sources.Register("javascript", c.ConfigMeta{
		APIType:    "Javascript",
		Properties: []c.ConfigProperty{},
		SkipConfig: true,
	})
	c.Sources.Register("node", c.ConfigMeta{
		APIType:    "Node",
		Properties: []c.ConfigProperty{},
		SkipConfig: true,
	})
	c.Sources.Register("reactnative", c.ConfigMeta{
		APIType:    "ReactNative",
		Properties: []c.ConfigProperty{},
		SkipConfig: true,
	})
	c.Sources.Register("ruby", c.ConfigMeta{
		APIType:    "Ruby",
		Properties: []c.ConfigProperty{},
		SkipConfig: true,
	})
	c.Sources.Register("webhook", c.ConfigMeta{
		APIType:    "webhook",
		Properties: []c.ConfigProperty{},
		SkipConfig: true,
	})
	c.Sources.Register("webhook_shopify", c.ConfigMeta{
		APIType:    "Shopify",
		Properties: []c.ConfigProperty{},
		SkipConfig: true,
	})
	c.Sources.Register("python", c.ConfigMeta{
		APIType:    "Python",
		Properties: []c.ConfigProperty{},
		SkipConfig: true,
	})
	c.Sources.Register("php", c.ConfigMeta{
		APIType:    "PHP",
		Properties: []c.ConfigProperty{},
		SkipConfig: true,
	})

	c.Sources.Register("dotnet", c.ConfigMeta{
		APIType:    "DotNet",
		Properties: []c.ConfigProperty{},
		SkipConfig: true,
	})

	c.Sources.Register("flutter", c.ConfigMeta{
		APIType:    "Flutter",
		Properties: []c.ConfigProperty{},
		SkipConfig: true,
	})

	c.Sources.Register("customerio", c.ConfigMeta{
		APIType:    "Customerio",
		Properties: []c.ConfigProperty{},
		SkipConfig: true,
	})

	c.Sources.Register("facebook_lead_ads", c.ConfigMeta{
		APIType:    "facebook_lead_ads",
		Properties: []c.ConfigProperty{},
		SkipConfig: true,
	})

	c.Sources.Register("adjust", c.ConfigMeta{
		APIType:    "Adjust",
		Properties: []c.ConfigProperty{},
		SkipConfig: true,
	})
}
