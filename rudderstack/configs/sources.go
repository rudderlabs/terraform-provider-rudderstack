package configs

func init() {
	Sources.Register("http", ConfigMeta{
		APIType:        "HTTP",
		Properties:     []ConfigProperty{},
		OptionalConfig: true,
		TestConfigs:    []TestConfig{emptyTestConfig},
	})
	Sources.Register("javascript", ConfigMeta{
		APIType:        "Javascript",
		Properties:     []ConfigProperty{},
		OptionalConfig: true,
		TestConfigs:    []TestConfig{emptyTestConfig},
	})
}
