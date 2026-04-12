package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"android", "androidKotlin", "ios", "iosSwift", "web", "unity", "amp", "cloud", "warehouse", "reactnative", "flutter", "cordova", "shopify"}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)

	properties := []c.ConfigProperty{
		c.Simple("datasetId", "dataset_id"),
		c.Simple("accessToken", "access_token"),
		c.Simple("actionSource", "action_source"),
		c.Simple("limitedDataUSage", "limited_data_usage", c.SkipZeroValue),
		c.Simple("testDestination", "test_destination", c.SkipZeroValue),
		c.Simple("testEventCode", "test_event_code", c.SkipZeroValue),
		c.Simple("removeExternalId", "remove_external_id", c.SkipZeroValue),
		c.ArrayWithObjects("eventsToEvents", "events_to_events", map[string]interface{}{
			"from": "from",
			"to":   "to",
		}),
		c.ArrayWithObjects("blacklistPiiProperties", "blacklist_pii_properties", map[string]interface{}{
			"blacklistPiiProperties": "property",
			"blacklistPiiHash":       "hash",
		}),
		c.ArrayWithObjects("whitelistPiiProperties", "whitelist_pii_properties", map[string]interface{}{
			"whitelistPiiProperties": "property",
		}),
	}

	properties = append(properties, commonProperties...)

	schema := map[string]*schema.Schema{
		"dataset_id": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Your Dataset ID, from the snippet created on the Facebook Dataset creation page.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"access_token": {
			Type:             schema.TypeString,
			Required:         true,
			Sensitive:        true,
			Description:      "Your Business Access token from your Business Account. Required for cloud-mode.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,500})$"),
		},
		"action_source": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "website",
			Description:      "Choose the fallback action_source value you want to set if action_source is not found in event properties.",
			ValidateDiagFunc: c.StringMatchesRegexp("^(website|email|app|phone_call|chat|physical_store|system_generated|other)$"),
		},
		"limited_data_usage": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable Limited Data Usage.",
		},
		"test_destination": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable this setting if you are using this destination for testing purposes.",
		},
		"test_event_code": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Your test event code from your Facebook Datasets dashboard. Required if Test Destination flag is turned ON.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
		},
		"remove_external_id": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Turn this on to send neither userId nor anonymousId as external_id.",
		},
		"events_to_events": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Map RudderStack events to Facebook standard events.",
			ConfigMode:  schema.SchemaConfigModeAttr,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"from": {
						Type:             schema.TypeString,
						Required:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
					},
					"to": {
						Type:             schema.TypeString,
						Required:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(ViewContent|Search|AddToCart|AddToWishlist|InitiateCheckout|AddPaymentInfo|Purchase|PageView|Lead|CompleteRegistration|Contact|CustomizeProduct|Donate|FindLocation|Schedule|StartTrial|SubmitApplication|Subscribe|)$"),
					},
				},
			},
		},
		"blacklist_pii_properties": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Enter the PII properties to be denylisted.",
			ConfigMode:  schema.SchemaConfigModeAttr,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"property": {
						Type:             schema.TypeString,
						Required:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
					},
					"hash": {
						Type:     schema.TypeBool,
						Required: true,
					},
				},
			},
		},
		"whitelist_pii_properties": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Enter the PII properties to be allowlisted.",
			ConfigMode:  schema.SchemaConfigModeAttr,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"property": {
						Type:             schema.TypeString,
						Required:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
					},
				},
			},
		},
	}

	for key, value := range commonSchema {
		schema[key] = value
	}

	c.Destinations.Register("facebook_conversions", c.ConfigMeta{
		APIType:      "FACEBOOK_CONVERSIONS",
		Properties:   properties,
		ConfigSchema: schema,
	})
}
