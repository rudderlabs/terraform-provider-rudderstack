package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceSqlServer(t *testing.T) {
	cmt.AssertDestination(t, "sql_server", []c.TestConfig{
		{
			TerraformCreate: `
				host               = "mssql.example.com"
				database           = "mydb"
				user               = "myuser"
				password           = "mypassword"
				use_rudder_storage = true
			`,
			APICreate: `{
				"host": "mssql.example.com",
				"database": "mydb",
				"user": "myuser",
				"password": "mypassword",
				"port": "1433",
				"sslMode": "disable",
				"syncFrequency": "30",
				"useRudderStorage": true
			}`,
			TerraformUpdate: `
				host               = "mssql.example.com"
				database           = "mydb"
				user               = "myuser"
				password           = "mypassword"
				port               = "1433"
				namespace          = "myschema"
				ssl_mode           = "true"
				sync_frequency     = "60"
				sync_start_at      = "09:00"
				use_rudder_storage = false
				bucket_provider    = "S3"
				bucket_name        = "my-bucket"
				access_key_id      = "my-key-id"
				access_key         = "my-secret"
			`,
			APIUpdate: `{
				"host": "mssql.example.com",
				"database": "mydb",
				"user": "myuser",
				"password": "mypassword",
				"port": "1433",
				"namespace": "myschema",
				"sslMode": "true",
				"syncFrequency": "60",
				"syncStartAt": "09:00",
				"useRudderStorage": false,
				"bucketProvider": "S3",
				"bucketName": "my-bucket",
				"accessKeyID": "my-key-id",
				"accessKey": "my-secret"
			}`,
		},
	})
}
