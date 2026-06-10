// Package cm provides unit-test helpers backed by testify mocks (no real API calls).
//
// # TestConfig convention for accounts
//
// AssertAccount reads testConfigs[0] and expects the following fields:
//
//   - TerraformCreate: raw HCL attribute/block lines to embed inside the `config { … }` block
//     for the "create" step. Example: `credentials = "..."\nproject = "my-proj"`.
//
//   - TerraformUpdate: same for the "update" step (name changes to "example-updated").
//
//   - APICreate: full JSON body of the Create request as the API receives it:
//     `{"name":"example","accountDefinitionName":"SOURCE_BIGQUERY","options":{…},"secret":{…}}`.
//     - "name" MUST be "example" (the harness hard-codes that name for the create step).
//     - "accountDefinitionName" MUST match configs.Accounts.Entries()[accountType].APIType.
//     - "options" is the non-secret account options object.
//     - "secret" is the secret/credentials object.
//
//   - APIUpdate: same structure, but "name" MUST be "example-updated" (the harness hard-codes
//     that name for the update step). "accountDefinitionName" is still present for documentation
//     but is NOT sent on Update (the mock only matches options/secret + name).
//
// The mock Get returns options from APICreate/APIUpdate (the secret is intentionally absent —
// the API never echoes it back; the resource preserves it from prior state).
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
	"github.com/stretchr/testify/mock"
	"github.com/tidwall/gjson"

	"github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

