package rudderstack

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Sources -
type Source struct {
	ID        types.String	        `tfsdk:"id"`
	Name      types.String	        `tfsdk:"name"`
	Type      types.String	        `tfsdk:"type"`
	CreatedAt types.String	        `tfsdk:"created_at"`
	UpdatedAt types.String	        `tfsdk:"updated_at"`

	Config    SourceConfig	        `tfsdk:"config"`
}

// Config for source.
type SourceConfig struct {
	ID        int 			        `tfsdk:"id"`
}

// Destinations -
type Destination struct {
	ID        types.String	        `tfsdk:"id"`
	Name      types.String	        `tfsdk:"name"`
	Type      types.String	        `tfsdk:"type"`
	CreatedAt types.String	        `tfsdk:"created_at"`
	UpdatedAt types.String	        `tfsdk:"updated_at"`

	Config    DestinationConfig	    `tfsdk:"config"`
}

// Config for destination.
type DestinationConfig struct {
	ID        int 			        `tfsdk:"id"`
}

type Connection struct {
	ID             types.String     `tfsdk:"id"`
	SourceID       types.String     `tfsdk:"source_id"`
	DestinationID  types.String     `tfsdk:"destination_id"`
}