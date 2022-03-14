package rudderstack

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudderlabs/rudder-api-go/client"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func resourceSource(cm configs.ConfigMeta) *schema.Resource {
	return &schema.Resource{
		Schema:        resourceSourceSchema(cm),
		CreateContext: resourceSourceCreate(cm),
		ReadContext:   resourceSourceRead(cm),
		UpdateContext: resourceSourceUpdate(cm),
		DeleteContext: resourceSourceDelete(cm),
	}
}

func resourceSourceSchema(cm configs.ConfigMeta) map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
		"config": {
			Type:     schema.TypeList,
			Optional: cm.OptionalConfig,
			Required: !cm.OptionalConfig,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: cm.ConfigSchema,
			},
		},
	}
}

func resourceSourceCreate(cm configs.ConfigMeta) schema.CreateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		c, ok := m.(*Client)
		if !ok {
			return diag.FromErr(fmt.Errorf("API client is not configured"))
		}

		source := &client.Source{}
		populateSourceFromState(cm, source, d)

		source, err := c.Sources.Create(ctx, source)
		if err != nil {
			return diag.FromErr(fmt.Errorf("could not create source: %w", err))
		}

		d.SetId(source.ID)

		return resourceSourceRead(cm)(ctx, d, m)
	}
}

func resourceSourceRead(cm configs.ConfigMeta) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		c, ok := m.(*Client)
		if !ok {
			return diag.FromErr(fmt.Errorf("API client is not configured"))
		}

		id := d.Id()

		source, err := c.Sources.Get(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}

		storeSourceToState(cm, source, d)

		return diag.Diagnostics{}
	}
}

func resourceSourceUpdate(cm configs.ConfigMeta) schema.UpdateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		c, ok := m.(*Client)
		if !ok {
			return diag.FromErr(fmt.Errorf("API client is not configured"))
		}

		source := &client.Source{}
		populateSourceFromState(cm, source, d)

		source, err := c.Sources.Update(ctx, source)
		if err != nil {
			return diag.FromErr(fmt.Errorf("could not create source: %w", err))
		}

		d.SetId(source.ID)

		return resourceSourceRead(cm)(ctx, d, m)
	}
}

func resourceSourceDelete(cm configs.ConfigMeta) schema.DeleteContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		c, ok := m.(*Client)
		if !ok {
			return diag.FromErr(fmt.Errorf("API client is not configured"))
		}

		if err := c.Sources.Delete(ctx, d.Id()); err != nil {
			return diag.FromErr(err)
		}

		d.SetId("")
		return diag.Diagnostics{}
	}
}

func populateSourceFromState(cm configs.ConfigMeta, source *client.Source, d *schema.ResourceData) {
	source.ID = d.Id()
	source.Type = cm.APIType
	source.Name = d.Get("name").(string)
	source.IsEnabled = d.Get("enabled").(bool)
	source.Config, _ = cm.ParseResourceData(d)
}

func storeSourceToState(cm configs.ConfigMeta, source *client.Source, d *schema.ResourceData) error {
	d.SetId(source.ID)
	if err := d.Set("name", source.Name); err != nil {
		return err
	}
	if err := d.Set("enabled", source.IsEnabled); err != nil {
		return err
	}
	if err := cm.StoreResourceData(source.Config, d); err != nil {
		return err
	}
	return nil
}
