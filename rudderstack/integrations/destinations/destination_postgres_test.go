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
				host = "test-host"
				database = "test-database"
				user = "test-user"
				password = "test-password"
				port = "test-port"
				ssl_mode = "disable"
				sync_frequency = "30"
				use_rudder_storage = true
			`,
			APICreate: `{
				"host": "test-host",
				"database": "test-database",
				"user": "test-user",
				"password": "test-password",
				"port": "test-port",
				"sslMode": "disable",
				"syncFrequency": "30",
				"useRudderStorage": true
			}`,
			TerraformUpdate: `
				host = "test-host"
				database = "test-database"
				user = "test-user"
				password = "test-password"
				port = "test-port"
				ssl_mode = "verify-ca"
				sync_frequency = "60"
				use_rudder_storage = true
			`,
			APIUpdate: `{
				"host": "test-host",
				"database": "test-database",
				"user": "test-user",
				"password": "test-password",
				"port": "test-port",
				"sslMode": "verify-ca",
				"syncFrequency": "60",
				"useRudderStorage": true
			}`,
		},
	})
}
