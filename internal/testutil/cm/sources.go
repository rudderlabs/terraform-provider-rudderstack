package cm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/mock"

	"github.com/rudderlabs/rudder-iac/api/client"
	"github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

// sourceSettingsFromAPIJSON parses the APICreate/APIUpdate JSON and returns
// GeoEnrichmentEnabled and the API-struct Transient value for mock Source structs.
// The JSON "transient" key holds the actual API wire value (source.Transient),
// which is the inverse of the TF field temporarily_store_events_for_retries.
// Returns nil pointers when the JSON is empty or has no settings keys.
func sourceSettingsFromAPIJSON(apiJSON string) (geoEnabled *bool, apiTransient *bool) {
	apiJSON = strings.TrimSpace(apiJSON)
	if apiJSON == "" || apiJSON == "{}" {
		return nil, nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(apiJSON), &m); err != nil {
		return nil, nil
	}
	_, hasGeo := m["geoEnrichmentEnabled"]
	_, hasTrans := m["transient"]
	if !hasGeo && !hasTrans {
		return nil, nil
	}
	var geo, trans bool
	if v, ok := m["geoEnrichmentEnabled"].(bool); ok {
		geo = v
	}
	if v, ok := m["transient"].(bool); ok {
		// transient in JSON = actual source.Transient wire value (inverse of temporarily_store_events_for_retries)
		trans = v
	}
	return &geo, &trans
}

// cmSourceHCL builds the HCL for a source resource. When configBlock is empty
// the resource has no extra blocks. Otherwise configBlock is embedded raw (the
// caller is responsible for passing a well-formed block such as "settings { … }").
func cmSourceHCL(source, name, configBlock string) string {
	provider := `
		provider "rudderstack" {
			access_token = "some-access-token"
		}`
	if configBlock == "" {
		return fmt.Sprintf(`%s

		resource "rudderstack_source_%s" "example" {
			name = %q
		}
		`, provider, source, name)
	}
	return fmt.Sprintf(`%s

	resource "rudderstack_source_%s" "example" {
		name = %q
		%s
	}
	`, provider, source, name, configBlock)
}

