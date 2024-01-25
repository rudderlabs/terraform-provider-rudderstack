package destinations

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func GetConfigMetaForGenericConsentManagement(supportedSourceTypes []string) ([]c.ConfigProperty, map[string]*schema.Schema) {
	consentManagementProperties := []c.ConfigProperty{}
	consentManagementSchema := map[string]*schema.Schema{}

	if len(supportedSourceTypes) != 0 && supportedSourceTypes != nil {
		consent_management_terraform_key := "consent_management"
		consent_management_elements_schema := map[string]*schema.Schema{}

		// Create property and schema for each source type
		for _, sourceType := range supportedSourceTypes {
			consentManagementProperties = append(consentManagementProperties, c.ArrayWithObjects(fmt.Sprintf("consentManagement.%s", sourceType), fmt.Sprintf("%s.0.%s", consent_management_terraform_key, sourceType), map[string]interface{}{
				"provider": "provider",
				"resolutionStrategy": "resolution_strategy",
				"consents": c.APINestedObject{
					TerraformKey: "consents",
					NestedKey: 	"consent",
				},
			}))

			consent_management_elements_schema[sourceType] = &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"provider": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The provider name.",
						},
						"resolution_strategy": {
							Type:        schema.TypeString,
							Optional:    true,
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
			Description: "Specify consent IDs for each CMP.",
			Elem: &schema.Resource{
				Schema: consent_management_elements_schema,
			},
		};
	}

	return consentManagementProperties, consentManagementSchema
}

func GetConfigMetaForOneTrustConsents(supportedSourceTypes []string) ([]c.ConfigProperty, map[string]*schema.Schema) {
	oneTrustConsentsProperties := []c.ConfigProperty{}
	oneTrustConsentsSchema := map[string]*schema.Schema{}

	if len(supportedSourceTypes) != 0 && supportedSourceTypes != nil {
		onetrust_terraform_key := "onetrust_cookie_categories"
		onetrust_elements_schema := make(map[string]*schema.Schema)

		// Create property and schema for each source type
		for _, sourceType := range supportedSourceTypes {
			oneTrustConsentsProperties = append(oneTrustConsentsProperties, c.ArrayWithStrings(fmt.Sprintf("oneTrustCookieCategories.%s", sourceType), "oneTrustCookieCategory", fmt.Sprintf("%s.0.%s", onetrust_terraform_key, sourceType)))

			onetrust_elements_schema[sourceType] = &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			}
		}

		oneTrustConsentsSchema[onetrust_terraform_key] = &schema.Schema{
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Specify OneTrust category IDs.",
			Elem: &schema.Resource{
				Schema: onetrust_elements_schema,
			},
		};
	}

	return oneTrustConsentsProperties, oneTrustConsentsSchema
}

func GetCommonConfigMeta(supportedSourceTypes []string) ([]c.ConfigProperty, map[string]*schema.Schema) {
	commonProperties := []c.ConfigProperty{}
	commonSchema := map[string]*schema.Schema{}

	oneTrustConsentFieldProperties, oneTrustConsentFieldSchema := GetConfigMetaForOneTrustConsents(supportedSourceTypes)

	commonProperties = append(commonProperties, oneTrustConsentFieldProperties...)

	for key, value := range oneTrustConsentFieldSchema {
		commonSchema[key] = value
	}

	gcmFieldProperties, gcmFieldSchema := GetConfigMetaForGenericConsentManagement(supportedSourceTypes)

	commonProperties = append(commonProperties, gcmFieldProperties...)

	for key, value := range gcmFieldSchema {
		commonSchema[key] = value
	}

	return commonProperties, commonSchema
}