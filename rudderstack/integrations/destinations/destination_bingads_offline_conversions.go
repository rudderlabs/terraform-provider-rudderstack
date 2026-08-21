package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"warehouse"}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)

	properties := []c.ConfigProperty{
		c.Simple("rudderAccountId", "rudder_account_id"),
		c.Simple("customerAccountId", "customer_account_id"),
		c.Simple("customerId", "customer_id"),
		c.Simple("isHashRequired", "is_hash_required"),
		c.Simple("connectionMode.warehouse", "connection_mode.0.warehouse", c.SkipZeroValue),
	}

	properties = append(properties, commonProperties...)

	destinationSchema := map[string]*schema.Schema{
		"rudder_account_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The RudderStack account ID for OAuth-based event delivery.",
		},
		"customer_account_id": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Your Bing Ads customer account ID.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^[0-9]+$"),
		},
		"customer_id": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Your Bing Ads customer ID.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^[0-9]+$"),
		},
		"is_hash_required": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable this setting if you are not sending hashed data.",
		},
		"connection_mode": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Configure the connection mode for Bing Ads Offline Conversions.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"warehouse": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$"),
					},
				},
			},
		},
	}

	for key, value := range commonSchema {
		destinationSchema[key] = value
	}

	c.Destinations.Register("bingads_offline_conversions", c.ConfigMeta{
		APIType:      "BINGADS_OFFLINE_CONVERSIONS",
		Version:      1,
		Properties:   properties,
		ConfigSchema: destinationSchema,
	})
}
