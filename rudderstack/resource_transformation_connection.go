package rudderstack

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/rudderlabs/rudder-iac/api/client"
	"github.com/rudderlabs/rudder-iac/api/client/transformations"
)

func resourceTransformationConnection() *schema.Resource {
	return &schema.Resource{
		Description: "Connects an existing (published) transformation to an existing destination. " +
			"A destination can have only one connected transformation at a time; connecting a " +
			"transformation to a destination replaces any transformation already connected to it. " +
			"The transformation must be published before it can be connected.",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "The composite ID of the connection, in the form " +
					"`<transformation_id>:<destination_id>`.",
			},
			"transformation_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the transformation to connect. Must be published.",
			},
			"destination_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the destination to connect the transformation to.",
			},
		},
		CreateContext: resourceTransformationConnectionCreate,
		ReadContext:   resourceTransformationConnectionRead,
		DeleteContext: resourceTransformationConnectionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceTransformationConnectionImportState,
		},
	}
}

// transformationConnectionID builds the synthetic resource ID for a
// transformation-to-destination connection.
func transformationConnectionID(transformationID, destinationID string) string {
	return fmt.Sprintf("%s:%s", transformationID, destinationID)
}

// parseTransformationConnectionID splits a synthetic resource ID back into its
// transformation and destination IDs.
func parseTransformationConnectionID(id string) (transformationID, destinationID string, err error) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf(
			"invalid transformation connection ID %q, expected format <transformation_id>:<destination_id>",
			id,
		)
	}
	return parts[0], parts[1], nil
}

func resourceTransformationConnectionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(*Client)
	if !ok {
		return diag.FromErr(fmt.Errorf("API client is not configured"))
	}

	transformationID := d.Get("transformation_id").(string)
	destinationID := d.Get("destination_id").(string)

	_, err := c.TransformationConnections.ConnectToDestination(ctx, transformationID, &transformations.ConnectToDestinationRequest{
		DestinationID: destinationID,
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf(
			"could not connect transformation %q to destination %q: %w",
			transformationID, destinationID, err,
		))
	}

	d.SetId(transformationConnectionID(transformationID, destinationID))

	return resourceTransformationConnectionRead(ctx, d, m)
}

func resourceTransformationConnectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(*Client)
	if !ok {
		return diag.FromErr(fmt.Errorf("API client is not configured"))
	}

	transformationID, destinationID, err := parseTransformationConnectionID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	transformation, err := c.TransformationConnections.GetTransformation(ctx, transformationID)
	if err != nil {
		var apiErr *client.APIError
		if errors.As(err, &apiErr) && apiErr.HTTPStatusCode == 404 {
			// The transformation no longer exists, so the connection is gone.
			d.SetId("")
			return diag.Diagnostics{}
		}
		return diag.FromErr(err)
	}

	// The connection only exists if the destination is still in the
	// transformation's list of connected destinations.
	if !transformationConnectedToDestination(transformation, destinationID) {
		d.SetId("")
		return diag.Diagnostics{}
	}

	if err := d.Set("transformation_id", transformationID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("destination_id", destinationID); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func resourceTransformationConnectionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(*Client)
	if !ok {
		return diag.FromErr(fmt.Errorf("API client is not configured"))
	}

	transformationID, destinationID, err := parseTransformationConnectionID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = c.TransformationConnections.DisconnectFromDestination(ctx, transformationID, &transformations.ConnectToDestinationRequest{
		DestinationID: destinationID,
	})
	if err != nil {
		var apiErr *client.APIError
		if errors.As(err, &apiErr) && apiErr.HTTPStatusCode == 404 {
			// Already gone.
			d.SetId("")
			return diag.Diagnostics{}
		}
		return diag.FromErr(fmt.Errorf(
			"could not disconnect transformation %q from destination %q: %w",
			transformationID, destinationID, err,
		))
	}

	d.SetId("")
	return diag.Diagnostics{}
}

func resourceTransformationConnectionImportState(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	transformationID, destinationID, err := parseTransformationConnectionID(d.Id())
	if err != nil {
		return nil, err
	}

	if err := d.Set("transformation_id", transformationID); err != nil {
		return nil, err
	}
	if err := d.Set("destination_id", destinationID); err != nil {
		return nil, err
	}

	diagnostics := resourceTransformationConnectionRead(ctx, d, m)
	if diagnostics.HasError() {
		for _, diagnostic := range diagnostics {
			if diagnostic.Severity == diag.Error {
				return nil, fmt.Errorf("could not import transformation connection: %s", diagnostic.Summary)
			}
		}
	}

	if d.Id() == "" {
		return nil, fmt.Errorf(
			"transformation %q is not connected to destination %q",
			transformationID, destinationID,
		)
	}

	return []*schema.ResourceData{d}, nil
}

func transformationConnectedToDestination(transformation *transformations.Transformation, destinationID string) bool {
	if transformation == nil {
		return false
	}
	for _, dest := range transformation.Destinations {
		if dest.ID == destinationID {
			return true
		}
	}
	return false
}
