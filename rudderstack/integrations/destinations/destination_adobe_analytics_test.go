package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceAdobeAnalytics(t *testing.T) {
	cmt.AssertDestination(t, "adobe_analytics", []c.TestConfig{
		{
			TerraformCreate: `
				report_suite_ids = "id001, id002"
			`,
			APICreate: `{
				"reportSuiteIds": "id001, id002"
			}`,
			TerraformUpdate: `
				report_suite_ids = "id001, id003"
			  `,
			APIUpdate: `{
				"reportSuiteIds": "id001, id003"
			  }`,
		},
	})
}