func AssertSource(t *testing.T, source string, testConfigs []configs.TestConfig) {
	cm := configs.Sources.Entries()[source]
	cfg := testConfigs[0]

	createSettingsJSON := cfg.APICreateSettings
	if createSettingsJSON == "" {
		createSettingsJSON = cfg.APICreate
	}
	updateSettingsJSON := cfg.APIUpdateSettings
	if updateSettingsJSON == "" {
		updateSettingsJSON = cfg.APIUpdate
	}
	createGeo, createTransientAPI := sourceSettingsFromAPIJSON(createSettingsJSON)
	updateGeo, updateTransientAPI := sourceSettingsFromAPIJSON(updateSettingsJSON)

	sources := &mockSourcesService{}

	sources.On("Create", mock.Anything, &client.Source{
		Type:                 cm.APIType,
		Name:                 "example",
		IsEnabled:            true,
		GeoEnrichmentEnabled: createGeo,
		Transient:            createTransientAPI,
	}).Return(&client.Source{
		ID:                   "some-id",
		Type:                 cm.APIType,
		Name:                 "example",
		IsEnabled:            true,
		WriteKey:             "some-write-key",
		GeoEnrichmentEnabled: createGeo,
		Transient:            createTransientAPI,
		CreatedAt:            testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt:            testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
	}, nil)

	updateName := "example-updated"
	if cfg.TerraformUpdate == cfg.TerraformCreate {
		updateName = "example-updated"
	}

	sources.On("Update", mock.Anything, &client.Source{
		ID:                   "some-id",
		Type:                 cm.APIType,
		Name:                 updateName,
		IsEnabled:            true,
		GeoEnrichmentEnabled: updateGeo,
		Transient:            updateTransientAPI,
	}).Return(&client.Source{
		ID:                   "some-id",
		Type:                 cm.APIType,
		Name:                 updateName,
		IsEnabled:            true,
		GeoEnrichmentEnabled: updateGeo,
		Transient:            updateTransientAPI,
		CreatedAt:            testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt:            testutil.TimePtr(time.Date(2010, 2, 2, 3, 4, 5, 0, time.UTC)),
	}, nil)

	sources.On("Get", mock.Anything, "some-id").Return(&client.Source{
		ID:                   "some-id",
		Type:                 cm.APIType,
		Name:                 "example",
		IsEnabled:            true,
		WriteKey:             "some-write-key",
		GeoEnrichmentEnabled: createGeo,
		Transient:            createTransientAPI,
		CreatedAt:            testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt:            testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
	}, nil).Times(3)

	sources.On("Get", mock.Anything, "some-id").Return(&client.Source{
		ID:                   "some-id",
		Type:                 cm.APIType,
		Name:                 updateName,
		IsEnabled:            true,
		WriteKey:             "some-write-key",
		GeoEnrichmentEnabled: updateGeo,
		Transient:            updateTransientAPI,
		CreatedAt:            testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt:            testutil.TimePtr(time.Date(2010, 2, 2, 3, 4, 5, 0, time.UTC)),
	}, nil).Twice()

	sources.On("Delete", mock.Anything, "some-id").Return(nil)

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"rudderstack": func() (*schema.Provider, error) {
				return rudderstack.NewWithConfigureClientFunc(func(_ context.Context, d *schema.ResourceData) (*rudderstack.Client, diag.Diagnostics) {
					return &rudderstack.Client{
						Sources: sources,
					}, diag.Diagnostics{}
				}), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: cmSourceHCL(source, "example", cfg.TerraformCreate),
				Check: func(state *terraform.State) error {
					resources := state.RootModule().Resources
					resource, ok := resources[fmt.Sprintf("rudderstack_source_%s.example", source)]
					if !ok {
						return fmt.Errorf("resource not found in state")
					}
					attributes := resource.Primary.Attributes
					if c, ok := attributes["created_at"]; !ok || c != "2010-01-02T03:04:05Z" {
						return fmt.Errorf("created_at was not set properly in state")
					}
					if c, ok := attributes["updated_at"]; !ok || c != "2010-01-02T03:04:05Z" {
						return fmt.Errorf("update_at was not set properly in state")
					}
					if c, ok := attributes["write_key"]; !ok || c != "some-write-key" {
						return fmt.Errorf("write_key was not set properly in state")
					}
					return nil
				},
			},
			{
				Config: cmSourceHCL(source, updateName, cfg.TerraformUpdate),
				Check: func(state *terraform.State) error {
					resources := state.RootModule().Resources
					resource, ok := resources[fmt.Sprintf("rudderstack_source_%s.example", source)]
					if !ok {
						return fmt.Errorf("resource not found in state")
					}
					attributes := resource.Primary.Attributes
					if c, ok := attributes["created_at"]; !ok || c != "2010-01-02T03:04:05Z" {
						return fmt.Errorf("created_at was not set properly in state")
					}
					if c, ok := attributes["updated_at"]; !ok || c != "2010-02-02T03:04:05Z" {
						return fmt.Errorf("update_at was not set properly in state")
					}
					return nil
				},
			},
		},
	})
}

type mockSourcesService struct {
	mock.Mock
}

func (m *mockSourcesService) Get(ctx context.Context, id string) (*client.Source, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*client.Source), args.Error(1)
}

func (m *mockSourcesService) Create(ctx context.Context, source *client.Source) (*client.Source, error) {
	args := m.Called(ctx, source)
	return args.Get(0).(*client.Source), args.Error(1)
}

func (m *mockSourcesService) Update(ctx context.Context, source *client.Source) (*client.Source, error) {
	args := m.Called(ctx, source)
	return args.Get(0).(*client.Source), args.Error(1)
}

func (m *mockSourcesService) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
