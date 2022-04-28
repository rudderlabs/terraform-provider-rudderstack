package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourcePostgres(t *testing.T) {
	cmt.AssertDestination(t, "postgres", []c.TestConfig{
		{
			TerraformCreate: `
				host = "example.com"
				port = "5432"
				user = "example-user"
				password = "example-password"
				namespace = "example-namespace"
				ssl_mode = "disable"
				database = "example-database"
				use_rudder_storage = true
				s3 {
					bucket_name = "some-bucket-name"
					access_key_id = "some-access-key-id"
					access_key = "some-access-key"
				}
				`,
			APICreate: `{
				"host": "example.com",
				"port": "5432",
				"user": "example-user",
				"password": "example-password",
				"database": "example-database",
				"namespace": "example-namespace",
				"sslMode": "disable",
				"useRudderStorage": true,
				"bucketProvider": "S3",
				"bucketName": "some-bucket-name",
				"accessKeyID": "some-access-key-id",
				"accessKey": "some-access-key"
			}`,
			TerraformUpdate: `
				host = "example.com"
				port = "5432"
				user = "example-user-updated"
				password = "example-password"
				namespace = "example-namespace"
				ssl_mode = "disable"
				database = "example-database"
				use_rudder_storage = true
				gcs {
					bucket_name = "some-bucket-name"
					credentials = "some-credentials"
				}
				`,
			APIUpdate: `{
				"host": "example.com",
				"port": "5432",
				"user": "example-user-updated",
				"password": "example-password",
				"database": "example-database",
				"namespace": "example-namespace",
				"sslMode": "disable",
				"useRudderStorage": true,
				"bucketProvider": "GCS",
				"bucketName": "some-bucket-name",
				"credentials": "some-credentials"
			}`,
		},
	})
}
