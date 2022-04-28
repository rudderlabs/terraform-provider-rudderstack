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
	"github.com/rudderlabs/rudder-api-go/client"
	"github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
	_ "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/integrations"
	"github.com/stretchr/testify/mock"
)

func TestDestinationResource(t *testing.T) {
	for destination, cm := range configs.Destinations.Entries() {
		destinations := &mockDestinationsService{}

		destinations.On("Create", mock.Anything, mock.MatchedBy(func(d *client.Destination) bool {
			return d.Type == cm.APIType &&
				d.ID == "" &&
				d.Name == "example" &&
				d.IsEnabled &&
				testutil.JSONEq(string(d.Config), cm.TestConfigs[0].APICreate)
		})).Return(&client.Destination{
			ID:        "some-id",
			Type:      cm.APIType,
			Name:      "example",
			IsEnabled: true,
			Config:    json.RawMessage(cm.TestConfigs[0].APICreate),
			CreatedAt: timePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
			UpdatedAt: timePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		}, nil)

		destinations.On("Update", mock.Anything, mock.MatchedBy(func(d *client.Destination) bool {
			return d.Type == cm.APIType &&
				d.ID == "some-id" &&
				d.Name == "example-updated" &&
				d.IsEnabled &&
				testutil.JSONEq(string(d.Config), cm.TestConfigs[0].APIUpdate)
		})).Return(&client.Destination{
			ID:        "some-id",
			Type:      cm.APIType,
			Name:      "example-updated",
			IsEnabled: true,
			Config:    json.RawMessage(cm.TestConfigs[0].APIUpdate),
			CreatedAt: timePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
			UpdatedAt: timePtr(time.Date(2010, 2, 2, 3, 4, 5, 0, time.UTC)),
		}, nil)

		destinations.On("Get", mock.Anything, "some-id").Return(&client.Destination{
			ID:        "some-id",
			Type:      cm.APIType,
			Name:      "example",
			IsEnabled: true,
			Config:    json.RawMessage(cm.TestConfigs[0].APICreate),
			CreatedAt: timePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
			UpdatedAt: timePtr(time.Date(2010, 2, 2, 3, 4, 5, 0, time.UTC)),
		}, nil).Times(3)

		destinations.On("Get", mock.Anything, "some-id").Return(&client.Destination{
			ID:        "some-id",
			Type:      cm.APIType,
			Name:      "example-updated",
			IsEnabled: true,
			Config:    json.RawMessage(cm.TestConfigs[0].APIUpdate),
			CreatedAt: timePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
			UpdatedAt: timePtr(time.Date(2010, 2, 2, 3, 4, 5, 0, time.UTC)),
		}, nil).Twice()

		destinations.On("Delete", mock.Anything, "some-id").Return(nil)

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
					`, destination, cm.TestConfigs[0].TerraformCreate),
					Check: func(state *terraform.State) error {
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
					`, destination, cm.TestConfigs[0].TerraformUpdate),
					Check: func(state *terraform.State) error {
						return nil
					},
				},
			},
		})
	}
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
