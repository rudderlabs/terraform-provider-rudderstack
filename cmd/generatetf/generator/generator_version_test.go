package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
	// treated the same as an unsupported apiType by the caller (skipped) —
	// this provider relies on the API always reporting a real version
	// (INT-6489); it does not coerce an absent/zero version to v1.
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
