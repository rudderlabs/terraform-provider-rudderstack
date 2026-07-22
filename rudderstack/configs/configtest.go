package configs

type TestConfig struct {
	TerraformCreate string
	APICreate       string
	TerraformUpdate string
	APIUpdate       string
	// APIReadRedactedFields lists top-level destination config fields that are sent
	// on create/update but are not returned by the live API on subsequent reads.
	APIReadRedactedFields []string
	// ImportStateVerifyIgnore lists Terraform state paths that cannot be reconstructed
	// during import, typically because the live API redacts write-only config values.
	ImportStateVerifyIgnore []string
	// APICreateSettings and APIUpdateSettings hold the expected settings JSON for
	// source-level fields (GeoEnrichmentEnabled, Transient) that are asserted separately from Config.
	APICreateSettings string
	APIUpdateSettings string
}

var EmptyTestConfig = TestConfig{TerraformCreate: "", APICreate: "{}", TerraformUpdate: "", APIUpdate: "{}"}
