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
				role = "example-role"
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
				consent_management {
					web = [
						{
							provider = "oneTrust"
							consents = ["one_web", "two_web", "three_web"]
							resolution_strategy = ""
						},
						{
							provider = "ketch"
							consents = ["one_web", "two_web", "three_web"]
							resolution_strategy = ""
						},
						{
							provider = "custom"
							resolution_strategy = "and"
							consents = ["one_web", "two_web", "three_web"]
						}
					]
					android = [{
						provider = "ketch"
						consents = ["one_android", "two_android", "three_android"]
						resolution_strategy = ""
					}]
					ios = [{
						provider = "custom"
						resolution_strategy = "and"
						consents = ["one_ios", "two_ios", "three_ios"]
					}]
					unity = [{
						provider = "custom"
						resolution_strategy = "or"
						consents = ["one_unity", "two_unity", "three_unity"]
					}]
					reactnative = [{
						provider = "custom"
						resolution_strategy = "and"
						consents = ["one_reactnative", "two_reactnative", "three_reactnative"]
					}]
					flutter = [{
						provider = "custom"
						resolution_strategy = "and"
						consents = ["one_flutter", "two_flutter", "three_flutter"]
					}]
					cordova = [{
						provider = "custom"
						resolution_strategy = "and"
						consents = ["one_cordova", "two_cordova", "three_cordova"]
					}]
					amp = [{
						provider = "custom"
						resolution_strategy = "and"
						consents = ["one_amp", "two_amp", "three_amp"]
					}]
					cloud = [{
						provider = "custom"
						resolution_strategy = "and"
						consents = ["one_cloud", "two_cloud", "three_cloud"]
					}]
					cloud_source = [{
						provider = "custom"
						resolution_strategy = "and"
						consents = ["one_cloud_source", "two_cloud_source", "three_cloud_source"]
					}]
					shopify = [{
						provider = "custom"
						resolution_strategy = "and"
						consents = ["one_shopify", "two_shopify", "three_shopify"]
					}]
				}
			`,
			APIUpdate: `{
				"account": "example-account",
				"database": "example-database",
				"warehouse": "example-warehouse",
				"user": "example-user",
				"password": "example-password",
				"role": "example-role",
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
				"consentManagement": {
					"web": [
						{
							"provider": "oneTrust",
							"resolutionStrategy": "",
							"consents": [
								{
									"consent": "one_web"
								},
								{
									"consent": "two_web"
								},
								{
									"consent": "three_web"
								}
							]
						},
						{
							"provider": "ketch",
							"resolutionStrategy": "",
							"consents": [
								{
									"consent": "one_web"
								},
								{
									"consent": "two_web"
								},
								{
									"consent": "three_web"
								}
							]
						},
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{
									"consent": "one_web"
								},
								{
									"consent": "two_web"
								},
								{
									"consent": "three_web"
								}
							]
						}
					],
					"android": [
						{
							"provider": "ketch",
							"resolutionStrategy": "",
							"consents": [
								{
									"consent": "one_android"
								},
								{
									"consent": "two_android"
								},
								{
									"consent": "three_android"
								}
							]
						}
					],
					"ios": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{
									"consent": "one_ios"
								},
								{
									"consent": "two_ios"
								},
								{
									"consent": "three_ios"
								}
							]
						}
					],
					"unity": [
						{
							"provider": "custom",
							"resolutionStrategy": "or",
							"consents": [
								{
									"consent": "one_unity"
								},
								{
									"consent": "two_unity"
								},
								{
									"consent": "three_unity"
								}
							]
						}
					],
					"reactnative": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{
									"consent": "one_reactnative"
								},
								{
									"consent": "two_reactnative"
								},
								{
									"consent": "three_reactnative"
								}
							]
						}
					],
					"flutter": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{
									"consent": "one_flutter"
								},
								{
									"consent": "two_flutter"
								},
								{
									"consent": "three_flutter"
								}
							]
						}
					],
					"cordova": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{
									"consent": "one_cordova"
								},
								{
									"consent": "two_cordova"
								},
								{
									"consent": "three_cordova"
								}
							]
						}
					],
					"amp": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{
									"consent": "one_amp"
								},
								{
									"consent": "two_amp"
								},
								{
									"consent": "three_amp"
								}
							]
						}
					],
					"cloud": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{
									"consent": "one_cloud"
								},
								{
									"consent": "two_cloud"
								},
								{
									"consent": "three_cloud"
								}
							]
						}
					],
					"cloudSource": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{
									"consent": "one_cloud_source"
								},
								{
									"consent": "two_cloud_source"
								},
								{
									"consent": "three_cloud_source"
								}
							]
						}
					],
					"shopify": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{
									"consent": "one_shopify"
								},
								{
									"consent": "two_shopify"
								},
								{
									"consent": "three_shopify"
								}
							]
						}
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
