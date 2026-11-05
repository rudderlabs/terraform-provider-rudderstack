package generator

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rudderlabs/rudder-iac/api/client"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

// These tests exercise configMeta and configMetaByVersion directly (whitebox,
// package generator) against small, self-contained registries — not the real,
// process-global configs.Destinations/configs.Sources populated by every
// integration's init().

func TestConfigMeta_FirstMatchOnAPIType(t *testing.T) {
	entries := map[string]configs.ConfigMeta{
		"javascript": {APIType: "Javascript"},
		"http":       {APIType: "HTTP"},
	}

	terraformType, cm := configMeta(entries, "HTTP")
	require.NotNil(t, cm)
	assert.Equal(t, "http", terraformType)
	assert.Equal(t, "HTTP", cm.APIType)
}

func TestConfigMeta_NoMatch(t *testing.T) {
	entries := map[string]configs.ConfigMeta{
		"http": {APIType: "HTTP"},
	}

	terraformType, cm := configMeta(entries, "Unknown")
	assert.Nil(t, cm)
	assert.Empty(t, terraformType)
}

func TestConfigMetaByVersion_MatchesExactAPITypeAndVersion(t *testing.T) {
	entries := map[string]configs.ConfigMeta{
		"braze":    {APIType: "BRAZE", Version: 1},
		"braze_v2": {APIType: "BRAZE", Version: 2},
	}

	t.Run("v1", func(t *testing.T) {
		terraformType, cm := configMetaByVersion(entries, "BRAZE", 1)
		require.NotNil(t, cm)
		assert.Equal(t, "braze", terraformType)
		assert.Equal(t, 1, cm.Version)
	})

	t.Run("v2", func(t *testing.T) {
		terraformType, cm := configMetaByVersion(entries, "BRAZE", 2)
		require.NotNil(t, cm)
		assert.Equal(t, "braze_v2", terraformType)
		assert.Equal(t, 2, cm.Version)
	})
}

func TestConfigMetaByVersion_MismatchedVersionIsSkipped(t *testing.T) {
	// A destination reporting a version with no matching registry entry is
	// skipped by the caller — this provider relies on the API always reporting
	// a real version (INT-6489); it does not coerce an absent/zero version to v1.
	entries := map[string]configs.ConfigMeta{
		"braze": {APIType: "BRAZE", Version: 1},
	}

	terraformType, cm := configMetaByVersion(entries, "BRAZE", 3)
	assert.Nil(t, cm)
	assert.Empty(t, terraformType)
}

func TestConfigMetaByVersion_ZeroVersionDoesNotCoerceToV1(t *testing.T) {
	entries := map[string]configs.ConfigMeta{
		"braze": {APIType: "BRAZE", Version: 1},
	}

	// dst.Version == 0 (e.g. an old API record that hasn't been backfilled by
	// the backend yet) does not match the v1 entry: no implicit coercion.
	terraformType, cm := configMetaByVersion(entries, "BRAZE", 0)
	assert.Nil(t, cm)
	assert.Empty(t, terraformType)
}

func TestConfigMetaByVersion_UnknownAPIType(t *testing.T) {
	entries := map[string]configs.ConfigMeta{
		"braze": {APIType: "BRAZE", Version: 1},
	}

	terraformType, cm := configMetaByVersion(entries, "UNKNOWN", 1)
	assert.Nil(t, cm)
	assert.Empty(t, terraformType)
}

func TestUnsupportedDestinationReason(t *testing.T) {
	entries := map[string]configs.ConfigMeta{
		"braze":    {APIType: "BRAZE", Version: 1},
		"braze_v2": {APIType: "BRAZE", Version: 2},
	}

	t.Run("unknown api type", func(t *testing.T) {
		got := unsupportedDestinationReason(entries, "UNKNOWN", 1)
		assert.Equal(t, "type 'UNKNOWN' is not supported", got)
	})

	t.Run("unsupported nonzero version", func(t *testing.T) {
		got := unsupportedDestinationReason(entries, "BRAZE", 3)
		assert.Equal(t, "type 'BRAZE' version 3 is not supported (available versions: 1, 2)", got)
	})

	t.Run("zero version", func(t *testing.T) {
		got := unsupportedDestinationReason(entries, "BRAZE", 0)
		assert.Equal(t, "type 'BRAZE' version 0 is not supported (available versions: 1, 2)", got)
	})
}

