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
				ssl_mode = "test-ssl_mode"
				sync_frequency = "test-sync_frequency"
				use_rudder_storage = true
			`,
			APICreate: `{
				"host": "test-host",
				"database": "test-database",
				"user": "test-user",
				"password": "test-password",
				"port": "test-port",
				"sslMode": "test-ssl_mode",
				"syncFrequency": "test-sync_frequency",
				"useRudderStorage": true
			}`,
			TerraformUpdate: `
				host = "test-host"
				database = "test-database"
				user = "test-user"
				password = "test-password"
				port = "test-port"
				ssl_mode = "test-ssl_mode"
				sync_frequency = "test-sync_frequency"
				use_rudder_storage = true
				exclude_window {
					exclude_window_start_time = "test-exclude_window_start_time"
					exclude_window_end_time = "test-exclude_window_end_time"
				}
				onetrust_cookie_categories = ["c001"]
			`,
			APIUpdate: `{
				"host": "test-host",
				"database": "test-database",
				"user": "test-user",
				"password": "test-password",
				"port": "test-port",
				"sslMode": "test-ssl_mode",
				"syncFrequency": "test-sync_frequency",
				"useRudderStorage": true,
				"excludeWindow": {
					"excludeWindowStartTime": "test-exclude_window_start_time",
					"excludeWindowEndTime": "test-exclude_window_end_time"
				},
				"oneTrustCookieCategories": ["c001"]
			}`,
		},
	})
}
