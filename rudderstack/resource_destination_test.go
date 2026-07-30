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

func TestResourceDestinationReadFiltersStaleConfigFields(t *testing.T) {
	cm := testDestinationStaleFieldConfigMeta()
	destinations := &stubDestinationsService{
		getFunc: func(_ context.Context, id string) (*client.Destination, error) {
			assert.Equal(t, "dst-id", id)
			return &client.Destination{
				ID:        "dst-id",
				Type:      cm.APIType,
				Version:   cm.Version,
				Name:      "example",
				IsEnabled: true,
				Config: json.RawMessage(`{
					"known": "kept",
					"staleTopLevel": "dropped",
					"nested": { "value": "kept-nested", "staleNested": "dropped" },
					"objects": [{
						"known": "kept-object",
						"staleObject": "dropped",
						"nestedValues": [
							{ "value": "kept-list-value", "staleListValue": "dropped" }
						]
					}]
				}`),
			}, nil
		},
	}
	d := schema.TestResourceDataRaw(t, resourceDestinationSchema(cm), map[string]interface{}{
		"name":    "before-read",
		"enabled": true,
	})
	d.SetId("dst-id")

	diags := resourceDestinationRead(cm)(context.Background(), d, &Client{Destinations: destinations})
	require.False(t, diags.HasError(), "unexpected diagnostics: %#v", diags)

	assert.Equal(t, "example", d.Get("name"))
	assert.Equal(t, "kept", d.Get("config.0.known"))
	assert.Equal(t, "kept-nested", d.Get("config.0.nested.0.value"))
	objects, ok := d.Get("config.0.objects").([]interface{})
	require.True(t, ok)
	require.Len(t, objects, 1)
	object, ok := objects[0].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "kept-object", object["known"])
	assert.Equal(t, []interface{}{"kept-list-value"}, object["nested_values"])
	assert.NotContains(t, object, "stale_object")

	attributes := d.State().Attributes
	assert.NotContains(t, attributes, "config.0.stale_top_level")
	assert.NotContains(t, attributes, "config.0.nested.0.stale_nested")
	assert.NotContains(t, attributes, "config.0.objects.0.stale_object")
}

func TestResourceDestinationReadPreservesWriteOnlyConfigFields(t *testing.T) {
	cm := testDestinationWriteOnlyConfigMeta()
	destinations := &stubDestinationsService{
		getFunc: func(_ context.Context, id string) (*client.Destination, error) {
			assert.Equal(t, "dst-id", id)
			return &client.Destination{
				ID:        "dst-id",
				Type:      cm.APIType,
				Version:   cm.Version,
				Name:      "example",
				IsEnabled: true,
				Config: json.RawMessage(`{
					"apiUrl": "https://example.com",
					"headers": [{ "from": "x-header", "to": "" }]
				}`),
			}, nil
		},
	}
	d := schema.TestResourceDataRaw(t, resourceDestinationSchema(cm), map[string]interface{}{
		"name":    "before-read",
		"enabled": true,
		"config": []interface{}{
			map[string]interface{}{
				"api_key": "existing-api-key",
				"api_url": "https://old.example.com",
				"headers": []interface{}{
					map[string]interface{}{"from": "x-header", "to": "existing-header-value"},
				},
				"stale": "should-not-survive",
			},
		},
	})
	d.SetId("dst-id")

	diags := resourceDestinationRead(cm)(context.Background(), d, &Client{Destinations: destinations})
	require.False(t, diags.HasError(), "unexpected diagnostics: %#v", diags)

	assert.Equal(t, "https://example.com", d.Get("config.0.api_url"))
	assert.Equal(t, "existing-api-key", d.Get("config.0.api_key"))
	headers, ok := d.Get("config.0.headers").([]interface{})
	require.True(t, ok)
	require.Len(t, headers, 1)
	header, ok := headers[0].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "x-header", header["from"])
	assert.Equal(t, "existing-header-value", header["to"])
	assert.Empty(t, d.Get("config.0.stale"))
}

