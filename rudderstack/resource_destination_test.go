package rudderstack

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rudderlabs/rudder-iac/api/client"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/integrations/destinations"
)

func TestResourceDestinationConsentManagementDuplicateProvider(t *testing.T) {
	_, consentSchema := destinations.GetConfigMetaForGenericConsentManagement([]string{"web", "android"})

	cm := configs.ConfigMeta{
		APIType:      "TEST",
		ConfigSchema: consentSchema,
	}

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"rudderstack": func() (*schema.Provider, error) {
				return &schema.Provider{
					ConfigureContextFunc: func(_ context.Context, _ *schema.ResourceData) (interface{}, diag.Diagnostics) {
						return nil, nil
					},
					ResourcesMap: map[string]*schema.Resource{
						"rudderstack_destination_test": resourceDestination(cm),
					},
				}, nil
			},
		},
		Steps: []resource.TestStep{
			{
				PlanOnly: true,
				Config: `
					resource "rudderstack_destination_test" "example" {
						name = "test-destination"

						config {
							consent_management {
								web = [
									{
										provider = "oneTrust"
										consents = ["a"]
										resolution_strategy = ""
									},
									{
										provider = "oneTrust"
										consents = ["b"]
										resolution_strategy = ""
									}
								]
							}
						}
					}
				`,
				ExpectError: regexp.MustCompile(`duplicate consent_management provider "oneTrust" configured for source type "web"`),
			},
		},
	})
}

func TestResourceDestinationConsentManagementRejectsEmptyConsents(t *testing.T) {
	_, consentSchema := destinations.GetConfigMetaForGenericConsentManagement([]string{"web", "android"})

	cm := configs.ConfigMeta{
		APIType:      "TEST",
		ConfigSchema: consentSchema,
	}

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"rudderstack": func() (*schema.Provider, error) {
				return &schema.Provider{
					ConfigureContextFunc: func(_ context.Context, _ *schema.ResourceData) (interface{}, diag.Diagnostics) {
						return nil, nil
					},
					ResourcesMap: map[string]*schema.Resource{
						"rudderstack_destination_test": resourceDestination(cm),
					},
				}, nil
			},
		},
		Steps: []resource.TestStep{
			{
				PlanOnly: true,
				Config: `
					resource "rudderstack_destination_test" "example" {
						name = "test-destination"

						config {
							consent_management {
								web = [
									{
										provider = "oneTrust"
										consents = []
										resolution_strategy = ""
									}
								]
							}
						}
					}
				`,
				ExpectError: regexp.MustCompile(`consents requires 1 item\s+minimum`),
			},
		},
	})
}

func TestResourceDestinationConsentManagementRejectsBlankConsentValue(t *testing.T) {
	_, consentSchema := destinations.GetConfigMetaForGenericConsentManagement([]string{"web", "android"})

	cm := configs.ConfigMeta{
		APIType:      "TEST",
		ConfigSchema: consentSchema,
	}

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"rudderstack": func() (*schema.Provider, error) {
				return &schema.Provider{
					ConfigureContextFunc: func(_ context.Context, _ *schema.ResourceData) (interface{}, diag.Diagnostics) {
						return nil, nil
					},
					ResourcesMap: map[string]*schema.Resource{
						"rudderstack_destination_test": resourceDestination(cm),
					},
				}, nil
			},
		},
		Steps: []resource.TestStep{
			{
				PlanOnly: true,
				Config: `
					resource "rudderstack_destination_test" "example" {
						name = "test-destination"

						config {
							consent_management {
								web = [
									{
										provider = "oneTrust"
										consents = ["   "]
										resolution_strategy = ""
									}
								]
							}
						}
					}
				`,
				ExpectError: regexp.MustCompile(`empty string or whitespace`),
			},
		},
	})
}

