package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	c.Destinations.Register("kafka", c.ConfigMeta{
		APIType: "KAFKA",
		Properties: []c.ConfigProperty{
			c.Simple("hostName", "host_name"),
			c.Simple("port", "port"),
			c.Simple("topic", "topic"),
			c.Simple("sslEnabled", "ssl_enabled", c.SkipZeroValue),
			c.Simple("caCertificate", "ca_certificate", c.SkipZeroValue),
			c.Simple("useSASL", "use_sasl", c.SkipZeroValue),
			c.Simple("saslType", "sasl_type", c.SkipZeroValue),
			c.Simple("username", "username", c.SkipZeroValue),
			c.Simple("password", "password", c.SkipZeroValue),
			c.Simple("convertToAvro", "convert_to_avro", c.SkipZeroValue),
			c.ArrayWithObjects("avroSchema", "avro_schema", map[string]string{
				"schemaId": "schema_id",
				"schema":   "schema",
			}),
			c.Simple("embedAvroSchemaID", "embed_avro_schema_id", c.SkipZeroValue),
			c.Simple("enableMultiTopic", "enable_multi_topic", c.SkipZeroValue),
			c.ArrayWithObjects("eventTypeToTopicMap", "event_type_to_topic_map", map[string]string{
				"from": "from",
				"to":   "to",
			}),
			c.ArrayWithObjects("eventToTopicMap", "event_to_topic_map", map[string]string{
				"from": "from",
				"to":   "to",
			}),
			c.ArrayWithStrings("oneTrustCookieCategories", "oneTrustCookieCategory", "onetrust_cookie_categories"),
		},
		ConfigSchema: map[string]*schema.Schema{
			"host_name": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "The host name of your Kafka service.",
			},
			"port": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "The port number associated with the Kafka service.",
			},
			"topic": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "The topic name in your Kafka service where the data will be sent.",
			},
			"ssl_enabled": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:		  true,
				Description:      "Whether to enable SSL for the Kafka service.",
			},
			"ca_certificate": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "The CA certificate for the Kafka service.",
			},
			"use_sasl": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:		  false,
				Description:      "Whether to use SASL for the Kafka service.",
			},
			"sasl_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "The SASL type for the Kafka service.",

			},
			"username": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "The username for the Kafka service.",

			},
			"password": {
				Type:             schema.TypeString,
				Optional:         true,
				Sensitive:		  true,
				Description:      "The password for the Kafka service.",

			},
			"convert_to_avro": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:		  false,
				Description:      "Whether to convert the data to Avro format.",
			},
			"avro_schema": {
				Type:             schema.TypeList,
				Optional:         true,
				Description:      "The Avro schema for the Kafka service.",
				ConfigMode:  schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"schema_id": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "The schema ID for the Avro schema.",
						},
						"schema": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "The schema for the Avro schema.",
						},
					},
				},
			},
			"embed_avro_schema_id": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:		  false,
				Description:      "Whether to embed the Avro schema ID.",
			},
			"enable_multi_topic": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:		  false,
				Description:      "Whether to enable multi-topic.",
			},
			"event_type_to_topic_map": {
				Type:             schema.TypeList,
				Optional:         true,
				Description:      "The event type to topic map for the Kafka service.",
				ConfigMode:  schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "The event type to topic map from.",
						},
						"to": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "The event type to topic map to.",
						},
					},
				},
			},
			"event_to_topic_map": {
				Type:             schema.TypeList,
				Optional:         true,
				Description:      "The event to topic map for the Kafka service.",
				ConfigMode:  schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "The event to topic map from.",
						},
						"to": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "The event to topic map to.",
						},
					},
				},
			},
			"onetrust_cookie_categories" : {
				Type:             schema.TypeList,
				Optional:         true,
				Description:      "Specify the OneTrust category name for mapping the OneTrust consent settings to RudderStack's consent purposes.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
			},
			},
		},
	})
}
