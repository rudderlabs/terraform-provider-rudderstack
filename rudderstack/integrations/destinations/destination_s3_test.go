package destinations_test

import (
	"testing"

	"github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceS3(t *testing.T) {
	testutil.AssertDestination(t, "s3", []c.TestConfig{
		{
			TerraformCreate: `
				bucket_name = "bucket"
			`,
			APICreate: `{
				"bucketName": "bucket"
			}`,
			TerraformUpdate: `
				bucket_name = "bucket"

				prefix        = "prefix"
				access_key_id = "..."
				access_key    = "..."
			
				enable_sse    = true
			`,
			APIUpdate: `{
				"bucketName": "bucket",
				"prefix": "prefix",
				"accessKeyID": "...",
				"accessKey": "...",
				"enableSSE": true
			}`,
		},
	})
}
