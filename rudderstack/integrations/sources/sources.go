package configs

import c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"

func init() {
	c.Sources.Register("http", c.ConfigMeta{
		APIType:        "HTTP",
		Properties:     []c.ConfigProperty{},
		OptionalConfig: true,
		TestConfigs:    []c.TestConfig{c.EmptyTestConfig},
	})
	c.Sources.Register("javascript", c.ConfigMeta{
		APIType:        "Javascript",
		Properties:     []c.ConfigProperty{},
		OptionalConfig: true,
		TestConfigs:    []c.TestConfig{c.EmptyTestConfig},
	})
}
