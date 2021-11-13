package rudderstack

import (
    "context"
    "strings"
    // "strconv"
    // "time"
    // "log"
    //"math/big"

    "github.com/hashicorp/terraform-plugin-framework/diag"
    "github.com/hashicorp/terraform-plugin-framework/tfsdk"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/rudderlabs/cp-client-go"
)

type resourceConnectionType struct{}

// Source Resource schema
func (r resourceConnectionType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
    return tfsdk.Schema{
        Attributes: map[string]tfsdk.Attribute{
            "id": {
                Type:     types.StringType,
                Computed: true,
            },
            "source_id": {
                Type:     types.StringType,
                Required: true,
            },
            "destination_id": {
                Type:     types.StringType,
                Required: true,
            },
        },
    }, nil
}

// New resource instance
func (r resourceConnectionType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
    return resourceConnection{
        p:   *(p.(*provider)),
    }, nil
}

type resourceConnection struct {
    p provider
}

func NewConnection(clientConnection *rudderclient.Connection) (Connection) {
    return Connection{
        ID                  : types.String{Value: clientConnection.ID},
        SourceID            : types.String{Value: clientConnection.SourceID},
        DestinationID       : types.String{Value: clientConnection.DestinationID},
    }
}

func (sdkConnection Connection) ToClient() rudderclient.Connection {
    return rudderclient.Connection {
        ID                  : sdkConnection.ID.Value,
        SourceID              : sdkConnection.SourceID.Value,
        DestinationID       : sdkConnection.DestinationID.Value,
    }
}

// Create a new resource
func (r resourceConnection) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
    if !r.p.configured {
        resp.Diagnostics.AddError(
            "Provider not configured",
            "The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
        )
        return
    }

    // Retrieve values from plan
    var plan Connection
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Convert terraform object to REST API Client object.
    clientConnection := plan.ToClient()

    // Create new source
    createdConnection, err := r.p.client.CreateConnection(clientConnection)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error creating connection",
            "Could not create connection, unexpected error: "+err.Error(),
        )
        return
    }

    state := NewConnection(createdConnection)
    
    diags = resp.State.Set(ctx, state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

// Read resource information
func (r resourceConnection) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
    // Get current state
    var state Connection
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Get src/dst ID from current state.
    connectionId := state.ID.Value

    // Get current value of source from API.
    connection, err := r.p.client.GetConnection(connectionId)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error reading connection",
            "Could not read connection (connectionId) " + connectionId + ": " + err.Error(),
        )
        return
    }

    state = NewConnection(connection)

    // Set state with updated value.
    diags = resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

// Update resource
func (r resourceConnection) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
    // Get plan values
    var plan Connection
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
            return
    }

    // Get current state
    var state Connection
    diags = req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
            return
    }

    // Convert terraform object to REST API Client object.
    clientConnection := plan.ToClient()

    // Get connection ID from current state.
    connectionID := state.ID.Value

    // Get current value of connection from API.
    connection, err := r.p.client.UpdateConnection(connectionID, clientConnection)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error updating connection",
            "Could not update connectionID "+connectionID+": "+err.Error(),
        )
        return
    }

    // Set state with updated value.
    state = NewConnection(connection)
    diags = resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

// ImportState resource
func (r resourceConnection) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
    var diags diag.Diagnostics

    // Get src/dst ID from import request.
    idFields := strings.Fields(req.ID)
    if (len(idFields) != 2) {
        resp.Diagnostics.AddError(
            "Error reading import request",
            "Could not read (sourceId, destinationId) for connection import " + req.ID,
        )
        return
    }
    sourceId := idFields[0]
    destinationId := idFields[1]

    // Get current connection value from API.
    connections, err := r.p.client.FilterConnections(sourceId, destinationId)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error filtering connection",
            "Could not read connection (sourceId, destinationId) = (" + sourceId + ","+ destinationId + ") : " + err.Error(),
        )
        return
    }

    if len(connections) != 1 {
        resp.Diagnostics.AddError(
            "No matching connection found",
            "Could not find any connection matching (src,dst) = (" + sourceId + ","+ destinationId + "): " + err.Error(),
        )
        return
    }

    state := NewConnection(&connections[0])

    // Set state with updated value.
    diags = resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

// Delete resource
func (r resourceConnection) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
    // Get current state
    var state Connection
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Get ID from current state.
    connectionId := state.ID.Value

    // Delete source via API.
    err := r.p.client.DeleteConnection(connectionId)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error deleting connection",
            "Could not read connectionId "+connectionId+": "+err.Error(),
        )
        return
    }

    // Set state.
    resp.State.RemoveResource(ctx)
}
