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

func TestWarehouseAccountResource(t *testing.T) {
	accounts := &mockAccountsService{}

	accounts.On("Create", mock.Anything, &client.AccountCreateInput{
		Name:                  "my-bq-account",
		AccountDefinitionName: "bigquery",
		Options: map[string]interface{}{
			"project":     "my-project",
			"credentials": "secret-json",
		},
	}).Return(&client.Account{
		ID:         "acct-id",
		Name:       "my-bq-account",
		Definition: &client.AccountDefinition{Type: "bigquery"},
		Options:    map[string]interface{}{"project": "my-project", "credentials": "secret-json"},
		CreatedAt:  testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt:  testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
	}, nil)

	accounts.On("Update", mock.Anything, "acct-id", &client.AccountUpdateInput{
		Name: "my-bq-account-updated",
		Options: map[string]interface{}{
			"project":     "my-project-2",
			"credentials": "secret-json-2",
		},
	}).Return(&client.Account{
		ID:         "acct-id",
		Name:       "my-bq-account-updated",
		Definition: &client.AccountDefinition{Type: "bigquery"},
		Options:    map[string]interface{}{"project": "my-project-2", "credentials": "secret-json-2"},
		CreatedAt:  testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt:  testutil.TimePtr(time.Date(2010, 2, 2, 3, 4, 5, 0, time.UTC)),
	}, nil)

	accounts.On("Get", mock.Anything, "acct-id").Return(&client.Account{
		ID:         "acct-id",
		Name:       "my-bq-account",
		Definition: &client.AccountDefinition{Type: "bigquery"},
		Options:    map[string]interface{}{"project": "my-project", "credentials": "secret-json"},
		CreatedAt:  testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt:  testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
	}, nil).Times(3)

	accounts.On("Get", mock.Anything, "acct-id").Return(&client.Account{
		ID:         "acct-id",
		Name:       "my-bq-account-updated",
		Definition: &client.AccountDefinition{Type: "bigquery"},
		Options:    map[string]interface{}{"project": "my-project-2", "credentials": "secret-json-2"},
		CreatedAt:  testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt:  testutil.TimePtr(time.Date(2010, 2, 2, 3, 4, 5, 0, time.UTC)),
	}, nil).Twice()

	accounts.On("Delete", mock.Anything, "acct-id").Return(nil)

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"rudderstack": func() (*schema.Provider, error) {
				return NewWithConfigureClientFunc(func(_ context.Context, d *schema.ResourceData) (*Client, diag.Diagnostics) {
					return &Client{
						Accounts: accounts,
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

					resource "rudderstack_warehouse_account" "example" {
						name                    = "my-bq-account"
						account_definition_name = "bigquery"
						options = {
							project     = "my-project"
							credentials = "secret-json"
						}
					}
				`,
				Check: func(state *terraform.State) error {
					resources := state.RootModule().Resources
					resource, ok := resources["rudderstack_warehouse_account.example"]
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

					resource "rudderstack_warehouse_account" "example" {
						name                    = "my-bq-account-updated"
						account_definition_name = "bigquery"
						options = {
							project     = "my-project-2"
							credentials = "secret-json-2"
						}
					}
				`,
				Check: func(state *terraform.State) error {
					resources := state.RootModule().Resources
					resource, ok := resources["rudderstack_warehouse_account.example"]
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

type mockAccountsService struct {
	mock.Mock
}

func (m *mockAccountsService) Get(ctx context.Context, id string) (*client.Account, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*client.Account), args.Error(1)
}

func (m *mockAccountsService) Create(ctx context.Context, input *client.AccountCreateInput) (*client.Account, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*client.Account), args.Error(1)
}

func (m *mockAccountsService) Update(ctx context.Context, id string, input *client.AccountUpdateInput) (*client.Account, error) {
	args := m.Called(ctx, id, input)
	return args.Get(0).(*client.Account), args.Error(1)
}

func (m *mockAccountsService) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
