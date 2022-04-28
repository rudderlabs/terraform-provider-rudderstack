package configs

type TestConfig struct {
	TerraformCreate string
	APICreate       string
	TerraformUpdate string
	APIUpdate       string
}

var EmptyTestConfig = TestConfig{TerraformCreate: "", APICreate: "{}", TerraformUpdate: "", APIUpdate: "{}"}
