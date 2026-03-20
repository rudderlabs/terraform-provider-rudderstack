package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceS3(t *testing.T) {
	cmt.AssertDestination(t, "s3", []c.TestConfig{
		{
			TerraformCreate: `
				bucket_name = "bucket"
				key_based_authentication {}
			`,
			APICreate: `{
				"bucketName": "bucket",
				"roleBasedAuth": false
			}`,
			TerraformUpdate: `
				bucket_name = "bucket"
				prefix      = "prefix"
				enable_sse  = true

				key_based_authentication {
					access_key_id = "my-key-id"
					access_key    = "my-secret"
				}
			`,
			APIUpdate: `{
				"bucketName": "bucket",
				"prefix": "prefix",
				"enableSSE": true,
				"roleBasedAuth": false,
				"accessKeyID": "my-key-id",
				"accessKey": "my-secret"
			}`,
		},
		{
			TerraformCreate: `
				bucket_name = "bucket"
				role_based_authentication {
					i_am_role_arn = "arn:aws:iam::123456789012:role/MyRole"
				}
			`,
			APICreate: `{
				"bucketName": "bucket",
				"roleBasedAuth": true,
				"iamRoleARN": "arn:aws:iam::123456789012:role/MyRole"
			}`,
			TerraformUpdate: `
				bucket_name = "bucket"
				prefix      = "prefix"
				enable_sse  = true
				role_based_authentication {
					i_am_role_arn = "arn:aws:iam::123456789012:role/MyRole"
				}
			`,
			APIUpdate: `{
				"bucketName": "bucket",
				"prefix": "prefix",
				"enableSSE": true,
				"roleBasedAuth": true,
				"iamRoleARN": "arn:aws:iam::123456789012:role/MyRole"
			}`,
		},
	})
}
