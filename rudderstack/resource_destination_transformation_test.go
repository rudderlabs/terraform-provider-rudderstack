package rudderstack

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/mock"

	"github.com/rudderlabs/rudder-iac/api/client"
	"github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

// A synthetic destination type with a trivial (empty-mapping) config schema, so
// the test exercises the transformation_id lifecycle without a real
// integration's config. No Properties => config round-trips as "{}".
func transformationTestProviderFactory(dests DestinationsService) func() (*schema.Provider, error) {
	cm := configs.ConfigMeta{APIType: "TEST", SkipConfig: true}
	return func() (*schema.Provider, error) {
		return &schema.Provider{
			ConfigureContextFunc: func(_ context.Context, _ *schema.ResourceData) (interface{}, diag.Diagnostics) {
				return &Client{Destinations: dests}, diag.Diagnostics{}
			},
			ResourcesMap: map[string]*schema.Resource{
				"rudderstack_destination_test": resourceDestination(cm),
			},
		}, nil
	}
}

func destWithName() *client.Destination {
	return &client.Destination{
		ID:        "dest-1",
		Type:      "TEST",
		Name:      "test-destination",
		IsEnabled: true,
		CreatedAt: testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt: testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
	}
}

func transformationConfig(id string) string {
	return `
		resource "rudderstack_destination_test" "example" {
			name              = "test-destination"
			transformation_id = "` + id + `"
		}
	`
}

// Create attaches the transformation, Update re-attaches a different one, and
// Destroy disconnects before deleting — the destination-side reconciliation.
func TestResourceDestinationTransformationLink(t *testing.T) {
	dests := &mockDestinationsService{}

	// Create + Read (x2: post-create read, then pre-update refresh).
	dests.On("Create", mock.Anything, mock.Anything).Return(destWithName(), nil).Once()
	dests.On("ConnectTransformation", mock.Anything, "dest-1", "tPub1").
		Return(&client.DestinationTransformation{DestinationID: "dest-1", TransformationID: "tPub1"}, nil).Once()
	dests.On("Get", mock.Anything, "dest-1").Return(destWithName(), nil)
	dests.On("GetTransformation", mock.Anything, "dest-1").
		Return(&client.DestinationTransformation{DestinationID: "dest-1", TransformationID: "tPub1"}, nil).Times(3)

	// Update: swap tPub1 -> tPub2.
	dests.On("Update", mock.Anything, mock.Anything).Return(destWithName(), nil).Once()
	dests.On("ConnectTransformation", mock.Anything, "dest-1", "tPub2").
		Return(&client.DestinationTransformation{DestinationID: "dest-1", TransformationID: "tPub2"}, nil).Once()
	dests.On("GetTransformation", mock.Anything, "dest-1").
		Return(&client.DestinationTransformation{DestinationID: "dest-1", TransformationID: "tPub2"}, nil).Times(2)

	// Destroy: disconnect then delete.
	dests.On("DisconnectTransformation", mock.Anything, "dest-1").Return(nil).Once()
	dests.On("Delete", mock.Anything, "dest-1").Return(nil).Once()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"rudderstack": transformationTestProviderFactory(dests),
		},
		Steps: []resource.TestStep{
			{
				Config: transformationConfig("tPub1"),
				Check:  checkTransformationID("tPub1"),
			},
			{
				Config: transformationConfig("tPub2"),
				Check:  checkTransformationID("tPub2"),
			},
		},
	})

	dests.AssertExpectations(t)
}

func checkTransformationID(want string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		return resource.TestCheckResourceAttr(
			"rudderstack_destination_test.example", "transformation_id", want,
		)(state)
	}
}

type mockDestinationsService struct {
	mock.Mock
}

func (m *mockDestinationsService) Create(ctx context.Context, d *client.Destination) (*client.Destination, error) {
	args := m.Called(ctx, d)
	return args.Get(0).(*client.Destination), args.Error(1)
}

func (m *mockDestinationsService) Get(ctx context.Context, id string) (*client.Destination, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*client.Destination), args.Error(1)
}

func (m *mockDestinationsService) Update(ctx context.Context, d *client.Destination) (*client.Destination, error) {
	args := m.Called(ctx, d)
	return args.Get(0).(*client.Destination), args.Error(1)
}

func (m *mockDestinationsService) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockDestinationsService) GetTransformation(ctx context.Context, destinationID string) (*client.DestinationTransformation, error) {
	args := m.Called(ctx, destinationID)
	t, _ := args.Get(0).(*client.DestinationTransformation)
	return t, args.Error(1)
}

func (m *mockDestinationsService) ConnectTransformation(ctx context.Context, destinationID, transformationID string) (*client.DestinationTransformation, error) {
	args := m.Called(ctx, destinationID, transformationID)
	t, _ := args.Get(0).(*client.DestinationTransformation)
	return t, args.Error(1)
}

func (m *mockDestinationsService) DisconnectTransformation(ctx context.Context, destinationID string) error {
	return m.Called(ctx, destinationID).Error(0)
}
