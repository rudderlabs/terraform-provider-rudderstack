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

// A NESTED redacted secret is preserved from state on read (top-level scalars are
// handled by DiffSuppressFunc instead — see TestTopLevelSecretUsesDiffSuppress).
func TestStoreDestinationToState_PreservesNestedRedactedSecret(t *testing.T) {
	cm := configs.ConfigMeta{
		APIType: "TEST",
		ConfigSchema: map[string]*schema.Schema{
			"s3": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"bucket":     {Type: schema.TypeString, Optional: true},
					"access_key": {Type: schema.TypeString, Optional: true, Sensitive: true},
				}},
			},
		},
		Properties: []configs.ConfigProperty{
			configs.Simple("bucketName", "s3.0.bucket"),
			configs.Simple("accessKey", "s3.0.access_key"),
		},
	}

	d := schema.TestResourceDataRaw(t, resourceDestinationSchema(cm), map[string]interface{}{
		"name":    "test-destination",
		"enabled": true,
		"config": []interface{}{map[string]interface{}{
			"s3": []interface{}{map[string]interface{}{"bucket": "b", "access_key": "shh"}},
		}},
	})

	// Redacted response: accessKey absent, bucketName present.
	dest := &client.Destination{ID: "d1", Name: "test-destination", Config: json.RawMessage(`{"bucketName":"b"}`)}
	require.NoError(t, storeDestinationToState(cm, dest, d))

	assert.Equal(t, "shh", d.Get("config.0.s3.0.access_key"), "nested redacted secret must be preserved")
	assert.Equal(t, "b", d.Get("config.0.s3.0.bucket"))
}

// A top-level scalar secret is handled by DiffSuppressFunc (empty->set suppressed,
// real change not), and excluded from the read-preserve path.
func TestTopLevelSecretUsesDiffSuppress(t *testing.T) {
	cm := configs.ConfigMeta{
		APIType: "TEST",
		ConfigSchema: map[string]*schema.Schema{
			"api_secret": {Type: schema.TypeString, Optional: true, Sensitive: true},
		},
		Properties: []configs.ConfigProperty{configs.Simple("apiSecret", "api_secret")},
	}

	resourceDestinationSchema(cm) // attaches the suppress func

	sup := cm.ConfigSchema["api_secret"].DiffSuppressFunc
	require.NotNil(t, sup, "top-level scalar secret must get a DiffSuppressFunc")
	assert.True(t, sup("config.0.api_secret", "", "shh", nil), "redacted empty->configured must be suppressed")
	assert.False(t, sup("config.0.api_secret", "old", "new", nil), "a real secret change must not be suppressed")

	assert.True(t, isDiffSuppressedSecret(cm, "api_secret"), "top-level scalar secret is owned by DiffSuppressFunc, not preserve")
}
