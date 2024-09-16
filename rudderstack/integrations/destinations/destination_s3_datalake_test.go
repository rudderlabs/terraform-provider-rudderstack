package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceS3Datalake(t *testing.T) {
	cmt.AssertDestination(t, "s3_datalake", []c.TestConfig{
		{
			TerraformCreate: `
				bucket_name = "bucket"
				use_glue    = false

				sync {
				  frequency = "30"
				}
			`,
			APICreate: `{
				"bucketName": "bucket",
				"useGlue": false,
				"syncFrequency": "30"
			}`,
			TerraformUpdate: `
				bucket_name = "bucket"

				namespace	  = "namespace"
				prefix        = "prefix"
				access_key_id = "..."
				access_key    = "..."
			
				enable_sse = true
				use_glue   = true
				region     = "region"

				sync {
					frequency = "30"
					start_at  = "10:00"
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
				"bucketName": "bucket",
				"namespace": "namespace",
				"prefix": "prefix",
				"accessKeyID": "...",
				"accessKey": "...",
				"enableSSE": true,
				"useGlue": true,
				"region": "region",
				"syncFrequency": "30",
				"syncStartAt": "10:00",
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
