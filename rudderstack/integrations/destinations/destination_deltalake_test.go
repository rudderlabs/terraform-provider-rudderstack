package destinations_test

import (
	"testing"

	acc "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/acc"
	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

var deltalakeTestConfigs = []c.TestConfig{
	{
		TerraformCreate: `
				host = "example-host.cloud.databricks.com"
				port = "443"
				path = "/sql/1.0/endpoints/example-endpoint"
				token = "example-token"
				bucket_provider = "S3"
				sync {
					frequency = "30"
				}
				s3 {
					bucket_name = "example-bucket-name"
					access_key_id = "example-access-key-id"
					access_key = "example-access-key"
				}
			`,
		APICreate: `{
				"host": "example-host.cloud.databricks.com",
				"port": "443",
				"path": "/sql/1.0/endpoints/example-endpoint",
				"token": "example-token",
				"useOauth": false,
				"syncFrequency": "30",
				"skipTracksTable": false,
				"skipUsersTable": true,
				"preferAppend": true,
				"underscoreDivideNumbers": false,
				"allowUsersContextTraits": false,
				"enableExternalLocation": false,
				"useRudderStorage": false,
				"cleanupObjectStorageFiles": false,
				"bucketProvider": "S3",
				"bucketName": "example-bucket-name",
				"enableSSE": false,
				"useSTSTokens": false,
				"accessKeyID": "example-access-key-id",
				"accessKey": "example-access-key",
				"roleBasedAuth": false
			}`,
		TerraformUpdate: `
				host = "example-host.cloud.databricks.com"
				port = "443"
				path = "/sql/1.0/endpoints/example-endpoint"
				token = "example-token"
				catalog = "example-catalog"
				namespace = "example-namespace"
				prefix = "example-prefix"
				bucket_provider = "S3"
				sync {
					frequency = "60"
					start_at                  = "10:00"
					exclude_window_start_time = "11:00"
					exclude_window_end_time   = "12:00"
				}
				skip_tracks_table = true
				skip_users_table = false
				prefer_append = false
				cleanup_object_storage_files = true
				s3 {
					bucket_name = "example-bucket-name"
					access_key_id = "example-access-key-id"
					access_key = "example-access-key"
					enable_sse = true
					use_sts_tokens = true
				}
				connection_mode {
					web = "cloud"
					android = "cloud"
					android_kotlin = "cloud"
					ios = "cloud"
					ios_swift = "cloud"
					unity = "cloud"
					amp = "cloud"
					cloud = "cloud"
					cloud_source = "cloud"
					reactnative = "cloud"
					flutter = "cloud"
					cordova = "cloud"
					shopify = "cloud"
				}
				consent_management {
					web = [
						{
							provider = "oneTrust"
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
						resolution_strategy = "or"
						consents = ["one_ios", "two_ios", "three_ios"]
					}]
				}
			`,
		APIUpdate: `{
				"host": "example-host.cloud.databricks.com",
				"port": "443",
				"path": "/sql/1.0/endpoints/example-endpoint",
				"token": "example-token",
				"useOauth": false,
				"catalog": "example-catalog",
				"namespace": "example-namespace",
				"prefix": "example-prefix",
				"syncFrequency": "60",
				"syncStartAt": "10:00",
				"excludeWindow": {
					"excludeWindowStartTime": "11:00",
					"excludeWindowEndTime": "12:00"
				},
				"skipTracksTable": true,
				"skipUsersTable": false,
				"preferAppend": false,
				"underscoreDivideNumbers": false,
				"allowUsersContextTraits": false,
				"enableExternalLocation": false,
				"useRudderStorage": false,
				"cleanupObjectStorageFiles": true,
				"bucketProvider": "S3",
				"bucketName": "example-bucket-name",
				"enableSSE": true,
				"useSTSTokens": true,
				"accessKeyID": "example-access-key-id",
				"accessKey": "example-access-key",
				"roleBasedAuth": false,
				"connectionMode": {
					"web": "cloud",
					"android": "cloud",
					"androidKotlin": "cloud",
					"ios": "cloud",
					"iosSwift": "cloud",
					"unity": "cloud",
					"amp": "cloud",
					"cloud": "cloud",
					"cloudSource": "cloud",
					"reactnative": "cloud",
					"flutter": "cloud",
					"cordova": "cloud",
					"shopify": "cloud"
				},
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
							"resolutionStrategy": "or",
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
					]
				}
			}`,
	},
}

func TestDestinationResourceDeltalake(t *testing.T) {
	cmt.AssertDestination(t, "deltalake", deltalakeTestConfigs)
}

func TestAccDestinationDeltalake(t *testing.T) {
	acc.AccAssertDestination(t, "deltalake", deltalakeTestConfigs)
}

// S3 staging with IAM role-based authentication instead of access keys.
func TestDestinationResourceDeltalakeWithRoleBasedAuth(t *testing.T) {
	cmt.AssertDestination(t, "deltalake", []c.TestConfig{
		{
			TerraformCreate: `
				host = "example-host.cloud.databricks.com"
				port = "443"
				path = "/sql/1.0/endpoints/example-endpoint"
				token = "example-token"
				bucket_provider = "S3"
				sync {
					frequency = "30"
				}
				s3 {
					bucket_name = "example-bucket-name"
					use_sts_tokens = true
					role_based_authentication {
						i_am_role_arn = "arn:aws:iam::123456789012:role/example"
					}
				}
			`,
			APICreate: `{
				"host": "example-host.cloud.databricks.com",
				"port": "443",
				"path": "/sql/1.0/endpoints/example-endpoint",
				"token": "example-token",
				"useOauth": false,
				"syncFrequency": "30",
				"skipTracksTable": false,
				"skipUsersTable": true,
				"preferAppend": true,
				"underscoreDivideNumbers": false,
				"allowUsersContextTraits": false,
				"enableExternalLocation": false,
				"useRudderStorage": false,
				"cleanupObjectStorageFiles": false,
				"bucketProvider": "S3",
				"bucketName": "example-bucket-name",
				"enableSSE": false,
				"useSTSTokens": true,
				"iamRoleARN": "arn:aws:iam::123456789012:role/example",
				"roleBasedAuth": true
			}`,
			TerraformUpdate: `
				host = "example-host.cloud.databricks.com"
				port = "443"
				path = "/sql/1.0/endpoints/example-endpoint"
				token = "example-token"
				bucket_provider = "S3"
				sync {
					frequency = "30"
				}
				s3 {
					bucket_name = "example-bucket-name"
					use_sts_tokens = true
					enable_sse = true
					role_based_authentication {
						i_am_role_arn = "arn:aws:iam::123456789012:role/example-updated"
					}
				}
			`,
			APIUpdate: `{
				"host": "example-host.cloud.databricks.com",
				"port": "443",
				"path": "/sql/1.0/endpoints/example-endpoint",
				"token": "example-token",
				"useOauth": false,
				"syncFrequency": "30",
				"skipTracksTable": false,
				"skipUsersTable": true,
				"preferAppend": true,
				"underscoreDivideNumbers": false,
				"allowUsersContextTraits": false,
				"enableExternalLocation": false,
				"useRudderStorage": false,
				"cleanupObjectStorageFiles": false,
				"bucketProvider": "S3",
				"bucketName": "example-bucket-name",
				"enableSSE": true,
				"useSTSTokens": true,
				"iamRoleARN": "arn:aws:iam::123456789012:role/example-updated",
				"roleBasedAuth": true
			}`,
		},
	})
}

func TestDestinationResourceDeltalakeWithGCS(t *testing.T) {
	cmt.AssertDestination(t, "deltalake", []c.TestConfig{
		{
			TerraformCreate: `
				host = "example-host.cloud.databricks.com"
				port = "443"
				path = "/sql/1.0/endpoints/example-endpoint"
				token = "example-token"
				bucket_provider = "GCS"
				sync {
					frequency = "180"
				}
				gcs {
					bucket_name = "example-bucket-name"
					credentials = "example-credentials"
				}
			`,
			APICreate: `{
				"host": "example-host.cloud.databricks.com",
				"port": "443",
				"path": "/sql/1.0/endpoints/example-endpoint",
				"token": "example-token",
				"useOauth": false,
				"syncFrequency": "180",
				"skipTracksTable": false,
				"skipUsersTable": true,
				"preferAppend": true,
				"underscoreDivideNumbers": false,
				"allowUsersContextTraits": false,
				"enableExternalLocation": false,
				"useRudderStorage": false,
				"cleanupObjectStorageFiles": false,
				"bucketProvider": "GCS",
				"bucketName": "example-bucket-name",
				"credentials": "example-credentials"
			}`,
			TerraformUpdate: `
				host = "example-host.cloud.databricks.com"
				port = "443"
				path = "/sql/1.0/endpoints/example-endpoint"
				token = "example-token"
				bucket_provider = "GCS"
				sync {
					frequency = "360"
				}
				gcs {
					bucket_name = "example-bucket-name-updated"
					credentials = "example-credentials"
				}
			`,
			APIUpdate: `{
				"host": "example-host.cloud.databricks.com",
				"port": "443",
				"path": "/sql/1.0/endpoints/example-endpoint",
				"token": "example-token",
				"useOauth": false,
				"syncFrequency": "360",
				"skipTracksTable": false,
				"skipUsersTable": true,
				"preferAppend": true,
				"underscoreDivideNumbers": false,
				"allowUsersContextTraits": false,
				"enableExternalLocation": false,
				"useRudderStorage": false,
				"cleanupObjectStorageFiles": false,
				"bucketProvider": "GCS",
				"bucketName": "example-bucket-name-updated",
				"credentials": "example-credentials"
			}`,
		},
	})
}

// Azure Blob staging, moving from account-key to SAS-token authentication.
func TestDestinationResourceDeltalakeWithAzureBlob(t *testing.T) {
	cmt.AssertDestination(t, "deltalake", []c.TestConfig{
		{
			TerraformCreate: `
				host = "example-host.cloud.databricks.com"
				port = "443"
				path = "/sql/1.0/endpoints/example-endpoint"
				token = "example-token"
				bucket_provider = "AZURE_BLOB"
				sync {
					frequency = "30"
				}
				azure_blob {
					container_name = "example-container"
					account_name = "example-account-name"
					account_key = "example-account-key"
				}
			`,
			APICreate: `{
				"host": "example-host.cloud.databricks.com",
				"port": "443",
				"path": "/sql/1.0/endpoints/example-endpoint",
				"token": "example-token",
				"useOauth": false,
				"syncFrequency": "30",
				"skipTracksTable": false,
				"skipUsersTable": true,
				"preferAppend": true,
				"underscoreDivideNumbers": false,
				"allowUsersContextTraits": false,
				"enableExternalLocation": false,
				"useRudderStorage": false,
				"cleanupObjectStorageFiles": false,
				"bucketProvider": "AZURE_BLOB",
				"containerName": "example-container",
				"accountName": "example-account-name",
				"accountKey": "example-account-key",
				"useSASTokens": false,
				"enableHierarchicalNamespace": false
			}`,
			TerraformUpdate: `
				host = "example-host.cloud.databricks.com"
				port = "443"
				path = "/sql/1.0/endpoints/example-endpoint"
				token = "example-token"
				bucket_provider = "AZURE_BLOB"
				sync {
					frequency = "30"
				}
				azure_blob {
					container_name = "example-container"
					account_name = "example-account-name"
					use_sas_tokens = true
					sas_token = "example-sas-token"
					enable_hierarchical_namespace = true
				}
			`,
			APIUpdate: `{
				"host": "example-host.cloud.databricks.com",
				"port": "443",
				"path": "/sql/1.0/endpoints/example-endpoint",
				"token": "example-token",
				"useOauth": false,
				"syncFrequency": "30",
				"skipTracksTable": false,
				"skipUsersTable": true,
				"preferAppend": true,
				"underscoreDivideNumbers": false,
				"allowUsersContextTraits": false,
				"enableExternalLocation": false,
				"useRudderStorage": false,
				"cleanupObjectStorageFiles": false,
				"bucketProvider": "AZURE_BLOB",
				"containerName": "example-container",
				"accountName": "example-account-name",
				"useSASTokens": true,
				"sasToken": "example-sas-token",
				"enableHierarchicalNamespace": true
			}`,
		},
	})
}

// M2M OAuth instead of a personal access token.
func TestDestinationResourceDeltalakeWithOauth(t *testing.T) {
	cmt.AssertDestination(t, "deltalake", []c.TestConfig{
		{
			TerraformCreate: `
				host = "example-host.cloud.databricks.com"
				port = "443"
				path = "/sql/1.0/endpoints/example-endpoint"
				use_oauth = true
				oauth_client_id = "example-client-id"
				oauth_client_secret = "example-client-secret"
				bucket_provider = "S3"
				use_rudder_storage = true
				sync {
					frequency = "30"
				}
			`,
			APICreate: `{
				"host": "example-host.cloud.databricks.com",
				"port": "443",
				"path": "/sql/1.0/endpoints/example-endpoint",
				"useOauth": true,
				"oauthClientID": "example-client-id",
				"oauthClientSecret": "example-client-secret",
				"syncFrequency": "30",
				"skipTracksTable": false,
				"skipUsersTable": true,
				"preferAppend": true,
				"underscoreDivideNumbers": false,
				"allowUsersContextTraits": false,
				"enableExternalLocation": false,
				"useRudderStorage": true,
				"cleanupObjectStorageFiles": false,
				"bucketProvider": "S3"
			}`,
			TerraformUpdate: `
				host = "example-host.cloud.databricks.com"
				port = "443"
				path = "/sql/1.0/endpoints/example-endpoint"
				use_oauth = true
				oauth_client_id = "example-client-id-updated"
				oauth_client_secret = "example-client-secret"
				bucket_provider = "S3"
				use_rudder_storage = true
				sync {
					frequency = "30"
				}
			`,
			APIUpdate: `{
				"host": "example-host.cloud.databricks.com",
				"port": "443",
				"path": "/sql/1.0/endpoints/example-endpoint",
				"useOauth": true,
				"oauthClientID": "example-client-id-updated",
				"oauthClientSecret": "example-client-secret",
				"syncFrequency": "30",
				"skipTracksTable": false,
				"skipUsersTable": true,
				"preferAppend": true,
				"underscoreDivideNumbers": false,
				"allowUsersContextTraits": false,
				"enableExternalLocation": false,
				"useRudderStorage": true,
				"cleanupObjectStorageFiles": false,
				"bucketProvider": "S3"
			}`,
		},
	})
}

// enable_external_location is a boolean on the wire: rudder-server reads the
// key with GetBoolDestinationConfig and the destination definition types it as
// a boolean. The definition's conditional branch that requires
// externalLocation compares against the string "true", so it never matches a
// boolean and the pairing is not enforced server-side.
func TestDestinationResourceDeltalakeWithExternalLocation(t *testing.T) {
	cmt.AssertDestination(t, "deltalake", []c.TestConfig{
		{
			TerraformCreate: `
				host = "example-host.cloud.databricks.com"
				port = "443"
				path = "/sql/1.0/endpoints/example-endpoint"
				token = "example-token"
				bucket_provider = "S3"
				enable_external_location = true
				external_location = "/mnt/example/delta"
				use_rudder_storage = true
				sync {
					frequency = "30"
				}
			`,
			APICreate: `{
				"host": "example-host.cloud.databricks.com",
				"port": "443",
				"path": "/sql/1.0/endpoints/example-endpoint",
				"token": "example-token",
				"useOauth": false,
				"syncFrequency": "30",
				"skipTracksTable": false,
				"skipUsersTable": true,
				"preferAppend": true,
				"underscoreDivideNumbers": false,
				"allowUsersContextTraits": false,
				"enableExternalLocation": true,
				"externalLocation": "/mnt/example/delta",
				"useRudderStorage": true,
				"cleanupObjectStorageFiles": false,
				"bucketProvider": "S3"
			}`,
			TerraformUpdate: `
				host = "example-host.cloud.databricks.com"
				port = "443"
				path = "/sql/1.0/endpoints/example-endpoint"
				token = "example-token"
				bucket_provider = "S3"
				enable_external_location = true
				external_location = "/mnt/example/delta-updated"
				use_rudder_storage = true
				sync {
					frequency = "30"
				}
			`,
			APIUpdate: `{
				"host": "example-host.cloud.databricks.com",
				"port": "443",
				"path": "/sql/1.0/endpoints/example-endpoint",
				"token": "example-token",
				"useOauth": false,
				"syncFrequency": "30",
				"skipTracksTable": false,
				"skipUsersTable": true,
				"preferAppend": true,
				"underscoreDivideNumbers": false,
				"allowUsersContextTraits": false,
				"enableExternalLocation": true,
				"externalLocation": "/mnt/example/delta-updated",
				"useRudderStorage": true,
				"cleanupObjectStorageFiles": false,
				"bucketProvider": "S3"
			}`,
		},
	})
}

// The two rs-immutable flags must round-trip when set explicitly, and both
// booleans must be sent even when false.
func TestDestinationResourceDeltalakeWithImmutableFlags(t *testing.T) {
	cmt.AssertDestination(t, "deltalake", []c.TestConfig{
		{
			TerraformCreate: `
				host = "example-host.cloud.databricks.com"
				port = "443"
				path = "/sql/1.0/endpoints/example-endpoint"
				token = "example-token"
				bucket_provider = "S3"
				use_rudder_storage = true
				underscore_divide_numbers = true
				allow_users_context_traits = true
				sync {
					frequency = "30"
				}
			`,
			APICreate: `{
				"host": "example-host.cloud.databricks.com",
				"port": "443",
				"path": "/sql/1.0/endpoints/example-endpoint",
				"token": "example-token",
				"useOauth": false,
				"syncFrequency": "30",
				"skipTracksTable": false,
				"skipUsersTable": true,
				"preferAppend": true,
				"underscoreDivideNumbers": true,
				"allowUsersContextTraits": true,
				"enableExternalLocation": false,
				"useRudderStorage": true,
				"cleanupObjectStorageFiles": false,
				"bucketProvider": "S3"
			}`,
			TerraformUpdate: `
				host = "example-host.cloud.databricks.com"
				port = "443"
				path = "/sql/1.0/endpoints/example-endpoint"
				token = "example-token-updated"
				bucket_provider = "S3"
				use_rudder_storage = true
				underscore_divide_numbers = true
				allow_users_context_traits = true
				sync {
					frequency = "30"
				}
			`,
			APIUpdate: `{
				"host": "example-host.cloud.databricks.com",
				"port": "443",
				"path": "/sql/1.0/endpoints/example-endpoint",
				"token": "example-token-updated",
				"useOauth": false,
				"syncFrequency": "30",
				"skipTracksTable": false,
				"skipUsersTable": true,
				"preferAppend": true,
				"underscoreDivideNumbers": true,
				"allowUsersContextTraits": true,
				"enableExternalLocation": false,
				"useRudderStorage": true,
				"cleanupObjectStorageFiles": false,
				"bucketProvider": "S3"
			}`,
		},
	})
}
