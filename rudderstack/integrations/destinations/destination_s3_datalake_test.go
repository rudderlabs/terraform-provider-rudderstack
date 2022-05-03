package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceS3Datalake(t *testing.T) {
	cmt.AssertDestination(t, "s3_datalake", []c.TestConfig{
		{
			TerraformCreate: `
				bucket_name = "bucket"
				use_glue    = false

				sync {
				  frequency = "30"
				}
			`,
			APICreate: `{
				"bucketName": "bucket",
				"useGlue": false,
				"syncFrequency": "30"
			}`,
			TerraformUpdate: `
				bucket_name = "bucket"

				prefix        = "prefix"
				access_key_id = "..."
				access_key    = "..."
			
				enable_sse = true
				use_glue   = true
				region     = "region"

				sync {
					frequency = "30"
					start_at  = "10:00"
				}
			`,
			APIUpdate: `{
				"bucketName": "bucket",
				"prefix": "prefix",
				"accessKeyID": "...",
				"accessKey": "...",
				"enableSSE": true,
				"useGlue": true,
				"region": "region",
				"syncFrequency": "30",
				"syncStartAt": "10:00"
			}`,
		},
	})
}
