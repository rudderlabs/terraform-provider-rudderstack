package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceSnowflake(t *testing.T) {
	cmt.AssertDestination(t, "snowflake", []c.TestConfig{
		{
			TerraformCreate: `
				account = "..."
				database = "..."
				warehouse = "..."
				user = "..."
				password = "..."
				use_rudder_storage = true
				sync {
					frequency = "30"
				}
			`,
			APICreate: `{
				"account": "...",
				"database": "...",
				"warehouse": "...",
				"user": "...",
				"password": "...",
				"syncFrequency": "30",
				"useRudderStorage": true,
				"additionalProperties": true
			}`,
			TerraformUpdate: `
				account = "..."
				database = "..."
				warehouse = "..."
				user = "..."
				password = "..."
				namespace = "..."
				sync {
					frequency = "60"
					start_at                  = "10:00"
					exclude_window_start_time = "11:00"
					exclude_window_end_time   = "12:00"
				}
				json_paths = "..."
				use_rudder_storage = false
				cloud_provider = "AWS"
				s3 {
					bucket_name = "..."
					access_key_id = "..."
					access_key = "..."
					enable_sse = true
				}
				prefix = "..."
			`,
			APIUpdate: `{
				"account": "...",
				"database": "...",
				"warehouse": "...",
				"user": "...",
				"password": "...",
				"namespace": "...",
				"syncFrequency": "60",
				"useRudderStorage": false,
				"additionalProperties": true,
				"syncStartAt": "10:00",
				"excludeWindow": {
						"excludeWindowStartTime": "11:00",
						"excludeWindowEndTime": "12:00"
				},
				"jsonPaths": "...",
				"cloudProvider": "AWS",
				"bucketName": "...",
				"accessKeyID": "...",
				"accessKey": "...",
				"enableSSE": true,
				"prefix": "..."
		}`,
		},
	})
}
