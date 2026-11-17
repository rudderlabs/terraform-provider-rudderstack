package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

// warehouseImmutableFlags maps allowUsersContextTraits and underscoreDivideNumbers
// for the warehouse destinations. Both are "rs-immutable" upstream: an update must
// send back exactly what the control plane stored at create time, and destinations
// created before the flags existed have nothing stored, so the update must omit them.
// They are Computed strings ("true", "false" or "") rather than bools because
// d.Get cannot tell a null bool from false, and sending false for a flag that was
// never stored is rejected as a modification. Users cannot set them, so upgrading
// the provider never plans a change of an existing warehouse destination.
func warehouseImmutableFlags() ([]c.ConfigProperty, map[string]*schema.Schema) {
	properties := []c.ConfigProperty{
		storedBool("allowUsersContextTraits", "allow_users_context_traits"),
		storedBool("underscoreDivideNumbers", "underscore_divide_numbers"),
	}
	description := `Read-only. The value RudderStack stored when the destination was created ("true" or "false"), empty if none is stored. Immutable.`
	schemas := map[string]*schema.Schema{
		"allow_users_context_traits": {Type: schema.TypeString, Computed: true, Description: description},
		"underscore_divide_numbers":  {Type: schema.TypeString, Computed: true, Description: description},
	}
	return properties, schemas
}

// storedBool round-trips an API bool through a "true"/"false" state string,
// writing nothing on either side when the other side has no value.
func storedBool(apiKey, terraformKey string) c.ConfigProperty {
	return c.ConfigProperty{
		ToStateFunc: func(state, config string) (string, error) {
			if v := gjson.Get(config, apiKey); v.IsBool() {
				return sjson.Set(state, terraformKey, v.String())
			}
			return state, nil
		},
		FromStateFunc: func(config, state string) (string, error) {
			switch v := gjson.Get(state, terraformKey).String(); v {
			case "true", "false":
				return sjson.Set(config, apiKey, v == "true")
			}
			return config, nil
		},
	}
}
