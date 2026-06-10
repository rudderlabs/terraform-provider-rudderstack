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

	"github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestAccountResource(t *testing.T) {
	cm := testAccountConfigMeta()
	accounts := &mockAccountsService{}

	createdAt := testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC))
	updatedAtV1 := testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC))
	updatedAtV2 := testutil.TimePtr(time.Date(2010, 2, 2, 3, 4, 5, 0, time.UTC))

	accounts.On("Create", mock.Anything, mock.MatchedBy(func(req *CreateAccountRequest) bool {
		return req.Name == "example" &&
			req.AccountDefinitionName == "SOURCE_BIGQUERY" &&
			testutil.JSONEq(string(req.Options), `{"project":"example-project","location":"US"}`) &&
			testutil.JSONEq(string(req.Secret), `{"credentials":"create-credentials"}`)
	})).Return(&Account{
		ID:   "some-id",
		Name: "example",
		Definition: AccountTypeInfo{
			Name:     "BigQuery",
			Type:     "warehouse",
			Category: "source",
		},
		Options:   json.RawMessage(`{"project":"example-project","location":"US"}`),
		CreatedAt: createdAt,
		UpdatedAt: updatedAtV1,
	}, nil)

	accounts.On("Update", mock.Anything, "some-id", mock.MatchedBy(func(req *UpdateAccountRequest) bool {
		return req.Name == "example-updated" &&
			testutil.JSONEq(string(req.Options), `{"project":"example-updated-project","location":"EU"}`) &&
			testutil.JSONEq(string(req.Secret), `{"credentials":"update-credentials"}`)
	})).Return(&Account{
		ID:   "some-id",
		Name: "example-updated",
		Definition: AccountTypeInfo{
			Name:     "BigQuery",
			Type:     "warehouse",
			Category: "source",
		},
		Options:   json.RawMessage(`{"project":"example-updated-project","location":"EU"}`),
		CreatedAt: createdAt,
		UpdatedAt: updatedAtV2,
	}, nil)

	accounts.On("Get", mock.Anything, "some-id").Return(&Account{
		ID:   "some-id",
		Name: "example",
		Definition: AccountTypeInfo{
			Name:     "BigQuery",
			Type:     "warehouse",
			Category: "source",
		},
		Options:   json.RawMessage(`{"project":"example-project","location":"US"}`),
		CreatedAt: createdAt,
		UpdatedAt: updatedAtV1,
	}, nil).Times(3)

	accounts.On("Get", mock.Anything, "some-id").Return(&Account{
		ID:   "some-id",
		Name: "example-updated",
		Definition: AccountTypeInfo{
			Name:     "BigQuery",
			Type:     "warehouse",
			Category: "source",
		},
		Options:   json.RawMessage(`{"project":"example-updated-project","location":"EU"}`),
		CreatedAt: createdAt,
		UpdatedAt: updatedAtV2,
	}, nil).Twice()

	accounts.On("Delete", mock.Anything, "some-id").Return(nil)

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"rudderstack": func() (*schema.Provider, error) {
				return testProviderWithAccountResource(cm, accounts), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccountResourceConfig("example", "example-project", "US", "create-credentials"),
				Check:  checkAccountState("example-project", "US", "create-credentials", "2010-01-02T03:04:05Z", "2010-01-02T03:04:05Z"),
			},
			{
				Config: testAccountResourceConfig("example-updated", "example-updated-project", "EU", "update-credentials"),
				Check:  checkAccountState("example-updated-project", "EU", "update-credentials", "2010-01-02T03:04:05Z", "2010-02-02T03:04:05Z"),
			},
		},
	})

	accounts.AssertExpectations(t)
}

func TestAccountResourceImportState(t *testing.T) {
	cm := testAccountConfigMeta()
	accounts := &mockAccountsService{}

	createdAt := testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC))
	updatedAt := testutil.TimePtr(time.Date(2010, 1, 2, 3, 4, 5, 0, time.UTC))

	accounts.On("Create", mock.Anything, mock.AnythingOfType("*rudderstack.CreateAccountRequest")).Return(&Account{
		ID:        "some-id",
		Name:      "example",
		Options:   json.RawMessage(`{"project":"example-project","location":"US"}`),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil)
	accounts.On("Get", mock.Anything, "some-id").Return(&Account{
		ID:        "some-id",
		Name:      "example",
		Options:   json.RawMessage(`{"project":"example-project","location":"US"}`),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil).Times(3)
	accounts.On("Delete", mock.Anything, "some-id").Return(nil)

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"rudderstack": func() (*schema.Provider, error) {
				return testProviderWithAccountResource(cm, accounts), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccountResourceConfig("example", "example-project", "US", "create-credentials"),
			},
			{
				Config:            testAccountResourceConfig("example", "example-project", "US", "create-credentials"),
				ResourceName:      "rudderstack_account_mock.example",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"config.0.credentials",
				},
			},
		},
	})

	accounts.AssertExpectations(t)
}

