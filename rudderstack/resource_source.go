package rudderstack

import (
	"context"
	// "strconv"
	"time"
	// "log"
	//"math/big"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/rudderlabs/cp-client-go"
)

type resourceSourceType struct{}

// Source Resource schema
func (r resourceSourceType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
func (r resourceSourceType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceSource{
		p:   *(p.(*provider)),
	}, nil
}

type resourceSource struct {
	p provider
}

func (clientSource rudderclient.Source) ToSdk() Source {
	return Source{
		ID        	        : types.String{Value: clientSource.ID},
		Name      			: types.String{Value: clientSource.Name},
		Type      			: types.String{Value: clientSource.Type},
		CreatedAt 			: types.String{Value: string(clientSource.CreatedAt.Format(time.RFC850))},
		UpdatedAt 			: types.String{Value: string(clientSource.UpdatedAt.Format(time.RFC850))},
	
		Config    			: SourceConfig{
			ID        : clientSource.Config.ID,
		},
	}
}

func (sdkSource Source) ToClient() rudderclient.Source {
	return rudderclient.Source {
		ID      		    : sdkSource.ID.Value,
		Name      			: sdkSource.Name.Value,
		Type      			: sdkSource.Type.Value,
		Config    			: rudderclient.SourceConfig {
			ID        : sdkSource.Config.ID,
		},
	}
}

// Create a new resource
func (r resourceSource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	// Retrieve values from plan
	var plan Source
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert terraform object to REST API Client object.
	clientSource := plan.ToClient()

	// Create new source
	createdSource, err := r.p.client.CreateSource(clientSource)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating source",
			"Could not create source, unexpected error: "+err.Error(),
		)
		return
	}

	state := createdSource.ToSdk()

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r resourceSource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	// Get current state
	var state Source
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get source ID from current state.
	sourceID := state.ID.Value

	// Get current value of source from API.
	source, err := r.p.client.GetSource(sourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading source",
			"Could not read sourceID "+sourceID+": "+err.Error(),
		)
		return
	}

	state = source.ToSdk()

	// Set state with updated value.
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update resource
func (r resourceSource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
}

// ImportState resource
func (r resourceSource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStateNotImplemented(ctx, "", resp)
}

// Delete resource
func (r resourceSource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	// Get current state
	var state Source
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get source ID from current state.
	sourceID := state.ID.Value

	// Delete source via API.
	err := r.p.client.DeleteSource(sourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting source",
			"Could not read sourceID "+sourceID+": "+err.Error(),
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
