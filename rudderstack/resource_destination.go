package rudderstack

import (
    "context"
    // "strconv"
    "strings"
    // "time"
    "log"
    // "math/big"

    "github.com/hashicorp/terraform-plugin-framework/diag"
    "github.com/hashicorp/terraform-plugin-framework/tfsdk"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/rudderlabs/cp-client-go"
)

type resourceDestinationType struct{}

// Destination Resource schema
func (r resourceDestinationType) GetSchema(context context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
            "allow_same_name": {
                Type:     types.BoolType,
                Optional: true,
            },
            /* Not config. Cause problems when server updates them.
            "created_at": {
                Type:     types.StringType,
                Computed: true,
            },
            "updated_at": {
                Type:     types.StringType,
                Computed: true,
            },*/
            "config": GetConfigJsonObjectAttributeSchema(context),
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

func NewDestination(clientDestination *rudderclient.Destination, allowSameName types.Bool) (Destination) {
    // log.Println("SDK dest config creation started.")
    // if (clientDestination.Config == nil) {
    //     log.Println("Got Client dest config nil.")
    // } else if (len(clientDestination.Config) == 0) {
    //     log.Println("Got Client dest config empty.")
    // }

    newConfig := ConvertApiClientConfigToTerraform(&clientDestination.Config)
    // log.Println("New SDK config gGenerated.", newConfig)

    retval := Destination{
        ID                        : types.String{Value: clientDestination.ID},
        Name                      : types.String{Value: clientDestination.Name},
        Type                      : types.String{Value: clientDestination.Type},
        AllowSameName             : allowSameName,
        /* Not config. Cause problems when server updates them.
        CreatedAt                 : types.String{Value: string(clientDestination.CreatedAt.Format(time.RFC850))},
        UpdatedAt                 : types.String{Value: string(clientDestination.UpdatedAt.Format(time.RFC850))},
        */
        Config                    : newConfig,
    }
    // log.Println("SDK dest config created.", newConfig.ObjectPropertiesMap)
    return retval
}

func (sdkDestination Destination) ToClient() rudderclient.Destination {
    // log.Println("Client dest config creation started.")
    retval := rudderclient.Destination {
        ID                        : sdkDestination.ID.Value,
        Name                      : sdkDestination.Name.Value,
        Type                      : sdkDestination.Type.Value,
        Config                    : sdkDestination.Config.ObjectPropertiesMap.TerraformToApiClient(),
    }
    // log.Println("Client dest config created.")
    return retval
}

// Create a new resource
func (r resourceDestination) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
    // log.Println("SDK dest creation started.")
    if !r.p.configured {
        resp.Diagnostics.AddError(
            "Provider not configured",
            "The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
        )
        return
    }

    // log.Println("SDK dest creation in progress 1.")
    // Retrieve values from plan
    var plan Destination
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Convert terraform object to REST API Client object.
    clientDestination := plan.ToClient()

    existingDestinations, err := r.p.client.FilterDestinations(clientDestination.Type, clientDestination.Name)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error filtering for existing destinations",
            "Could not create destination, unexpected error: "+err.Error(),
        )
        return
    }

    state := Destination{}
    allowSameName := !plan.AllowSameName.Null && plan.AllowSameName.Value

    if len(existingDestinations) > 0 {
        logStr := "Found "
        if len(existingDestinations) > 1 {
            logStr += "more than one pre-existing destinations"
        } else {
            logStr += "a pre-existing destination"
        }
        logStr += " with matching type/name."

        if (allowSameName) {
            logStr = "Warning: " + logStr + " Creating another one."
            resp.Diagnostics.AddWarning(
                "Anomaly creating destination",
                logStr,
            )
        } else {
            logStr = "Error: " + logStr + " Giving up.\n"
            logStr += "Fix by following one of the options below:\n"
            logStr += "1) Invoke ImportState command to import the upstream resource into local terraform state.\n"
            logStr += "2) Resolve name conflict by changing the name of either upstream or downstream resource.\n"
            logStr += "3) Force same name by setting ForceSameName=true in the terraform resource config."
            resp.Diagnostics.AddError(
                "Anomaly creating destination",
                logStr,
            )
        }
        log.Println(logStr,
            "type=", clientDestination.Type,
            "name=", clientDestination.Name,
            "existing=", existingDestinations)
    } else {
        log.Println("Existing destination not found, type=", clientDestination.Type, ", name=", clientDestination.Name)
    }

    if len(existingDestinations) == 0 || allowSameName {
        // Create new destination. 
        createdDestination, err2 := r.p.client.CreateDestination(clientDestination)
        if err2 != nil {
            resp.Diagnostics.AddError(
                "Error creating destination",
                "Could not create destination, unexpected error: "+err2.Error(),
            )
            return
        }

        state = NewDestination(createdDestination, plan.AllowSameName)

        diags = resp.State.Set(ctx, state)
        resp.Diagnostics.Append(diags...)
        if resp.Diagnostics.HasError() {
            return
        }
    }
}

// Read resource information
func (r resourceDestination) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
    // log.Println("SDK dest read started.")
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
    // log.Println("SDK dest read value=", destination.Config)

    state = NewDestination(destination, state.AllowSameName)
    // log.Println("SDK dest value being set=", state.Config)

    // Set state with updated value.
    diags = resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

// Update resource
func (r resourceDestination) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
    // log.Println("SDK dest update started.")
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
    state = NewDestination(destination, plan.AllowSameName)
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
            "Could not read (destinationType, destinationName) for destination import " + req.ID,
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

    state := NewDestination(&destinations[0], types.Bool{Null: true})

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
            "Could not delete destinationID "+destinationID+": "+err.Error(),
        )
        return
    }

    // Remove state.
    resp.State.RemoveResource(ctx)
}
