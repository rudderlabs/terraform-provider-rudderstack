package sources

import c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"

func init() {
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
}
