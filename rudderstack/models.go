package rudderstack

import (
    "github.com/hashicorp/terraform-plugin-framework/types"
)

// Sources -
type Source struct {
    ID                     types.String                    `tfsdk:"id"`
    Name                   types.String                    `tfsdk:"name"`
    Type                   types.String                    `tfsdk:"type"`
    CreatedAt              types.String                    `tfsdk:"created_at"`
    UpdatedAt              types.String                    `tfsdk:"updated_at"`

    Config                 *EncapsulatedConfigObject       `tfsdk:"config"`
}

// Destinations -
type Destination struct {
    ID                     types.String                    `tfsdk:"id"`
    Name                   types.String                    `tfsdk:"name"`
    Type                   types.String                    `tfsdk:"type"`
    CreatedAt              types.String                    `tfsdk:"created_at"`
    UpdatedAt              types.String                    `tfsdk:"updated_at"`

    Config                 *EncapsulatedConfigObject       `tfsdk:"config"`
}

type Connection struct {
    ID                     types.String                    `tfsdk:"id"`
    SourceID               types.String                    `tfsdk:"source_id"`
    DestinationID          types.String                    `tfsdk:"destination_id"`
}

type EncapsulatedConfigObject struct {
    ObjectPropertiesMap    ObjectPropertiesMap             `tfsdk:"object"`
}

type ObjectPropertiesMap   map[string]SingleObjectProperty

type SingleObjectProperty struct {
    StrValue               types.String                    `tfsdk:"str"`
    NumValue               types.Number                    `tfsdk:"num"`
    BoolValue              types.Bool                      `tfsdk:"bool"`
    ObjectValue            *ObjectPropertiesMap            `tfsdk:"object"`
    ObjectsListValue       *[]EncapsulatedConfigObject     `tfsdk:"objects_list"`
}
