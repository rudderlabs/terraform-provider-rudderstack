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

func resourceRETLConnection() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the RETL source.",
			},
			"destination_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the destination.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "An enabled connection allows data to be transferred.",
			},
			"schedule": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Sync schedule configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Schedule type: \"basic\" or \"cron\".",
						},
						"every": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Sync interval in minutes (for basic schedule type).",
						},
						"cron": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Cron expression (for cron schedule type).",
						},
					},
				},
			},
			"sync_behaviour": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sync behaviour, e.g. \"upsert\" or \"mirror\".",
			},
			"event": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Event name for the sync.",
			},
			"event_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Event type, e.g. \"track\" or \"identify\".",
			},
			"mappings": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "JSON-encoded field mappings.",
			},
			"constants": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "JSON-encoded constant values.",
			},
			"query": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SQL query for the sync.",
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
		CreateContext: resourceRETLConnectionCreate,
		ReadContext:   resourceRETLConnectionRead,
		UpdateContext: resourceRETLConnectionUpdate,
		DeleteContext: resourceRETLConnectionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRETLConnectionImportState,
		},
	}
}

func resourceRETLConnectionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(*Client)
	if !ok {
		return diag.FromErr(fmt.Errorf("API client is not configured"))
	}

	connection := &client.RETLConnection{}
	if err := populateRETLConnectionFromState(connection, d); err != nil {
		return diag.FromErr(err)
	}

	connection, err := c.RETLConnections.Create(ctx, connection)
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not create RETL connection: %w", err))
	}

	d.SetId(connection.ID)

	return resourceRETLConnectionRead(ctx, d, m)
}

func resourceRETLConnectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(*Client)
	if !ok {
		return diag.FromErr(fmt.Errorf("API client is not configured"))
	}

	id := d.Id()

	connection, err := c.RETLConnections.Get(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := storeRETLConnectionToState(connection, d); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func resourceRETLConnectionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(*Client)
	if !ok {
		return diag.FromErr(fmt.Errorf("API client is not configured"))
	}

	connection := &client.RETLConnection{}
	if err := populateRETLConnectionFromState(connection, d); err != nil {
		return diag.FromErr(err)
	}

	connection, err := c.RETLConnections.Update(ctx, connection)
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not update RETL connection: %w", err))
	}

	d.SetId(connection.ID)

	return resourceRETLConnectionRead(ctx, d, m)
}

func resourceRETLConnectionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(*Client)
	if !ok {
		return diag.FromErr(fmt.Errorf("API client is not configured"))
	}

	if err := c.RETLConnections.Delete(ctx, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diag.Diagnostics{}
}

func resourceRETLConnectionImportState(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	diagnostics := resourceRETLConnectionRead(ctx, d, m)
	if diagnostics.HasError() {
		for _, diagnostic := range diagnostics {
			if diagnostic.Severity == diag.Error {
				return nil, fmt.Errorf("could not import RETL connection: %s", diagnostic.Summary)
			}
		}
	}
	return []*schema.ResourceData{d}, nil
}

