package rudderstack

import (
	"context"
	// "strconv"
	"time"
	// "log"
	// "math/big"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/rudderlabs/cp-client-go"
)

type resourceDestinationType struct{}

// Destination Resource schema
func (r resourceDestinationType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"name": {
				Type:     types.StringType,
				Required: true,
			},
			"type": {
				Type:     types.StringType,
				Required: true,
			},
			"created_at": {
				Type:     types.StringType,
				Computed: true,
			},
			"updated_at": {
				Type:     types.StringType,
				Computed: true,
			},
			"config": {
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"id": {
						Type:     types.NumberType,
						Computed: true,
						Optional: true,
					},
				}),
				Optional: true,
			},
		},
	}, nil
}

// New resource instance
func (r resourceDestinationType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceDestination{
		p:   *(p.(*provider)),
	}, nil
}

type resourceDestination struct {
	p provider
}

// Create a new resource
func (r resourceDestination) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	// Retrieve values from plan
	var plan Destination
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert terraform object to REST API Client object.
	clientDestination := rudderclient.Destination {
		Name      : plan.Name.Value,
		Type      : plan.Type.Value,
		Config    : rudderclient.DestinationConfig {
		},
	}

	// Create new destination
	createdDestination, err := r.p.client.CreateDestination(clientDestination)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating destination",
			"Could not create destination, unexpected error: "+err.Error(),
		)
		return
	}

	state := Destination{
		ID        : types.String{Value: createdDestination.ID},
		Name      : types.String{Value: createdDestination.Name},
		Type      : types.String{Value: createdDestination.Type},
		CreatedAt : types.String{Value: string(createdDestination.CreatedAt.Format(time.RFC850))},
		UpdatedAt : types.String{Value: string(createdDestination.UpdatedAt.Format(time.RFC850))},
	
		Config    : DestinationConfig{
			ID        : createdDestination.Config.ID,
		},
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r resourceDestination) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	// Get current state
	var state Destination
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get destination ID from current state.
	destinationID := state.ID.Value

	// Get current value of destination from API.
	destination, err := r.p.client.GetDestination(destinationID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading destination",
			"Could not read destinationID "+destinationID+": "+err.Error(),
		)
		return
	}

	state = Destination{
		ID        : types.String{Value: destination.ID},
		Name      : types.String{Value: destination.Name},
		Type      : types.String{Value: destination.Type},
		CreatedAt : types.String{Value: string(destination.CreatedAt.Format(time.RFC850))},
		UpdatedAt : types.String{Value: string(destination.UpdatedAt.Format(time.RFC850))},
	
		Config    : DestinationConfig{},
	}

	// Set state with updated value.
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update resource
func (r resourceDestination) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
}

// ImportState resource
func (r resourceDestination) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStateNotImplemented(ctx, "", resp)
}

// Delete resource
func (r resourceDestination) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	// Get current state
	var state Destination
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get destination ID from current state.
	destinationID := state.ID.Value

	// Delete destination via API.
	err := r.p.client.DeleteDestination(destinationID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting destination",
			"Could not read destinationID "+destinationID+": "+err.Error(),
		)
		return
	}

	// Set state.
	diags = resp.State.Set(ctx, nil)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
