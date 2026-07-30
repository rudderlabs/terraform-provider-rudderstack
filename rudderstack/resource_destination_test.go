package rudderstack

import (
	"context"
	"encoding/json"
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
	assert.Equal(t, 2, destination.Version)
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

	assert.Equal(t, 0, destination.Version)
}

func TestStoreDestinationToStatePreservesPrunedSensitiveAndEmptyConfigValues(t *testing.T) {
	cm := configs.ConfigMeta{
		APIType: "TEST",
		ConfigSchema: map[string]*schema.Schema{
			"api_secret": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"labels": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"headers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
					},
				},
			},
			"key_based_authentication": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_key_id": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"access_key": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
					},
				},
			},
			"non_secret": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Properties: []configs.ConfigProperty{
			configs.Simple("apiSecret", "api_secret"),
			configs.Simple("url", "url"),
			configs.ArrayWithStrings("labels", "value", "labels"),
			configs.ArrayWithObjects("headers", "headers", map[string]any{
				"name":  "name",
				"value": "value",
			}),
			configs.ArrayWithObjects("keyBasedAuthentication", "key_based_authentication", map[string]any{
				"accessKeyID": "access_key_id",
				"accessKey":   "access_key",
			}),
			configs.Simple("nonSecret", "non_secret"),
		},
	}

	d := schema.TestResourceDataRaw(t, resourceDestinationSchema(cm), map[string]interface{}{
		"name":    "test-destination",
		"enabled": true,
		"config": []interface{}{
			map[string]interface{}{
				"api_secret": "prior-api-secret",
				"url":        "https://old.example.com",
				"labels":     []interface{}{},
				"headers": []interface{}{
					map[string]interface{}{
						"name":  "Authorization",
						"value": "prior-header-secret",
					},
				},
				"key_based_authentication": []interface{}{
					map[string]interface{}{
						"access_key_id": "prior-access-key-id",
						"access_key":    "prior-access-key",
					},
				},
				"non_secret": "prior-non-secret",
			},
		},
	})

	destination := &client.Destination{
		ID:        "destination-id",
		Name:      "test-destination",
		IsEnabled: true,
		Config: json.RawMessage(`{
			"url": "https://new.example.com",
			"nonSecret": "api-non-secret",
			"headers": [
				{
					"name": "Authorization",
					"value": ""
				}
			]
		}`),
	}

	require.NoError(t, storeDestinationToState(cm, destination, d))

	assert.Equal(t, "prior-api-secret", d.Get("config.0.api_secret"))
	assert.Equal(t, "https://new.example.com", d.Get("config.0.url"))
	assert.Equal(t, 0, d.Get("config.0.labels.#"))
	assert.Equal(t, "Authorization", d.Get("config.0.headers.0.name"))
	assert.Equal(t, "prior-header-secret", d.Get("config.0.headers.0.value"))
	assert.Equal(t, "prior-access-key-id", d.Get("config.0.key_based_authentication.0.access_key_id"))
	assert.Equal(t, "prior-access-key", d.Get("config.0.key_based_authentication.0.access_key"))
	assert.Equal(t, "api-non-secret", d.Get("config.0.non_secret"))
}
