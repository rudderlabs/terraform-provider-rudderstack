package accounts_test

import (
	"testing"

	sdkschema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	_ "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/accounts"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestPostgresAccountRegistration(t *testing.T) {
	entries := configs.Accounts.Entries()
	cm, ok := entries["postgres"]
	require.True(t, ok, "postgres must be registered in the Accounts registry")

	require.Equal(t, "SOURCE_POSTGRES", cm.APIType, "APIType must be SOURCE_POSTGRES")

	schema := cm.ConfigSchema
	require.NotNil(t, schema, "ConfigSchema must not be nil")

	for _, field := range []string{"host", "dbname", "user", "port"} {
		s, ok := schema[field]
		require.True(t, ok, "ConfigSchema must contain %q", field)
		require.True(t, s.Required, "%q must be Required", field)
		require.False(t, s.Sensitive, "%q must not be Sensitive", field)
		require.Equal(t, sdkschema.TypeString, s.Type, "%q must be a string", field)
	}

	sslMode, ok := schema["ssl_mode"]
	require.True(t, ok, "ConfigSchema must contain 'ssl_mode'")
	require.True(t, sslMode.Optional, "'ssl_mode' must be Optional")
	require.False(t, sslMode.Required, "'ssl_mode' must not be Required")

	password, ok := schema["password"]
	require.True(t, ok, "ConfigSchema must contain 'password'")
	require.True(t, password.Required, "'password' must be Required")
	require.True(t, password.Sensitive, "'password' must be Sensitive")

	// host, dbname, user, port, ssl_mode, password
	require.Len(t, cm.Properties, 6, "expected 6 property mappings")
}
