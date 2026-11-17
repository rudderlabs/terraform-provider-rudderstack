package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

// warehouseImmutableFlags maps allowUsersContextTraits and underscoreDivideNumbers
// for the warehouse destinations. Both are "rs-immutable" upstream and the
// control plane stores a value for them at create time, so an update that omits
// them is rejected as a modification ("Field ... is immutable and cannot be
// modified"). Keeping them Computed-only reads the stored value into state and
// sends it back unchanged on every update. Users cannot set them, so upgrading
// the provider never plans a change, let alone a replacement, of an existing
// warehouse destination.
func warehouseImmutableFlags() ([]c.ConfigProperty, map[string]*schema.Schema) {
	properties := []c.ConfigProperty{
		c.Simple("allowUsersContextTraits", "allow_users_context_traits"),
		c.Simple("underscoreDivideNumbers", "underscore_divide_numbers"),
	}
	schemas := map[string]*schema.Schema{
		"allow_users_context_traits": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Read-only. Set by RudderStack when the destination is created and immutable afterwards.",
		},
		"underscore_divide_numbers": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Read-only. Set by RudderStack when the destination is created and immutable afterwards.",
		},
	}
	return properties, schemas
}
