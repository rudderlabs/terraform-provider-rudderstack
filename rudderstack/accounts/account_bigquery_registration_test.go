package accounts_test

import (
	"testing"

	sdkschema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
	_ "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/accounts"
)

func TestBigQueryAccountRegistration(t *testing.T) {
	entries := configs.Accounts.Entries()
	cm, ok := entries["bigquery"]
	require.True(t, ok, "bigquery must be registered in the Accounts registry")

	require.Equal(t, "SOURCE_BIGQUERY", cm.APIType, "APIType must be SOURCE_BIGQUERY")

	schema := cm.ConfigSchema
	require.NotNil(t, schema, "ConfigSchema must not be nil")

	// project: Required
	project, ok := schema["project"]
	require.True(t, ok, "ConfigSchema must contain 'project'")
	require.True(t, project.Required, "'project' must be Required")

	// location: Optional
	location, ok := schema["location"]
	require.True(t, ok, "ConfigSchema must contain 'location'")
	require.True(t, location.Optional, "'location' must be Optional")
	require.False(t, location.Required, "'location' must not be Required")

	// credentials: Required and Sensitive
	credentials, ok := schema["credentials"]
	require.True(t, ok, "ConfigSchema must contain 'credentials'")
	require.True(t, credentials.Required, "'credentials' must be Required")
	require.True(t, credentials.Sensitive, "'credentials' must be Sensitive")

	// Field types
	require.Equal(t, sdkschema.TypeString, project.Type, "'project' must be a string")
	require.Equal(t, sdkschema.TypeString, location.Type, "'location' must be a string")
	require.Equal(t, sdkschema.TypeString, credentials.Type, "'credentials' must be a string")

	// API field mappings (Properties) must be wired for project, location and credentials.
	require.Len(t, cm.Properties, 3, "expected project, location and credentials property mappings")
}