func TestResourceDestinationConsentManagementAllowsDistinctAndPerSourceType(t *testing.T) {
	_, consentSchema := destinations.GetConfigMetaForGenericConsentManagement([]string{"web", "android"})

	cm := configs.ConfigMeta{
		APIType:      "TEST",
		ConfigSchema: consentSchema,
	}

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"rudderstack": func() (*schema.Provider, error) {
				return &schema.Provider{
					ConfigureContextFunc: func(_ context.Context, _ *schema.ResourceData) (interface{}, diag.Diagnostics) {
						return nil, nil
					},
					ResourcesMap: map[string]*schema.Resource{
						"rudderstack_destination_test": resourceDestination(cm),
					},
				}, nil
			},
		},
		Steps: []resource.TestStep{
			{
				PlanOnly: true,
				Config: `
					resource "rudderstack_destination_test" "example" {
						name = "test-destination"

						config {
							consent_management {
								web = [
									{
										provider = "oneTrust"
										consents = ["a"]
										resolution_strategy = ""
									},
									{
										provider = "ketch"
										consents = ["b"]
										resolution_strategy = ""
									}
								]
								android = [
									{
										provider = "oneTrust"
										consents = ["c"]
										resolution_strategy = ""
									}
								]
							}
						}
					}
				`,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestHTTPV2BatchSizeConditionalValidation(t *testing.T) {
	cm := configs.Destinations.Entries()["http_v2"]
	for _, tc := range []struct {
		name        string
		format      string
		batchSize   string
		expectError bool
	}{
		{"default JSON requires batch size", "", "", true},
		{"JSON requires batch size", "JSON", "", true},
		{"XML requires batch size", "XML", "", true},
		{"JSON with size", "JSON", "5000", false},
		{"XML with size", "XML", "1", false},
		{"FORM without size", "FORM", "", false},
		{"FORM with stale size", "FORM", "5000", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config := `resource "rudderstack_destination_http_v2" "example" {
				name = "test"
				config {
					api_url = "https://example.com"
					request_body_template = "{\"event\": $$.events[0]}"
				`
			if tc.format != "" {
				config += fmt.Sprintf("\nformat = %q", tc.format)
			}
			if tc.batchSize != "" {
				config += fmt.Sprintf("\nmax_batch_size = %q", tc.batchSize)
			}
			config += "\n}\n}\n"

			step := resource.TestStep{PlanOnly: true, Config: config, ExpectNonEmptyPlan: !tc.expectError}
			if tc.expectError {
				step.ExpectError = regexp.MustCompile("max_batch_size is required for JSON and XML")
			}
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: map[string]func() (*schema.Provider, error){
					"rudderstack": func() (*schema.Provider, error) {
						return &schema.Provider{
							ConfigureContextFunc: func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) { return nil, nil },
							ResourcesMap:         map[string]*schema.Resource{"rudderstack_destination_http_v2": resourceDestination(cm)},
						}, nil
					},
				},
				Steps: []resource.TestStep{step},
			})
		})
	}
}

func TestHTTPDestinationReadRejectsVersionMismatch(t *testing.T) {
	for _, tc := range []struct {
		name            string
		resourceVersion int
		remoteVersion   int64
	}{
		{"v1 resource after migration", 1, 2},
		{"v2 resource pointing at v1", 2, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cm := configs.Destinations.Entries()["http"]
			if tc.resourceVersion == 2 {
				cm = configs.Destinations.Entries()["http_v2"]
			}
			d := schema.TestResourceDataRaw(t, resourceDestinationSchema(cm), map[string]interface{}{"name": "original"})
			d.SetId("existing-id")
			service := &versionedDestinationService{destination: &client.Destination{
				ID: "existing-id", Type: "HTTP", Version: tc.remoteVersion, Name: "migrated",
				Config: json.RawMessage(`{"apiUrl":"https://example.com","requestBodyTemplate":"{}","format":"FORM"}`),
			}}

			diagnostics := resourceDestinationRead(cm)(context.Background(), d, &Client{Destinations: service})
			require.True(t, diagnostics.HasError())
			assert.Contains(t, diagnostics[0].Summary, fmt.Sprintf("version %d, but this resource manages version %d", tc.remoteVersion, tc.resourceVersion))
			assert.Equal(t, "existing-id", d.Id(), "a mismatch must not schedule a destroy/recreate")
			assert.Equal(t, "original", d.Get("name"), "do not read v2 config into the v1 schema")
		})
	}
}

type versionedDestinationService struct {
	DestinationsService
	destination *client.Destination
}

