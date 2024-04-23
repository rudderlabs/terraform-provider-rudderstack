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
				account = "example-account"
				database = "example-database"
				warehouse = "example-warehouse"
				user = "example-user"
				password = "example-password"
				use_rudder_storage = true
				sync {
					frequency = "30"
				}
			`,
			APICreate: `{
				"account": "example-account",
				"database": "example-database",
				"warehouse": "example-warehouse",
				"user": "example-user",
				"password": "example-password",
				"syncFrequency": "30",
				"useRudderStorage": true,
				"additionalProperties": true
			}`,
			TerraformUpdate: `
				account = "example-account"
				database = "example-database"
				warehouse = "example-warehouse"
				user = "example-user"
				password = "example-password"
				use_rudder_storage = false
				sync {
					frequency = "60"
					start_at                  = "10:00"
					exclude_window_start_time = "11:00"
					exclude_window_end_time   = "12:00"
				}
				namespace = "example-namespace"
				json_paths = "./example-paths"
				prefix = "example-prefix"
				s3 {
					bucket_name = "example-bucket-name"
					access_key_id = "example-access-key-id"
					access_key = "example-access-key"
					enable_sse = true
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
					warehouse = ["one", "two", "three"]
					shopify = ["one", "two", "three"]
				}
			`,
			APIUpdate: `{
				"account": "example-account",
				"database": "example-database",
				"warehouse": "example-warehouse",
				"user": "example-user",
				"password": "example-password",
				"namespace": "example-namespace",
				"syncFrequency": "60",
				"syncStartAt": "10:00",
				"excludeWindow": {
					"excludeWindowStartTime": "11:00",
					"excludeWindowEndTime": "12:00"
				},
				"useRudderStorage": false,
				"additionalProperties": true,
				"jsonPaths": "./example-paths",
				"cloudProvider": "AWS",
				"prefix": "example-prefix",
        		"bucketName": "example-bucket-name",
        		"accessKeyID": "example-access-key-id",
        		"accessKey": "example-access-key",
        		"enableSSE": true,
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
					"warehouse": [
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

func TestDestinationResourceSnowflakeWithGCP(t *testing.T) {
	cmt.AssertDestination(t, "snowflake", []c.TestConfig{
		{
			TerraformCreate: `
				account = "example-account"
				database = "example-database"
				warehouse = "example-warehouse"
				user = "example-user"
				password = "example-password"
				use_rudder_storage = true
				sync {
					frequency = "30"
				}
			`,
			APICreate: `{
				"account": "example-account",
				"database": "example-database",
				"warehouse": "example-warehouse",
				"user": "example-user",
				"password": "example-password",
				"syncFrequency": "30",
				"useRudderStorage": true,
				"additionalProperties": true
			}`,
			TerraformUpdate: `
				account = "example-account"
				database = "example-database"
				warehouse = "example-warehouse"
				user = "example-user"
				password = "example-password"
				use_rudder_storage = false
				sync {
					frequency = "60"
					start_at                  = "10:00"
					exclude_window_start_time = "11:00"
					exclude_window_end_time   = "12:00"
				}
				namespace = "example-namespace"
				json_paths = "./example-paths"
				prefix = "example-prefix"
				gcp {
					bucket_name = "example-bucket-name"
					credentials = "example-credentials"
					storage_integration = "example-storage"      
				}
			`,
			APIUpdate: `{
				"account": "example-account",
				"database": "example-database",
				"warehouse": "example-warehouse",
				"user": "example-user",
				"password": "example-password",
				"namespace": "example-namespace",
				"syncFrequency": "60",
				"syncStartAt": "10:00",
				"excludeWindow": {
					"excludeWindowStartTime": "11:00",
					"excludeWindowEndTime": "12:00"
				},
				"useRudderStorage": false,
				"additionalProperties": true,
				"jsonPaths": "./example-paths",
				"cloudProvider": "GCP",
				"prefix": "example-prefix",
        "bucketName": "example-bucket-name",
				"credentials": "example-credentials",
				"storageIntegration": "example-storage"
			}`,
		},
	})
}

func TestDestinationResourceSnowflakeWithAzure(t *testing.T) {
	cmt.AssertDestination(t, "snowflake", []c.TestConfig{
		{
			TerraformCreate: `
				account = "example-account"
				database = "example-database"
				warehouse = "example-warehouse"
				user = "example-user"
				password = "example-password"
				use_rudder_storage = true
				sync {
					frequency = "30"
				}
			`,
			APICreate: `{
				"account": "example-account",
				"database": "example-database",
				"warehouse": "example-warehouse",
				"user": "example-user",
				"password": "example-password",
				"syncFrequency": "30",
				"useRudderStorage": true,
				"additionalProperties": true
			}`,
			TerraformUpdate: `
				account = "example-account"
				database = "example-database"
				warehouse = "example-warehouse"
				user = "example-user"
				password = "example-password"
				use_rudder_storage = false
				sync {
					frequency = "60"
					start_at                  = "10:00"
					exclude_window_start_time = "11:00"
					exclude_window_end_time   = "12:00"
				}
				namespace = "example-namespace"
				json_paths = "./example-paths"
				prefix = "example-prefix"
				azure {
					container_name = "example-container-name"
					account_name = "example-account-name"
					account_key = "example-account-key"
					storage_integration = "example-storage" 
				}
			`,
			APIUpdate: `{
				"account": "example-account",
				"database": "example-database",
				"warehouse": "example-warehouse",
				"user": "example-user",
				"password": "example-password",
				"namespace": "example-namespace",
				"syncFrequency": "60",
				"syncStartAt": "10:00",
				"excludeWindow": {
					"excludeWindowStartTime": "11:00",
					"excludeWindowEndTime": "12:00"
				},
				"useRudderStorage": false,
				"additionalProperties": true,
				"jsonPaths": "./example-paths",
				"cloudProvider": "AZURE",
				"containerName": "example-container-name",
				"accountName": "example-account-name",
				"accountKey": "example-account-key",
				"storageIntegration": "example-storage",
				"prefix": "example-prefix"
			}`,
		},
	})
}
