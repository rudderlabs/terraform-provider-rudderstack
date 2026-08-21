package destinations_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	acc "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/acc"
	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

var customerioTestConfigs = []c.TestConfig{
	{
		TerraformCreate: `
				site_id = "cd820c1b31d8f2696f3b"
				api_key = "cg044d23bc1beb3031c5"
				datacenter = "US"
				api_version = "v2"
				user_id_mapping = "id"

				connection_mode {
					web       = "device"
					android   = "cloud"
					ios       = "cloud"
					cloud     = "cloud"
					warehouse = "cloud"
				}

				use_native_sdk {
					web     = true
					android = true
					ios     = true
				}

				send_page_name_in_sdk {
					web = true
				}

				data_use_in_app {
					web = true
				}

				auto_track_device_attributes {
					android = true
					ios     = false
				}

				background_queue_min_number_of_tasks {
					android = "10"
				}

				background_queue_seconds_delay {
					android = "30"
				}
			`,
		APICreate: `{
				"siteID": "cd820c1b31d8f2696f3b",
				"apiKey": "cg044d23bc1beb3031c5",
				"datacenter": "US",
				"apiVersion": "v2",
				"userIdMapping": "id",
				"connectionMode": {
					"web": "device",
					"android": "cloud",
					"ios": "cloud",
					"cloud": "cloud",
					"warehouse": "cloud"
				},
				"useNativeSDK": {
					"web": true,
					"android": true,
					"ios": true
				},
				"sendPageNameInSDK": {
					"web": true
				},
				"dataUseInApp": {
					"web": true
				},
				"autoTrackDeviceAttributes": {
					"android": true,
					"ios": false
				},
				"backgroundQueueMinNumberOfTasks": {
					"android": "10"
				},
				"backgroundQueueSecondsDelay": {
					"android": "30"
				}
			}`,
		TerraformUpdate: `
				site_id = "cd820c1b31d8f2696f3b"
				api_key = "cg044d23bc1beb3031c5"
				datacenter = "EU"
				api_version = "v2"
				user_id_mapping = "id"
				device_token_event_name = "name"

				event_filtering {
					blacklist = [ "one", "two", "three" ]
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
				"siteID": "cd820c1b31d8f2696f3b",
				"apiKey": "cg044d23bc1beb3031c5",
				"datacenter": "EU",
				"apiVersion": "v2",
				"userIdMapping": "id",
				"deviceTokenEventName": "name",
				"eventFilteringOption": "blacklistedEvents",
				"blacklistedEvents": [{
					"eventName": "one"
				}, {
					"eventName": "two"
				}, {
					"eventName": "three"
				}],
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
}

func TestDestinationResourceCustomerIO(t *testing.T) {
	cmt.AssertDestination(t, "customerio", customerioTestConfigs)
}

var customerioAcceptanceTestConfigs = customerioConfigsWithoutAPIVersionFields(customerioTestConfigs)

func customerioConfigsWithoutAPIVersionFields(configs []c.TestConfig) []c.TestConfig {
	out := make([]c.TestConfig, len(configs))
	copy(out, configs)

	for i := range out {
		out[i].TerraformCreate = removeCustomerIOAPIVersionFieldsFromTerraform(out[i].TerraformCreate)
		out[i].TerraformUpdate = removeCustomerIOAPIVersionFieldsFromTerraform(out[i].TerraformUpdate)
		out[i].APICreate = removeCustomerIOAPIVersionFieldsFromAPI(out[i].APICreate)
		out[i].APIUpdate = removeCustomerIOAPIVersionFieldsFromAPI(out[i].APIUpdate)
	}

	return out
}

func removeCustomerIOAPIVersionFieldsFromTerraform(config string) string {
	config = strings.ReplaceAll(config, "\n\t\t\t\tapi_version = \"v2\"", "")
	config = strings.ReplaceAll(config, "\n\t\t\t\tuser_id_mapping = \"id\"", "")
	return config
}

func removeCustomerIOAPIVersionFieldsFromAPI(config string) string {
	config = strings.ReplaceAll(config, "\n\t\t\t\t\"apiVersion\": \"v2\",", "")
	config = strings.ReplaceAll(config, "\n\t\t\t\t\"userIdMapping\": \"id\",", "")
	return config
}

func TestCustomerIODefaultAPIVersionStateToAPI(t *testing.T) {
	cm := c.Destinations.Entries()["customerio"]

	api, err := cm.StateToAPI(`{
		"api_version": "v1",
		"user_id_mapping": ""
	}`)
	require.NoError(t, err)

	assert.JSONEq(t, `{}`, api)
}

func TestCustomerIODefaultAPIVersionAPIToState(t *testing.T) {
	cm := c.Destinations.Entries()["customerio"]

	state, err := cm.APIToState(`{}`)
	require.NoError(t, err)

	assert.JSONEq(t, `{
		"api_version": "v1"
	}`, state)
}

func TestAccDestinationCustomerIO(t *testing.T) {
	acc.AccAssertDestination(t, "customerio", customerioAcceptanceTestConfigs)
}
