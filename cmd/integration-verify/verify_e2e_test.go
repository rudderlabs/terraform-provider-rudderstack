package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/integrations"
)

// newMockAPI creates an httptest server that simulates the RudderStack API
// for verifying onboarded integrations.
func newMockAPI(t *testing.T, kind, id string, config json.RawMessage) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := fmt.Sprintf("/%ss/%s", kind, id)
		if r.URL.Path != expectedPath || r.Method != http.MethodGet {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		var resp interface{}
		switch kind {
		case "destination":
			resp = map[string]interface{}{
				"destination": map[string]interface{}{
					"id": id, "name": "test", "type": "WEBHOOK", "enabled": true,
					"config": json.RawMessage(config),
				},
			}
		case "source":
			resp = map[string]interface{}{
				"source": map[string]interface{}{
					"id": id, "name": "test", "type": "HTTP", "enabled": true,
					"config": json.RawMessage(config),
				},
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
}

func TestE2E_InvalidState(t *testing.T) {
	t.Setenv("RUDDERSTACK_ACCESS_TOKEN", "test-token")
	err := verifyFromState([]byte(`{}`), "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no resources found")
}

func TestE2E_NoRudderstackResource(t *testing.T) {
	stateJSON := buildStateJSON(t, stateResource{
		Type:   "aws_s3_bucket",
		Name:   "example",
		ID:     "bucket-1",
		TFName: "my-bucket",
	})

	t.Setenv("RUDDERSTACK_ACCESS_TOKEN", "test-token")
	err := verifyFromState(stateJSON, "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no matching")
}

func TestE2E_MissingAccessToken(t *testing.T) {
	stateJSON := buildStateJSON(t, stateResource{
		Type:   "rudderstack_destination_webhook",
		Name:   "test",
		ID:     "dest-1",
		TFName: "my-webhook",
		Config: map[string]interface{}{"webhook_url": "https://example.com"},
	})

	t.Setenv("RUDDERSTACK_ACCESS_TOKEN", "")
	err := verifyFromState(stateJSON, "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "RUDDERSTACK_ACCESS_TOKEN")
}

func TestE2E_DestinationMatch(t *testing.T) {
	apiConfig := `{
		"webhookUrl": "https://example.com/hook",
		"webhookMethod": "POST"
	}`
	server := newMockAPI(t, "destination", "dest-ok", json.RawMessage(apiConfig))
	defer server.Close()

	stateJSON := buildStateJSON(t, stateResource{
		Type:   "rudderstack_destination_webhook",
		Name:   "test",
		ID:     "dest-ok",
		TFName: "my-webhook",
		Config: map[string]interface{}{
			"webhook_url":    "https://example.com/hook",
			"webhook_method": "POST",
		},
	})

	t.Setenv("RUDDERSTACK_ACCESS_TOKEN", "test-token")
	t.Setenv("RUDDERSTACK_API_URL", server.URL)

	err := verifyFromState(stateJSON, "")
	require.NoError(t, err)
}

func TestE2E_DestinationMismatch(t *testing.T) {
	apiConfig := `{
		"webhookUrl": "https://wrong.com/hook",
		"webhookMethod": "POST"
	}`
	server := newMockAPI(t, "destination", "dest-bad", json.RawMessage(apiConfig))
	defer server.Close()

	stateJSON := buildStateJSON(t, stateResource{
		Type:   "rudderstack_destination_webhook",
		Name:   "test",
		ID:     "dest-bad",
		TFName: "my-webhook",
		Config: map[string]interface{}{
			"webhook_url":    "https://example.com/hook",
			"webhook_method": "POST",
		},
	})

	t.Setenv("RUDDERSTACK_ACCESS_TOKEN", "test-token")
	t.Setenv("RUDDERSTACK_API_URL", server.URL)

	err := verifyFromState(stateJSON, "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "verification failed")
}

func TestE2E_SourceMatch(t *testing.T) {
	server := newMockAPI(t, "source", "src-ok", json.RawMessage(`{}`))
	defer server.Close()

	stateJSON := buildStateJSON(t, stateResource{
		Type:   "rudderstack_source_http",
		Name:   "test",
		ID:     "src-ok",
		TFName: "my-source",
	})

	t.Setenv("RUDDERSTACK_ACCESS_TOKEN", "test-token")
	t.Setenv("RUDDERSTACK_API_URL", server.URL)

	err := verifyFromState(stateJSON, "")
	require.NoError(t, err)
}

func TestE2E_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"internal error"}`))
	}))
	defer server.Close()

	stateJSON := buildStateJSON(t, stateResource{
		Type:   "rudderstack_destination_webhook",
		Name:   "test",
		ID:     "dest-err",
		TFName: "my-webhook",
		Config: map[string]interface{}{"webhook_url": "https://example.com"},
	})

	t.Setenv("RUDDERSTACK_ACCESS_TOKEN", "test-token")
	t.Setenv("RUDDERSTACK_API_URL", server.URL)

	err := verifyFromState(stateJSON, "")
	require.Error(t, err)
}

func TestE2E_FacebookPixelFullConfig(t *testing.T) {
	apiConfig := `{
		"pixelId": "facebook-pixel-id-123",
		"accessToken": "fb-access-token-456",
		"standardPageCall": true,
		"valueFieldIdentifier": "properties.price",
		"advancedMapping": true,
		"testDestination": true,
		"testEventCode": "TEST123",
		"eventsToEvents": [
			{"from": "signup", "to": "CompleteRegistration"},
			{"from": "purchase", "to": "Purchase"}
		],
		"eventCustomProperties": [
			{"eventCustomProperties": "customProp1"},
			{"eventCustomProperties": "customProp2"}
		],
		"blacklistPiiProperties": [
			{"blacklistPiiProperties": "email", "blacklistPiiHash": true},
			{"blacklistPiiProperties": "phone", "blacklistPiiHash": false}
		],
		"whitelistPiiProperties": [
			{"whitelistPiiProperties": "firstName"},
			{"whitelistPiiProperties": "lastName"}
		],
		"categoryToContent": [
			{"from": "games", "to": "product"}
		],
		"legacyConversionPixelId": {"from": "ViewContent", "to": "pixel-789"},
		"useNativeSDK": {"web": true},
		"blacklistedEvents": [
			{"eventName": "checkout_started"},
			{"eventName": "product_viewed"}
		],
		"eventFilteringOption": "blacklistedEvents"
	}`
	server := newMockAPI(t, "destination", "dest-fb", json.RawMessage(apiConfig))
	defer server.Close()

	stateJSON := buildStateJSON(t, stateResource{
		Type:   "rudderstack_destination_facebook_pixel",
		Name:   "test",
		ID:     "dest-fb",
		TFName: "e2e-test-facebook-pixel",
		Config: map[string]interface{}{
			"pixel_id":               "facebook-pixel-id-123",
			"access_token":           "fb-access-token-456",
			"standard_page_call":     true,
			"value_field_identifier": "properties.price",
			"advanced_mapping":       true,
			"test_destination":       true,
			"test_event_code":        "TEST123",
			"events_to_events": []interface{}{
				map[string]interface{}{"from": "signup", "to": "CompleteRegistration"},
				map[string]interface{}{"from": "purchase", "to": "Purchase"},
			},
			"event_custom_properties": []interface{}{"customProp1", "customProp2"},
			"blacklist_pii_properties": []interface{}{
				map[string]interface{}{"property": "email", "hash": true},
				map[string]interface{}{"property": "phone", "hash": false},
			},
			"whitelist_pii_properties": []interface{}{
				map[string]interface{}{"property": "firstName"},
				map[string]interface{}{"property": "lastName"},
			},
			"category_to_content": []interface{}{
				map[string]interface{}{"from": "games", "to": "product"},
			},
			"legacy_conversion_pixel_id": []interface{}{
				map[string]interface{}{"from": "ViewContent", "to": "pixel-789"},
			},
			"use_native_sdk": []interface{}{
				map[string]interface{}{"web": true},
			},
			"event_filtering": []interface{}{
				map[string]interface{}{
					"blacklist": []interface{}{"checkout_started", "product_viewed"},
				},
			},
		},
	})

	t.Setenv("RUDDERSTACK_ACCESS_TOKEN", "test-token")
	t.Setenv("RUDDERSTACK_API_URL", server.URL)

	err := verifyFromState(stateJSON, "")
	require.NoError(t, err)
}

func TestE2E_TargetResourceFlag(t *testing.T) {
	apiConfig := `{
		"webhookUrl": "https://second.com",
		"webhookMethod": "GET"
	}`
	server := newMockAPI(t, "destination", "dest-2nd", json.RawMessage(apiConfig))
	defer server.Close()

	stateJSON := buildStateJSON(t,
		stateResource{
			Type:   "rudderstack_destination_webhook",
			Name:   "first",
			ID:     "dest-1st",
			TFName: "first",
			Config: map[string]interface{}{
				"webhook_url":    "https://first.com",
				"webhook_method": "POST",
			},
		},
		stateResource{
			Type:   "rudderstack_destination_webhook",
			Name:   "second",
			ID:     "dest-2nd",
			TFName: "second",
			Config: map[string]interface{}{
				"webhook_url":    "https://second.com",
				"webhook_method": "GET",
			},
		},
	)

	t.Setenv("RUDDERSTACK_ACCESS_TOKEN", "test-token")
	t.Setenv("RUDDERSTACK_API_URL", server.URL)

	err := verifyFromState(stateJSON, "second")
	require.NoError(t, err)
}
