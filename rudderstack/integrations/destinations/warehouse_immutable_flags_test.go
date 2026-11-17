package destinations_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

// The control plane rejects an update that changes an immutable flag, and it
// counts an omitted flag as a change. A destination created as `true` must
// therefore send `true` back on every update, read through state.
func TestWarehouseImmutableFlagsRoundTrip(t *testing.T) {
	entries := c.Destinations.Entries()
	for _, name := range []string{"bigquery", "postgres", "redshift", "s3_datalake", "snowflake"} {
		t.Run(name, func(t *testing.T) {
			cm, ok := entries[name]
			require.True(t, ok, "destination %q not registered", name)

			state, err := cm.APIToState(`{"allowUsersContextTraits": true, "underscoreDivideNumbers": true}`)
			require.NoError(t, err)

			api, err := cm.StateToAPI(state)
			require.NoError(t, err)
			require.True(t, gjson.Get(api, "allowUsersContextTraits").Bool(), api)
			require.True(t, gjson.Get(api, "underscoreDivideNumbers").Bool(), api)
		})
	}
}
