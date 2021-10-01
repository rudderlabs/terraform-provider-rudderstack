package rudderstack

import (
    "context"
    // "strconv"
	"strings"
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

func NewDestination(clientDestination *rudderclient.Destination) (Destination) {
    return Destination{
        ID                        : types.String{Value: clientDestination.ID},
        Name                      : types.String{Value: clientDestination.Name},
        Type                      : types.String{Value: clientDestination.Type},
        CreatedAt                 : types.String{Value: string(clientDestination.CreatedAt.Format(time.RFC850))},
        UpdatedAt                 : types.String{Value: string(clientDestination.UpdatedAt.Format(time.RFC850))},
    
        Config                    : DestinationConfig{
            ID        : clientDestination.Config.ID,
        },
    }
}

func (sdkDestination Destination) ToClient() rudderclient.Destination {
    return rudderclient.Destination {
        ID                    : sdkDestination.ID.Value,
        Name                  : sdkDestination.Name.Value,
        Type                  : sdkDestination.Type.Value,
        Config                : rudderclient.DestinationConfig {
            ID        : sdkDestination.Config.ID,
        },
    }
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
    clientDestination := plan.ToClient()

    // Create new destination
    createdDestination, err := r.p.client.CreateDestination(clientDestination)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error creating destination",
            "Could not create destination, unexpected error: "+err.Error(),
        )
        return
    }

    state := NewDestination(createdDestination)

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

    state = NewDestination(destination)

    // Set state with updated value.
    diags = resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

// Update resource
func (r resourceDestination) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
    // Get plan values
    var plan Destination
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
            return
    }

    // Get current state
    var state Destination
    diags = req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
            return
    }

    // Convert terraform object to REST API Client object.
    clientDestination := plan.ToClient()

    // Get destination ID from current state.
    destinationID := state.ID.Value

    // Get current value of destination from API.
    destination, err := r.p.client.UpdateDestination(destinationID, clientDestination)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error updating destination",
            "Could not update destinationID "+destinationID+": "+err.Error(),
        )
        return
    }

    // Set state with updated value.
    state = NewDestination(destination)
    diags = resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

// ImportState resource
func (r resourceDestination) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	var diags diag.Diagnostics

	// Get destination type/name from import request.
    idFields := strings.Fields(req.ID)
	destinationType := ""
	destinationName := ""
	if (len(idFields) == 1) {
		destinationName = idFields[0]
	} else if (len(idFields) == 2) {
		destinationType = idFields[0]
		destinationName = idFields[1]
	} else {
        resp.Diagnostics.AddError(
            "Error reading import request",
            "Could not read (destinationType, destinationName) for connection import " + req.ID,
        )
        return
	}

    // Get current value of destination from API.
    destinations, err := r.p.client.FilterDestinations(destinationType, destinationName)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error filtering destination",
            "Could not filter destinations by import request " + req.ID + ": "+err.Error(),
        )
        return
    }

	if len(destinations) != 1 {
        resp.Diagnostics.AddError(
            "No matching destination found",
            "Number of destinations matching import request ==" + req.ID + " is " + string(len(destinations)) + "!= 1: " + err.Error(),
        )
        return
	}

    state := NewDestination(&destinations[0])

    // Set state with updated value.
    diags = resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
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
