package rudderstack

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/mock"

	"github.com/rudderlabs/rudder-api-go/client"
	"github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil"
	_ "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/integrations"
)

func TestRETLDestinationResource(t *testing.T) {
	destinations := &mockDestinationsServiceForRETL{}

	destinations.On("Create", mock.Anything, &client.Destination{
		Name:      "my-retl-dest",
		Type:      "google_ads",
		IsEnabled: true,
		Config:    json.RawMessage(`{"customerId":"123","loginCustomerId":"456"}`),
	}).Return(&client.Destination{
		ID:        "dest-id",
		Name:      "my-retl-dest",
		Type:      "google_ads",
		IsEnabled: true,
		Config:    json.RawMessage(`{"customerId":"123","loginCustomerId":"456"}`),
		CreatedAt: testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt: testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
	}, nil)

	destinations.On("Update", mock.Anything, &client.Destination{
		ID:        "dest-id",
		Name:      "my-retl-dest-updated",
		Type:      "google_ads",
		IsEnabled: true,
		Config:    json.RawMessage(`{"customerId":"789","loginCustomerId":"012"}`),
	}).Return(&client.Destination{
		ID:        "dest-id",
		Name:      "my-retl-dest-updated",
		Type:      "google_ads",
		IsEnabled: true,
		Config:    json.RawMessage(`{"customerId":"789","loginCustomerId":"012"}`),
		CreatedAt: testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt: testutil.TimePtr(time.Date(2010, 2, 2, 3, 4, 5, 0, time.UTC)),
	}, nil)

	destinations.On("Get", mock.Anything, "dest-id").Return(&client.Destination{
		ID:        "dest-id",
		Name:      "my-retl-dest",
		Type:      "google_ads",
		IsEnabled: true,
		Config:    json.RawMessage(`{"customerId":"123","loginCustomerId":"456"}`),
		CreatedAt: testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt: testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
	}, nil).Times(3)

	destinations.On("Get", mock.Anything, "dest-id").Return(&client.Destination{
		ID:        "dest-id",
		Name:      "my-retl-dest-updated",
		Type:      "google_ads",
		IsEnabled: true,
		Config:    json.RawMessage(`{"customerId":"789","loginCustomerId":"012"}`),
		CreatedAt: testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt: testutil.TimePtr(time.Date(2010, 2, 2, 3, 4, 5, 0, time.UTC)),
	}, nil).Twice()

	destinations.On("Delete", mock.Anything, "dest-id").Return(nil)

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"rudderstack": func() (*schema.Provider, error) {
				return NewWithConfigureClientFunc(func(_ context.Context, d *schema.ResourceData) (*Client, diag.Diagnostics) {
					return &Client{
						Destinations: destinations,
					}, diag.Diagnostics{}
				}), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: `
					provider "rudderstack" {
						access_token = "some-access-token"
						workspace_id = "some-workspace-id"
					}

					resource "rudderstack_retl_destination" "example" {
						name    = "my-retl-dest"
						type    = "google_ads"
						config  = "{\"customerId\":\"123\",\"loginCustomerId\":\"456\"}"
					}
				`,
				Check: func(state *terraform.State) error {
					resources := state.RootModule().Resources
					resource, ok := resources["rudderstack_retl_destination.example"]
					if !ok {
						return fmt.Errorf("resource not found in state")
					}
					attributes := resource.Primary.Attributes
					if c, ok := attributes["created_at"]; !ok || c != "2010-01-02T03:04:05Z" {
						return fmt.Errorf("created_at was not set properly in state")
					}
					if c, ok := attributes["updated_at"]; !ok || c != "2010-01-02T03:04:05Z" {
						return fmt.Errorf("updated_at was not set properly in state")
					}
					return nil
				},
			},
			{
				Config: `
					provider "rudderstack" {
						access_token = "some-access-token"
						workspace_id = "some-workspace-id"
					}

					resource "rudderstack_retl_destination" "example" {
						name    = "my-retl-dest-updated"
						type    = "google_ads"
						config  = "{\"customerId\":\"789\",\"loginCustomerId\":\"012\"}"
					}
				`,
				Check: func(state *terraform.State) error {
					resources := state.RootModule().Resources
					resource, ok := resources["rudderstack_retl_destination.example"]
					if !ok {
						return fmt.Errorf("resource not found in state")
					}
					attributes := resource.Primary.Attributes
					if c, ok := attributes["created_at"]; !ok || c != "2010-01-02T03:04:05Z" {
						return fmt.Errorf("created_at was not set properly in state")
					}
					if c, ok := attributes["updated_at"]; !ok || c != "2010-02-02T03:04:05Z" {
						return fmt.Errorf("updated_at was not set properly in state")
					}
					return nil
				},
			},
		},
	})
}

type mockDestinationsServiceForRETL struct {
	mock.Mock
}

func (m *mockDestinationsServiceForRETL) Get(ctx context.Context, id string) (*client.Destination, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*client.Destination), args.Error(1)
}

func (m *mockDestinationsServiceForRETL) Create(ctx context.Context, destination *client.Destination) (*client.Destination, error) {
	args := m.Called(ctx, destination)
	return args.Get(0).(*client.Destination), args.Error(1)
}

func (m *mockDestinationsServiceForRETL) Update(ctx context.Context, destination *client.Destination) (*client.Destination, error) {
	args := m.Called(ctx, destination)
	return args.Get(0).(*client.Destination), args.Error(1)
}

func (m *mockDestinationsServiceForRETL) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