func populateRETLConnectionFromState(connection *client.RETLConnection, d *schema.ResourceData) error {
	connection.ID = d.Id()
	connection.SourceID = d.Get("source_id").(string)
	connection.DestinationID = d.Get("destination_id").(string)
	connection.IsEnabled = d.Get("enabled").(bool)

	// Build the config object from schedule and other fields
	sourceConfig := make(map[string]interface{})

	// Parse schedule block
	if scheduleList, ok := d.GetOk("schedule"); ok {
		schedules := scheduleList.([]interface{})
		if len(schedules) > 0 && schedules[0] != nil {
			scheduleMap := schedules[0].(map[string]interface{})
			schedule := make(map[string]interface{})

			if v, ok := scheduleMap["type"]; ok {
				schedule["type"] = v.(string)
			}
			if v, ok := scheduleMap["every"]; ok && v.(int) > 0 {
				schedule["every"] = v.(int)
			}
			if v, ok := scheduleMap["cron"]; ok && v.(string) != "" {
				schedule["cron"] = v.(string)
			}

			sourceConfig["schedule"] = schedule
		}
	}

	if v, ok := d.GetOk("sync_behaviour"); ok {
		sourceConfig["syncBehaviour"] = v.(string)
	}
	if v, ok := d.GetOk("event"); ok {
		sourceConfig["event"] = v.(string)
	}
	if v, ok := d.GetOk("event_type"); ok {
		sourceConfig["eventType"] = v.(string)
	}
	if v, ok := d.GetOk("query"); ok {
		sourceConfig["query"] = v.(string)
	}

	if v, ok := d.GetOk("mappings"); ok && v.(string) != "" {
		var mappings interface{}
		if err := json.Unmarshal([]byte(v.(string)), &mappings); err != nil {
			return fmt.Errorf("could not parse mappings JSON: %w", err)
		}
		sourceConfig["mappings"] = mappings
	}

	if v, ok := d.GetOk("constants"); ok && v.(string) != "" {
		var constants interface{}
		if err := json.Unmarshal([]byte(v.(string)), &constants); err != nil {
			return fmt.Errorf("could not parse constants JSON: %w", err)
		}
		sourceConfig["constants"] = constants
	}

	config := map[string]interface{}{
		"source": sourceConfig,
	}

	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}
	connection.Config = json.RawMessage(configJSON)

	return nil
}

func storeRETLConnectionToState(connection *client.RETLConnection, d *schema.ResourceData) error {
	d.SetId(connection.ID)

	if err := d.Set("source_id", connection.SourceID); err != nil {
		return err
	}
	if err := d.Set("destination_id", connection.DestinationID); err != nil {
		return err
	}
	if err := d.Set("enabled", connection.IsEnabled); err != nil {
		return err
	}

	if connection.CreatedAt != nil {
		createdAt := connection.CreatedAt.Format(time.RFC3339)
		if err := d.Set("created_at", createdAt); err != nil {
			return err
		}
	}
	if connection.UpdatedAt != nil {
		updatedAt := connection.UpdatedAt.Format(time.RFC3339)
		if err := d.Set("updated_at", updatedAt); err != nil {
			return err
		}
	}

	// Parse config JSON and extract fields into state
	if len(connection.Config) > 0 {
		var config map[string]interface{}
		if err := json.Unmarshal(connection.Config, &config); err != nil {
			return err
		}

		if sourceConfig, ok := config["source"].(map[string]interface{}); ok {
			if schedule, ok := sourceConfig["schedule"].(map[string]interface{}); ok {
				scheduleBlock := map[string]interface{}{
					"type": "",
					"every": 0,
					"cron":  "",
				}
				if v, ok := schedule["type"]; ok {
					scheduleBlock["type"] = v
				}
				if v, ok := schedule["every"]; ok {
					// JSON numbers decode as float64
					switch n := v.(type) {
					case float64:
						scheduleBlock["every"] = int(n)
					case int:
						scheduleBlock["every"] = n
					}
				}
				if v, ok := schedule["cron"]; ok {
					scheduleBlock["cron"] = v
				}

				if err := d.Set("schedule", []interface{}{scheduleBlock}); err != nil {
					return err
				}
			}

			if v, ok := sourceConfig["syncBehaviour"]; ok {
				if err := d.Set("sync_behaviour", v); err != nil {
					return err
				}
			}
			if v, ok := sourceConfig["event"]; ok {
				if err := d.Set("event", v); err != nil {
					return err
				}
			}
			if v, ok := sourceConfig["eventType"]; ok {
				if err := d.Set("event_type", v); err != nil {
					return err
				}
			}
			if v, ok := sourceConfig["query"]; ok {
				if err := d.Set("query", v); err != nil {
					return err
				}
			}

			if v, ok := sourceConfig["mappings"]; ok {
				mappingsJSON, err := json.Marshal(v)
				if err != nil {
					return err
				}
				if err := d.Set("mappings", string(mappingsJSON)); err != nil {
					return err
				}
			}

			if v, ok := sourceConfig["constants"]; ok {
				constantsJSON, err := json.Marshal(v)
				if err != nil {
					return err
				}
				if err := d.Set("constants", string(constantsJSON)); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
