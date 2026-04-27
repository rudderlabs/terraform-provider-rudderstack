package retl

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/rudderlabs/rudder-iac/api/client/retl"
)

// ResourceTable returns the schema for `rudderstack_retl_source_table`.
func ResourceTable() *schema.Resource {
	return buildResource(typeAdapter{
		sourceType:      retl.TableSourceType,
		configSchema:    tableConfigSchema,
		marshalConfig:   tableMarshalConfig,
		unmarshalConfig: tableUnmarshalConfig,
	}, "A RETL source backed by a warehouse table (e.g. snowflake / bigquery / postgres).")
}

func tableConfigSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"primary_key": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Column used as the primary key for change tracking.",
		},
		"schema": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Warehouse schema (or dataset) the table lives in.",
		},
		"table": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Warehouse table name.",
		},
	}
}

func tableMarshalConfig(cfg map[string]interface{}) ([]byte, error) {
	return json.Marshal(retl.RETLTableConfig{
		PrimaryKey: stringField(cfg, "primary_key"),
		Schema:     stringField(cfg, "schema"),
		Table:      stringField(cfg, "table"),
	})
}

func tableUnmarshalConfig(raw []byte) ([]map[string]interface{}, error) {
	var cfg retl.RETLTableConfig
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &cfg); err != nil {
			return nil, err
		}
	}
	return []map[string]interface{}{{
		"primary_key": cfg.PrimaryKey,
		"schema":      cfg.Schema,
		"table":       cfg.Table,
	}}, nil
}
