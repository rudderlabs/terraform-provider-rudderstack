package configs

type TestConfig struct {
	TerraformCreate string
	APICreate       string
	TerraformUpdate string
	APIUpdate       string
	// APICreateSettings and APIUpdateSettings hold the expected settings JSON for
	// source-level fields (GeoEnrichmentEnabled, Transient) that are asserted separately from Config.
	APICreateSettings string
	APIUpdateSettings string
}

var EmptyTestConfig = TestConfig{TerraformCreate: "", APICreate: "{}", TerraformUpdate: "", APIUpdate: "{}"}