func (s *versionedDestinationService) Get(_ context.Context, _ string) (*client.Destination, error) {
	return s.destination, nil
}

func TestHTTPV2BatchSizeUnknownFormat(t *testing.T) {
	cm := configs.Destinations.Entries()["http_v2"]
	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"rudderstack": func() (*schema.Provider, error) {
				return &schema.Provider{
					ConfigureContextFunc: func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) { return nil, nil },
					ResourcesMap:         map[string]*schema.Resource{"rudderstack_destination_http_v2": resourceDestination(cm)},
				}, nil
			},
		},
		Steps: []resource.TestStep{{
			PlanOnly: true,
			Config: `resource "rudderstack_destination_http_v2" "dependency" {
				name = "dependency"
				config {
					api_url = "https://example.com"
					request_body_template = "{}"
					max_batch_size = "1"
				}
			}
			resource "rudderstack_destination_http_v2" "example" {
				name = "example"
				config {
					api_url = "https://example.com"
					request_body_template = "{}"
					format = rudderstack_destination_http_v2.dependency.id
				}
			}`,
			ExpectNonEmptyPlan: true,
		}},
	})
}

func TestPopulateDestinationFromState_SetsVersionFromConfigMeta(t *testing.T) {
	cm := configs.ConfigMeta{
		APIType:    "TEST",
		Version:    2,
		SkipConfig: true,
	}

	d := schema.TestResourceDataRaw(t, resourceDestinationSchema(cm), map[string]interface{}{
		"name":    "test-destination",
		"enabled": true,
	})

	destination := &client.Destination{}
	require.NoError(t, populateDestinationFromState(cm, destination, d))

	assert.Equal(t, "TEST", destination.Type)
	assert.Equal(t, int64(2), destination.Version)
}

func TestPopulateDestinationFromState_NoVersionSetsZero(t *testing.T) {
	// A destination registered without an explicit Version can't actually
	// exist in the real Destinations registry (Register panics on Version==0,
	// see rudderstack/configs/registries.go), but populateDestinationFromState
	// itself has no opinion on that — it just forwards cm.Version verbatim.
	// This documents that behavior for a directly-constructed ConfigMeta, as
	// used by other tests in this file (e.g. consent-management tests above).
	cm := configs.ConfigMeta{
		APIType:    "TEST",
		SkipConfig: true,
	}

	d := schema.TestResourceDataRaw(t, resourceDestinationSchema(cm), map[string]interface{}{
		"name":    "test-destination",
		"enabled": true,
	})

	destination := &client.Destination{}
	require.NoError(t, populateDestinationFromState(cm, destination, d))

	assert.Equal(t, int64(0), destination.Version)
}

// The backend redacts Sensitive fields from responses. storeDestinationToState
// must keep the secret already in state rather than clear it, else the field
// shows perpetual plan drift.
func TestStoreDestinationToState_PreservesRedactedSecret(t *testing.T) {
	cm := configs.ConfigMeta{
		APIType: "TEST",
		ConfigSchema: map[string]*schema.Schema{
			"api_key":    {Type: schema.TypeString, Optional: true},
			"api_secret": {Type: schema.TypeString, Optional: true, Sensitive: true},
		},
		Properties: []configs.ConfigProperty{
			configs.Simple("apiKey", "api_key"),
			configs.Simple("apiSecret", "api_secret"),
		},
	}

	d := schema.TestResourceDataRaw(t, resourceDestinationSchema(cm), map[string]interface{}{
		"name":    "test-destination",
		"enabled": true,
		"config": []interface{}{
			map[string]interface{}{"api_key": "pub", "api_secret": "shh"},
		},
	})

	// Simulate a redacted response: apiSecret absent, apiKey present.
	dest := &client.Destination{ID: "d1", Name: "test-destination", Config: json.RawMessage(`{"apiKey":"pub"}`)}
	require.NoError(t, storeDestinationToState(cm, dest, d))

	assert.Equal(t, "shh", d.Get("config.0.api_secret"), "redacted secret must be preserved from state")
	assert.Equal(t, "pub", d.Get("config.0.api_key"))
}
