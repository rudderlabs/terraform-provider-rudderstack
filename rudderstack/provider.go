package rudderstack

import (
	"context"
	"os"
	// "log"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/rudderlabs/cp-client-go"
)

var stderr = os.Stderr

func New() tfsdk.Provider {
	return &provider{}
}

type provider struct {
	configured bool
	client     *rudderclient.Client
}

// GetSchema
func (p *provider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"host": {
				Type:     types.StringType,
				Optional: true,
				Computed: true,
			},
			"token": {
				Type:      types.StringType,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
			"schema_host": {
				Type:     types.StringType,
				Optional: true,
				Computed: true,
			},
			"schema_token": {
				Type:      types.StringType,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
		},
	}, nil
}

// Provider schema struct
type providerData struct {
	WorkspaceHost  types.String `tfsdk:"host"`
	WorkspaceToken types.String `tfsdk:"token"`
	SchemaHost    types.String `tfsdk:"schema_host"`
	SchemaToken   types.String `tfsdk:"schema_token"`
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	// Retrieve provider data from configuration
	var config providerData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// User must specify a workspace host
	var workspaceHost string
	if config.WorkspaceHost.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddError(
			"Unable to create client",
			"Cannot use unknown value as workspace host",
		)
		return
	}

	if config.WorkspaceHost.Null {
		workspaceHost = os.Getenv("RUDDERSTACK_HOST")
	} else {
		workspaceHost = config.WorkspaceHost.Value
	}

	if workspaceHost == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find workspace host",
			"Workspace Host cannot be an empty string",
		)
		return
	}

	// User must provide a workspace token to the provider
	var workspaceToken string
	if config.WorkspaceToken.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddError(
			"Unable to create client",
			"Cannot use unknown value as workspace token",
		)
		return
	}

	if config.WorkspaceToken.Null {
		workspaceToken = os.Getenv("RUDDERSTACK_TOKEN")
	} else {
		workspaceToken = config.WorkspaceToken.Value
	}

	if workspaceToken == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find workspace token",
			"Workspace token cannot be an empty string",
		)
		return
	}

	// User must specify a schema host
	var schemaHost string
	if config.SchemaHost.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddError(
			"Unable to create client",
			"Cannot use unknown value as schema host",
		)
		return
	}

	if config.SchemaHost.Null {
		schemaHost = os.Getenv("RUDDERSTACK_SCHEMA_HOST")
	} else {
		schemaHost = config.SchemaHost.Value
	}

	if schemaHost == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find schema host",
			"Schema Host cannot be an empty string",
		)
		return
	}

	// User must provide a schema token to the provider
	var schemaToken string
	if config.SchemaToken.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddError(
			"Unable to create client",
			"Cannot use unknown value as schema token",
		)
		return
	}

	if config.SchemaToken.Null {
		schemaToken = os.Getenv("RUDDERSTACK_SCHEMA_TOKEN")
	} else {
		schemaToken = config.SchemaToken.Value
	}

	if schemaToken == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find schema token",
			"Schema token cannot be an empty string",
		)
		return
	}

	// Create a new HashiCups client and set it to the provider client
	c, err := rudderclient.NewClient(&workspaceHost, &workspaceToken, &schemaHost, &schemaToken)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create client",
			"Unable to create rudderstack client:\n\n"+err.Error(),
		)
		return
	}

	p.client = c
	p.configured = true
}

// GetResources - Defines provider resources
func (p *provider) GetResources(_ context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		"rudderstack_source": resourceSourceType{},
		"rudderstack_destination": resourceDestinationType{},
		"rudderstack_connection": resourceConnectionType{},
	}, nil
}

// GetDataSources - Defines provider data sources
func (p *provider) GetDataSources(_ context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{}, nil
}
