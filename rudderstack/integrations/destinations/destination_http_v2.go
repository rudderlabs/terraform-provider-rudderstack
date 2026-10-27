package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"android", "androidKotlin", "ios", "iosSwift", "web", "unity", "amp", "cloud", "warehouse", "reactnative", "flutter", "cordova", "shopify"}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)

	properties := []c.ConfigProperty{
		c.Simple("apiUrl", "api_url"),
		c.Simple("auth", "auth"),
		c.Simple("username", "username", c.SkipZeroValue),
		c.Simple("password", "password", c.SkipZeroValue),
		c.Simple("bearerToken", "bearer_token", c.SkipZeroValue),
		c.Simple("apiKeyName", "api_key_name", c.SkipZeroValue),
		c.Simple("apiKeyValue", "api_key_value", c.SkipZeroValue),
		c.Simple("method", "method"),
		c.Simple("format", "format"),
		c.Simple("requestBodyTemplate", "request_body_template", c.SkipZeroValue),
		c.Simple("supportedMessageTypes", "supported_message_types", c.SkipZeroValue),
		c.ArrayWithObjects("queryParams", "query_params", map[string]interface{}{"to": "to", "from": "from"}),
		c.ArrayWithObjects("headers", "headers", map[string]interface{}{"to": "to", "from": "from"}),
		c.ArrayWithObjects("pathParams", "path_params", map[string]interface{}{"path": "path"}),
		c.Simple("isBatchingEnabled", "is_batching_enabled"),
		c.Simple("maxBatchSize", "max_batch_size", c.SkipZeroValue),
		c.ArrayWithStrings("whitelistedEvents", "eventName", "event_filtering.0.whitelist"),
		c.ArrayWithStrings("blacklistedEvents", "eventName", "event_filtering.0.blacklist"),
		c.Discriminator("eventFilteringOption", c.DiscriminatorValues{
			"event_filtering.0.whitelist": "whitelistedEvents",
			"event_filtering.0.blacklist": "blacklistedEvents",
		}),
	}
	for _, sourceType := range supportedSourceTypes {
		properties = append(properties, c.Simple("connectionMode."+sourceType, "connection_mode.0."+camelToSnake(sourceType), c.SkipZeroValue))
	}
	properties = append(properties, commonProperties...)

	// Keep the v2 schema independent from HTTP v1: the v1 field mappings and
	// JSONPath validators must not constrain JSONata request expressions.
	configSchema := map[string]*schema.Schema{
		"api_url": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Enter the base URL for your HTTP endpoint.",
			ValidateDiagFunc: c.ValidateAll(
				c.StringMatchesRegexp(`^(https?://)([a-zA-Z0-9-]{1,63}\.)+[a-zA-Z]{2,}(:(6553[0-5]|655[0-2][0-9]|65[0-4][0-9]{2}|6[0-4][0-9]{3}|[1-5]\d{4}|[1-9]\d{1,3}))?(/.*)?$`),
				c.StringNotMatchesRegexp(`^https?://[^/]*\.ngrok\.io(?:[:/].*)?$`),
				c.StringNotMatchesRegexp(`^https?://(?:localhost[^/:]*|[^/]*\.localhost[^/:]*)(?:[:/].*)?$`),
			),
		},
		"auth": {
			Type: schema.TypeString, Optional: true, Default: "noAuth",
			Description:      "Select the authentication method used for the HTTP request.",
			ValidateDiagFunc: c.StringMatchesRegexp("^(noAuth|basicAuth|bearerTokenAuth|apiKeyAuth)$"),
		},
		"username": {
			Type: schema.TypeString, Optional: true,
			Description:      "Username for basic authentication.",
			ValidateDiagFunc: c.StringMatchesRegexp("^(.{1,100})$"),
		},
		"password": {
			Type: schema.TypeString, Optional: true, Sensitive: true,
			Description:      "Password for basic authentication.",
			ValidateDiagFunc: c.StringMatchesRegexp("^(.{0,100})$"),
		},
		"bearer_token": {
			Type: schema.TypeString, Optional: true, Sensitive: true,
			Description:  "Bearer token for authorization.",
			ValidateFunc: validation.StringLenBetween(1, 2048),
		},
		"api_key_name": {
			Type: schema.TypeString, Optional: true,
			Description:      "Header name for API key authentication.",
			ValidateDiagFunc: c.StringMatchesRegexp(`^(\S{1,100})$`),
		},
		"api_key_value": {
			Type: schema.TypeString, Optional: true, Sensitive: true,
			Description:      "Value for API key authentication.",
			ValidateDiagFunc: c.StringMatchesRegexp("^(.{1,100})$"),
		},
		"method": {
			Type: schema.TypeString, Optional: true, Default: "POST",
			Description:      "HTTP method to use when sending requests.",
			ValidateDiagFunc: c.StringMatchesRegexp("^(POST|PUT|PATCH|GET|DELETE)$"),
		},
		"format": {
			Type: schema.TypeString, Optional: true, Default: "JSON",
			Description:      "HTTP v2 request body format (JSON only).",
			ValidateDiagFunc: c.StringMatchesRegexp("^(JSON)$"),
		},
		"request_body_template": {
			Type: schema.TypeString, Required: true,
			Description:      "JSONata template for the complete JSON request body, evaluated over $$.events. Verify the expression before saving.",
			ValidateDiagFunc: c.StringMatchesRegexp(`(?s)\S`),
		},
		"supported_message_types": {
			Type: schema.TypeList, Optional: true,
			Description: "Message types accepted by this destination (no duplicates).",
			Elem:        &schema.Schema{Type: schema.TypeString, ValidateDiagFunc: c.StringMatchesRegexp("^(track|identify|page|screen|group|alias)$")},
		},
		"headers": {
			Type: schema.TypeList, Optional: true, ConfigMode: schema.SchemaConfigModeAttr,
			Description: "Custom header names and JSONata expressions or literal values.",
			Elem: &schema.Resource{Schema: map[string]*schema.Schema{
				"to":   {Type: schema.TypeString, Optional: true, ValidateDiagFunc: c.StringMatchesRegexp("^[A-Za-z0-9!#%&'*+.^_`|~-]{0,100}$")},
				"from": {Type: schema.TypeString, Optional: true},
			}},
		},
		"query_params": {
			Type: schema.TypeList, Optional: true, ConfigMode: schema.SchemaConfigModeAttr,
			Description: "Query parameter names and JSONata expressions or literal values.",
			Elem: &schema.Resource{Schema: map[string]*schema.Schema{
				"to":   {Type: schema.TypeString, Optional: true, ValidateDiagFunc: c.StringMatchesRegexp("^[A-Za-z0-9!#%&'*+.^_`|~\\- ]{0,100}$")},
				"from": {Type: schema.TypeString, Optional: true},
			}},
		},
		"path_params": {
			Type: schema.TypeList, Optional: true, ConfigMode: schema.SchemaConfigModeAttr,
			Description: "Ordered URL path segments as JSONata expressions or literal values.",
			Elem: &schema.Resource{Schema: map[string]*schema.Schema{
				"path": {Type: schema.TypeString, Optional: true},
			}},
		},
		"is_batching_enabled": {
			Type: schema.TypeBool, Optional: true, Default: false,
			Description: "Enable batching for JSON requests.",
		},
		"max_batch_size": {
			Type: schema.TypeString, Optional: true,
			Description:      "Maximum number of events per request, from 1 to 5,000.",
			ValidateDiagFunc: c.StringMatchesRegexp(`^([1-9]|[1-9][0-9]{1,2}|[1-4][0-9]{3}|5000)$`),
		},
		"event_filtering": {
			Type: schema.TypeList, MaxItems: 1, Optional: true,
			Description: "Allowlist or denylist client-side event names.",
			Elem: &schema.Resource{Schema: map[string]*schema.Schema{
				"whitelist": {Type: schema.TypeList, Optional: true, ExactlyOneOf: []string{"config.0.event_filtering.0.whitelist", "config.0.event_filtering.0.blacklist"}, Elem: &schema.Schema{Type: schema.TypeString}},
				"blacklist": {Type: schema.TypeList, Optional: true, ExactlyOneOf: []string{"config.0.event_filtering.0.whitelist", "config.0.event_filtering.0.blacklist"}, Elem: &schema.Schema{Type: schema.TypeString}},
			}},
		},
	}

	connectionModes := map[string]*schema.Schema{}
	for _, sourceType := range supportedSourceTypes {
		connectionModes[camelToSnake(sourceType)] = &schema.Schema{Type: schema.TypeString, Optional: true, ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$")}
	}
	configSchema["connection_mode"] = &schema.Schema{
		Type: schema.TypeList, MaxItems: 1, Optional: true,
		Description: "Connection mode for each supported source type (cloud only).",
		Elem:        &schema.Resource{Schema: connectionModes},
	}
	for key, value := range commonSchema {
		configSchema[key] = value
	}

	c.Destinations.Register("http_v2", c.ConfigMeta{
		APIType: "HTTP", Version: 2,
		Properties: properties, ConfigSchema: configSchema,
	})
}
