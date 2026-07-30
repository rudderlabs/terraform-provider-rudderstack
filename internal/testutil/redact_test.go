package testutil

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

// RedactedAPIConfigKeys must discover the API-side key of a Sensitive terraform
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

	got := RedactedAPIConfigKeys(t, cm)
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

	got := RedactedAPIConfigKeys(t, cm)
	assert.True(t, got["accessKey"], "nested sensitive field's API key should be redacted")
	assert.True(t, got["accessKeyID"], "nested sensitive field's API key should be redacted")
	assert.False(t, got["bucketName"], "non-sensitive nested field should not be redacted")
}

func TestRedactedAPIConfigKeys_NoSensitiveFields(t *testing.T) {
	cm := configs.ConfigMeta{
		ConfigSchema: map[string]*schema.Schema{"api_key": {Type: schema.TypeString}},
		Properties:   []configs.ConfigProperty{configs.Simple("apiKey", "api_key")},
	}
	assert.Empty(t, RedactedAPIConfigKeys(t, cm))
}

// ConfigHasRedactedSecret drives the perpetual diff: only true when the config
// actually sets a redacted key to a non-empty value.
func TestConfigHasRedactedSecret(t *testing.T) {
	redacted := map[string]bool{"apiSecret": true, "headers": true}

	assert.True(t, ConfigHasRedactedSecret(`{"apiKey":"x","apiSecret":"shh"}`, redacted),
		"a set scalar secret counts")
	assert.False(t, ConfigHasRedactedSecret(`{"apiKey":"x"}`, redacted),
		"no redacted key set -> no diff (e.g. secret only set on update)")
	assert.False(t, ConfigHasRedactedSecret(`{"apiSecret":""}`, redacted),
		"an empty redacted value does not count")
	assert.True(t, ConfigHasRedactedSecret(`{"headers":[{"from":"a","to":"b"}]}`, redacted),
		"a non-empty redacted collection counts")
	assert.False(t, ConfigHasRedactedSecret(`{"headers":[]}`, redacted),
		"an empty redacted collection does not count")
}
