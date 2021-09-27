package rudderclient

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Source Definition -
type SourceDefinition struct {
	ID        types.String 				`json:"id,omitempty"`
	Name      types.String 				`json:"name"`
	Category  types.String 				`json:"category"`
	CreatedAt time.Time    				`json:"createdAt"`
	UpdatedAt time.Time    				`json:"updatedAt"`

	Config SourceDefinitionConfig		`json:"config"`
}

// Destination Definition -
type DestinationDefinition struct {
	ID        types.String 				`json:"id,omitempty"`
	Name      types.String 				`json:"name"`
	Category  types.String 				`json:"category"`
	CreatedAt time.Time    				`json:"createdAt"`
	UpdatedAt time.Time    				`json:"updatedAt"`

	Config DestinationDefinitionConfig	`json:"config"`
}

// Sources -
type Source struct {
	ID        string 					`json:"id,omitempty"`
	Name      string 					`json:"name"`
	Type      string	 				`json:"type,omitempty"`
	CreatedAt time.Time    				`json:"createdAt"`
	UpdatedAt time.Time    				`json:"updatedAt"`

	Config    SourceConfig 				`json:"config"`
}

// Destinations -
type Destination struct {
	ID        types.String 				`json:"id,omitempty" tfsdk:"id"`
	Name      types.String 				`json:"name" tfsdk:"name"`
	Type      types.String 				`json:"type" tfsdk:"type"`
	CreatedAt time.Time    				`json:"createdAt" tfsdk:"created_at"`
	UpdatedAt time.Time    				`json:"updatedAt" tfsdk:"updated_at"`

	Config    DestinationConfig 		`json:"config" tfsdk:"config"`
}

type SourceConfig struct {
	ID        int		 				`json:"id,omitempty"`
}

type DestinationConfig struct {
}

type SourceDefinitionConfig struct {
}

type DestinationDefinitionConfig struct {
}
