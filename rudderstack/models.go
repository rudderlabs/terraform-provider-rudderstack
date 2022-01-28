package rudderstack

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Sources -
type Source struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Type          types.String `tfsdk:"type"`
	AllowSameName types.Bool   `tfsdk:"allow_same_name"`
	/* Time stamps are not config. Server side updates cause problems.
	   CreatedAt              types.String                    `tfsdk:"created_at"`
	   UpdatedAt              types.String                    `tfsdk:"updated_at"`
	*/

	Config    *EncapsulatedConfigObject `tfsdk:"config"`
	IsEnabled types.Bool                `tfsdk:"enabled"`
}

// Destinations -
type Destination struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Type          types.String `tfsdk:"type"`
	AllowSameName types.Bool   `tfsdk:"allow_same_name"`

	/* Time stamps are not config. Server side updates cause problems.
	   CreatedAt              types.String                    `tfsdk:"created_at"`
	   UpdatedAt              types.String                    `tfsdk:"updated_at"`
	*/

	Config    *EncapsulatedConfigObject `tfsdk:"config"`
	IsEnabled types.Bool                `tfsdk:"enabled"`
}

type Connection struct {
	ID            types.String `tfsdk:"id"`
	SourceID      types.String `tfsdk:"source_id"`
	DestinationID types.String `tfsdk:"destination_id"`
	IsEnabled     types.Bool   `tfsdk:"enabled"`
}

type EncapsulatedConfigObject struct {
	JsonObjectMap JsonObjectMap `tfsdk:"object"`
}

type JsonObjectMap map[string]JsonElement

type JsonElement struct {
	StrValue  types.String          `tfsdk:"str"`
	NumValue  types.Number          `tfsdk:"num"`
	BoolValue types.Bool            `tfsdk:"bool"`
	ObjectValue *JsonObjectMap      `tfsdk:"object"`
	ListValue *[]JsonElement        `tfsdk:"list"`
	ObjectsListValue *[]JsonElement `tfsdk:"objects_list"`
}
