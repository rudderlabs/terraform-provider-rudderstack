package destinations_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

// The control plane rejects an update that changes an immutable flag, and it
// counts both a flipped value and false-for-undefined as a change. Updates must
// send back exactly what the API stored, and nothing when it stored nothing.
func TestWarehouseImmutableFlagsRoundTrip(t *testing.T) {
	cases := map[string]struct {
		api  string
		want string // "" means the key is absent on both sides
	}{
		"stored true":  {`{"allowUsersContextTraits": true, "underscoreDivideNumbers": true}`, "true"},
		"stored false": {`{"allowUsersContextTraits": false, "underscoreDivideNumbers": false}`, "false"},
		"not stored":   {`{}`, ""},
	}
	entries := c.Destinations.Entries()
	for _, name := range []string{"bigquery", "postgres", "redshift", "s3_datalake", "snowflake"} {
		cm, ok := entries[name]
		require.True(t, ok, "destination %q not registered", name)
		for caseName, tc := range cases {
			t.Run(name+"/"+caseName, func(t *testing.T) {
				state, err := cm.APIToState(tc.api)
				require.NoError(t, err)
				api, err := cm.StateToAPI(state)
				require.NoError(t, err)

				for apiKey, tfKey := range map[string]string{
					"allowUsersContextTraits": "allow_users_context_traits",
					"underscoreDivideNumbers": "underscore_divide_numbers",
				} {
					s, a := gjson.Get(state, tfKey), gjson.Get(api, apiKey)
					if tc.want == "" {
						require.False(t, s.Exists(), state)
						require.False(t, a.Exists(), api)
						continue
					}
					require.Equal(t, tc.want, s.String(), state)
					require.True(t, a.IsBool(), api)
					require.Equal(t, tc.want == "true", a.Bool(), api)
				}
			})
		}
	}
}
