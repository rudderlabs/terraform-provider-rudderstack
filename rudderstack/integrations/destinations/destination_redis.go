package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"web", "android", "ios", "unity", "reactnative", "flutter", "cordova", "amp", "cloud", "warehouse", "shopify"}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)

	properties := []c.ConfigProperty{
		c.Simple("address", "address"),
		c.Simple("password", "password", c.SkipZeroValue),
		c.Simple("clusterMode", "cluster_mode"),
		c.Simple("secure", "secure"),
		c.Simple("prefix", "prefix", c.SkipZeroValue),
		c.Simple("database", "database", c.SkipZeroValue),
		c.Simple("caCertificate", "ca_certificate", c.SkipZeroValue),
		c.Simple("skipVerify", "skip_verify"),
	}

	properties = append(properties, commonProperties...)

	schema := map[string]*schema.Schema{
		"address": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Enter the address associated with your Redis cluster.",
			// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"password": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "Enter the password associated with your Redis user.",
			// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"cluster_mode": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Use this setting to enable the Redis cluster mode.",
		},
		"secure": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Enable this setting if you want to send the data to Redis via SSL.",
		},
		"prefix": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "By default, RudderStack stores user traits with the key user:<user_id>. An extra prefix can be added in the destination configuration to distinguish all RudderStack-stored keys with a prefix.",
			// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"database": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "RudderStack stores the user traits in the default database of the Redis instance. A different database inside the Redis instance can be configured using this setting.",
			// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"ca_certificate": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Enter the certificate which needs to be verified while establishing a secure connection. Skip setting this if Root CA of your server can be verified with any client, e.g. AWS Elasticache.",
			// ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"skip_verify": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Enable this setting to skip the client's verification of the server's certificate chain and host name.",
		},
	}

	for key, value := range commonSchema {
		schema[key] = value
	}

	c.Destinations.Register("redis", c.ConfigMeta{
		APIType:      "REDIS",
		Properties:   properties,
		ConfigSchema: schema,
	})
}
