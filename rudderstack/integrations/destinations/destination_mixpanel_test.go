package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceMixPanel(t *testing.T) {
	cmt.AssertDestination(t, "mixpanel", []c.TestConfig{
		{
			TerraformCreate: `
				token = "avasaffav1241"
				data_residency = "us"
				persistence = "none"
			`,
			APICreate: `{
				"token": "avasaffav1241",
				"dataResidency": "us",
				"persistence: "none"
			}"
			`,
			TerraformUpdate: `
				data_residency = "us"
				people = true
				set_all_traits_by_default = true
				use_native_sdk {
					web = true
				}
			`,
			APIUpdate: `{
				"dataResidency": "us",
				"people": true,
				"setAllTraitsByDefault": true,
				"useNativeSDK": {
				  "web": true
				}
			}
			`,
		},
	})
}
