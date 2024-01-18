package destinations

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func GetConfigMetaForOneTrustConsents(supportedSourceTypes []string) ([]c.ConfigProperty, map[string]*schema.Schema) {
	onetrust_terraform_key := "onetrust_cookie_categories"

	oneTrustConsentFieldProperties := []c.ConfigProperty{}
	onetrust_elements_schema := make(map[string]*schema.Schema)

	for _, sourceType := range supportedSourceTypes {
		oneTrustConsentFieldProperties = append(oneTrustConsentFieldProperties, c.ArrayWithStrings(fmt.Sprintf("oneTrustCookieCategories.%s", sourceType), "oneTrustCookieCategory", fmt.Sprintf("%s.0.%s", onetrust_terraform_key, sourceType)))

		onetrust_elements_schema[sourceType] = &schema.Schema{
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		}
	}

	oneTrustConsentFieldSchema := map[string]*schema.Schema{
		onetrust_terraform_key: {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Specify OneTrust category IDs.",
			Elem: &schema.Resource{
				Schema: onetrust_elements_schema,
			},
		},
	}

	return oneTrustConsentFieldProperties, oneTrustConsentFieldSchema
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
