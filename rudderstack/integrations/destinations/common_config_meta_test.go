// write unit tests for the common_config_meta.go file

package destinations

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/stretchr/testify/require"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
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
			supportedSourceTypes: []string{"web", "android", "ios", "cloudSource"},
			expectedProperties: []c.ConfigProperty{
				c.ArrayWithObjects("consentManagement.web", "consent_management.0.web", map[string]interface{}{
					"provider":           "provider",
					"resolutionStrategy": "resolution_strategy",
					"consents":           c.APINestedObject{TerraformKey: "consents", NestedKey: "consent"},
				}),
				c.ArrayWithObjects("consentManagement.android", "consent_management.0.android", map[string]interface{}{
					"provider":           "provider",
					"resolutionStrategy": "resolution_strategy",
					"consents":           c.APINestedObject{TerraformKey: "consents", NestedKey: "consent"},
				}),
				c.ArrayWithObjects("consentManagement.ios", "consent_management.0.ios", map[string]interface{}{
					"provider":           "provider",
					"resolutionStrategy": "resolution_strategy",
					"consents":           c.APINestedObject{TerraformKey: "consents", NestedKey: "consent"},
				}),
				c.ArrayWithObjects("consentManagement.cloudSource", "consent_management.0.cloud_source", map[string]interface{}{
					"provider":           "provider",
					"resolutionStrategy": "resolution_strategy",
					"consents":           c.APINestedObject{TerraformKey: "consents", NestedKey: "consent"},
				}),
			},
			expectedSchema: map[string]*schema.Schema{
				"consent_management": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: "Allows you to specify consent configuration data for multiple providers for each source type.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"web": {
								Type:        schema.TypeList,
								Optional:    true,
								ConfigMode:  schema.SchemaConfigModeAttr,
								Description: "Allows you to specify consent configuration data for multiple providers.",
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"provider": {
											Type:     schema.TypeString,
											Required: true,
											ValidateFunc: validation.StringInSlice([]string{
												"oneTrust",
												"ketch",
												"custom",
											}, false),
											Description: "The provider name.",
										},
										"resolution_strategy": {
											Type:     schema.TypeString,
											Optional: true,
											ValidateFunc: validation.StringInSlice([]string{
												"and",
												"or",
												"",
											}, false),
											Description: "The resolution strategy for the provider.",
										},
										"consents": {
											Type:        schema.TypeList,
											Required:    true,
											Description: "The list of consent IDs for the provider.",
											Elem: &schema.Schema{
												Type: schema.TypeString,
											},
										},
									},
								},
							},
							"android": {
								Type:        schema.TypeList,
								Optional:    true,
								ConfigMode:  schema.SchemaConfigModeAttr,
								Description: "Allows you to specify consent configuration data for multiple providers.",
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"provider": {
											Type:     schema.TypeString,
											Required: true,
											ValidateFunc: validation.StringInSlice([]string{
												"oneTrust",
												"ketch",
												"custom",
											}, false),
											Description: "The provider name.",
										},
										"resolution_strategy": {
											Type:     schema.TypeString,
											Optional: true,
											ValidateFunc: validation.StringInSlice([]string{
												"and",
												"or",
												"",
											}, false),
											Description: "The resolution strategy for the provider.",
										},
										"consents": {
											Type:        schema.TypeList,
											Required:    true,
											Description: "The list of consent IDs for the provider.",
											Elem: &schema.Schema{
												Type: schema.TypeString,
											},
										},
									},
								},
							},
							"ios": {
								Type:        schema.TypeList,
								Optional:    true,
								ConfigMode:  schema.SchemaConfigModeAttr,
								Description: "Allows you to specify consent configuration data for multiple providers.",
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"provider": {
											Type:     schema.TypeString,
											Required: true,
											ValidateFunc: validation.StringInSlice([]string{
												"oneTrust",
												"ketch",
												"custom",
											}, false),
											Description: "The provider name.",
										},
										"resolution_strategy": {
											Type:     schema.TypeString,
											Optional: true,
											ValidateFunc: validation.StringInSlice([]string{
												"and",
												"or",
												"",
											}, false),
											Description: "The resolution strategy for the provider.",
										},
										"consents": {
											Type:        schema.TypeList,
											Required:    true,
											Description: "The list of consent IDs for the provider.",
											Elem: &schema.Schema{
												Type: schema.TypeString,
											},
										},
									},
								},
							},
							"cloudSource": {
								Type:        schema.TypeList,
								Optional:    true,
								ConfigMode:  schema.SchemaConfigModeAttr,
								Description: "Allows you to specify consent configuration data for multiple providers.",
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"provider": {
											Type:     schema.TypeString,
											Required: true,
											ValidateFunc: validation.StringInSlice([]string{
												"oneTrust",
												"ketch",
												"custom",
											}, false),
											Description: "The provider name.",
										},
										"resolution_strategy": {
											Type:     schema.TypeString,
											Optional: true,
											ValidateFunc: validation.StringInSlice([]string{
												"and",
												"or",
												"",
											}, false),
											Description: "The resolution strategy for the provider.",
										},
										"consents": {
											Type:        schema.TypeList,
											Required:    true,
											Description: "The list of consent IDs for the provider.",
											Elem: &schema.Schema{
												Type: schema.TypeString,
											},
										},
									},
								},
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
				c.ArrayWithObjects("consentManagement.web", "consent_management.0.web", map[string]interface{}{
					"provider":           "provider",
					"resolutionStrategy": "resolution_strategy",
					"consents":           c.APINestedObject{TerraformKey: "consents", NestedKey: "consent"},
				}),
			},
			expectedSchema: map[string]*schema.Schema{
				"consent_management": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: "Allows you to specify consent configuration data for multiple providers for each source type.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"web": {
								Type:        schema.TypeList,
								Optional:    true,
								ConfigMode:  schema.SchemaConfigModeAttr,
								Description: "Allows you to specify consent configuration data for multiple providers.",
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"provider": {
											Type:     schema.TypeString,
											Required: true,
											ValidateFunc: validation.StringInSlice([]string{
												"oneTrust",
												"ketch",
												"custom",
											}, false),
											Description: "The provider name.",
										},
										"resolution_strategy": {
											Type:     schema.TypeString,
											Optional: true,
											ValidateFunc: validation.StringInSlice([]string{
												"and",
												"or",
												"",
											}, false),
											Description: "The resolution strategy for the provider.",
										},
										"consents": {
											Type:        schema.TypeList,
											Required:    true,
											Description: "The list of consent IDs for the provider.",
											Elem: &schema.Schema{
												Type: schema.TypeString,
											},
										},
									},
								},
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

			require.True(t, len(actualProperties) == len(tc.expectedProperties))

			require.True(t, len(actualSchema) == len(tc.expectedSchema))
		})
	}
}
