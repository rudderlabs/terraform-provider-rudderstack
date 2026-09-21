package configs

type TestConfig struct {
	TerraformCreate string
	APICreate       string
	TerraformUpdate string
	APIUpdate       string
	// APIResponseOptionalFields lists top-level config fields that are sent to the
	// API but may be omitted from destination read responses. When present in a
	// response, their values are still verified.
	APIResponseOptionalFields []string
	// APICreateSettings and APIUpdateSettings hold the expected settings JSON for
	// source-level fields (GeoEnrichmentEnabled, Transient) that are asserted separately from Config.
	APICreateSettings string
	APIUpdateSettings string
}

var EmptyTestConfig = TestConfig{TerraformCreate: "", APICreate: "{}", TerraformUpdate: "", APIUpdate: "{}"}
