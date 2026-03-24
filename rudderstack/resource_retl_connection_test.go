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

func TestRETLConnectionResource(t *testing.T) {
	retlConnections := &mockRETLConnectionsService{}

	createConfig := json.RawMessage(`{"source":{"event":"Purchase","eventType":"track","schedule":{"every":30,"type":"basic"},"syncBehaviour":"upsert"}}`)
	createReturnConfig := json.RawMessage(`{"source":{"event":"Purchase","eventType":"track","schedule":{"every":30,"type":"basic"},"syncBehaviour":"upsert"}}`)
	updateConfig := json.RawMessage(`{"source":{"event":"Login","eventType":"identify","schedule":{"cron":"0 */2 * * *","type":"cron"},"syncBehaviour":"mirror"}}`)
	updateReturnConfig := json.RawMessage(`{"source":{"event":"Login","eventType":"identify","schedule":{"cron":"0 */2 * * *","type":"cron"},"syncBehaviour":"mirror"}}`)

	retlConnections.On("Create", mock.Anything, &client.RETLConnection{
		SourceID:      "retl-src-id",
		DestinationID: "dest-id",
		IsEnabled:     true,
		Config:        createConfig,
	}).Return(&client.RETLConnection{
		ID:            "retl-conn-id",
		SourceID:      "retl-src-id",
		DestinationID: "dest-id",
		IsEnabled:     true,
		Config:        createReturnConfig,
		CreatedAt:     testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt:     testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
	}, nil)

	retlConnections.On("Update", mock.Anything, &client.RETLConnection{
		ID:            "retl-conn-id",
		SourceID:      "retl-src-id",
		DestinationID: "dest-id",
		IsEnabled:     true,
		Config:        updateConfig,
	}).Return(&client.RETLConnection{
		ID:            "retl-conn-id",
		SourceID:      "retl-src-id",
		DestinationID: "dest-id",
		IsEnabled:     true,
		Config:        updateReturnConfig,
		CreatedAt:     testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt:     testutil.TimePtr(time.Date(2010, 2, 2, 3, 4, 5, 0, time.UTC)),
	}, nil)

	retlConnections.On("Get", mock.Anything, "retl-conn-id").Return(&client.RETLConnection{
		ID:            "retl-conn-id",
		SourceID:      "retl-src-id",
		DestinationID: "dest-id",
		IsEnabled:     true,
		Config:        createReturnConfig,
		CreatedAt:     testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt:     testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
	}, nil).Times(3)

	retlConnections.On("Get", mock.Anything, "retl-conn-id").Return(&client.RETLConnection{
		ID:            "retl-conn-id",
		SourceID:      "retl-src-id",
		DestinationID: "dest-id",
		IsEnabled:     true,
		Config:        updateReturnConfig,
		CreatedAt:     testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt:     testutil.TimePtr(time.Date(2010, 2, 2, 3, 4, 5, 0, time.UTC)),
	}, nil).Twice()

	retlConnections.On("Delete", mock.Anything, "retl-conn-id").Return(nil)

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"rudderstack": func() (*schema.Provider, error) {
				return NewWithConfigureClientFunc(func(_ context.Context, d *schema.ResourceData) (*Client, diag.Diagnostics) {
					return &Client{
						RETLConnections: retlConnections,
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

					resource "rudderstack_retl_connection" "example" {
						source_id      = "retl-src-id"
						destination_id = "dest-id"

						schedule {
							type  = "basic"
							every = 30
						}

						sync_behaviour = "upsert"
						event          = "Purchase"
						event_type     = "track"
					}
				`,
				Check: func(state *terraform.State) error {
					resources := state.RootModule().Resources
					resource, ok := resources["rudderstack_retl_connection.example"]
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

					resource "rudderstack_retl_connection" "example" {
						source_id      = "retl-src-id"
						destination_id = "dest-id"

						schedule {
							type = "cron"
							cron = "0 */2 * * *"
						}

						sync_behaviour = "mirror"
						event          = "Login"
						event_type     = "identify"
					}
				`,
				Check: func(state *terraform.State) error {
					resources := state.RootModule().Resources
					resource, ok := resources["rudderstack_retl_connection.example"]
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

type mockRETLConnectionsService struct {
	mock.Mock
}

func (m *mockRETLConnectionsService) Get(ctx context.Context, id string) (*client.RETLConnection, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*client.RETLConnection), args.Error(1)
}

func (m *mockRETLConnectionsService) Create(ctx context.Context, connection *client.RETLConnection) (*client.RETLConnection, error) {
	args := m.Called(ctx, connection)
	return args.Get(0).(*client.RETLConnection), args.Error(1)
}

func (m *mockRETLConnectionsService) Update(ctx context.Context, connection *client.RETLConnection) (*client.RETLConnection, error) {
	args := m.Called(ctx, connection)
	return args.Get(0).(*client.RETLConnection), args.Error(1)
}

func (m *mockRETLConnectionsService) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
