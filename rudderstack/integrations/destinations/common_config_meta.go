package destinations

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

// A function to covert a string from camelCase to snake_case
func camelToSnake(s string) string {
	var res string
	for i, v := range s {
		if 'A' <= v && v <= 'Z' {
			if i != 0 {
				res += "_"
			}
			res += string(v + 32)
		} else {
			res += string(v)
		}
	}
	return res
}

func GetConfigMetaForGenericConsentManagement(supportedSourceTypes []string) ([]c.ConfigProperty, map[string]*schema.Schema) {
	consentManagementProperties := []c.ConfigProperty{}
	consentManagementSchema := map[string]*schema.Schema{}

	if len(supportedSourceTypes) != 0 && supportedSourceTypes != nil {
		consent_management_terraform_key := "consent_management"
		consent_management_elements_schema := map[string]*schema.Schema{}

		// Create property and schema for each source type
		for _, sourceType := range supportedSourceTypes {
			validSourceType := camelToSnake(sourceType)

			consentManagementProperties = append(consentManagementProperties, c.ArrayWithObjects(fmt.Sprintf("consentManagement.%s", sourceType), fmt.Sprintf("%s.0.%s", consent_management_terraform_key, validSourceType), map[string]interface{}{
				"provider":           "provider",
				"resolutionStrategy": "resolution_strategy",
				"consents": c.APINestedObject{
					TerraformKey: "consents",
					NestedKey:    "consent",
				},
			}))

			consent_management_elements_schema[validSourceType] = &schema.Schema{
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
								"iubenda",
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
			}
		}

		consentManagementSchema[consent_management_terraform_key] = &schema.Schema{
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Allows you to specify consent configuration data for multiple providers for each source type.",
			Elem: &schema.Resource{
				Schema: consent_management_elements_schema,
			},
		}
	}

	return consentManagementProperties, consentManagementSchema
}

func GetCommonConfigMeta(supportedSourceTypes []string) ([]c.ConfigProperty, map[string]*schema.Schema) {
	commonProperties := []c.ConfigProperty{}
	commonSchema := map[string]*schema.Schema{}

	gcmFieldProperties, gcmFieldSchema := GetConfigMetaForGenericConsentManagement(supportedSourceTypes)

	commonProperties = append(commonProperties, gcmFieldProperties...)

	for key, value := range gcmFieldSchema {
		commonSchema[key] = value
	}

	return commonProperties, commonSchema
}