func TestResourceDestinationUpdateDoesNotRoundTripStaleConfigFields(t *testing.T) {
	cm := testDestinationStaleFieldConfigMeta()
	getCalls := 0
	destinations := &stubDestinationsService{
		getFunc: func(_ context.Context, id string) (*client.Destination, error) {
			assert.Equal(t, "dst-id", id)
			getCalls++
			name := "example"
			if getCalls > 1 {
				name = "example-updated"
			}
			return &client.Destination{
				ID:        "dst-id",
				Type:      cm.APIType,
				Version:   cm.Version,
				Name:      name,
				IsEnabled: true,
				Config: json.RawMessage(`{
					"known": "kept",
					"staleTopLevel": "dropped",
					"nested": { "value": "kept-nested", "staleNested": "dropped" },
					"objects": [{
						"known": "kept-object",
						"staleObject": "dropped",
						"nestedValues": [{ "value": "kept-list-value", "staleListValue": "dropped" }]
					}]
				}`),
			}, nil
		},
		updateFunc: func(_ context.Context, destination *client.Destination) (*client.Destination, error) {
			assert.Equal(t, "dst-id", destination.ID)
			assert.Equal(t, cm.APIType, destination.Type)
			assert.Equal(t, cm.Version, destination.Version)
			assert.Equal(t, "example-updated", destination.Name)
			assert.JSONEq(t, `{
				"known": "kept",
				"nested": { "value": "kept-nested" },
				"objects": [{
					"known": "kept-object",
					"nestedValues": [{ "value": "kept-list-value" }]
				}]
			}`, string(destination.Config))
			return &client.Destination{
				ID:        "dst-id",
				Type:      cm.APIType,
				Version:   cm.Version,
				Name:      "example-updated",
				IsEnabled: true,
				Config:    destination.Config,
			}, nil
		},
	}
	d := schema.TestResourceDataRaw(t, resourceDestinationSchema(cm), map[string]interface{}{
		"name":    "before-read",
		"enabled": true,
	})
	d.SetId("dst-id")

	diags := resourceDestinationRead(cm)(context.Background(), d, &Client{Destinations: destinations})
	require.False(t, diags.HasError(), "unexpected diagnostics on read: %#v", diags)
	require.NoError(t, d.Set("name", "example-updated"))

	diags = resourceDestinationUpdate(cm)(context.Background(), d, &Client{Destinations: destinations})
	require.False(t, diags.HasError(), "unexpected diagnostics on update: %#v", diags)
	assert.Equal(t, 2, getCalls)
}

func testDestinationStaleFieldConfigMeta() configs.ConfigMeta {
	return configs.ConfigMeta{
		APIType: "TEST_DESTINATION",
		Version: 1,
		ConfigSchema: map[string]*schema.Schema{
			"known": {
				Type:     schema.TypeString,
				Required: true,
			},
			"nested": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"objects": {
				Type:       schema.TypeList,
				Optional:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"known": {
							Type:     schema.TypeString,
							Required: true,
						},
						"nested_values": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
		Properties: []configs.ConfigProperty{
			configs.Simple("known", "known"),
			configs.Simple("nested.value", "nested.0.value"),
			configs.ArrayWithObjects("objects", "objects", map[string]interface{}{
				"known": "known",
				"nestedValues": configs.APINestedObject{
					TerraformKey: "nested_values",
					NestedKey:    "value",
				},
			}),
		},
	}
}

func testDestinationWriteOnlyConfigMeta() configs.ConfigMeta {
	return configs.ConfigMeta{
		APIType: "TEST_DESTINATION",
		Version: 1,
		ConfigSchema: map[string]*schema.Schema{
			"api_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"api_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"headers": {
				Type:      schema.TypeList,
				Optional:  true,
				Sensitive: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from": {
							Type:     schema.TypeString,
							Required: true,
						},
						"to": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
		Properties: []configs.ConfigProperty{
			configs.Simple("apiKey", "api_key"),
			configs.Simple("apiUrl", "api_url"),
			configs.Simple("headers", "headers"),
		},
	}
}

type stubDestinationsService struct {
	createFunc func(context.Context, *client.Destination) (*client.Destination, error)
	getFunc    func(context.Context, string) (*client.Destination, error)
	updateFunc func(context.Context, *client.Destination) (*client.Destination, error)
	deleteFunc func(context.Context, string) error
}

func (s *stubDestinationsService) Create(ctx context.Context, destination *client.Destination) (*client.Destination, error) {
	return s.createFunc(ctx, destination)
}

func (s *stubDestinationsService) Get(ctx context.Context, id string) (*client.Destination, error) {
	return s.getFunc(ctx, id)
}

func (s *stubDestinationsService) Update(ctx context.Context, destination *client.Destination) (*client.Destination, error) {
	return s.updateFunc(ctx, destination)
}

func (s *stubDestinationsService) Delete(ctx context.Context, id string) error {
	return s.deleteFunc(ctx, id)
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