func TestGenerateTerraform_LogsVersionMismatch(t *testing.T) {
	buf := captureLogger(t)

	destinations := []client.Destination{
		{ID: "id-braze", Name: "braze", Type: "BRAZE", Version: 0},
	}

	data, err := GenerateTerraform(nil, destinations, nil, nil, nil)
	require.NoError(t, err)
	assert.Empty(t, string(data))
	assert.Contains(t, buf.String(), "could not generate resource block for destination 'id-braze': type 'BRAZE' version 0 is not supported (available versions: 1)")
}

func TestGenerateImportScript_LogsVersionMismatch(t *testing.T) {
	buf := captureLogger(t)

	destinations := []client.Destination{
		{ID: "id-braze", Name: "braze", Type: "BRAZE", Version: 3},
	}

	data, err := GenerateImportScript(nil, destinations, nil, nil, nil)
	require.NoError(t, err)
	assert.Empty(t, string(data))
	assert.Contains(t, buf.String(), "skipping destination 'id-braze': type 'BRAZE' version 3 is not supported (available versions: 1)")
}

func TestGenerateHTTPVersionsAndBodyFormats(t *testing.T) {
	destinations := []client.Destination{
		{ID: "http-v1", Name: "old", Type: "HTTP", Version: 1, Config: json.RawMessage(`{"apiUrl":"https://example.com/old","auth":"noAuth","method":"POST","format":"JSON","isBatchingEnabled":false,"isDefaultMapping":true}`)},
		{ID: "http-xml", Name: "xml", Type: "HTTP", Version: 2, Config: json.RawMessage(`{"apiUrl":"https://example.com/xml","auth":"noAuth","method":"POST","format":"XML","xmlRootKey":"Envelope","requestBodyTemplate":"{\"event\": $$.events}","maxBatchSize":"10"}`)},
		{ID: "http-form", Name: "form", Type: "HTTP", Version: 2, Config: json.RawMessage(`{"apiUrl":"https://example.com/form","auth":"noAuth","method":"POST","format":"FORM","requestBodyTemplate":"{\"email\": $$.events[0].traits.email}"}`)},
	}

	hcl, err := GenerateTerraform(nil, destinations, nil, nil, nil)
	require.NoError(t, err)
	hclText := string(hcl)
	assert.Contains(t, hclText, `resource "rudderstack_destination_http" "dst_http-v1"`)
	assert.Contains(t, hclText, `resource "rudderstack_destination_http_v2" "dst_http-xml"`)
	formResource := `resource "rudderstack_destination_http_v2" "dst_http-form"`
	assert.Contains(t, hclText, formResource)
	assert.Contains(t, hclText, `xml_root_key`)
	assert.Contains(t, hclText, `"Envelope"`)
	assert.Contains(t, hclText, `max_batch_size`)
	assert.NotContains(t, hclText, `is_batching_enabled = true`)
	// FORM has no batch size or XML root key in the generated HCL.
	form := hclText[strings.Index(hclText, formResource):]
	assert.Regexp(t, `(?m)^\s*format\s+= "FORM"$`, form)
	assert.NotContains(t, form, `max_batch_size`)
	assert.NotContains(t, form, `xml_root_key`)

	imports, err := GenerateImportScript(nil, destinations, nil, nil, nil)
	require.NoError(t, err)
	assert.Contains(t, string(imports), `terraform import "rudderstack_destination_http.dst_http-v1" "http-v1"`)
	assert.Contains(t, string(imports), `terraform import "rudderstack_destination_http_v2.dst_http-xml" "http-xml"`)
	assert.Contains(t, string(imports), `terraform import "rudderstack_destination_http_v2.dst_http-form" "http-form"`)
}

func captureLogger(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	prev := logger
	logger = log.New(&buf, "", 0)
	t.Cleanup(func() { logger = prev })
	return &buf
}
