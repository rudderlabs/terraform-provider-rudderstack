package rudderstack

import (
    "context"
    "strings"
    "strconv"
    // "time"
    "log"
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
            "enabled": {
                Type:     types.BoolType,
                Computed: true,
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
        IsEnabled           : types.Bool{Value: clientConnection.IsEnabled},
    }
}

func (sdkConnection Connection) ToClient() rudderclient.Connection {
    return rudderclient.Connection {
        ID                  : sdkConnection.ID.Value,
        SourceID            : sdkConnection.SourceID.Value,
        DestinationID       : sdkConnection.DestinationID.Value,
        IsEnabled           : sdkConnection.IsEnabled.Value,
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

    existingConnections, err := r.p.client.FilterConnections(clientConnection.SourceID, clientConnection.DestinationID)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error filtering for existing connections",
            "Could not create connection, unexpected error: "+err.Error(),
        )
        return
    }

    state := Connection{}
    if len(existingConnections) == 0 {
        log.Println("Existing connection not found, src=", clientConnection.SourceID, ", dst=", clientConnection.DestinationID)
        // Create new connection. 
        createdConnection, err2 := r.p.client.CreateConnection(clientConnection)
        if err2 != nil {
            resp.Diagnostics.AddError(
                "Error creating connection",
                "Could not create connection, unexpected error: "+err2.Error(),
            )
            return
        }

        state = NewConnection(createdConnection)
    } else if len(existingConnections) == 1 {
        log.Println("ReUsing existing connection")
        // Use existing connection.
        state = NewConnection(&existingConnections[0])
    } else {
        resp.Diagnostics.AddError(
            "Error creating connection",
            "Filtered connection for source/destination but got too many connections: ",
        )
        return
    }

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

    if connection.IsDeleted {
        resp.Diagnostics.AddError(
            "Target connection deleted",
            "Target connection has been marked as deleted. Could not read connection (connectionId) " + connectionId + ": " + err.Error(),
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

    helpStr := "Either specify {connectionGuid} or {sourceId},{destinationId} for connection import."
    connection, getErr := r.p.client.GetConnection(req.ID)

    if (getErr != nil) {
        // Get src/dst ID from import request.
        idFields := strings.SplitN(req.ID, ",", 2)
        if (len(idFields) != 2) {
            resp.Diagnostics.AddError(
                "Error parsing connection import request",
                "Could not parse " + req.ID + ". " + helpStr,
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
                "Could not filter conections by import request " + req.ID + ". Error:"+err.Error() + ". " + helpStr,
            )
            return
        }

        if len(connections) != 1 {
            resp.Diagnostics.AddError(
                "No matching connection found",
                "Number of connections matching import request '" + req.ID + "' is " + strconv.Itoa(len(connections)) + " != 1. " + helpStr,
            )
            return
        }

        connection = &connections[0]
    }

    state := NewConnection(connection)

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
