// write unit tests for the common_config_meta.go file

package destinations

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
	"github.com/stretchr/testify/require"
)

func TestGetCommonConfigMeta(t *testing.T) {
	testCases := []struct {
		description          string
		supportedSourceTypes []string
		expectedProperties   []c.ConfigProperty
		expectedSchema       map[string]*schema.Schema
	}{
		{
			description:          "Valid list of supported source types",
			supportedSourceTypes: []string{"web", "android", "ios"},
			expectedProperties: []c.ConfigProperty{
				c.ArrayWithStrings("oneTrustCookieCategories.web", "oneTrustCookieCategory", "onetrust_cookie_categories.0.web"),
				c.ArrayWithStrings("oneTrustCookieCategories.android", "oneTrustCookieCategory", "onetrust_cookie_categories.0.android"),
				c.ArrayWithStrings("oneTrustCookieCategories.ios", "oneTrustCookieCategory", "onetrust_cookie_categories.0.ios"),
			},
			expectedSchema: map[string]*schema.Schema{
				"onetrust_cookie_categories": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: "Specify OneTrust category IDs.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"web": {
								Type:        schema.TypeList,
								Optional:    true,
								Elem:        &schema.Schema{Type: schema.TypeString},
							},
							"android": {
								Type:        schema.TypeList,
								Optional:    true,
								Elem:        &schema.Schema{Type: schema.TypeString},
							},
							"ios": {
								Type:        schema.TypeList,
								Optional:    true,
								Elem:        &schema.Schema{Type: schema.TypeString},
							},
						},
					},
				},
			},
		},
		{
			description:          "Empty list of supported source types",
			supportedSourceTypes: []string{},
			expectedProperties:   []c.ConfigProperty{},
			expectedSchema:       map[string]*schema.Schema{},
		},
		{
			description:          "Nil supported source types",
			supportedSourceTypes: nil,
			expectedProperties:   []c.ConfigProperty{},
			expectedSchema:       map[string]*schema.Schema{},
		},
		{
			description:          "A single supported source type",
			supportedSourceTypes: []string{"web"},
			expectedProperties: []c.ConfigProperty{
				c.ArrayWithStrings("oneTrustCookieCategories.web", "oneTrustCookieCategory", "onetrust_cookie_categories.0.web"),
			},
			expectedSchema: map[string]*schema.Schema{
				"onetrust_cookie_categories": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: "Specify OneTrust category IDs.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"web": {
								Type:        schema.TypeList,
								Optional:    true,
								Elem:        &schema.Schema{Type: schema.TypeString},
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actualProperties, actualSchema := GetCommonConfigMeta(tc.supportedSourceTypes)
			require.EqualValues(t, tc.expectedSchema, actualSchema)
			require.True(t, len(actualProperties) == len(tc.expectedProperties))
		})
	}
}
