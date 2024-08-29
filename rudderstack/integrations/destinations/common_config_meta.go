package destinations

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

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

	return commonProperties, commonSchema
}
