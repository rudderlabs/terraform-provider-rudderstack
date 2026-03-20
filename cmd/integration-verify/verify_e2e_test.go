package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
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

func writeTF(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "main.tf")
	require.NoError(t, os.WriteFile(path, []byte(content), 0600))
	return path
}

func TestE2E_MissingArgs(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run([]string{}, &stdout, &stderr)
	assert.Equal(t, 1, code)
	assert.Contains(t, stderr.String(), "Usage:")
}

func TestE2E_MissingFile(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run([]string{"-file", "", "-id", "some-id"}, &stdout, &stderr)
	assert.Equal(t, 1, code)
	assert.Contains(t, stderr.String(), "Usage:")
}

func TestE2E_MissingID(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run([]string{"-file", "/tmp/some.tf", "-id", ""}, &stdout, &stderr)
	assert.Equal(t, 1, code)
	assert.Contains(t, stderr.String(), "Usage:")
}

func TestE2E_InvalidTFFile(t *testing.T) {
	path := writeTF(t, `this is not valid HCL {{{`)

	var stdout, stderr bytes.Buffer

	t.Setenv("RUDDERSTACK_ACCESS_TOKEN", "test-token")
	code := run([]string{"-file", path, "-id", "some-id"}, &stdout, &stderr)
	assert.Equal(t, 1, code)
	assert.Contains(t, stderr.String(), "Error:")
}

func TestE2E_NoRudderstackResource(t *testing.T) {
	path := writeTF(t, `
resource "aws_s3_bucket" "example" {
  bucket = "my-bucket"
}
`)
	var stdout, stderr bytes.Buffer

	t.Setenv("RUDDERSTACK_ACCESS_TOKEN", "test-token")
	code := run([]string{"-file", path, "-id", "some-id"}, &stdout, &stderr)
	assert.Equal(t, 1, code)
	assert.Contains(t, stderr.String(), "no matching")
}

func TestE2E_MissingAccessToken(t *testing.T) {
	path := writeTF(t, `
resource "rudderstack_destination_webhook" "test" {
  name = "my-webhook"
  config {
    webhook_url = "https://example.com"
  }
}
`)
	var stdout, stderr bytes.Buffer

	t.Setenv("RUDDERSTACK_ACCESS_TOKEN", "")
	code := run([]string{"-file", path, "-id", "dest-1"}, &stdout, &stderr)
	assert.Equal(t, 1, code)
	assert.Contains(t, stderr.String(), "RUDDERSTACK_ACCESS_TOKEN")
}

func TestE2E_DestinationMatch(t *testing.T) {
	apiConfig := `{
		"webhookUrl": "https://example.com/hook",
		"webhookMethod": "POST",
		"extraField": "not-in-tf-but-thats-ok"
	}`
	server := newMockAPI(t, "destination", "dest-ok", json.RawMessage(apiConfig))
	defer server.Close()

	path := writeTF(t, `
resource "rudderstack_destination_webhook" "test" {
  name = "my-webhook"
  config {
    webhook_url    = "https://example.com/hook"
    webhook_method = "POST"
  }
}
`)
	var stdout, stderr bytes.Buffer

	t.Setenv("RUDDERSTACK_ACCESS_TOKEN", "test-token")
	t.Setenv("RUDDERSTACK_API_URL", server.URL)

	code := run([]string{"-file", path, "-id", "dest-ok"}, &stdout, &stderr)
	assert.Equal(t, 0, code)
	assert.Contains(t, stdout.String(), "PASS")
	assert.Empty(t, stderr.String())
}

func TestE2E_DestinationMismatch(t *testing.T) {
	apiConfig := `{
		"webhookUrl": "https://wrong.com/hook",
		"webhookMethod": "POST"
	}`
	server := newMockAPI(t, "destination", "dest-bad", json.RawMessage(apiConfig))
	defer server.Close()

	path := writeTF(t, `
resource "rudderstack_destination_webhook" "test" {
  name = "my-webhook"
  config {
    webhook_url    = "https://example.com/hook"
    webhook_method = "POST"
  }
}
`)
	var stdout, stderr bytes.Buffer

	t.Setenv("RUDDERSTACK_ACCESS_TOKEN", "test-token")
	t.Setenv("RUDDERSTACK_API_URL", server.URL)

	code := run([]string{"-file", path, "-id", "dest-bad"}, &stdout, &stderr)
	assert.Equal(t, 1, code)
	assert.Contains(t, stdout.String(), "FAIL")
	assert.Contains(t, stdout.String(), "Differences")
}

func TestE2E_SourceMatch(t *testing.T) {
	server := newMockAPI(t, "source", "src-ok", json.RawMessage(`{}`))
	defer server.Close()

	path := writeTF(t, `
resource "rudderstack_source_http" "test" {
  name = "my-source"
}
`)
	var stdout, stderr bytes.Buffer

	t.Setenv("RUDDERSTACK_ACCESS_TOKEN", "test-token")
	t.Setenv("RUDDERSTACK_API_URL", server.URL)

	code := run([]string{"-file", path, "-id", "src-ok"}, &stdout, &stderr)
	assert.Equal(t, 0, code)
	assert.Contains(t, stdout.String(), "PASS")
}

func TestE2E_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"internal error"}`))
	}))
	defer server.Close()

	path := writeTF(t, `
resource "rudderstack_destination_webhook" "test" {
  name = "my-webhook"
  config {
    webhook_url = "https://example.com"
  }
}
`)
	var stdout, stderr bytes.Buffer

	t.Setenv("RUDDERSTACK_ACCESS_TOKEN", "test-token")
	t.Setenv("RUDDERSTACK_API_URL", server.URL)

	code := run([]string{"-file", path, "-id", "dest-err"}, &stdout, &stderr)
	assert.Equal(t, 1, code)
	assert.Contains(t, stderr.String(), "Error:")
}

func TestE2E_TargetResourceFlag(t *testing.T) {
	apiConfig := `{
		"webhookUrl": "https://second.com",
		"webhookMethod": "GET"
	}`
	server := newMockAPI(t, "destination", "dest-2nd", json.RawMessage(apiConfig))
	defer server.Close()

	path := writeTF(t, `
resource "rudderstack_destination_webhook" "first" {
  name = "first"
  config {
    webhook_url    = "https://first.com"
    webhook_method = "POST"
  }
}

resource "rudderstack_destination_webhook" "second" {
  name = "second"
  config {
    webhook_url    = "https://second.com"
    webhook_method = "GET"
  }
}
`)
	var stdout, stderr bytes.Buffer

	t.Setenv("RUDDERSTACK_ACCESS_TOKEN", "test-token")
	t.Setenv("RUDDERSTACK_API_URL", server.URL)

	code := run([]string{"-file", path, "-id", "dest-2nd", "-resource", "second"}, &stdout, &stderr)
	assert.Equal(t, 0, code)
	assert.Contains(t, stdout.String(), "PASS")
}
