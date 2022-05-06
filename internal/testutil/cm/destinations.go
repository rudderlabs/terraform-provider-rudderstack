package cm

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
	"github.com/rudderlabs/rudder-api-go/client"
	"github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
	"github.com/stretchr/testify/mock"
)

func AssertDestination(t *testing.T, destination string, testConfigs []configs.TestConfig) {
	cm := configs.Destinations.Entries()[destination]
	destinations := &mockDestinationsService{}

	destinations.On("Create", mock.Anything, mock.MatchedBy(func(d *client.Destination) bool {
		return d.Type == cm.APIType &&
			d.ID == "" &&
			d.Name == "example" &&
			d.IsEnabled &&
			testutil.JSONEq(string(d.Config), testConfigs[0].APICreate)
	})).Return(&client.Destination{
		ID:        "some-id",
		Type:      cm.APIType,
		Name:      "example",
		IsEnabled: true,
		Config:    json.RawMessage(testConfigs[0].APICreate),
		CreatedAt: testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt: testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
	}, nil)

	destinations.On("Update", mock.Anything, mock.MatchedBy(func(d *client.Destination) bool {
		return d.Type == cm.APIType &&
			d.ID == "some-id" &&
			d.Name == "example-updated" &&
			d.IsEnabled &&
			testutil.JSONEq(string(d.Config), testConfigs[0].APIUpdate)
	})).Return(&client.Destination{
		ID:        "some-id",
		Type:      cm.APIType,
		Name:      "example-updated",
		IsEnabled: true,
		Config:    json.RawMessage(testConfigs[0].APIUpdate),
		CreatedAt: testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt: testutil.TimePtr(time.Date(2010, 2, 2, 3, 4, 5, 0, time.UTC)),
	}, nil)

	destinations.On("Get", mock.Anything, "some-id").Return(&client.Destination{
		ID:        "some-id",
		Type:      cm.APIType,
		Name:      "example",
		IsEnabled: true,
		Config:    json.RawMessage(testConfigs[0].APICreate),
		CreatedAt: testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt: testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
	}, nil).Times(3)

	destinations.On("Get", mock.Anything, "some-id").Return(&client.Destination{
		ID:        "some-id",
		Type:      cm.APIType,
		Name:      "example-updated",
		IsEnabled: true,
		Config:    json.RawMessage(testConfigs[0].APIUpdate),
		CreatedAt: testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt: testutil.TimePtr(time.Date(2010, 2, 2, 3, 4, 5, 0, time.UTC)),
	}, nil).Twice()

	destinations.On("Delete", mock.Anything, "some-id").Return(nil)

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"rudderstack": func() (*schema.Provider, error) {
				return rudderstack.NewWithConfigureClientFunc(func(_ context.Context, d *schema.ResourceData) (*rudderstack.Client, diag.Diagnostics) {
					return &rudderstack.Client{
						Destinations: destinations,
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

					resource "rudderstack_destination_%s" "example" {
						name = "example"
						config {
							%s
						}
					}
				`, destination, testConfigs[0].TerraformCreate),
				Check: func(state *terraform.State) error {
					resources := state.RootModule().Resources
					resource, ok := resources[fmt.Sprintf("rudderstack_destination_%s.example", destination)]
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
					return nil
				},
			},
			{
				Config: fmt.Sprintf(`
					provider "rudderstack" {
						access_token = "some-access-token"
					}

					resource "rudderstack_destination_%s" "example" {
						name = "example-updated"
						config {
							%s
						}
					}
				`, destination, testConfigs[0].TerraformUpdate),
				Check: func(state *terraform.State) error {
					resources := state.RootModule().Resources
					resource, ok := resources[fmt.Sprintf("rudderstack_destination_%s.example", destination)]
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

type mockDestinationsService struct {
	mock.Mock
}

func (m *mockDestinationsService) Get(ctx context.Context, id string) (*client.Destination, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*client.Destination), args.Error(1)
}

func (m *mockDestinationsService) Create(ctx context.Context, destination *client.Destination) (*client.Destination, error) {
	args := m.Called(ctx, destination)
	return args.Get(0).(*client.Destination), args.Error(1)
}

func (m *mockDestinationsService) Update(ctx context.Context, destination *client.Destination) (*client.Destination, error) {
	args := m.Called(ctx, destination)
	return args.Get(0).(*client.Destination), args.Error(1)
}

func (m *mockDestinationsService) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
