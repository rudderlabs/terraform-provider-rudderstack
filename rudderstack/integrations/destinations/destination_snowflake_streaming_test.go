package destinations_test

import (
	"testing"

	acc "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/acc"
	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

var snowflakeStreamingTestConfigs = []c.TestConfig{
	{
		TerraformCreate: `
				account   = "example-account"
				database  = "example-database"
				warehouse = "example-warehouse"
				user      = "example-user"
				private_key = "-----BEGIN PRIVATE KEY-----\nabc\n-----END PRIVATE KEY-----"
				namespace = "example_namespace"
			`,
		APICreate: `{
				"account": "example-account",
				"database": "example-database",
				"warehouse": "example-warehouse",
				"user": "example-user",
				"privateKey": "-----BEGIN PRIVATE KEY-----\nabc\n-----END PRIVATE KEY-----",
				"namespace": "example_namespace",
				"skipTracksTable": false,
				"enableIceberg": false,
				"underscoreDivideNumbers": false,
				"allowUsersContextTraits": false
			}`,
		TerraformUpdate: `
				account   = "updated-account"
				database  = "updated-database"
				warehouse = "updated-warehouse"
				user      = "updated-user"
				role      = "updated-role"
				private_key = "-----BEGIN PRIVATE KEY-----\nupdated\n-----END PRIVATE KEY-----"
				private_key_passphrase = "updated-passphrase"
				namespace = "updated_namespace"
				skip_tracks_table = true
				json_paths = "event.properties.a,event.properties.b"
				enable_iceberg = true
				external_volume = "EXT_VOLUME"
				underscore_divide_numbers = true
				allow_users_context_traits = true
				connection_mode {
					android = "cloud"
					android_kotlin = "cloud"
					ios = "cloud"
					ios_swift = "cloud"
					web = "cloud"
					unity = "cloud"
					amp = "cloud"
					reactnative = "cloud"
					cloud = "cloud"
					cloud_source = "cloud"
					flutter = "cloud"
					cordova = "cloud"
					shopify = "cloud"
				}
				one_trust_cookie_categories {
					web = [{ one_trust_cookie_category = "ot_web" }]
					cloud_source = [{ one_trust_cookie_category = "ot_cloud_source" }]
				}
				ketch_consent_purposes {
					android = [{ purpose = "purpose_android" }]
					shopify = [{ purpose = "purpose_shopify" }]
				}
				consent_management {
					web = [{
						provider = "custom"
						resolution_strategy = "and"
						consents = ["one_web", "two_web"]
					}]
					cloud_source = [{
						provider = "ketch"
						resolution_strategy = ""
						consents = ["one_cloud_source"]
					}]
				}
			`,
		APIUpdate: `{
				"account": "updated-account",
				"database": "updated-database",
				"warehouse": "updated-warehouse",
				"user": "updated-user",
				"role": "updated-role",
				"privateKey": "-----BEGIN PRIVATE KEY-----\nupdated\n-----END PRIVATE KEY-----",
				"privateKeyPassphrase": "updated-passphrase",
				"namespace": "updated_namespace",
				"skipTracksTable": true,
				"jsonPaths": "event.properties.a,event.properties.b",
				"enableIceberg": true,
				"externalVolume": "EXT_VOLUME",
				"underscoreDivideNumbers": true,
				"allowUsersContextTraits": true,
				"connectionMode": {
					"android": "cloud",
					"androidKotlin": "cloud",
					"ios": "cloud",
					"iosSwift": "cloud",
					"web": "cloud",
					"unity": "cloud",
					"amp": "cloud",
					"reactnative": "cloud",
					"cloud": "cloud",
					"cloudSource": "cloud",
					"flutter": "cloud",
					"cordova": "cloud",
					"shopify": "cloud"
				},
				"oneTrustCookieCategories": {
					"web": [{"oneTrustCookieCategory": "ot_web"}],
					"cloudSource": [{"oneTrustCookieCategory": "ot_cloud_source"}]
				},
				"ketchConsentPurposes": {
					"android": [{"purpose": "purpose_android"}],
					"shopify": [{"purpose": "purpose_shopify"}]
				},
				"consentManagement": {
					"web": [
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{"consent": "one_web"},
								{"consent": "two_web"}
							]
						}
					],
					"cloudSource": [
						{
							"provider": "ketch",
							"resolutionStrategy": "",
							"consents": [
								{"consent": "one_cloud_source"}
							]
						}
					]
				}
			}`,
	},
}

func TestDestinationResourceSnowflakeStreaming(t *testing.T) {
	cmt.AssertDestination(t, "snowflake_streaming", snowflakeStreamingTestConfigs)
}

func TestAccDestinationSnowflakeStreaming(t *testing.T) {
	acc.AccAssertDestination(t, "snowflake_streaming", snowflakeStreamingTestConfigs)
}
