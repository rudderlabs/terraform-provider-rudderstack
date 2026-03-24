package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractResourceType(t *testing.T) {
	tests := []struct {
		input        string
		expectedKind string
		expectedType string
		expectErr    bool
	}{
		{"rudderstack_destination_webhook", "destination", "webhook", false},
		{"rudderstack_destination_google_analytics", "destination", "google_analytics", false},
		{"rudderstack_source_shopify", "source", "shopify", false},
		{"rudderstack_source_http", "source", "http", false},
		{"rudderstack_connection_foo", "", "", true},
		{"aws_s3_bucket", "", "", true},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			kind, intType, err := ExtractResourceType(tc.input)
			if tc.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedKind, kind)
			assert.Equal(t, tc.expectedType, intType)
		})
	}
}



// buildStateJSON constructs a terraform show -json output for testing.
func buildStateJSON(t *testing.T, resources ...stateResource) []byte {
	t.Helper()

	type tfResource struct {
		Type   string                 `json:"type"`
		Name   string                 `json:"name"`
		Values map[string]interface{} `json:"values"`
	}

	var tfResources []tfResource
	for _, r := range resources {
		values := map[string]interface{}{
			"id":   r.ID,
			"name": r.TFName,
		}
		if r.Config != nil {
			values["config"] = []interface{}{r.Config}
		}
		tfResources = append(tfResources, tfResource{
			Type:   r.Type,
			Name:   r.Name,
			Values: values,
		})
	}

	state := map[string]interface{}{
		"values": map[string]interface{}{
			"root_module": map[string]interface{}{
				"resources": tfResources,
			},
		},
	}

	data, err := json.Marshal(state)
	require.NoError(t, err)
	return data
}

type stateResource struct {
	Type   string                 // terraform resource type, e.g. "rudderstack_destination_webhook"
	Name   string                 // terraform resource label, e.g. "test"
	ID     string                 // resource ID
	TFName string                 // the name attribute in values
	Config map[string]interface{} // config block content (nil for no config)
}

func TestParseTerraformState(t *testing.T) {
	t.Run("parses destination resource", func(t *testing.T) {
		stateJSON := buildStateJSON(t, stateResource{
			Type:   "rudderstack_destination_webhook",
			Name:   "test",
			ID:     "dest-123",
			TFName: "my-webhook",
			Config: map[string]interface{}{
				"webhook_url":    "https://example.com",
				"webhook_method": "POST",
			},
		})

		resources, err := ParseTerraformState(stateJSON, "")
		require.NoError(t, err)
		require.Len(t, resources, 1)

		r := resources[0]
		assert.Equal(t, "destination", r.Kind)
		assert.Equal(t, "webhook", r.IntegrationType)
		assert.Equal(t, "rudderstack_destination_webhook", r.ResourceType)
		assert.Equal(t, "test", r.Name)
		assert.Equal(t, "dest-123", r.ResourceID)
		assert.Contains(t, r.ConfigState, "webhook_url")
	})

	t.Run("parses source resource", func(t *testing.T) {
		stateJSON := buildStateJSON(t, stateResource{
			Type:   "rudderstack_source_http",
			Name:   "test",
			ID:     "src-456",
			TFName: "my-source",
		})

		resources, err := ParseTerraformState(stateJSON, "")
		require.NoError(t, err)
		require.Len(t, resources, 1)

		r := resources[0]
		assert.Equal(t, "source", r.Kind)
		assert.Equal(t, "http", r.IntegrationType)
		assert.Equal(t, "src-456", r.ResourceID)
		assert.Equal(t, "{}", r.ConfigState)
	})

	t.Run("filters by target resource name", func(t *testing.T) {
		stateJSON := buildStateJSON(t,
			stateResource{
				Type:   "rudderstack_destination_webhook",
				Name:   "first",
				ID:     "dest-1",
				TFName: "first",
				Config: map[string]interface{}{"webhook_url": "https://first.com"},
			},
			stateResource{
				Type:   "rudderstack_destination_webhook",
				Name:   "second",
				ID:     "dest-2",
				TFName: "second",
				Config: map[string]interface{}{"webhook_url": "https://second.com"},
			},
		)

		resources, err := ParseTerraformState(stateJSON, "second")
		require.NoError(t, err)
		require.Len(t, resources, 1)
		assert.Equal(t, "second", resources[0].Name)
		assert.Equal(t, "dest-2", resources[0].ResourceID)
	})

	t.Run("skips non-rudderstack resources", func(t *testing.T) {
		stateJSON := buildStateJSON(t,
			stateResource{
				Type:   "aws_s3_bucket",
				Name:   "example",
				ID:     "bucket-1",
				TFName: "my-bucket",
			},
			stateResource{
				Type:   "rudderstack_destination_webhook",
				Name:   "test",
				ID:     "dest-1",
				TFName: "my-webhook",
				Config: map[string]interface{}{"webhook_url": "https://example.com"},
			},
		)

		resources, err := ParseTerraformState(stateJSON, "")
		require.NoError(t, err)
		require.Len(t, resources, 1)
		assert.Equal(t, "rudderstack_destination_webhook", resources[0].ResourceType)
	})

	t.Run("returns error for no matching resource", func(t *testing.T) {
		stateJSON := buildStateJSON(t, stateResource{
			Type:   "aws_s3_bucket",
			Name:   "example",
			ID:     "bucket-1",
			TFName: "my-bucket",
		})

		_, err := ParseTerraformState(stateJSON, "")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no matching")
	})

	t.Run("returns error for empty state", func(t *testing.T) {
		_, err := ParseTerraformState([]byte(`{}`), "")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no resources found")
	})

	t.Run("returns multiple resources", func(t *testing.T) {
		stateJSON := buildStateJSON(t,
			stateResource{
				Type:   "rudderstack_destination_webhook",
				Name:   "first",
				ID:     "dest-1",
				TFName: "first",
				Config: map[string]interface{}{"webhook_url": "https://first.com"},
			},
			stateResource{
				Type:   "rudderstack_source_http",
				Name:   "second",
				ID:     "src-1",
				TFName: "second",
			},
		)

		resources, err := ParseTerraformState(stateJSON, "")
		require.NoError(t, err)
		require.Len(t, resources, 2)
		assert.Equal(t, "destination", resources[0].Kind)
		assert.Equal(t, "source", resources[1].Kind)
	})
}
