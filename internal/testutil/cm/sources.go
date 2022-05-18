package cm

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/rudderlabs/rudder-api-go/client"
	"github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
	"github.com/stretchr/testify/mock"
)

func AssertSource(t *testing.T, source string, testConfigs []configs.TestConfig) {
	cm := configs.Sources.Entries()[source]
	sources := &mockSourcesService{}

	sources.On("Create", mock.Anything, &client.Source{
		Type:      cm.APIType,
		Name:      "example",
		IsEnabled: true,
	}).Return(&client.Source{
		ID:        "some-id",
		Type:      cm.APIType,
		Name:      "example",
		IsEnabled: true,
		WriteKey:  "some-write-key",
		CreatedAt: testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt: testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
	}, nil)

	sources.On("Update", mock.Anything, &client.Source{
		ID:        "some-id",
		Type:      cm.APIType,
		Name:      "example-updated",
		IsEnabled: true,
	}).Return(&client.Source{
		ID:        "some-id",
		Type:      cm.APIType,
		Name:      "example-updated",
		IsEnabled: true,
		CreatedAt: testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt: testutil.TimePtr(time.Date(2010, 2, 2, 3, 4, 5, 0, time.UTC)),
	}, nil)

	sources.On("Get", mock.Anything, "some-id").Return(&client.Source{
		ID:        "some-id",
		Type:      cm.APIType,
		Name:      "example",
		IsEnabled: true,
		WriteKey:  "some-write-key",
		CreatedAt: testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt: testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
	}, nil).Times(3)

	sources.On("Get", mock.Anything, "some-id").Return(&client.Source{
		ID:        "some-id",
		Type:      cm.APIType,
		Name:      "example-updated",
		IsEnabled: true,
		WriteKey:  "some-write-key",
		CreatedAt: testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt: testutil.TimePtr(time.Date(2010, 2, 2, 3, 4, 5, 0, time.UTC)),
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
				Config: fmt.Sprintf(`
					provider "rudderstack" {
						access_token = "some-access-token"
					}

					resource "rudderstack_source_%s" "example" {
						name = "example"
					}
				`, source),
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
				Config: fmt.Sprintf(`
					provider "rudderstack" {
						access_token = "some-access-token"
					}

					resource "rudderstack_source_%s" "example" {
						name = "example-updated"
					}
				`, source),
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
