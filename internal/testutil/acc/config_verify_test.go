package acc

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

// redactedAPIConfigKeys must discover the API-side key of a Sensitive terraform
// field by running it through the real state->API transform.
func TestRedactedAPIConfigKeys(t *testing.T) {
	cm := configs.ConfigMeta{
		ConfigSchema: map[string]*schema.Schema{
			"api_secret": {Type: schema.TypeString, Sensitive: true},
			"api_key":    {Type: schema.TypeString},
		},
		Properties: []configs.ConfigProperty{
			configs.Simple("apiSecret", "api_secret"),
			configs.Simple("apiKey", "api_key"),
		},
	}

	got := redactedAPIConfigKeys(t, cm)
	assert.True(t, got["apiSecret"], "sensitive field's API key should be marked redacted")
	assert.False(t, got["apiKey"], "non-sensitive field's API key should not be marked redacted")
}

// A Sensitive field nested inside a block must still be discovered (its flat
// API key derived through the transform).
func TestRedactedAPIConfigKeys_Nested(t *testing.T) {
	cm := configs.ConfigMeta{
		ConfigSchema: map[string]*schema.Schema{
			"s3": {
				Type: schema.TypeList,
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"access_key":    {Type: schema.TypeString, Sensitive: true},
					"access_key_id": {Type: schema.TypeString, Sensitive: true},
					"bucket":        {Type: schema.TypeString},
				}},
			},
		},
		Properties: []configs.ConfigProperty{
			configs.Simple("accessKey", "s3.0.access_key"),
			configs.Simple("accessKeyID", "s3.0.access_key_id"),
			configs.Simple("bucketName", "s3.0.bucket"),
		},
	}

	got := redactedAPIConfigKeys(t, cm)
	assert.True(t, got["accessKey"], "nested sensitive field's API key should be redacted")
	assert.True(t, got["accessKeyID"], "nested sensitive field's API key should be redacted")
	assert.False(t, got["bucketName"], "non-sensitive nested field should not be redacted")
}

// A block whose fields are all Sensitive collapses on import and is ignored
// wholesale; a block with a non-secret field ignores only its sensitive leaf.
func TestSensitiveImportIgnorePaths(t *testing.T) {
	cm := configs.ConfigMeta{
		ConfigSchema: map[string]*schema.Schema{
			"key_based_authentication": {
				Type: schema.TypeList,
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"access_key":    {Type: schema.TypeString, Sensitive: true},
					"access_key_id": {Type: schema.TypeString, Sensitive: true},
				}},
			},
			"s3": {
				Type: schema.TypeList,
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"access_key": {Type: schema.TypeString, Sensitive: true},
					"bucket":     {Type: schema.TypeString},
				}},
			},
			"api_secret": {Type: schema.TypeString, Sensitive: true},
		},
	}

	got := cm.SensitiveImportIgnorePaths()
	assert.Contains(t, got, "key_based_authentication", "all-secret block ignored wholesale")
	assert.Contains(t, got, "s3.0.access_key", "partial block ignores only the secret leaf")
	assert.NotContains(t, got, "s3", "partial block must not be ignored wholesale")
	assert.Contains(t, got, "api_secret")
}

func TestRedactedAPIConfigKeys_NoSensitiveFields(t *testing.T) {
	cm := configs.ConfigMeta{
		ConfigSchema: map[string]*schema.Schema{"api_key": {Type: schema.TypeString}},
		Properties:   []configs.ConfigProperty{configs.Simple("apiKey", "api_key")},
	}
	assert.Empty(t, redactedAPIConfigKeys(t, cm))
}

// A redacted secret absent from the response is fine; any other missing field is not.
func TestCompareConfig_SkipsRedactedSecret(t *testing.T) {
	redacted := map[string]bool{"apiSecret": true}

	require.NoError(t, compareConfig(
		json.RawMessage(`{"apiKey":"abc"}`),
		`{"apiKey":"abc","apiSecret":"shh"}`,
		redacted,
	), "redacted apiSecret missing from response must not fail")

	err := compareConfig(
		json.RawMessage(`{"apiKey":"abc"}`),
		`{"apiKey":"abc","residencyServer":"standard"}`,
		redacted,
	)
	require.Error(t, err, "a non-redacted missing field must still fail")
	assert.Contains(t, err.Error(), "residencyServer")

	// Redacted fields aren't verified at all — the backend may return them
	// blanked in place (e.g. a Sensitive list emptied), so present-but-different
	// is tolerated too.
	require.NoError(t, compareConfig(
		json.RawMessage(`{"apiKey":"abc","apiSecret":""}`),
		`{"apiKey":"abc","apiSecret":"right"}`,
		redacted,
	), "a redacted field returned blanked must not fail")
}
