package rudderstack

import (
    "context"
    "strconv"
    "strings"
    //"time"
    "log"
    //"encoding/json"
    //"math/big"

    "github.com/hashicorp/terraform-plugin-framework/diag"
    "github.com/hashicorp/terraform-plugin-framework/tfsdk"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/rudderlabs/cp-client-go"
)

type resourceSourceType struct{}

// Source Resource schema
func (r resourceSourceType) GetSchema(context context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
            },
            */
            "config": GetConfigJsonObjectAttributeSchema(context),
            "enabled": {
                Type:     types.BoolType,
                Computed: true,
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

func NewSource(clientSource *rudderclient.Source, allowSameName types.Bool) (Source) {
    newConfig := ConvertApiClientConfigToTerraform(&clientSource.Config)

    return Source{
        ID                        : types.String{Value: clientSource.ID},
        Name                      : types.String{Value: clientSource.Name},
        Type                      : types.String{Value: clientSource.Type},
        AllowSameName             : allowSameName,
        /* Not config. Cause problems when server updates them.
        CreatedAt                 : types.String{Value: string(clientSource.CreatedAt.Format(time.RFC850))},
        UpdatedAt                 : types.String{Value: string(clientSource.UpdatedAt.Format(time.RFC850))},
        */
        Config                    : newConfig,
        IsEnabled                 : types.Bool{Value: clientSource.IsEnabled},
    }
}

func (sdkSource Source) TerraformToApiClient() rudderclient.Source {
    return rudderclient.Source {
        ID                        : sdkSource.ID.Value,
        Name                      : sdkSource.Name.Value,
        Type                      : sdkSource.Type.Value,
        Config                    : sdkSource.Config.ObjectPropertiesMap.TerraformToApiClient(),
        IsEnabled                 : sdkSource.IsEnabled.Value,
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
    clientSource := plan.TerraformToApiClient()

    existingSources, err := r.p.client.FilterSources(clientSource.Type, clientSource.Name)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error filtering for existing sources",
            "Could not create source, unexpected error: "+err.Error(),
        )
        return
    }

    state := Source{}
    allowSameName := !plan.AllowSameName.Null && plan.AllowSameName.Value

    if len(existingSources) > 0 {
        logStr := "Found "
        if len(existingSources) > 1 {
            logStr += "more than one pre-existing sources"
        } else {
            logStr += "a pre-existing source"
        }
        logStr += " with matching type/name."

        if (allowSameName) {
            logStr = "Warning: " + logStr + " Creating another one."
            resp.Diagnostics.AddWarning(
                "Anomaly creating source",
                logStr,
            )
        } else {
            logStr = "Error: " + logStr + " Giving up.\n"
            logStr += "Fix by following one of the options below:\n"
            logStr += "1) Invoke ImportState command to import the upstream resource into local terraform state.\n"
            logStr += "2) Resolve name conflict by changing the name of either upstream or downstream resource.\n"
            logStr += "3) Force same name by setting allow_same_name=true in the terraform resource config."
            resp.Diagnostics.AddError(
                "Anomaly creating source",
                logStr,
            )
        }
        log.Println(logStr,
            "type=", clientSource.Type,
            "name=", clientSource.Name,
            "existing=", existingSources)
    } else {
        log.Println("No existing source with same type/name found, type=", clientSource.Type, ", name=", clientSource.Name)
    }

    if len(existingSources) == 0 || allowSameName {
        // Create new source. 
        createdSource, err2 := r.p.client.CreateSource(clientSource)
        if err2 != nil {
            resp.Diagnostics.AddError(
                "Error creating source",
                "Could not create source, unexpected error: "+err2.Error(),
            )
            return
        }

        state = NewSource(createdSource, plan.AllowSameName)
        state.AllowSameName = plan.AllowSameName

        diags = resp.State.Set(ctx, state)
        resp.Diagnostics.Append(diags...)
        if resp.Diagnostics.HasError() {
            return
        }
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

    if source.IsDeleted {
        resp.Diagnostics.AddError(
            "Target source deleted",
            "Target source has been marked as deleted. Could not read sourceID "+sourceID+": "+err.Error(),
        )
        return
    }

    state = NewSource(source, state.AllowSameName)

    // Set state with updated value.
    diags = resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

// Update resource
func (r resourceSource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
    // Get plan values
    var plan Source
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Get current state
    var state Source
    diags = req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Convert terraform object to REST API Client object.
    clientSource := plan.TerraformToApiClient()

    // Get source ID from current state.
    sourceID := state.ID.Value

    // Update source with current value of source.
    source, err := r.p.client.UpdateSource(sourceID, clientSource)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error updating source",
            "Could not update sourceID "+sourceID+": "+err.Error(),
        )
        return
    }

    // Set state with updated value.
    state = NewSource(source, plan.AllowSameName)
    diags = resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

// ImportState resource
func (r resourceSource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
    var diags diag.Diagnostics

    helpStr := "Either specify source ID or '{sourceType}/{sourceName}' for source import."
    // Get source type/name from import request..
    source, getErr := r.p.client.GetSource(req.ID)
    if getErr != nil {
        sourceType := ""
        sourceName := ""

        idFields := strings.SplitN(req.ID, "/", 2)
        if (len(idFields) == 1) {
            sourceName = idFields[0]
        } else if (len(idFields) == 2) {
            sourceType = idFields[0]
            sourceName = idFields[1]
        } else {
            resp.Diagnostics.AddError(
                "Error parsing source import request",
		"Error parsing '" + req.ID + "'. " + helpStr,
            )
            return
        }

        // Get current value of source from API.
        sources, err := r.p.client.FilterSources(sourceType, sourceName)
        if err != nil {
            resp.Diagnostics.AddError(
                "Error filtering source",
		"Could not filter sources by import request " + req.ID + ". Error:"+err.Error() + ". " + helpStr,
            )
            return
        }

        if len(sources) != 1 {
            resp.Diagnostics.AddError(
                "No matching source found",
                "Number of sources matching import request '" + req.ID + "' is " + strconv.Itoa(len(sources)) + " != 1. " + helpStr,
            )
            return
        }

	source = &sources[0]
    }

    state := NewSource(source, types.Bool{Null: true})

    // Set state with updated value.
    diags = resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
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

    // Remove state.
    resp.State.RemoveResource(ctx)
}
