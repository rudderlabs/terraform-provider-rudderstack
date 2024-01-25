package destinations_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestDestinationResourceGoogleSheets(t *testing.T) {
	cmt.AssertDestination(t, "google_sheets", []c.TestConfig{
		{
			TerraformCreate: `
				sheet_name = "sheet"
                credentials = "..."
                sheet_id = "123"
			`,
			APICreate: `{
				"sheetName": "sheet",
                 "credentials": "...",
                 "sheetId": "123"
			}`,
			TerraformUpdate: `
				sheet_name = "sheetName"
                credentials = "..."
                sheet_id = "1234"
                event_key_map = [
					{
						from = "a1"
						to   = "b1"
					},
				]
				onetrust_cookie_categories = ["one", "two", "three"]
			`,
			APIUpdate: `{
				"sheetName": "sheetName",
                "credentials": "...",
                "sheetId": "1234",
                 "eventKeyMap": [
					{
						"from": "a1",
						"to": "b1"
					}
				],
				"oneTrustCookieCategories": [
					{ "oneTrustCookieCategory": "one" },
					{ "oneTrustCookieCategory": "two" },
					{ "oneTrustCookieCategory": "three" }
				]
			}`,
		},
	})
}
