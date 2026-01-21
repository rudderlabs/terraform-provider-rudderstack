package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceBQStream(t *testing.T) {
	cmt.AssertDestination(t, "bqstream", []c.TestConfig{
		{
			TerraformCreate: `
				projectId = "test_project_id"
				datasetId = "test_dataset_id"
				tableId = "test_table_id"
				insertId = "test_insert_id"
				credentials = "test_credentials"
				connection_mode {
					web = "cloud"
				}
			`,
			APICreate: `{
				"projectId": "test_project_id",
				"datasetId": "test_dataset_id",
				"tableId": "test_table_id",
				"insertId": "test_insert_id",
				"credentials": "test_credentials",
				"connectionMode": {
					"web": "cloud"
				}
			}`,
			TerraformUpdate: `
				projectId = "updated_project_id"
				datasetId = "updated_dataset_id"
				tableId = "updated_table_id"
				insertId = "updated_insert_id"
				credentials = "updated_credentials"
				connection_mode {
					web = "cloud"
				}
				event_filtering {
					whitelist = [ "one", "two", "three" ]
				}
				connection_mode {
					android = "cloud"
					androidKotlin = "cloud"
					ios = "cloud"
					iosSwift = "cloud"
					unity = "cloud"
					amp = "cloud"
					reactnative = "cloud"
					flutter = "cloud"
					cordova = "cloud"
					shopify = "cloud"
					cloud = "cloud"
					warehouse = "cloud"
				}
				consent_management {
					android = [{
						provider = "ketch"
						consents = ["one_android", "two_android", "three_android"]
						resolution_strategy = ""
					}]
					androidKotlin = [{
						provider = "ketch"
						consents = ["one_androidKotlin", "two_androidKotlin", "three_androidKotlin"]
						resolution_strategy = ""
					}]
					ios = [{
						provider = "custom"
						resolution_strategy = "and"
						consents = ["one_ios", "two_ios", "three_ios"]
					}]
					iosSwift = [{
						provider = "custom"
						resolution_strategy = "and"
						consents = ["one_iosSwift", "two_iosSwift", "three_iosSwift"]
					}]
					unity = [{
						provider = "custom"
						resolution_strategy = "or"
						consents = ["one_unity", "two_unity", "three_unity"]
					}]
					amp = [{
						provider = "custom"
						resolution_strategy = "and"
						consents = ["one_amp", "two_amp", "three_amp"]
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
					cloud = [{
						provider = "custom"
						resolution_strategy = "and"
						consents = ["one_cloud", "two_cloud", "three_cloud"]
					}]
					warehouse = [{
						provider = "custom"
						resolution_strategy = "and"
						consents = ["one_warehouse", "two_warehouse", "three_warehouse"]
					}]
					shopify = [{
						provider = "custom"
						resolution_strategy = "and"
						consents = ["one_shopify", "two_shopify", "three_shopify"]
					}]
				}
			`,
			APIUpdate: `{
				"projectId": "updated_project_id",
				"datasetId": "updated_dataset_id",
				"tableId": "updated_table_id",
				"insertId": "updated_insert_id",
				"credentials": "updated_credentials",
				"connectionMode": {
					"web": "cloud"
				},
				"consentManagement": {
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
							"provider": "ketch",
							"resolutionStrategy": "",
							"consents": [
								{
									"consent": "one_androidKotlin"
								},
								{
									"consent": "two_androidKotlin"
								},
								{
									"consent": "three_androidKotlin"
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
									"consent": "one_iosSwift"
								},
								{
									"consent": "two_iosSwift"
								},
								{
									"consent": "three_iosSwift"
								}
							]
						}
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
					"warehouse": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{
									"consent": "one_warehouse"
								},
								{
									"consent": "two_warehouse"
								},
								{
									"consent": "three_warehouse"
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
