package rudderstack

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/mock"

	"github.com/rudderlabs/rudder-iac/api/client/transformations"
	_ "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/integrations"
)

func TestTransformationConnectionResource(t *testing.T) {
	connections := &mockTransformationConnectionsService{}

	connections.On("ConnectToDestination", mock.Anything, "transformation-id", &transformations.ConnectToDestinationRequest{
		DestinationID: "destination-id",
	}).Return(&transformations.Transformation{
		ID: "transformation-id",
		Destinations: []transformations.TransformationDestination{
			{ID: "destination-id", Name: "My destination", Enabled: true},
		},
	}, nil)

	connections.On("GetTransformation", mock.Anything, "transformation-id").Return(&transformations.Transformation{
		ID: "transformation-id",
		Destinations: []transformations.TransformationDestination{
			{ID: "destination-id", Name: "My destination", Enabled: true},
		},
	}, nil)

	connections.On("DisconnectFromDestination", mock.Anything, "transformation-id", &transformations.ConnectToDestinationRequest{
		DestinationID: "destination-id",
	}).Return(&transformations.Transformation{
		ID: "transformation-id",
	}, nil)

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"rudderstack": func() (*schema.Provider, error) {
				return NewWithConfigureClientFunc(func(_ context.Context, d *schema.ResourceData) (*Client, diag.Diagnostics) {
					return &Client{
						TransformationConnections: connections,
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

					resource "rudderstack_transformation_connection" "example" {
						transformation_id = "transformation-id"
						destination_id    = "destination-id"
					}
				`,
				Check: func(state *terraform.State) error {
					resources := state.RootModule().Resources
					res, ok := resources["rudderstack_transformation_connection.example"]
					if !ok {
						return fmt.Errorf("resource not found in state")
					}
					attributes := res.Primary.Attributes
					if id := res.Primary.ID; id != "transformation-id:destination-id" {
						return fmt.Errorf("unexpected id %q", id)
					}
					if v := attributes["transformation_id"]; v != "transformation-id" {
						return fmt.Errorf("transformation_id was not set properly in state")
					}
					if v := attributes["destination_id"]; v != "destination-id" {
						return fmt.Errorf("destination_id was not set properly in state")
					}
					return nil
				},
			},
		},
	})
}

// TestTransformationConnectionResourceDrift verifies that when the connection no
// longer exists upstream (the destination is absent from the transformation's
// destinations), the resource is removed from state on read.
func TestTransformationConnectionResourceDrift(t *testing.T) {
	connections := &mockTransformationConnectionsService{}

	connections.On("ConnectToDestination", mock.Anything, "transformation-id", &transformations.ConnectToDestinationRequest{
		DestinationID: "destination-id",
	}).Return(&transformations.Transformation{
		ID: "transformation-id",
		Destinations: []transformations.TransformationDestination{
			{ID: "destination-id", Name: "My destination", Enabled: true},
		},
	}, nil).Once()

	// First read (right after create) returns the connected destination.
	connections.On("GetTransformation", mock.Anything, "transformation-id").Return(&transformations.Transformation{
		ID: "transformation-id",
		Destinations: []transformations.TransformationDestination{
			{ID: "destination-id", Name: "My destination", Enabled: true},
		},
	}, nil).Once()

	// Subsequent read (refresh) shows the connection is gone.
	connections.On("GetTransformation", mock.Anything, "transformation-id").Return(&transformations.Transformation{
		ID:           "transformation-id",
		Destinations: nil,
	}, nil)

	config := `
		provider "rudderstack" {
			access_token = "some-access-token"
		}

		resource "rudderstack_transformation_connection" "example" {
			transformation_id = "transformation-id"
			destination_id    = "destination-id"
		}
	`

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"rudderstack": func() (*schema.Provider, error) {
				return NewWithConfigureClientFunc(func(_ context.Context, d *schema.ResourceData) (*Client, diag.Diagnostics) {
					return &Client{
						TransformationConnections: connections,
					}, diag.Diagnostics{}
				}), nil
			},
		},
		Steps: []resource.TestStep{
			{
				// After apply, the framework performs a refresh. The drifted
				// GetTransformation (no destinations) causes the resource to be
				// removed from state, producing a non-empty (re-create) plan.
				Config:             config,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

type mockTransformationConnectionsService struct {
	mock.Mock
}

func (m *mockTransformationConnectionsService) GetTransformation(ctx context.Context, id string) (*transformations.Transformation, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*transformations.Transformation), args.Error(1)
}

func (m *mockTransformationConnectionsService) ConnectToDestination(ctx context.Context, id string, req *transformations.ConnectToDestinationRequest) (*transformations.Transformation, error) {
	args := m.Called(ctx, id, req)
	return args.Get(0).(*transformations.Transformation), args.Error(1)
}

func (m *mockTransformationConnectionsService) DisconnectFromDestination(ctx context.Context, id string, req *transformations.ConnectToDestinationRequest) (*transformations.Transformation, error) {
	args := m.Called(ctx, id, req)
	return args.Get(0).(*transformations.Transformation), args.Error(1)
}
