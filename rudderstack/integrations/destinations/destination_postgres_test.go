package destinations_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	acc "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/acc"
	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

var postgresTestConfigs = []c.TestConfig{
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
				ssl_mode = "disable"
				sync_frequency = "60"
				use_rudder_storage = true
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
					android_kotlin = [{
						provider = "ketch"
						consents = ["one_android_kotlin", "two_android_kotlin", "three_android_kotlin"]
						resolution_strategy = ""
					}]
					android = [{
						provider = "ketch"
						consents = ["one_android", "two_android", "three_android"]
						resolution_strategy = ""
					}]
					ios_swift = [{
						provider = "custom"
						resolution_strategy = "and"
						consents = ["one_ios_swift", "two_ios_swift", "three_ios_swift"]
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
				"host": "test-host",
				"database": "test-database",
				"user": "test-user",
				"password": "test-password",
				"port": "test-port",
				"sslMode": "disable",
				"syncFrequency": "60",
				"useRudderStorage": true,
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
					"androidKotlin": [
						{
							"provider": "ketch",
							"resolutionStrategy": "",
							"consents": [
								{
									"consent": "one_android_kotlin"
								},
								{
									"consent": "two_android_kotlin"
								},
								{
									"consent": "three_android_kotlin"
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
					"iosSwift": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{
									"consent": "one_ios_swift"
								},
								{
									"consent": "two_ios_swift"
								},
								{
									"consent": "three_ios_swift"
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
}

func TestDestinationResourcePostgresS3AccessKeyStorage(t *testing.T) {
	cmt.AssertDestination(t, "postgres", []c.TestConfig{
		{
			TerraformCreate: `
				host = "test-host"
				database = "test-database"
				user = "test-user"
				password = "test-password"
				use_rudder_storage = false
				s3 {
					bucket_name = "test-s3-bucket"
					access_key_id = "test-access-key-id"
					access_key = "test-access-key"
				}
			`,
			APICreate: `{
				"host": "test-host",
				"database": "test-database",
				"user": "test-user",
				"password": "test-password",
				"port": "5432",
				"sslMode": "disable",
				"syncFrequency": "30",
				"useRudderStorage": false,
				"bucketProvider": "S3",
				"bucketName": "test-s3-bucket",
				"roleBasedAuth": false,
				"accessKeyID": "test-access-key-id",
				"accessKey": "test-access-key"
			}`,
			TerraformUpdate: `
				host = "test-host"
				database = "test-database"
				user = "test-user"
				password = "test-password"
				use_rudder_storage = false
				s3 {
					bucket_name = "test-s3-bucket"
					role_based_authentication {
						i_am_role_arn = "arn:aws:iam::123456789012:role/S3Access"
					}
				}
			`,
			APIUpdate: `{
				"host": "test-host",
				"database": "test-database",
				"user": "test-user",
				"password": "test-password",
				"port": "5432",
				"sslMode": "disable",
				"syncFrequency": "30",
				"useRudderStorage": false,
				"bucketProvider": "S3",
				"bucketName": "test-s3-bucket",
				"roleBasedAuth": true,
				"iamRoleARN": "arn:aws:iam::123456789012:role/S3Access"
			}`,
		},
	})
}

func TestDestinationResourcePostgresGCPStorage(t *testing.T) {
	cmt.AssertDestination(t, "postgres", []c.TestConfig{
		{
			TerraformCreate: `
				host = "test-host"
				database = "test-database"
				user = "test-user"
				password = "test-password"
				use_rudder_storage = false
				gcp {
					bucket_name = "test-gcs-bucket"
					credentials = "{\"type\":\"service_account\"}"
				}
			`,
			APICreate: `{
				"host": "test-host",
				"database": "test-database",
				"user": "test-user",
				"password": "test-password",
				"port": "5432",
				"sslMode": "disable",
				"syncFrequency": "30",
				"useRudderStorage": false,
				"bucketProvider": "GCS",
				"bucketName": "test-gcs-bucket",
				"credentials": "{\"type\":\"service_account\"}"
			}`,
			TerraformUpdate: `
				host = "test-host"
				database = "test-database"
				user = "test-user"
				password = "test-password"
				use_rudder_storage = true
			`,
			APIUpdate: `{
				"host": "test-host",
				"database": "test-database",
				"user": "test-user",
				"password": "test-password",
				"port": "5432",
				"sslMode": "disable",
				"syncFrequency": "30",
				"useRudderStorage": true
			}`,
		},
	})
}

func TestDestinationResourcePostgresAzureStorage(t *testing.T) {
	cmt.AssertDestination(t, "postgres", []c.TestConfig{
		{
			TerraformCreate: `
				host = "test-host"
				database = "test-database"
				user = "test-user"
				password = "test-password"
				use_rudder_storage = false
				azure {
					container_name = "test-container"
					account_name = "test-account"
					account_key = "test-account-key"
				}
			`,
			APICreate: `{
				"host": "test-host",
				"database": "test-database",
				"user": "test-user",
				"password": "test-password",
				"port": "5432",
				"sslMode": "disable",
				"syncFrequency": "30",
				"useRudderStorage": false,
				"bucketProvider": "AZURE_BLOB",
				"containerName": "test-container",
				"accountName": "test-account",
				"accountKey": "test-account-key",
				"useSASTokens": false
			}`,
			TerraformUpdate: `
				host = "test-host"
				database = "test-database"
				user = "test-user"
				password = "test-password"
				use_rudder_storage = false
				azure {
					container_name = "test-container"
					account_name = "test-account"
					use_sas_tokens = true
					sas_token = "test-sas-token"
				}
			`,
			APIUpdate: `{
				"host": "test-host",
				"database": "test-database",
				"user": "test-user",
				"password": "test-password",
				"port": "5432",
				"sslMode": "disable",
				"syncFrequency": "30",
				"useRudderStorage": false,
				"bucketProvider": "AZURE_BLOB",
				"containerName": "test-container",
				"accountName": "test-account",
				"useSASTokens": true,
				"sasToken": "test-sas-token"
			}`,
		},
	})
}

func TestDestinationResourcePostgresMinIOStorage(t *testing.T) {
	cmt.AssertDestination(t, "postgres", []c.TestConfig{
		{
			TerraformCreate: `
				host = "test-host"
				database = "test-database"
				user = "test-user"
				password = "test-password"
				use_rudder_storage = false
				minio {
					bucket_name = "test-minio-bucket"
					end_point = "minio.example.com:9000"
					access_key_id = "test-minio-access-key-id"
					secret_access_key = "test-minio-secret-access-key"
				}
			`,
			APICreate: `{
				"host": "test-host",
				"database": "test-database",
				"user": "test-user",
				"password": "test-password",
				"port": "5432",
				"sslMode": "disable",
				"syncFrequency": "30",
				"useRudderStorage": false,
				"bucketProvider": "MINIO",
				"bucketName": "test-minio-bucket",
				"endPoint": "minio.example.com:9000",
				"accessKeyID": "test-minio-access-key-id",
				"secretAccessKey": "test-minio-secret-access-key",
				"useSSL": true
			}`,
			TerraformUpdate: `
				host = "test-host"
				database = "test-database"
				user = "test-user"
				password = "test-password"
				use_rudder_storage = false
				minio {
					bucket_name = "test-minio-bucket"
					end_point = "minio.example.com:9000"
					access_key_id = "test-minio-access-key-id"
					secret_access_key = "test-minio-secret-access-key"
					use_ssl = false
				}
			`,
			APIUpdate: `{
				"host": "test-host",
				"database": "test-database",
				"user": "test-user",
				"password": "test-password",
				"port": "5432",
				"sslMode": "disable",
				"syncFrequency": "30",
				"useRudderStorage": false,
				"bucketProvider": "MINIO",
				"bucketName": "test-minio-bucket",
				"endPoint": "minio.example.com:9000",
				"accessKeyID": "test-minio-access-key-id",
				"secretAccessKey": "test-minio-secret-access-key",
				"useSSL": false
			}`,
		},
	})
}

func TestDestinationResourcePostgresValidatesStorageBlocks(t *testing.T) {
	tests := []struct {
		name        string
		config      string
		expectError string
	}{
		{
			name: "requires storage block when Rudder storage is disabled",
			config: `
				host = "test-host"
				database = "test-database"
				user = "test-user"
				password = "test-password"
				use_rudder_storage = false
			`,
			expectError: `exactly one of config\.0\.s3, config\.0\.gcp, config\.0\.azure, or config\.0\.minio must be configured when use_rudder_storage is false`,
		},
		{
			name: "rejects storage block when Rudder storage is enabled",
			config: `
				host = "test-host"
				database = "test-database"
				user = "test-user"
				password = "test-password"
				use_rudder_storage = true
				s3 {
					bucket_name = "test-s3-bucket"
					access_key_id = "test-access-key-id"
					access_key = "test-access-key"
				}
			`,
			expectError: `config\.0\.s3 cannot be configured when use_rudder_storage is true`,
		},
		{
			name: "requires Azure account key when SAS tokens are disabled",
			config: `
				host = "test-host"
				database = "test-database"
				user = "test-user"
				password = "test-password"
				use_rudder_storage = false
				azure {
					container_name = "test-container"
					account_name = "test-account"
				}
			`,
			expectError: `config\.0\.azure\.0\.account_key is required when use_sas_tokens is false`,
		},
		{
			name: "rejects Azure SAS token when SAS tokens are disabled",
			config: `
				host = "test-host"
				database = "test-database"
				user = "test-user"
				password = "test-password"
				use_rudder_storage = false
				azure {
					container_name = "test-container"
					account_name = "test-account"
					sas_token = "test-sas-token"
				}
			`,
			expectError: `config\.0\.azure\.0\.sas_token cannot be configured when use_sas_tokens is false`,
		},
		{
			name: "requires Azure SAS token when SAS tokens are enabled",
			config: `
				host = "test-host"
				database = "test-database"
				user = "test-user"
				password = "test-password"
				use_rudder_storage = false
				azure {
					container_name = "test-container"
					account_name = "test-account"
					use_sas_tokens = true
				}
			`,
			expectError: `config\.0\.azure\.0\.sas_token is required when use_sas_tokens is true`,
		},
		{
			name: "rejects Azure account key when SAS tokens are enabled",
			config: `
				host = "test-host"
				database = "test-database"
				user = "test-user"
				password = "test-password"
				use_rudder_storage = false
				azure {
					container_name = "test-container"
					account_name = "test-account"
					use_sas_tokens = true
					account_key = "test-account-key"
				}
			`,
			expectError: `config\.0\.azure\.0\.account_key cannot be configured when use_sas_tokens is true`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assertPostgresInvalidPlan(t, tc.config, regexp.MustCompile(tc.expectError))
		})
	}
}

func assertPostgresInvalidPlan(t *testing.T, config string, expectError *regexp.Regexp) {
	t.Helper()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"rudderstack": func() (*schema.Provider, error) {
				return rudderstack.NewWithConfigureClientFunc(func(_ context.Context, _ *schema.ResourceData) (*rudderstack.Client, diag.Diagnostics) {
					return &rudderstack.Client{}, diag.Diagnostics{}
				}), nil
			},
		},
		Steps: []resource.TestStep{
			{
				PlanOnly: true,
				Config: fmt.Sprintf(`
					provider "rudderstack" {
						access_token = "some-access-token"
					}

					resource "rudderstack_destination_postgres" "example" {
						name = "example"
						config {
							%s
						}
					}
				`, config),
				ExpectError: expectError,
			},
		},
	})
}

func TestDestinationResourcePostgres(t *testing.T) {
	cmt.AssertDestination(t, "postgres", postgresTestConfigs)
}

func TestAccDestinationPostgres(t *testing.T) {
	acc.AccAssertDestination(t, "postgres", postgresTestConfigs)
}
