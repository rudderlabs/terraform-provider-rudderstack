package rudderstack

import (
	"context"
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

func TestConnectionResource(t *testing.T) {
	connections := &mockConnectionsService{}

	connections.On("Create", mock.Anything, &client.Connection{
		SourceID:      "source-id",
		DestinationID: "destination-id",
		IsEnabled:     true,
	}).Return(&client.Connection{
		ID:            "some-id",
		SourceID:      "source-id",
		DestinationID: "destination-id",
		CreatedAt:     testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt:     testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
	}, nil)

	connections.On("Update", mock.Anything, &client.Connection{
		ID:            "some-id",
		SourceID:      "source-id-2",
		DestinationID: "destination-id-2",
		IsEnabled:     true,
	}).Return(&client.Connection{
		ID:            "some-id",
		SourceID:      "source-id-2",
		DestinationID: "destination-id-2",
		CreatedAt:     testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt:     testutil.TimePtr(time.Date(2010, 2, 2, 3, 4, 5, 0, time.UTC)),
	}, nil)

	connections.On("Get", mock.Anything, "some-id").Return(&client.Connection{
		ID:            "some-id",
		SourceID:      "source-id",
		DestinationID: "destination-id",
		IsEnabled:     true,
		CreatedAt:     testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt:     testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
	}, nil).Times(3)

	connections.On("Get", mock.Anything, "some-id").Return(&client.Connection{
		ID:            "some-id",
		SourceID:      "source-id-2",
		DestinationID: "destination-id-2",
		IsEnabled:     true,
		CreatedAt:     testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt:     testutil.TimePtr(time.Date(2010, 2, 2, 3, 4, 5, 0, time.UTC)),
	}, nil).Twice()

	connections.On("Delete", mock.Anything, "some-id").Return(nil)

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"rudderstack": func() (*schema.Provider, error) {
				return NewWithConfigureClientFunc(func(_ context.Context, d *schema.ResourceData) (*Client, diag.Diagnostics) {
					return &Client{
						Connections: connections,
					}, diag.Diagnostics{}
				}), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: `
					provider "rudderstack" {
						access_token = "some-access-token"
					}

					resource "rudderstack_connection" "example" {
						source_id      = "source-id"
						destination_id = "destination-id"
					}
				`,
				Check: func(state *terraform.State) error {
					resources := state.RootModule().Resources
					resource, ok := resources["rudderstack_connection.example"]
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
					}

					resource "rudderstack_connection" "example" {
						source_id      = "source-id-2"
						destination_id = "destination-id-2"
					}
				`,
				Check: func(state *terraform.State) error {
					resources := state.RootModule().Resources
					resource, ok := resources["rudderstack_connection.example"]
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

type mockConnectionsService struct {
	mock.Mock
}

func (m *mockConnectionsService) Get(ctx context.Context, id string) (*client.Connection, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*client.Connection), args.Error(1)
}

func (m *mockConnectionsService) Create(ctx context.Context, connection *client.Connection) (*client.Connection, error) {
	args := m.Called(ctx, connection)
	return args.Get(0).(*client.Connection), args.Error(1)
}

func (m *mockConnectionsService) Update(ctx context.Context, connection *client.Connection) (*client.Connection, error) {
	args := m.Called(ctx, connection)
	return args.Get(0).(*client.Connection), args.Error(1)
}

func (m *mockConnectionsService) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
