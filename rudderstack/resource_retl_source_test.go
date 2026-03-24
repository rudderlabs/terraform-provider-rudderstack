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

func TestRETLSourceResource(t *testing.T) {
	retlSources := &mockRETLSourcesService{}

	retlSources.On("Create", mock.Anything, &client.RETLSource{
		Name:                 "my-bq-source",
		SourceDefinitionName: "bigquery",
		SourceType:           "table",
		AccountID:            "acct-id",
		PrimaryKey:           "user_id",
		IsEnabled:            true,
		Config:               json.RawMessage(`{"schema":"my_dataset","table":"my_table"}`),
	}).Return(&client.RETLSource{
		ID:                   "retl-src-id",
		Name:                 "my-bq-source",
		SourceDefinitionName: "bigquery",
		SourceType:           "table",
		AccountID:            "acct-id",
		PrimaryKey:           "user_id",
		IsEnabled:            true,
		WriteKey:             "write-key-1",
		Config:               json.RawMessage(`{"schema":"my_dataset","table":"my_table"}`),
		CreatedAt:            testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt:            testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
	}, nil)

	retlSources.On("Update", mock.Anything, &client.RETLSource{
		ID:                   "retl-src-id",
		Name:                 "my-bq-source-updated",
		SourceDefinitionName: "bigquery",
		SourceType:           "table",
		AccountID:            "acct-id-2",
		PrimaryKey:           "email",
		IsEnabled:            true,
		Config:               json.RawMessage(`{"schema":"my_dataset_2","table":"my_table_2"}`),
	}).Return(&client.RETLSource{
		ID:                   "retl-src-id",
		Name:                 "my-bq-source-updated",
		SourceDefinitionName: "bigquery",
		SourceType:           "table",
		AccountID:            "acct-id-2",
		PrimaryKey:           "email",
		IsEnabled:            true,
		WriteKey:             "write-key-1",
		Config:               json.RawMessage(`{"schema":"my_dataset_2","table":"my_table_2"}`),
		CreatedAt:            testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt:            testutil.TimePtr(time.Date(2010, 2, 2, 3, 4, 5, 0, time.UTC)),
	}, nil)

	retlSources.On("Get", mock.Anything, "retl-src-id").Return(&client.RETLSource{
		ID:                   "retl-src-id",
		Name:                 "my-bq-source",
		SourceDefinitionName: "bigquery",
		SourceType:           "table",
		AccountID:            "acct-id",
		PrimaryKey:           "user_id",
		IsEnabled:            true,
		WriteKey:             "write-key-1",
		Config:               json.RawMessage(`{"schema":"my_dataset","table":"my_table"}`),
		CreatedAt:            testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt:            testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
	}, nil).Times(3)

	retlSources.On("Get", mock.Anything, "retl-src-id").Return(&client.RETLSource{
		ID:                   "retl-src-id",
		Name:                 "my-bq-source-updated",
		SourceDefinitionName: "bigquery",
		SourceType:           "table",
		AccountID:            "acct-id-2",
		PrimaryKey:           "email",
		IsEnabled:            true,
		WriteKey:             "write-key-1",
		Config:               json.RawMessage(`{"schema":"my_dataset_2","table":"my_table_2"}`),
		CreatedAt:            testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt:            testutil.TimePtr(time.Date(2010, 2, 2, 3, 4, 5, 0, time.UTC)),
	}, nil).Twice()

	retlSources.On("Delete", mock.Anything, "retl-src-id").Return(nil)

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"rudderstack": func() (*schema.Provider, error) {
				return NewWithConfigureClientFunc(func(_ context.Context, d *schema.ResourceData) (*Client, diag.Diagnostics) {
					return &Client{
						RETLSources: retlSources,
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

					resource "rudderstack_retl_source" "example" {
						name                   = "my-bq-source"
						source_definition_name = "bigquery"
						source_type            = "table"
						account_id             = "acct-id"
						primary_key            = "user_id"

						config {
							schema_name = "my_dataset"
							table_name  = "my_table"
						}
					}
				`,
				Check: func(state *terraform.State) error {
					resources := state.RootModule().Resources
					resource, ok := resources["rudderstack_retl_source.example"]
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
					if c, ok := attributes["write_key"]; !ok || c != "write-key-1" {
						return fmt.Errorf("write_key was not set properly in state")
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

					resource "rudderstack_retl_source" "example" {
						name                   = "my-bq-source-updated"
						source_definition_name = "bigquery"
						source_type            = "table"
						account_id             = "acct-id-2"
						primary_key            = "email"

						config {
							schema_name = "my_dataset_2"
							table_name  = "my_table_2"
						}
					}
				`,
				Check: func(state *terraform.State) error {
					resources := state.RootModule().Resources
					resource, ok := resources["rudderstack_retl_source.example"]
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

type mockRETLSourcesService struct {
	mock.Mock
}

func (m *mockRETLSourcesService) Get(ctx context.Context, id string) (*client.RETLSource, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*client.RETLSource), args.Error(1)
}

func (m *mockRETLSourcesService) Create(ctx context.Context, source *client.RETLSource) (*client.RETLSource, error) {
	args := m.Called(ctx, source)
	return args.Get(0).(*client.RETLSource), args.Error(1)
}

func (m *mockRETLSourcesService) Update(ctx context.Context, source *client.RETLSource) (*client.RETLSource, error) {
	args := m.Called(ctx, source)
	return args.Get(0).(*client.RETLSource), args.Error(1)
}

func (m *mockRETLSourcesService) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
