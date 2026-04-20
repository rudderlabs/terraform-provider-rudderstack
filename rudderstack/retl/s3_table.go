package retl

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/rudderlabs/rudder-iac/api/client/retl"
)

const s3SourceDefinitionName = "s3"

// ResourceS3Table returns the schema for `rudderstack_retl_source_s3_table`.
func ResourceS3Table() *schema.Resource {
	return buildResource(typeAdapter{
		sourceType:                retl.TableSourceType,
		fixedSourceDefinitionName: s3SourceDefinitionName,
		configSchema:              s3TableConfigSchema,
		marshalConfig:             s3TableMarshalConfig,
		unmarshalConfig:           s3TableUnmarshalConfig,
	}, "A RETL source backed by an S3 bucket.")
}

func s3TableConfigSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"bucket_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name of the S3 bucket containing the source data.",
		},
		"object_prefix": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Optional object key prefix used to scope the source within the bucket.",
		},
	}
}

func s3TableMarshalConfig(cfg map[string]interface{}) ([]byte, error) {
	return json.Marshal(retl.RETLS3TableConfig{
		BucketName:   stringField(cfg, "bucket_name"),
		ObjectPrefix: stringField(cfg, "object_prefix"),
	})
}

func s3TableUnmarshalConfig(raw []byte) ([]map[string]interface{}, error) {
	var cfg retl.RETLS3TableConfig
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &cfg); err != nil {
			return nil, err
		}
	}
	return []map[string]interface{}{{
		"bucket_name":   cfg.BucketName,
		"object_prefix": cfg.ObjectPrefix,
	}}, nil
}
