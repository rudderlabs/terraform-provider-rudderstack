package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceMixPanel(t *testing.T){
	cmt.AssertDestination(t, "mixpanel", []c.TestConfig{
		{
			TerraformCreate: ``,
			APICreate: ``,
			TerraformUpdate: ``,
			APIUpdate: ``,
		},
	})
}