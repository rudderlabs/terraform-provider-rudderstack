package destinations_test

import (
	"testing"

	"github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceBigQuery(t *testing.T) {
	testutil.AssertDestination(t, "bigquery", []c.TestConfig{
		{
			TerraformCreate: `
				project     = "project"
				bucket_name = "bucket"
				credentials = "..."
						
				sync {
					frequency = "30"
				}
			`,
			APICreate: `{
				"project": "project",
				"bucketName": "bucket",
				"credentials": "...",
				"syncFrequency": "30"
			}`,
			TerraformUpdate: `
				project     = "project"
				bucket_name = "bucket"
				credentials = "..."
			
				location  = "us-east1"
				prefix    = "prefix"
				namespace = "namespace"
			
				sync {
					frequency				  = "30"
					start_at   				  = "10:00"
					exclude_window_start_time = "11:00"
					exclude_window_end_time   = "12:00"
				}
			`,
			APIUpdate: `{
				"project": "project",
				"bucketName": "bucket",
				"credentials": "...",
				"location": "us-east1",
				"prefix": "prefix",
				"namespace": "namespace",
				"syncFrequency": "30",
				"syncStartAt": "10:00",
				"excludeWindow": {
					"excludeWindowStartTime": "11:00",
					"excludeWindowEndTime": "12:00"
				}
			}`,
		},
	})
}
