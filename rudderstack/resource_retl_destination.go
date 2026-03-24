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

func resourceRETLDestination() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Human readable name of the RETL destination.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Destination type identifier.",
			},
			"config": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "JSON-encoded destination configuration.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "An enabled destination allows data to be sent to it.",
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
		CreateContext: resourceRETLDestinationCreate,
		ReadContext:   resourceRETLDestinationRead,
		UpdateContext: resourceRETLDestinationUpdate,
		DeleteContext: resourceRETLDestinationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRETLDestinationImportState,
		},
	}
}

func resourceRETLDestinationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(*Client)
	if !ok {
		return diag.FromErr(fmt.Errorf("API client is not configured"))
	}

	destination := &client.Destination{}
	if err := populateRETLDestinationFromState(destination, d); err != nil {
		return diag.FromErr(err)
	}

	destination, err := c.Destinations.Create(ctx, destination)
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not create RETL destination: %w", err))
	}

	d.SetId(destination.ID)

	return resourceRETLDestinationRead(ctx, d, m)
}

func resourceRETLDestinationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(*Client)
	if !ok {
		return diag.FromErr(fmt.Errorf("API client is not configured"))
	}

	id := d.Id()

	destination, err := c.Destinations.Get(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := storeRETLDestinationToState(destination, d); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func resourceRETLDestinationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(*Client)
	if !ok {
		return diag.FromErr(fmt.Errorf("API client is not configured"))
	}

	destination := &client.Destination{}
	if err := populateRETLDestinationFromState(destination, d); err != nil {
		return diag.FromErr(err)
	}

	destination, err := c.Destinations.Update(ctx, destination)
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not update RETL destination: %w", err))
	}

	d.SetId(destination.ID)

	return resourceRETLDestinationRead(ctx, d, m)
}

func resourceRETLDestinationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

func resourceRETLDestinationImportState(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	diagnostics := resourceRETLDestinationRead(ctx, d, m)
	if diagnostics.HasError() {
		for _, diagnostic := range diagnostics {
			if diagnostic.Severity == diag.Error {
				return nil, fmt.Errorf("could not import RETL destination: %s", diagnostic.Summary)
			}
		}
	}
	return []*schema.ResourceData{d}, nil
}

func populateRETLDestinationFromState(destination *client.Destination, d *schema.ResourceData) error {
	destination.ID = d.Id()
	destination.Name = d.Get("name").(string)
	destination.Type = d.Get("type").(string)
	destination.IsEnabled = d.Get("enabled").(bool)

	configStr := d.Get("config").(string)
	destination.Config = json.RawMessage(configStr)

	return nil
}

func storeRETLDestinationToState(destination *client.Destination, d *schema.ResourceData) error {
	d.SetId(destination.ID)

	if err := d.Set("name", destination.Name); err != nil {
		return err
	}
	if err := d.Set("type", destination.Type); err != nil {
		return err
	}
	if err := d.Set("enabled", destination.IsEnabled); err != nil {
		return err
	}

	if len(destination.Config) > 0 {
		if err := d.Set("config", string(destination.Config)); err != nil {
			return err
		}
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

	return nil
}
