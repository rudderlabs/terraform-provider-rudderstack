package accounts_test

import (
	"testing"

	sdkschema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	_ "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/accounts"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestSnowflakeAccountRegistration(t *testing.T) {
	entries := configs.Accounts.Entries()
	cm, ok := entries["snowflake"]
	require.True(t, ok, "snowflake must be registered in the Accounts registry")

	require.Equal(t, "SOURCE_SNOWFLAKE", cm.APIType, "APIType must be SOURCE_SNOWFLAKE")

	schema := cm.ConfigSchema
	require.NotNil(t, schema, "ConfigSchema must not be nil")

	// Required, non-secret option fields.
	for _, field := range []string{"account", "dbname", "warehouse", "user", "authentication_type"} {
		s, ok := schema[field]
		require.True(t, ok, "ConfigSchema must contain %q", field)
		require.True(t, s.Required, "%q must be Required", field)
		require.False(t, s.Sensitive, "%q must not be Sensitive", field)
		require.Equal(t, sdkschema.TypeString, s.Type, "%q must be a string", field)
	}

	// role: Optional, non-secret.
	role, ok := schema["role"]
	require.True(t, ok, "ConfigSchema must contain 'role'")
	require.True(t, role.Optional, "'role' must be Optional")
	require.False(t, role.Required, "'role' must not be Required")

	// The three secrets: Optional (exactly one auth mode at a time) and Sensitive.
	for _, field := range []string{"password", "private_key", "private_key_passphrase"} {
		s, ok := schema[field]
		require.True(t, ok, "ConfigSchema must contain %q", field)
		require.True(t, s.Optional, "%q must be Optional", field)
		require.True(t, s.Sensitive, "%q must be Sensitive", field)
	}

	// account, dbname, warehouse, user, role, authentication_type, password, private_key, private_key_passphrase
	require.Len(t, cm.Properties, 9, "expected 9 property mappings")
}
