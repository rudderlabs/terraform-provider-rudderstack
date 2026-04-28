package retl

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/rudderlabs/rudder-iac/api/client/retl"
)

// ResourceModel returns the schema for `rudderstack_retl_source_model`.
func ResourceModel() *schema.Resource {
	return buildResource(typeAdapter{
		sourceType:      retl.ModelSourceType,
		configSchema:    modelConfigSchema,
		marshalConfig:   modelMarshalConfig,
		unmarshalConfig: modelUnmarshalConfig,
	}, "A RETL source backed by a SQL model (custom query).")
}

func modelConfigSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"primary_key": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Column used as the primary key for change tracking.",
		},
		"sql": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "SQL query that defines the model.",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Optional human-readable description of the model.",
		},
	}
}

func modelMarshalConfig(cfg map[string]any) (retl.RETLConfig, error) {
	return retl.RETLSQLModelConfig{
		PrimaryKey:  stringField(cfg, "primary_key"),
		Sql:         stringField(cfg, "sql"),
		Description: stringField(cfg, "description"),
	}, nil
}

func modelUnmarshalConfig(raw retl.RETLConfig) ([]map[string]interface{}, error) {
	cfg, err := retl.DecodeConfig[retl.RETLSQLModelConfig](raw)
	if err != nil {
		return nil, err
	}
	return []map[string]interface{}{{
		"primary_key": cfg.PrimaryKey,
		"sql":         cfg.Sql,
		"description": cfg.Description,
	}}, nil
}
