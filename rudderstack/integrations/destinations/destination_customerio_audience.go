package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"warehouse"}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)

	properties := []c.ConfigProperty{
		c.Simple("siteId", "site_id"),
		c.Simple("apiKey", "api_key"),
		c.Simple("appApiKey", "app_api_key"),
		c.Simple("region", "region"),
		c.Simple("connectionMode.warehouse", "connection_mode.0.warehouse", c.SkipZeroValue),
	}

	properties = append(properties, commonProperties...)

	s := map[string]*schema.Schema{
		"site_id": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Settings > Account Settings > API Credentials > Track APP Keys > Site ID",
			ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
		},
		"api_key": {
			Type:             schema.TypeString,
			Required:         true,
			Sensitive:        true,
			Description:      "Settings > Account Settings > API Credentials > Track APP Keys > API KEY of the corresponding Site ID",
			ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
		},
		"app_api_key": {
			Type:             schema.TypeString,
			Required:         true,
			Sensitive:        true,
			Description:      "Settings > Account Settings > API Credentials > APP API Keys > API KEY",
			ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^(.{1,100})$"),
		},
		"region": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Select your Customer.io Data Center",
			ValidateDiagFunc: c.StringMatchesRegexp("^(US|EU)$"),
		},
		"connection_mode": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Configure the connection mode for Customer.io Audience.",
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
		s[key] = value
	}

	c.Destinations.Register("customerio_audience", c.ConfigMeta{
		APIType:      "CUSTOMERIO_AUDIENCE",
		Properties:   properties,
		ConfigSchema: s,
	})
}
