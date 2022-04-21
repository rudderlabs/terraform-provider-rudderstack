package rudderstack

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudderlabs/rudder-api-go/client"
)

func resourceConnection() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"destination_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
		CreateContext: resourceConnectionCreate,
		ReadContext:   resourceConnectionRead,
		UpdateContext: resourceConnectionUpdate,
		DeleteContext: resourceConnectionDelete,
	}
}

func resourceConnectionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(*Client)
	if !ok {
		return diag.FromErr(fmt.Errorf("API client is not configured"))
	}

	connection := &client.Connection{}
	err := populateConnectionFromState(connection, d)
	if err != nil {
		return diag.FromErr(err)
	}

	connection, err = c.Connections.Create(ctx, connection)
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not create connection: %w", err))
	}

	d.SetId(connection.ID)

	return resourceConnectionRead(ctx, d, m)
}

func resourceConnectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(*Client)
	if !ok {
		return diag.FromErr(fmt.Errorf("API client is not configured"))
	}

	id := d.Id()

	connection, err := c.Connections.Get(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	storeConnectionToState(connection, d)

	return diag.Diagnostics{}
}

func resourceConnectionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(*Client)
	if !ok {
		return diag.FromErr(fmt.Errorf("API client is not configured"))
	}

	connection := &client.Connection{}
	err := populateConnectionFromState(connection, d)
	if err != nil {
		return diag.FromErr(err)
	}

	connection, err = c.Connections.Update(ctx, connection)
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not create source: %w", err))
	}

	d.SetId(connection.ID)

	return resourceConnectionRead(ctx, d, m)
}

func resourceConnectionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(*Client)
	if !ok {
		return diag.FromErr(fmt.Errorf("API client is not configured"))
	}

	if err := c.Connections.Delete(ctx, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diag.Diagnostics{}
}

func populateConnectionFromState(connection *client.Connection, d *schema.ResourceData) error {
	connection.ID = d.Id()
	connection.SourceID = d.Get("source_id").(string)
	connection.DestinationID = d.Get("destination_id").(string)
	connection.IsEnabled = d.Get("enabled").(bool)

	return nil
}

func storeConnectionToState(connection *client.Connection, d *schema.ResourceData) error {
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

	return nil
}
