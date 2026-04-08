package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"android", "androidKotlin", "ios", "iosSwift", "web", "unity", "amp", "cloud", "warehouse", "reactnative", "flutter", "cordova", "shopify"}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)
	httpURLPattern := `^(https?://)([a-zA-Z0-9-]{1,63}\.)+[a-zA-Z]{2,}(:(6553[0-5]|655[0-2][0-9]|65[0-4][0-9]{2}|6[0-4][0-9]{3}|[1-5]\d{4}|[1-9]\d{1,3}))?(/.*)?$`
	httpNgrokPattern := `^https?://[^/]*\.ngrok\.io(?:[:/].*)?$`
	httpLocalhostPattern := `^https?://(?:localhost[^/:]*|[^/]*\.localhost[^/:]*)(?:[:/].*)?$`

	properties := []c.ConfigProperty{
		c.Simple("apiUrl", "api_url"),
		c.Simple("auth", "auth"),
		c.Simple("username", "username", c.SkipZeroValue),
		c.Simple("password", "password", c.SkipZeroValue),
		c.Simple("bearerToken", "bearer_token", c.SkipZeroValue),
		c.Simple("apiKeyName", "api_key_name", c.SkipZeroValue),
		c.Simple("apiKeyValue", "api_key_value", c.SkipZeroValue),
		c.Simple("xmlRootKey", "xml_root_key", c.SkipZeroValue),
		c.Simple("method", "method"),
		c.Simple("format", "format"),
		c.ArrayWithObjects("propertiesMapping", "properties_mapping", map[string]interface{}{
			"to":   "to",
			"from": "from",
		}),
		c.ArrayWithObjects("queryParams", "query_params", map[string]interface{}{
			"to":   "to",
			"from": "from",
		}),
		c.ArrayWithObjects("headers", "headers", map[string]interface{}{
			"to":   "to",
			"from": "from",
		}),
		c.ArrayWithObjects("pathParams", "path_params", map[string]interface{}{
			"path": "path",
		}),
		c.Simple("isBatchingEnabled", "is_batching_enabled"),
		c.Simple("maxBatchSize", "max_batch_size", c.SkipZeroValue),
		c.ArrayWithStrings("whitelistedEvents", "eventName", "event_filtering.0.whitelist"),
		c.ArrayWithStrings("blacklistedEvents", "eventName", "event_filtering.0.blacklist"),
		c.Discriminator("eventFilteringOption", c.DiscriminatorValues{
			"event_filtering.0.whitelist": "whitelistedEvents",
			"event_filtering.0.blacklist": "blacklistedEvents",
		}),
		c.Simple("isDefaultMapping", "is_default_mapping"),
		c.Simple("connectionMode.android", "connection_mode.0.android", c.SkipZeroValue),
		c.Simple("connectionMode.androidKotlin", "connection_mode.0.android_kotlin", c.SkipZeroValue),
		c.Simple("connectionMode.ios", "connection_mode.0.ios", c.SkipZeroValue),
		c.Simple("connectionMode.iosSwift", "connection_mode.0.ios_swift", c.SkipZeroValue),
		c.Simple("connectionMode.web", "connection_mode.0.web", c.SkipZeroValue),
		c.Simple("connectionMode.unity", "connection_mode.0.unity", c.SkipZeroValue),
		c.Simple("connectionMode.amp", "connection_mode.0.amp", c.SkipZeroValue),
		c.Simple("connectionMode.cloud", "connection_mode.0.cloud", c.SkipZeroValue),
		c.Simple("connectionMode.warehouse", "connection_mode.0.warehouse", c.SkipZeroValue),
		c.Simple("connectionMode.reactnative", "connection_mode.0.reactnative", c.SkipZeroValue),
		c.Simple("connectionMode.flutter", "connection_mode.0.flutter", c.SkipZeroValue),
		c.Simple("connectionMode.cordova", "connection_mode.0.cordova", c.SkipZeroValue),
		c.Simple("connectionMode.shopify", "connection_mode.0.shopify", c.SkipZeroValue),
	}

	properties = append(properties, commonProperties...)

	schema := map[string]*schema.Schema{
		"api_url": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Enter the base URL for your HTTP endpoint.",
			ValidateDiagFunc: c.ValidateAll(
				c.StringMatchesRegexp(httpURLPattern),
				c.StringNotMatchesRegexp(httpNgrokPattern),
				c.StringNotMatchesRegexp(httpLocalhostPattern),
			),
		},
		"auth": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "noAuth",
			Description:      "Select the authentication method used for the HTTP request.",
			ValidateDiagFunc: c.StringMatchesRegexp("^(noAuth|basicAuth|bearerTokenAuth|apiKeyAuth)$"),
		},
		"username": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Enter the username for basic authentication.",
			ValidateDiagFunc: c.StringMatchesRegexp("^(.{1,100})$"),
		},
		"password": {
			Type:             schema.TypeString,
			Optional:         true,
			Sensitive:        true,
			Description:      "Enter the password for basic authentication.",
			ValidateDiagFunc: c.StringMatchesRegexp("^(.{0,100})$"),
		},
		"bearer_token": {
			Type:             schema.TypeString,
			Optional:         true,
			Sensitive:        true,
			Description:      "Enter the bearer token used for authorization.",
			ValidateDiagFunc: c.StringMatchesRegexp("^(.{1,255})$"),
		},
		"api_key_name": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Enter the header name for API key authentication.",
			ValidateDiagFunc: c.StringMatchesRegexp("^(\\S{1,100})$"),
		},
		"api_key_value": {
			Type:             schema.TypeString,
			Optional:         true,
			Sensitive:        true,
			Description:      "Enter the value for API key authentication.",
			ValidateDiagFunc: c.StringMatchesRegexp("^(.{1,100})$"),
		},
		"xml_root_key": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Enter the XML root key used as a common prefix for mapped fields.",
			ValidateDiagFunc: c.StringMatchesRegexp("^(.{0,100})$"),
		},
		"method": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "POST",
			Description:      "Select the HTTP method to use when sending requests.",
			ValidateDiagFunc: c.StringMatchesRegexp("^(POST|PUT|PATCH|GET|DELETE)$"),
		},
		"format": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "JSON",
			Description:      "Select the body format for the outgoing request.",
			ValidateDiagFunc: c.StringMatchesRegexp("^(JSON|XML|FORM)$"),
		},
		"properties_mapping": {
			Type:        schema.TypeList,
			Optional:    true,
			ConfigMode:  schema.SchemaConfigModeAttr,
			Description: "Map the outgoing request payload using JSON path keys and values from the RudderStack payload or constants.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"to": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(?:$|\\$(?:\\.|(\\.(\\w+|\\*)|\\[\\d+\\]|\\[('[^\\s']+'|\"[^\\s\"]+\")\\]|\\[\\*\\]|\\.\\w+\\(\\))*)$)"),
					},
					"from": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^\\$(?:\\.|(?:\\.\\.(\\w+|\\*)|\\.(\\w+|\\*)|\\[\\d+\\]|\\[('\\w+'|\"\\w+\")\\]|\\[\\*\\]|\\.\\w+\\(\\))*)$|^[A-Za-z0-9!#$%&'*+.^_`|~-]{0,100}$"),
					},
				},
			},
		},
		"query_params": {
			Type:        schema.TypeList,
			Optional:    true,
			ConfigMode:  schema.SchemaConfigModeAttr,
			Description: "Map query parameter keys to values from the RudderStack payload or constants.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"to": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^[A-Za-z0-9!#$%&'*+.^_`|~\\- ]{0,100}$"),
					},
					"from": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^\\$(?:\\.|(\\.(\\w+|\\*)|\\[\\d+\\]|\\[('\\w+'|\"\\w+\")\\]|\\[\\*\\]|\\.\\w+\\(\\))*)$|^[A-Za-z0-9!#$%&'*+.^_`|~\\- ]{0,100}$"),
					},
				},
			},
		},
		"headers": {
			Type:        schema.TypeList,
			Optional:    true,
			ConfigMode:  schema.SchemaConfigModeAttr,
			Description: "Build custom request headers using constants or values from the RudderStack payload.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"to": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^[A-Za-z0-9!#$%&'*+.^_`|~-]{0,100}$"),
					},
					"from": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^\\$(?:\\.|(\\.(\\w+|\\*)|\\[\\d+\\]|\\[('\\w+'|\"\\w+\")\\]|\\[\\*\\]|\\.\\w+\\(\\))*)$|^[A-Za-z0-9!#$%&'*+.^_`|~\\- /\\\\]{0,100}$"),
					},
				},
			},
		},
		"path_params": {
			Type:        schema.TypeList,
			Optional:    true,
			ConfigMode:  schema.SchemaConfigModeAttr,
			Description: "Enter path parameters in sequence using JSON path values or constants.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"path": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^\\$(?:\\.|(\\.(\\w+|\\*)|\\[\\d+\\]|\\[('\\w+'|\"\\w+\")\\]|\\[\\*\\]|\\.\\w+\\(\\))*)\\/?$|^[A-Za-z0-9!#$%&'*+.^_`|~-]{0,100}\\/?$"),
					},
				},
			},
		},
		"is_batching_enabled": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable batching for JSON payloads.",
		},
		"max_batch_size": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Specify the maximum number of events to include in a batch, up to 100.",
			ValidateDiagFunc: c.StringMatchesRegexp("^([1-9][0-9]{0,1}|100)$"),
		},
		"event_filtering": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Choose whether to allowlist or denylist client-side events.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"whitelist": {
						Type:         schema.TypeList,
						Optional:     true,
						Description:  "Enter the event names to be allowlisted.",
						ExactlyOneOf: []string{"config.0.event_filtering.0.whitelist", "config.0.event_filtering.0.blacklist"},
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"blacklist": {
						Type:         schema.TypeList,
						Optional:     true,
						Description:  "Enter the event names to be denylisted.",
						ExactlyOneOf: []string{"config.0.event_filtering.0.whitelist", "config.0.event_filtering.0.blacklist"},
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
		"is_default_mapping": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Send the event payload as-is when enabled.",
		},
		"connection_mode": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Configure the connection mode for each supported source type.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"android": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$"),
					},
					"android_kotlin": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$"),
					},
					"ios": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$"),
					},
					"ios_swift": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$"),
					},
					"web": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$"),
					},
					"unity": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$"),
					},
					"amp": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$"),
					},
					"cloud": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$"),
					},
					"warehouse": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$"),
					},
					"reactnative": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$"),
					},
					"flutter": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$"),
					},
					"cordova": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$"),
					},
					"shopify": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$"),
					},
				},
			},
		},
	}

	for key, value := range commonSchema {
		schema[key] = value
	}

	c.Destinations.Register("http", c.ConfigMeta{
		APIType:      "HTTP",
		Properties:   properties,
		ConfigSchema: schema,
	})
}
