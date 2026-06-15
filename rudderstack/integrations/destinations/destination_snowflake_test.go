package destinations_test

import (
	"testing"

	acc "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/acc"
	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

var snowflakeTestConfigs = []c.TestConfig{
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
				"useKeyPairAuth": false,
				"password": "example-password",
				"syncFrequency": "30",
				"skipTracksTable": false,
				"skipUsersTable": true,
				"preferAppend": true,
				"manualSync": false,
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
				skip_tracks_table = true
				skip_users_table = false
				prefer_append = false
				manual_sync = true
				json_paths = "./example-paths"
				prefix = "example-prefix"
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
				s3 {
					bucket_name = "example-bucket-name"
					access_key_id = "example-access-key-id"
					access_key = "example-access-key"
					enable_sse = true
				}
				one_trust_cookie_categories {
					web = [{ one_trust_cookie_category = "one_trust_web" }]
					android = [{ one_trust_cookie_category = "one_trust_android" }]
					ios = [{ one_trust_cookie_category = "one_trust_ios" }]
					unity = [{ one_trust_cookie_category = "one_trust_unity" }]
					amp = [{ one_trust_cookie_category = "one_trust_amp" }]
					cloud = [{ one_trust_cookie_category = "one_trust_cloud" }]
					reactnative = [{ one_trust_cookie_category = "one_trust_reactnative" }]
					cloud_source = [{ one_trust_cookie_category = "one_trust_cloud_source" }]
					flutter = [{ one_trust_cookie_category = "one_trust_flutter" }]
					cordova = [{ one_trust_cookie_category = "one_trust_cordova" }]
					shopify = [{ one_trust_cookie_category = "one_trust_shopify" }]
				}
				ketch_consent_purposes {
					web = [{ purpose = "ketch_web" }]
					android = [{ purpose = "ketch_android" }]
					ios = [{ purpose = "ketch_ios" }]
					unity = [{ purpose = "ketch_unity" }]
					amp = [{ purpose = "ketch_amp" }]
					cloud = [{ purpose = "ketch_cloud" }]
					reactnative = [{ purpose = "ketch_reactnative" }]
					cloud_source = [{ purpose = "ketch_cloud_source" }]
					flutter = [{ purpose = "ketch_flutter" }]
					cordova = [{ purpose = "ketch_cordova" }]
					shopify = [{ purpose = "ketch_shopify" }]
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
					android_kotlin = [{
						provider = "custom"
						resolution_strategy = "and"
						consents = ["one_android_kotlin", "two_android_kotlin", "three_android_kotlin"]
					}]
					ios = [{
						provider = "custom"
						resolution_strategy = "and"
						consents = ["one_ios", "two_ios", "three_ios"]
					}]
					ios_swift = [{
						provider = "custom"
						resolution_strategy = "or"
						consents = ["one_ios_swift", "two_ios_swift", "three_ios_swift"]
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
				"useKeyPairAuth": false,
				"password": "example-password",
				"role": "example-role",
				"syncFrequency": "60",
				"syncStartAt": "10:00",
				"excludeWindow": {
					"excludeWindowStartTime": "11:00",
					"excludeWindowEndTime": "12:00"
				},
				"skipTracksTable": true,
				"skipUsersTable": false,
				"preferAppend": false,
				"manualSync": true,
				"useRudderStorage": false,
				"additionalProperties": true,
				"jsonPaths": "./example-paths",
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
				"cloudProvider": "AWS",
				"roleBasedAuth": false,
				"storageIntegration": "",
				"prefix": "example-prefix",
	        		"bucketName": "example-bucket-name",
	        		"accessKeyID": "example-access-key-id",
	        		"accessKey": "example-access-key",
	        		"enableSSE": true,
				"oneTrustCookieCategories": {
					"web": [{"oneTrustCookieCategory": "one_trust_web"}],
					"android": [{"oneTrustCookieCategory": "one_trust_android"}],
					"ios": [{"oneTrustCookieCategory": "one_trust_ios"}],
					"unity": [{"oneTrustCookieCategory": "one_trust_unity"}],
					"amp": [{"oneTrustCookieCategory": "one_trust_amp"}],
					"cloud": [{"oneTrustCookieCategory": "one_trust_cloud"}],
					"reactnative": [{"oneTrustCookieCategory": "one_trust_reactnative"}],
					"cloudSource": [{"oneTrustCookieCategory": "one_trust_cloud_source"}],
					"flutter": [{"oneTrustCookieCategory": "one_trust_flutter"}],
					"cordova": [{"oneTrustCookieCategory": "one_trust_cordova"}],
					"shopify": [{"oneTrustCookieCategory": "one_trust_shopify"}]
				},
				"ketchConsentPurposes": {
					"web": [{"purpose": "ketch_web"}],
					"android": [{"purpose": "ketch_android"}],
					"ios": [{"purpose": "ketch_ios"}],
					"unity": [{"purpose": "ketch_unity"}],
					"amp": [{"purpose": "ketch_amp"}],
					"cloud": [{"purpose": "ketch_cloud"}],
					"reactnative": [{"purpose": "ketch_reactnative"}],
					"cloudSource": [{"purpose": "ketch_cloud_source"}],
					"flutter": [{"purpose": "ketch_flutter"}],
					"cordova": [{"purpose": "ketch_cordova"}],
					"shopify": [{"purpose": "ketch_shopify"}]
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
					"androidKotlin": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
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
					"iosSwift": [
						{
							"provider": "custom",
							"resolutionStrategy": "or",
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

func TestDestinationResourceSnowflake(t *testing.T) {
	cmt.AssertDestination(t, "snowflake", snowflakeTestConfigs)
}

func TestAccDestinationSnowflake(t *testing.T) {
	acc.AccAssertDestination(t, "snowflake", snowflakeTestConfigs)
}

func TestDestinationResourceSnowflakeWithKeyPairAuth(t *testing.T) {
	cmt.AssertDestination(t, "snowflake", []c.TestConfig{
		{
			TerraformCreate: `
				account = "example-account"
				database = "example-database"
				warehouse = "example-warehouse"
				user = "example-user"
				use_key_pair_auth = true
				private_key = "example-private-key"
				private_key_passphrase = "example-passphrase"
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
				"useKeyPairAuth": true,
				"privateKey": "-----BEGIN PRIVATE KEY-----\nexample-private-key\n-----END PRIVATE KEY-----",
				"privateKeyPassphrase": "example-passphrase",
				"syncFrequency": "30",
				"skipTracksTable": false,
				"skipUsersTable": true,
				"preferAppend": true,
				"manualSync": false,
				"useRudderStorage": true,
				"additionalProperties": true
			}`,
			TerraformUpdate: `
				account = "example-account"
				database = "example-database"
				warehouse = "example-warehouse"
				user = "example-user"
				use_key_pair_auth = true
				private_key = "example-private-key-updated"
				use_rudder_storage = false
				sync {
					frequency = "60"
				}
			`,
			APIUpdate: `{
				"account": "example-account",
				"database": "example-database",
				"warehouse": "example-warehouse",
				"user": "example-user",
				"useKeyPairAuth": true,
				"privateKey": "-----BEGIN PRIVATE KEY-----\nexample-private-key-updated\n-----END PRIVATE KEY-----",
				"syncFrequency": "60",
				"skipTracksTable": false,
				"skipUsersTable": true,
				"preferAppend": true,
				"manualSync": false,
				"useRudderStorage": false,
				"additionalProperties": true
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
				"useKeyPairAuth": false,
				"password": "example-password",
				"syncFrequency": "30",
				"skipTracksTable": false,
				"skipUsersTable": true,
				"preferAppend": true,
				"manualSync": false,
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
				"useKeyPairAuth": false,
				"password": "example-password",
				"syncFrequency": "60",
				"syncStartAt": "10:00",
				"excludeWindow": {
					"excludeWindowStartTime": "11:00",
					"excludeWindowEndTime": "12:00"
				},
				"skipTracksTable": false,
				"skipUsersTable": true,
				"preferAppend": true,
				"manualSync": false,
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
				"useKeyPairAuth": false,
				"password": "example-password",
				"syncFrequency": "30",
				"skipTracksTable": false,
				"skipUsersTable": true,
				"preferAppend": true,
				"manualSync": false,
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
				"useKeyPairAuth": false,
				"password": "example-password",
				"syncFrequency": "60",
				"syncStartAt": "10:00",
				"excludeWindow": {
					"excludeWindowStartTime": "11:00",
					"excludeWindowEndTime": "12:00"
				},
				"skipTracksTable": false,
				"skipUsersTable": true,
				"preferAppend": true,
				"manualSync": false,
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

func TestDestinationResourceSnowflakeWithPEMPrivateKey(t *testing.T) {
	cmt.AssertDestination(t, "snowflake", []c.TestConfig{
		{
			TerraformCreate: `
				account = "example-account"
				database = "example-database"
				warehouse = "example-warehouse"
				user = "example-user"
				use_key_pair_auth = true
				private_key = "-----BEGIN PRIVATE KEY-----\nexample-pem-key\n-----END PRIVATE KEY-----"
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
				"useKeyPairAuth": true,
				"privateKey": "-----BEGIN PRIVATE KEY-----\nexample-pem-key\n-----END PRIVATE KEY-----",
				"syncFrequency": "30",
				"skipTracksTable": false,
				"skipUsersTable": true,
				"preferAppend": true,
				"manualSync": false,
				"useRudderStorage": true,
				"additionalProperties": true
			}`,
			TerraformUpdate: `
				account = "example-account"
				database = "example-database"
				warehouse = "example-warehouse"
				user = "example-user"
				use_key_pair_auth = true
				private_key = "-----BEGIN ENCRYPTED PRIVATE KEY-----\nexample-encrypted-key\n-----END ENCRYPTED PRIVATE KEY-----"
				use_rudder_storage = true
				sync {
					frequency = "30"
				}
			`,
			APIUpdate: `{
				"account": "example-account",
				"database": "example-database",
				"warehouse": "example-warehouse",
				"user": "example-user",
				"useKeyPairAuth": true,
				"privateKey": "-----BEGIN ENCRYPTED PRIVATE KEY-----\nexample-encrypted-key\n-----END ENCRYPTED PRIVATE KEY-----",
				"syncFrequency": "30",
				"skipTracksTable": false,
				"skipUsersTable": true,
				"preferAppend": true,
				"manualSync": false,
				"useRudderStorage": true,
				"additionalProperties": true
			}`,
		},
	})
}

func TestDestinationResourceSnowflakeWithRoleBasedAuth(t *testing.T) {
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
				"useKeyPairAuth": false,
				"password": "example-password",
				"syncFrequency": "30",
				"skipTracksTable": false,
				"skipUsersTable": true,
				"preferAppend": true,
				"manualSync": false,
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
				}
				s3 {
					bucket_name = "example-bucket-name"
					role_based_authentication {
						i_am_role_arn = "arn:aws:iam::123456789012:role/S3Access"
					}
					storage_integration = "example-aws-int"
				}
			`,
			APIUpdate: `{
				"account": "example-account",
				"database": "example-database",
				"warehouse": "example-warehouse",
				"user": "example-user",
				"useKeyPairAuth": false,
				"password": "example-password",
				"syncFrequency": "60",
				"skipTracksTable": false,
				"skipUsersTable": true,
				"preferAppend": true,
				"manualSync": false,
				"useRudderStorage": false,
				"additionalProperties": true,
				"cloudProvider": "AWS",
				"bucketName": "example-bucket-name",
				"roleBasedAuth": true,
				"iamRoleARN": "arn:aws:iam::123456789012:role/S3Access",
				"storageIntegration": "example-aws-int"
			}`,
		},
	})
}
