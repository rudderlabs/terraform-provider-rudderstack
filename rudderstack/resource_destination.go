package rudderstack

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudderlabs/rudder-api-go/client"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func resourceDestination(cm configs.ConfigMeta) *schema.Resource {
	return &schema.Resource{
		Schema:        resourceDestinationSchema(cm),
		CreateContext: resourceDestinationCreate(cm),
		ReadContext:   resourceDestinationRead(cm),
		UpdateContext: resourceDestinationUpdate(cm),
		DeleteContext: resourceDestinationDelete(cm),
		Importer: &schema.ResourceImporter{
			StateContext: resourceDestinationImportState(cm),
		},
	}
}

func resourceDestinationSchema(cm configs.ConfigMeta) map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"updated_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}

	if !cm.SkipConfig {
		s["config"] = &schema.Schema{
			Type:     schema.TypeList,
			Optional: cm.SkipConfig,
			Required: !cm.SkipConfig,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: cm.ConfigSchema,
			},
		}
	}

	return s
}

func resourceDestinationCreate(cm configs.ConfigMeta) schema.CreateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		c, ok := m.(*Client)
		if !ok {
			return diag.FromErr(fmt.Errorf("API client is not configured"))
		}

		destination := &client.Destination{}
		err := populateDestinationFromState(cm, destination, d)
		if err != nil {
			return diag.FromErr(err)
		}

		destination, err = c.Destinations.Create(ctx, destination)
		if err != nil {
			return diag.FromErr(fmt.Errorf("could not create destination: %w", err))
		}

		d.SetId(destination.ID)

		return resourceDestinationRead(cm)(ctx, d, m)
	}
}

func resourceDestinationRead(cm configs.ConfigMeta) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		c, ok := m.(*Client)
		if !ok {
			return diag.FromErr(fmt.Errorf("API client is not configured"))
		}

		id := d.Id()

		destination, err := c.Destinations.Get(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}

		if err := storeDestinationToState(cm, destination, d); err != nil {
			return diag.FromErr(err)
		}

		return diag.Diagnostics{}
	}
}

func resourceDestinationUpdate(cm configs.ConfigMeta) schema.UpdateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		c, ok := m.(*Client)
		if !ok {
			return diag.FromErr(fmt.Errorf("API client is not configured"))
		}

		destination := &client.Destination{}
		err := populateDestinationFromState(cm, destination, d)
		if err != nil {
			return diag.FromErr(err)
		}

		destination, err = c.Destinations.Update(ctx, destination)
		if err != nil {
			return diag.FromErr(fmt.Errorf("could not create destination: %w", err))
		}

		d.SetId(destination.ID)

		return resourceDestinationRead(cm)(ctx, d, m)
	}
}

func resourceDestinationDelete(cm configs.ConfigMeta) schema.DeleteContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		c, ok := m.(*Client)
		if !ok {
			return diag.FromErr(fmt.Errorf("API client is not configured"))
		}

		if err := c.Destinations.Delete(ctx, d.Id()); err != nil {
			return diag.FromErr(err)
		}

		d.SetId("")
		return diag.Diagnostics{}
	}
}

func resourceDestinationImportState(cm configs.ConfigMeta) schema.StateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
		diagnostics := resourceDestinationRead(cm)(ctx, d, m)
		if diagnostics.HasError() {
			for _, diagnostic := range diagnostics {
				if diagnostic.Severity == diag.Error {
					return nil, fmt.Errorf("could not import connection: %s", diagnostic.Summary)
				}
			}
		}
		return []*schema.ResourceData{d}, nil
	}
}

func populateDestinationFromState(cm configs.ConfigMeta, destination *client.Destination, d *schema.ResourceData) error {
	destination.ID = d.Id()
	destination.Type = cm.APIType
	destination.Name = d.Get("name").(string)
	destination.IsEnabled = d.Get("enabled").(bool)

	if c := d.Get("config.0"); c != nil {
		state, err := json.Marshal(c)
		if err != nil {
			return err
		}
		apiConfig, err := cm.StateToAPI(string(state))
		if err != nil {
			return err
		}
		destination.Config = json.RawMessage(apiConfig)
	}

	return nil
}

func storeDestinationToState(cm configs.ConfigMeta, destination *client.Destination, d *schema.ResourceData) error {
	d.SetId(destination.ID)
	if err := d.Set("name", destination.Name); err != nil {
		return err
	}
	if err := d.Set("enabled", destination.IsEnabled); err != nil {
		return err
	}
	if destination.CreatedAt != nil {
		createdAt := destination.CreatedAt.Format(time.RFC3339)
		if err := d.Set("created_at", createdAt); err != nil {
			return err
		}
	}
	if destination.UpdatedAt != nil {
		updatedAt := destination.UpdatedAt.Format(time.RFC3339)
		if err := d.Set("updated_at", updatedAt); err != nil {
			return err
		}
	}

	state, err := cm.APIToState(string(destination.Config))
	if err != nil {
		return err
	}

	properties := make(map[string]interface{})
	if err := json.Unmarshal([]byte(state), &properties); err != nil {
		return err
	}

	if len(properties) > 0 {
		if err := d.Set("config", []interface{}{properties}); err != nil {
			return err
		}
	} else {
		if err := d.Set("config", []interface{}{}); err != nil {
			return err
		}
	}

	return nil
}
