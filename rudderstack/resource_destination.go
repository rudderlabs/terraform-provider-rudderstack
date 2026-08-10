package rudderstack

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/rudderlabs/rudder-iac/api/client"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/integrations/destinations"
)

func resourceDestination(cm configs.ConfigMeta) *schema.Resource {
	return &schema.Resource{
		Schema:        resourceDestinationSchema(cm),
		CreateContext: resourceDestinationCreate(cm),
		ReadContext:   resourceDestinationRead(cm),
		UpdateContext: resourceDestinationUpdate(cm),
		DeleteContext: resourceDestinationDelete(cm),
		CustomizeDiff: resourceDestinationCustomizeDiff(cm),
		Importer: &schema.ResourceImporter{
			StateContext: resourceDestinationImportState(cm),
		},
	}
}

func resourceDestinationCustomizeDiff(cm configs.ConfigMeta) schema.CustomizeDiffFunc {
	return func(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
		if cm.SkipConfig {
			return nil
		}

		consentManagement := d.Get("config.0.consent_management.0")
		if consentManagement == nil {
			return nil
		}

		return destinations.ValidateConsentManagementUniqueProviders(consentManagement)
	}
}

func resourceDestinationSchema(cm configs.ConfigMeta) map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Human readable name of the destination. The value has to be unique across all the destinations.",
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
		"transformation_id": {
			Type:     schema.TypeString,
			Optional: true,
			Description: "ID of a published transformation to attach to this destination. " +
				"The transformation must be published before it can be attached. A destination " +
				"can have at most one transformation. Omit or clear the value to detach.",
		},
	}

	if !cm.SkipConfig {
		s["config"] = &schema.Schema{
			Type:     schema.TypeList,
			Optional: cm.SkipConfig,
			Required: !cm.SkipConfig,
			Description: "Destination specific configuration. Check the nested block documenation " +
				"for more information.",
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

		// The transformation link is a separate endpoint, not part of the destination
		// create body — attach it after the destination exists. Mirrors rudder-iac's
		// destination handler.
		if newID := d.Get("transformation_id").(string); newID != "" {
			if err := reconcileTransformationLink(ctx, c, destination.ID, "", newID); err != nil {
				return diag.FromErr(err)
			}
		}

		return resourceDestinationRead(cm)(ctx, d, m)
	}
}

// reconcileTransformationLink brings the destination's attached transformation from
// oldID to newID. Empty string means "no transformation". Mirrors rudder-iac's
// syncTransformationLink. ponytail: one small helper beside the resource, no shared pkg.
func reconcileTransformationLink(ctx context.Context, c *Client, destinationID, oldID, newID string) error {
	if oldID == newID {
		return nil
	}
	if newID == "" {
		if err := c.Destinations.DisconnectTransformation(ctx, destinationID); err != nil {
			return fmt.Errorf("disconnecting transformation from destination %s: %w", destinationID, err)
		}
		return nil
	}
	if _, err := c.Destinations.ConnectTransformation(ctx, destinationID, newID); err != nil {
		return fmt.Errorf("connecting transformation %s to destination %s: %w", newID, destinationID, err)
	}
	return nil
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
			var apiErr *client.APIError
			if errors.As(err, &apiErr) && apiErr.HTTPStatusCode == 404 {
				d.SetId("")
				return diag.Diagnostics{}
			}
			return diag.FromErr(err)
		}

		if err := storeDestinationToState(cm, destination, d); err != nil {
			return diag.FromErr(err)
		}

		// GET /destinations/{id} does not embed the transformation link, so read it
		// separately. A not-found means nothing is attached (empty), and also lets
		// drift (out-of-band attach/detach) surface on the next plan.
		transformation, err := c.Destinations.GetTransformation(ctx, id)
		switch {
		case err == nil:
			if err := d.Set("transformation_id", transformation.TransformationID); err != nil {
				return diag.FromErr(err)
			}
		case errors.Is(err, client.ErrResourceNotFound):
			if err := d.Set("transformation_id", ""); err != nil {
				return diag.FromErr(err)
			}
		default:
			return diag.FromErr(fmt.Errorf("reading transformation link for destination %s: %w", id, err))
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
			return diag.FromErr(fmt.Errorf("could not update destination: %w", err))
		}

		d.SetId(destination.ID)

		if d.HasChange("transformation_id") {
			oldID, newID := d.GetChange("transformation_id")
			if err := reconcileTransformationLink(ctx, c, destination.ID, oldID.(string), newID.(string)); err != nil {
				return diag.FromErr(err)
			}
		}

		return resourceDestinationRead(cm)(ctx, d, m)
	}
}

func resourceDestinationDelete(cm configs.ConfigMeta) schema.DeleteContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		c, ok := m.(*Client)
		if !ok {
			return diag.FromErr(fmt.Errorf("API client is not configured"))
		}

		// Detach any linked transformation first, mirroring rudder-iac's destination
		// handler. Not-found is fine — the link may already be gone or cascade-deleted.
		if d.Get("transformation_id").(string) != "" {
			if err := c.Destinations.DisconnectTransformation(ctx, d.Id()); err != nil {
				var apiErr *client.APIError
				if !errors.Is(err, client.ErrResourceNotFound) && !(errors.As(err, &apiErr) && apiErr.HTTPStatusCode == 404) {
					return diag.FromErr(fmt.Errorf("disconnecting transformation from destination %s: %w", d.Id(), err))
				}
			}
		}

		if err := c.Destinations.Delete(ctx, d.Id()); err != nil {
			var apiErr *client.APIError
			if errors.As(err, &apiErr) && apiErr.HTTPStatusCode == 404 {
				d.SetId("")
				return diag.Diagnostics{}
			}
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
	// rudder-iac >= v0.19.0 widened Destination.Version to int64.
	destination.Version = int64(cm.Version)
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

	// SkipConfig destinations have no "config" attribute in their schema, so
	// there is nothing to write back (and d.Set("config", ...) would fail).
	if cm.SkipConfig {
		return nil
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
