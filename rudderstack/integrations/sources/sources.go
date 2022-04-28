package sources

import c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"

func init() {
	c.Sources.Register("http", c.ConfigMeta{
		APIType:    "HTTP",
		Properties: []c.ConfigProperty{},
		SkipConfig: true,
	})
	c.Sources.Register("javascript", c.ConfigMeta{
		APIType:    "Javascript",
		Properties: []c.ConfigProperty{},
		SkipConfig: true,
	})
}