// AssertAccount is the mock-backed unit test helper for account resources.
//
// It registers a mockAccountsService, sets up Create/Update/Get/Delete expectations,
// then runs resource.UnitTest with a provider that injects the mock client.
// Two test steps are run: create (name="example") and update (name="example-updated").
//
// See the package-level doc comment for the TestConfig convention.
func AssertAccount(t *testing.T, accountType string, testConfigs []configs.TestConfig) {
	t.Helper()

	cm := configs.Accounts.Entries()[accountType]
	cfg := testConfigs[0]

	// Extract options and secret from the full request body JSON.
	createOptionsJSON := gjson.Get(cfg.APICreate, "options").Raw
	if createOptionsJSON == "" {
		createOptionsJSON = "{}"
	}
	createSecretJSON := gjson.Get(cfg.APICreate, "secret").Raw
	if createSecretJSON == "" {
		createSecretJSON = "{}"
	}

	updateOptionsJSON := gjson.Get(cfg.APIUpdate, "options").Raw
	if updateOptionsJSON == "" {
		updateOptionsJSON = "{}"
	}
	updateSecretJSON := gjson.Get(cfg.APIUpdate, "secret").Raw
	if updateSecretJSON == "" {
		updateSecretJSON = "{}"
	}

	accounts := &mockAccountsService{}

	// Create: match name, accountDefinitionName, options, and secret.
	accounts.On("Create", mock.Anything, mock.MatchedBy(func(req *rudderstack.CreateAccountRequest) bool {
		return req.Name == "example" &&
			req.AccountDefinitionName == cm.APIType &&
			testutil.JSONEq(string(req.Options), createOptionsJSON) &&
			testutil.JSONEq(string(req.Secret), createSecretJSON)
	})).Return(&rudderstack.Account{
		ID:   "some-id",
		Name: "example",
		Definition: rudderstack.AccountDefinition{
			Name: cm.APIType,
			Type: cm.APIType,
		},
		Options:   json.RawMessage(createOptionsJSON),
		CreatedAt: testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt: testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
	}, nil)

	// Update: match id, name, options, and secret (accountDefinitionName is NOT sent on update).
	accounts.On("Update", mock.Anything, "some-id", mock.MatchedBy(func(req *rudderstack.UpdateAccountRequest) bool {
		return req.Name == "example-updated" &&
			testutil.JSONEq(string(req.Options), updateOptionsJSON) &&
			testutil.JSONEq(string(req.Secret), updateSecretJSON)
	})).Return(&rudderstack.Account{
		ID:   "some-id",
		Name: "example-updated",
		Definition: rudderstack.AccountDefinition{
			Name: cm.APIType,
			Type: cm.APIType,
		},
		Options:   json.RawMessage(updateOptionsJSON),
		CreatedAt: testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt: testutil.TimePtr(time.Date(2010, 2, 2, 3, 4, 5, 0, time.UTC)),
	}, nil)

	// Get after Create: called 3 times (post-create refresh, plan check before step 2, and
	// one more read in step 2's apply — matches the destination/source pattern of .Times(3)).
	accounts.On("Get", mock.Anything, "some-id").Return(&rudderstack.Account{
		ID:   "some-id",
		Name: "example",
		Definition: rudderstack.AccountDefinition{
			Name: cm.APIType,
			Type: cm.APIType,
		},
		// Secret intentionally absent — the API never returns it.
		Options:   json.RawMessage(createOptionsJSON),
		CreatedAt: testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt: testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
	}, nil).Times(3)

	// Get after Update: called twice (post-update refresh + destroy plan check).
	accounts.On("Get", mock.Anything, "some-id").Return(&rudderstack.Account{
		ID:   "some-id",
		Name: "example-updated",
		Definition: rudderstack.AccountDefinition{
			Name: cm.APIType,
			Type: cm.APIType,
		},
		// Secret intentionally absent — the API never returns it.
		Options:   json.RawMessage(updateOptionsJSON),
		CreatedAt: testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt: testutil.TimePtr(time.Date(2010, 2, 2, 3, 4, 5, 0, time.UTC)),
	}, nil).Twice()

	accounts.On("Delete", mock.Anything, "some-id").Return(nil)

	resourceTypeName := fmt.Sprintf("rudderstack_account_source_%s", accountType)
	resourceName := fmt.Sprintf("%s.example", resourceTypeName)

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

					resource "rudderstack_account_source_%s" "example" {
						name = "example"
						config {
							%s
						}
					}
				`, accountType, cfg.TerraformCreate),
				Check: func(state *terraform.State) error {
					resources := state.RootModule().Resources
					res, ok := resources[resourceName]
					if !ok {
						return fmt.Errorf("resource not found in state: %s", resourceName)
					}
					attributes := res.Primary.Attributes
					if c, ok := attributes["created_at"]; !ok || c != "2010-01-02T03:04:05Z" {
						return fmt.Errorf("created_at was not set properly in state: got %q", attributes["created_at"])
					}
					if c, ok := attributes["updated_at"]; !ok || c != "2010-01-02T03:04:05Z" {
						return fmt.Errorf("updated_at was not set properly in state: got %q", attributes["updated_at"])
					}
					return nil
				},
			},
			{
				Config: fmt.Sprintf(`
					provider "rudderstack" {
						access_token = "some-access-token"
					}

					resource "rudderstack_account_source_%s" "example" {
						name = "example-updated"
						config {
							%s
						}
					}
				`, accountType, cfg.TerraformUpdate),
				Check: func(state *terraform.State) error {
					resources := state.RootModule().Resources
					res, ok := resources[resourceName]
					if !ok {
						return fmt.Errorf("resource not found in state: %s", resourceName)
					}
					attributes := res.Primary.Attributes
					if c, ok := attributes["created_at"]; !ok || c != "2010-01-02T03:04:05Z" {
						return fmt.Errorf("created_at was not set properly in state: got %q", attributes["created_at"])
					}
					if c, ok := attributes["updated_at"]; !ok || c != "2010-02-02T03:04:05Z" {
						return fmt.Errorf("updated_at was not set properly in state: got %q", attributes["updated_at"])
					}
					return nil
				},
			},
		},
	})
}

// mockAccountsService is a testify-mock implementation of rudderstack.AccountsService.
// All return values are nil-safe: a nil first return is returned as-is rather than
// causing a nil-pointer panic via type assertion (mirroring resource_account_test.go).
type mockAccountsService struct {
	mock.Mock
}

func (m *mockAccountsService) Create(ctx context.Context, req *rudderstack.CreateAccountRequest) (*rudderstack.Account, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*rudderstack.Account), args.Error(1)
}

func (m *mockAccountsService) Get(ctx context.Context, id string) (*rudderstack.Account, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*rudderstack.Account), args.Error(1)
}

func (m *mockAccountsService) Update(ctx context.Context, id string, req *rudderstack.UpdateAccountRequest) (*rudderstack.Account, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*rudderstack.Account), args.Error(1)
}

func (m *mockAccountsService) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
