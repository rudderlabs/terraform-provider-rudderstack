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
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/accounts"
	"github.com/stretchr/testify/mock"
)

type TestConfigAPIPayload struct {
	Options string
	Secret  string
}

type AccountsTestConfig struct {
	TerraformCreate string
	APICreate       TestConfigAPIPayload
	TerraformUpdate string
	APIUpdate       TestConfigAPIPayload
}

var EmptyTestConfig = AccountsTestConfig{TerraformCreate: "", APICreate: TestConfigAPIPayload{Options: "{}", Secret: "{}"}, TerraformUpdate: "", APIUpdate: TestConfigAPIPayload{Options: "{}", Secret: "{}"}}

func AssertAccount(t *testing.T, category accounts.AccountCategory, account string, testConfigs []AccountsTestConfig) {
	cm := accounts.Accounts.Entries()[category][account]
	accounts := &mockAccountsService{}

	accounts.On("Create", mock.Anything, mock.MatchedBy(func(a *client.AccountWithSecret) bool {
		return a.Type == cm.APIType &&
			a.Category == string(category) &&
			a.ID == "" &&
			a.Name == "example" &&
			testutil.JSONEq(string(a.Options), testConfigs[0].APICreate.Options) &&
			testutil.JSONEq(string(a.Secret), testConfigs[0].APICreate.Secret)
	})).Return(&client.Account{
		ID:        "some-id",
		Type:      cm.APIType,
		Name:      "example",
		Options:   json.RawMessage(testConfigs[0].APICreate.Options),
		CreatedAt: testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt: testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
	}, nil)

	accounts.On("Update", mock.Anything, mock.MatchedBy(func(a *client.AccountWithSecret) bool {
		return a.Type == cm.APIType &&
			a.ID == "some-id" &&
			a.Name == "example-updated" &&
			testutil.JSONEq(string(a.Options), testConfigs[0].APIUpdate.Options) &&
			testutil.JSONEq(string(a.Secret), testConfigs[0].APIUpdate.Secret)
	})).Return(&client.Account{
		ID:        "some-id",
		Type:      cm.APIType,
		Name:      "example-updated",
		Options:   json.RawMessage(testConfigs[0].APIUpdate.Options),
		CreatedAt: testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt: testutil.TimePtr(time.Date(2010, 2, 2, 3, 4, 5, 0, time.UTC)),
	}, nil)

	accounts.On("Get", mock.Anything, "some-id").Return(&client.Account{
		ID:        "some-id",
		Type:      cm.APIType,
		Category:  string(category),
		Name:      "example",
		Options:   json.RawMessage(testConfigs[0].APICreate.Options),
		CreatedAt: testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt: testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
	}, nil).Times(1)

	accounts.On("Get", mock.Anything, "some-id").Return(&client.Account{
		ID:        "some-id",
		Type:      cm.APIType,
		Name:      "example-updated",
		Options:   json.RawMessage(testConfigs[0].APIUpdate.Options),
		CreatedAt: testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt: testutil.TimePtr(time.Date(2010, 2, 2, 3, 4, 5, 0, time.UTC)),
	}, nil).Twice()

	accounts.On("Delete", mock.Anything, "some-id").Return(nil)

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"rudderstack": func() (*schema.Provider, error) {
				return rudderstack.NewWithConfigureClientFunc(func(_ context.Context, d *schema.ResourceData) (*rudderstack.Client, diag.Diagnostics) {
					return &rudderstack.Client{
						Accounts: accounts,
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

					resource "rudderstack_account_%s_%s" "example" {
						name = "example"
						config {
							%s
						}
					}
				`, category, account, testConfigs[0].TerraformCreate),
				Check: func(state *terraform.State) error {
					resources := state.RootModule().Resources
					resource, ok := resources[fmt.Sprintf("rudderstack_account_%s_%s.example", category, account)]
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

					resource "rudderstack_account_%s_%s" "example" {
						name = "example-updated"
						config {
							%s
						}
					}
				`, category, account, testConfigs[0].TerraformUpdate),
				Check: func(state *terraform.State) error {
					resources := state.RootModule().Resources
					resource, ok := resources[fmt.Sprintf("rudderstack_account_%s_%s.example", category, account)]
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

type mockAccountsService struct {
	mock.Mock
}

func (m *mockAccountsService) Get(ctx context.Context, id string) (*client.Account, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*client.Account), args.Error(1)
}

func (m *mockAccountsService) Create(ctx context.Context, account *client.AccountWithSecret) (*client.Account, error) {
	args := m.Called(ctx, account)
	return args.Get(0).(*client.Account), args.Error(1)
}

func (m *mockAccountsService) Update(ctx context.Context, account *client.AccountWithSecret) (*client.Account, error) {
	args := m.Called(ctx, account)
	return args.Get(0).(*client.Account), args.Error(1)
}

func (m *mockAccountsService) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