func TestResourceAccountReadNotFound(t *testing.T) {
	cm := testAccountConfigMeta()
	accounts := &mockAccountsService{}
	accounts.On("Get", mock.Anything, "missing-id").Return((*Account)(nil), ErrAccountNotFound).Once()

	client := &Client{}
	setAccountsClientForTests(client, accounts)
	t.Cleanup(func() {
		setAccountsClientForTests(client, nil)
	})

	d := resourceAccount(cm).TestResourceData()
	d.SetId("missing-id")

	diagnostics := resourceAccountRead(cm)(context.Background(), d, client)
	if diagnostics.HasError() {
		t.Fatalf("expected no diagnostics, got: %v", diagnostics)
	}
	if d.Id() != "" {
		t.Fatalf("expected id to be cleared on not found, got %q", d.Id())
	}

	accounts.AssertExpectations(t)
}

func TestResourceAccountDeleteNotFound(t *testing.T) {
	cm := testAccountConfigMeta()
	accounts := &mockAccountsService{}
	accounts.On("Delete", mock.Anything, "missing-id").Return(ErrAccountNotFound).Once()

	client := &Client{}
	setAccountsClientForTests(client, accounts)
	t.Cleanup(func() {
		setAccountsClientForTests(client, nil)
	})

	d := resourceAccount(cm).TestResourceData()
	d.SetId("missing-id")

	diagnostics := resourceAccountDelete(cm)(context.Background(), d, client)
	if diagnostics.HasError() {
		t.Fatalf("expected no diagnostics, got: %v", diagnostics)
	}
	if d.Id() != "" {
		t.Fatalf("expected id to be cleared on not found, got %q", d.Id())
	}

	accounts.AssertExpectations(t)
}

func testProviderWithAccountResource(cm configs.ConfigMeta, accounts accountsAPI) *schema.Provider {
	p := NewWithConfigureClientFunc(func(_ context.Context, _ *schema.ResourceData) (*Client, diag.Diagnostics) {
		c := &Client{}
		setAccountsClientForTests(c, accounts)
		return c, nil
	})
	p.ResourcesMap["rudderstack_account_mock"] = resourceAccount(cm)
	return p
}

func setAccountsClientForTests(c *Client, svc accountsAPI) {
	if c == nil {
		return
	}
	if svc == nil {
		accountServiceByClient.Delete(c)
		return
	}
	accountServiceByClient.Store(c, svc)
}

func testAccountConfigMeta() configs.ConfigMeta {
	return configs.ConfigMeta{
		APIType: "SOURCE_BIGQUERY",
		ConfigSchema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
			},
			"location": {
				Type:     schema.TypeString,
				Required: true,
			},
			"credentials": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
		},
		Properties: []configs.ConfigProperty{
			configs.Simple("options.project", "project"),
			configs.Simple("options.location", "location"),
			configs.Simple("secret.credentials", "credentials", configs.SkipZeroValue),
		},
	}
}

func testAccountResourceConfig(name, project, location, credentials string) string {
	return fmt.Sprintf(`
		provider "rudderstack" {
			access_token = "some-access-token"
		}

		resource "rudderstack_account_mock" "example" {
			name = %q
			config {
				project = %q
				location = %q
				credentials = %q
			}
		}
	`, name, project, location, credentials)
}

func checkAccountState(project, location, credentials, createdAt, updatedAt string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		resources := state.RootModule().Resources
		res, ok := resources["rudderstack_account_mock.example"]
		if !ok {
			return fmt.Errorf("resource not found in state")
		}

		attributes := res.Primary.Attributes
		if got := attributes["config.0.project"]; got != project {
			return fmt.Errorf("project mismatch: got %q want %q", got, project)
		}
		if got := attributes["config.0.location"]; got != location {
			return fmt.Errorf("location mismatch: got %q want %q", got, location)
		}
		if got := attributes["config.0.credentials"]; got != credentials {
			return fmt.Errorf("credentials mismatch: got %q want %q", got, credentials)
		}
		if got := attributes["created_at"]; got != createdAt {
			return fmt.Errorf("created_at mismatch: got %q want %q", got, createdAt)
		}
		if got := attributes["updated_at"]; got != updatedAt {
			return fmt.Errorf("updated_at mismatch: got %q want %q", got, updatedAt)
		}

		return nil
	}
}

type mockAccountsService struct {
	mock.Mock
}

func (m *mockAccountsService) Create(ctx context.Context, req *CreateAccountRequest) (*Account, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*Account), args.Error(1)
}

func (m *mockAccountsService) Get(ctx context.Context, id string) (*Account, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*Account), args.Error(1)
}

func (m *mockAccountsService) Update(ctx context.Context, id string, req *UpdateAccountRequest) (*Account, error) {
	args := m.Called(ctx, id, req)
	return args.Get(0).(*Account), args.Error(1)
}

func (m *mockAccountsService) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
