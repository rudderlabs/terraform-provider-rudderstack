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
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
	"github.com/stretchr/testify/mock"
)

func TestSourceResource(t *testing.T) {
	for source, cm := range configs.Sources.Entries() {
		sources := &mockSourcesService{}

		sources.On("Create", mock.Anything, &client.Source{
			Type:      cm.APIType,
			Name:      "example",
			IsEnabled: true,
			Config:    json.RawMessage(cm.TestConfigs[0].APICreate),
		}).Return(&client.Source{
			ID:        "some-id",
			Type:      cm.APIType,
			Name:      "example",
			IsEnabled: true,
			Config:    json.RawMessage(cm.TestConfigs[0].APICreate),
			CreatedAt: timePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
			UpdatedAt: timePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		}, nil)

		sources.On("Update", mock.Anything, &client.Source{
			ID:        "some-id",
			Type:      cm.APIType,
			Name:      "example-updated",
			IsEnabled: true,
			Config:    json.RawMessage(cm.TestConfigs[0].APIUpdate),
		}).Return(&client.Source{
			ID:        "some-id",
			Type:      cm.APIType,
			Name:      "example-updated",
			IsEnabled: true,
			Config:    json.RawMessage(cm.TestConfigs[0].APIUpdate),
			CreatedAt: timePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
			UpdatedAt: timePtr(time.Date(2010, 2, 2, 3, 4, 5, 0, time.UTC)),
		}, nil)

		sources.On("Get", mock.Anything, "some-id").Return(&client.Source{
			ID:        "some-id",
			Type:      cm.APIType,
			Name:      "example",
			IsEnabled: true,
			Config:    json.RawMessage(cm.TestConfigs[0].APICreate),
			CreatedAt: timePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
			UpdatedAt: timePtr(time.Date(2010, 2, 2, 3, 4, 5, 0, time.UTC)),
		}, nil).Times(3)

		sources.On("Get", mock.Anything, "some-id").Return(&client.Source{
			ID:        "some-id",
			Type:      cm.APIType,
			Name:      "example-updated",
			IsEnabled: true,
			Config:    json.RawMessage(cm.TestConfigs[0].APIUpdate),
			CreatedAt: timePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
			UpdatedAt: timePtr(time.Date(2010, 2, 2, 3, 4, 5, 0, time.UTC)),
		}, nil).Twice()

		sources.On("Delete", mock.Anything, "some-id").Return(nil)

		resource.UnitTest(t, resource.TestCase{
			ProviderFactories: map[string]func() (*schema.Provider, error){
				"rudderstack": func() (*schema.Provider, error) {
					return NewWithConfigureClientFunc(func(_ context.Context, d *schema.ResourceData) (*Client, diag.Diagnostics) {
						return &Client{
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
							config {
								%s
							}
						}
					`, source, cm.TestConfigs[0].TerraformCreate),
					Check: func(state *terraform.State) error {
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
							config {
								%s
							}
						}
					`, source, cm.TestConfigs[0].TerraformUpdate),
					Check: func(state *terraform.State) error {
						return nil
					},
				},
			},
		})
	}
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

func timePtr(t time.Time) *time.Time {
	return &t
}
