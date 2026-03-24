package rudderstack

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/rudderlabs/rudder-api-go/client"
)

func resourceRETLSource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Human readable name of the RETL source.",
			},
			"source_definition_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Source definition name, e.g. \"bigquery\".",
			},
			"source_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Source type: \"table\" or \"model\".",
			},
			"account_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the warehouse account to use.",
			},
			"primary_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Primary key column name for the source.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "An enabled RETL source allows data to be read from it.",
			},
			"write_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The write key that identifies the source in RudderStack data plane.",
			},
			"config": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Source-specific configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"schema_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Dataset/schema name for table type sources.",
						},
						"table_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Table name for table type sources.",
						},
						"sql": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "SQL query for model type sources.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description for model type sources.",
						},
					},
				},
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Time when the resource was created, in ISO 8601 format.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Time when the resource was last updated, in ISO 8601 format.",
			},
		},
		CreateContext: resourceRETLSourceCreate,
		ReadContext:   resourceRETLSourceRead,
		UpdateContext: resourceRETLSourceUpdate,
		DeleteContext: resourceRETLSourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRETLSourceImportState,
		},
	}
}

func resourceRETLSourceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(*Client)
	if !ok {
		return diag.FromErr(fmt.Errorf("API client is not configured"))
	}

	source := &client.RETLSource{}
	if err := populateRETLSourceFromState(source, d); err != nil {
		return diag.FromErr(err)
	}

	source, err := c.RETLSources.Create(ctx, source)
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not create RETL source: %w", err))
	}

	d.SetId(source.ID)

	return resourceRETLSourceRead(ctx, d, m)
}

func resourceRETLSourceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(*Client)
	if !ok {
		return diag.FromErr(fmt.Errorf("API client is not configured"))
	}

	id := d.Id()

	source, err := c.RETLSources.Get(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := storeRETLSourceToState(source, d); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func resourceRETLSourceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(*Client)
	if !ok {
		return diag.FromErr(fmt.Errorf("API client is not configured"))
	}

	source := &client.RETLSource{}
	if err := populateRETLSourceFromState(source, d); err != nil {
		return diag.FromErr(err)
	}

	source, err := c.RETLSources.Update(ctx, source)
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not update RETL source: %w", err))
	}

	d.SetId(source.ID)

	return resourceRETLSourceRead(ctx, d, m)
}

func resourceRETLSourceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(*Client)
	if !ok {
		return diag.FromErr(fmt.Errorf("API client is not configured"))
	}

	if err := c.RETLSources.Delete(ctx, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diag.Diagnostics{}
}

func resourceRETLSourceImportState(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	diagnostics := resourceRETLSourceRead(ctx, d, m)
	if diagnostics.HasError() {
		for _, diagnostic := range diagnostics {
			if diagnostic.Severity == diag.Error {
				return nil, fmt.Errorf("could not import RETL source: %s", diagnostic.Summary)
			}
		}
	}
	return []*schema.ResourceData{d}, nil
}

func populateRETLSourceFromState(source *client.RETLSource, d *schema.ResourceData) error {
	source.ID = d.Id()
	source.Name = d.Get("name").(string)
	source.SourceDefinitionName = d.Get("source_definition_name").(string)
	source.SourceType = d.Get("source_type").(string)
	source.AccountID = d.Get("account_id").(string)
	source.IsEnabled = d.Get("enabled").(bool)

	if pk, ok := d.GetOk("primary_key"); ok {
		source.PrimaryKey = pk.(string)
	}

	if configList, ok := d.GetOk("config"); ok {
		configs := configList.([]interface{})
		if len(configs) > 0 && configs[0] != nil {
			configMap := configs[0].(map[string]interface{})
			apiConfig := make(map[string]interface{})

			if v, ok := configMap["schema_name"]; ok && v.(string) != "" {
				apiConfig["schema"] = v.(string)
			}
			if v, ok := configMap["table_name"]; ok && v.(string) != "" {
				apiConfig["table"] = v.(string)
			}
			if v, ok := configMap["sql"]; ok && v.(string) != "" {
				apiConfig["sql"] = v.(string)
			}
			if v, ok := configMap["description"]; ok && v.(string) != "" {
				apiConfig["description"] = v.(string)
			}

			configJSON, err := json.Marshal(apiConfig)
			if err != nil {
				return err
			}
			source.Config = json.RawMessage(configJSON)
		}
	}

	return nil
}

func storeRETLSourceToState(source *client.RETLSource, d *schema.ResourceData) error {
	d.SetId(source.ID)

	if err := d.Set("name", source.Name); err != nil {
		return err
	}
	if err := d.Set("source_definition_name", source.SourceDefinitionName); err != nil {
		return err
	}
	if err := d.Set("source_type", source.SourceType); err != nil {
		return err
	}
	if err := d.Set("account_id", source.AccountID); err != nil {
		return err
	}
	if err := d.Set("primary_key", source.PrimaryKey); err != nil {
		return err
	}
	if err := d.Set("enabled", source.IsEnabled); err != nil {
		return err
	}
	if err := d.Set("write_key", source.WriteKey); err != nil {
		return err
	}

	if source.CreatedAt != nil {
		createdAt := source.CreatedAt.Format(time.RFC3339)
		if err := d.Set("created_at", createdAt); err != nil {
			return err
		}
	}
	if source.UpdatedAt != nil {
		updatedAt := source.UpdatedAt.Format(time.RFC3339)
		if err := d.Set("updated_at", updatedAt); err != nil {
			return err
		}
	}

	// Parse config JSON back to terraform state
	if len(source.Config) > 0 {
		var apiConfig map[string]interface{}
		if err := json.Unmarshal(source.Config, &apiConfig); err != nil {
			return err
		}

		configBlock := map[string]interface{}{
			"schema_name": "",
			"table_name":  "",
			"sql":         "",
			"description": "",
		}

		if v, ok := apiConfig["schema"]; ok {
			configBlock["schema_name"] = v
		}
		if v, ok := apiConfig["table"]; ok {
			configBlock["table_name"] = v
		}
		if v, ok := apiConfig["sql"]; ok {
			configBlock["sql"] = v
		}
		if v, ok := apiConfig["description"]; ok {
			configBlock["description"] = v
		}

		if err := d.Set("config", []interface{}{configBlock}); err != nil {
			return err
		}
	}

	return nil
}
