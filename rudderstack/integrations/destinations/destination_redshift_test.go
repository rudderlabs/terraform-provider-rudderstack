package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceRedshift(t *testing.T) {
	cmt.AssertDestination(t, "redshift", []c.TestConfig{
		{
			TerraformCreate: `
				host = "example.com"
				port = "5432"
				user = "example-user"
				password = "example-password"
				namespace = "example-namespace"
				database = "example-database"
				use_rudder_storage = true

				sync {
					frequency = "30"
				}
			`,
			APICreate: `{
				"host": "example.com",
				"port": "5432",
				"user": "example-user",
				"password": "example-password",
				"database": "example-database",
				"namespace": "example-namespace",
				"useRudderStorage": true,
				"syncFrequency": "30"
			}`,
			TerraformUpdate: `
				host = "example.com"
				port = "5432"
				user = "example-user"
				password = "example-password"
				namespace = "example-namespace"
				enable_sse = true
				database = "example-database"
				use_rudder_storage = false

				sync {
					frequency = "30"
					start_at  = "10:00"
					exclude_window_start_time = "11:00"
					exclude_window_end_time   = "12:00"
				}

				s3 {
					bucket_name = "some-bucket-name"
					access_key_id = "some-access-key-id"
					access_key = "some-access-key"
				}
				onetrust_cookie_categories {
					web = ["one", "two", "three"]
					android = ["one", "two", "three"]
					ios = ["one", "two", "three"]
					unity = ["one", "two", "three"]
					reactnative = ["one", "two", "three"]
					flutter = ["one", "two", "three"]
					cordova = ["one", "two", "three"]
					amp = ["one", "two", "three"]
					cloud = ["one", "two", "three"]
					cloud_source = ["one", "two", "three"]
					shopify = ["one", "two", "three"]
				}
				`,
			APIUpdate: `{
				"host": "example.com",
				"port": "5432",
				"user": "example-user",
				"password": "example-password",
				"database": "example-database",
				"namespace": "example-namespace",
				"enableSSE": true,
				"useRudderStorage": false,
				"bucketName": "some-bucket-name",
				"accessKeyID": "some-access-key-id",
				"accessKey": "some-access-key",
				"syncFrequency": "30",
				"syncStartAt": "10:00",
				"excludeWindow": {
					"excludeWindowStartTime": "11:00",
					"excludeWindowEndTime": "12:00"	
				},
				"oneTrustCookieCategories": {
					"web": [
						{ "oneTrustCookieCategory": "one" },
						{ "oneTrustCookieCategory": "two" },
						{ "oneTrustCookieCategory": "three" }
					],
					"android": [
						{ "oneTrustCookieCategory": "one" },
						{ "oneTrustCookieCategory": "two" },
						{ "oneTrustCookieCategory": "three" }
					],
					"ios": [
						{ "oneTrustCookieCategory": "one" },
						{ "oneTrustCookieCategory": "two" },
						{ "oneTrustCookieCategory": "three" }
					],
					"unity": [
						{ "oneTrustCookieCategory": "one" },
						{ "oneTrustCookieCategory": "two" },
						{ "oneTrustCookieCategory": "three" }
					],
					"reactnative": [
						{ "oneTrustCookieCategory": "one" },
						{ "oneTrustCookieCategory": "two" },
						{ "oneTrustCookieCategory": "three" }
					],
					"flutter": [
						{ "oneTrustCookieCategory": "one" },
						{ "oneTrustCookieCategory": "two" },
						{ "oneTrustCookieCategory": "three" }
					],
					"cordova": [
						{ "oneTrustCookieCategory": "one" },
						{ "oneTrustCookieCategory": "two" },
						{ "oneTrustCookieCategory": "three" }
					],
					"amp": [
						{ "oneTrustCookieCategory": "one" },
						{ "oneTrustCookieCategory": "two" },
						{ "oneTrustCookieCategory": "three" }
					],
					"cloud": [
						{ "oneTrustCookieCategory": "one" },
						{ "oneTrustCookieCategory": "two" },
						{ "oneTrustCookieCategory": "three" }
					],
					"cloud_source": [
						{ "oneTrustCookieCategory": "one" },
						{ "oneTrustCookieCategory": "two" },
						{ "oneTrustCookieCategory": "three" }
					],
					"shopify": [
						{ "oneTrustCookieCategory": "one" },
						{ "oneTrustCookieCategory": "two" },
						{ "oneTrustCookieCategory": "three" }
					]
				}
			}`,
		},
	})
}
