package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceBigQuery(t *testing.T) {
	cmt.AssertDestination(t, "bigquery", []c.TestConfig{
		{
			TerraformCreate: `
				project     = "project"
				bucket_name = "bucket"
				credentials = "..."
				partition_column = "loaded_at"
				partition_type = "hour"
						
				sync {
					frequency = "30"
				}
			`,
			APICreate: `{
				"project": "project",
				"bucketName": "bucket",
				"credentials": "...",
				"skipTracksTable": false,
				"skipViews": false,
				"skipUsersTable": true,
				"partitionColumn": "loaded_at",
				"partitionType": "hour",
				"syncFrequency": "30"
			}`,
			TerraformUpdate: `
				project     = "project"
				bucket_name = "bucket"
				credentials = "..."
			
				location  = "us-east1"
				prefix    = "prefix"
				namespace = "namespace"
				skip_tracks_table = true
				skip_users_table = false
				skip_views = false
			
				sync {
					frequency				  = "30"
					start_at   				  = "10:00"
					exclude_window_start_time = "11:00"
					exclude_window_end_time   = "12:00"
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
				cleanup_object_storage_files = true
				json_paths = "event.properties.key1,event.properties.key2"
				connection_mode {
					web = "cloud"
					ios = "cloud"
					android = "cloud"
					reactnative = "cloud"
					unity = "cloud"
					amp = "cloud"
					flutter = "cloud"
					cordova = "cloud"
					shopify = "cloud"
					cloud = "cloud"
					cloud_source = "cloud"
				}
			`,
			APIUpdate: `{
				"project": "project",
				"bucketName": "bucket",
				"credentials": "...",
				"skipTracksTable": true,
				"skipViews": false,
				"skipUsersTable": false,
				"partitionColumn": "_PARTITIONTIME",
				"partitionType": "day",
				"location": "us-east1",
				"prefix": "prefix",
				"namespace": "namespace",
				"syncFrequency": "30",
				"syncStartAt": "10:00",
				"excludeWindow": {
					"excludeWindowStartTime": "11:00",
					"excludeWindowEndTime": "12:00"
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
				},
				"cleanupObjectStorageFiles": true,
				"jsonPaths": "event.properties.key1,event.properties.key2",
				"connectionMode": {
					"web": "cloud",
					"ios": "cloud",
					"android": "cloud",
					"reactnative": "cloud",
					"unity": "cloud",
					"amp": "cloud",
					"flutter": "cloud",
					"cordova": "cloud",
					"shopify": "cloud",
					"cloud": "cloud",
					"cloudSource": "cloud"
				}
			}`,
		},
	})
}
