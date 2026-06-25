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
	"github.com/stretchr/testify/mock"

	"github.com/rudderlabs/rudder-iac/api/client"
	"github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

// mockAccountsService is a testify-mock backed AccountsService.
type mockAccountsService struct {
	mock.Mock
}

func (m *mockAccountsService) Create(ctx context.Context, req *client.CreateAccountRequest) (*client.Account, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*client.Account), args.Error(1)
}

func (m *mockAccountsService) Get(ctx context.Context, id string) (*client.Account, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*client.Account), args.Error(1)
}

func (m *mockAccountsService) Update(ctx context.Context, id string, req *client.UpdateAccountRequest) (*client.Account, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*client.Account), args.Error(1)
}

func (m *mockAccountsService) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// testAccountConfigMeta builds a minimal ConfigMeta used only in tests.
// It maps:
//   - options.foo → TF config.foo   (option field)
//   - secret.bar  → TF config.bar   (secret/credential field)
func testAccountConfigMeta() configs.ConfigMeta {
	return configs.ConfigMeta{
		APIType: "SOURCE_TESTWH",
		ConfigSchema: map[string]*schema.Schema{
			"foo": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bar": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
		},
		Properties: []configs.ConfigProperty{
			configs.Simple("options.foo", "foo"),
			configs.Simple("secret.bar", "bar"),
		},
	}
}

// TestResourceAccountCreateReadUpdate drives the full create→read→update lifecycle using
// a mock AccountsService injected into the Client.
func TestResourceAccountCreateReadUpdate(t *testing.T) {
	cm := testAccountConfigMeta()
	accounts := &mockAccountsService{}

	createOptionsJSON := `{"foo":"hello"}`
	createSecretJSON := `{"bar":"s3cr3t"}` //nolint:gosec // test fixture, not a real credential
	updateOptionsJSON := `{"foo":"world"}`
	updateSecretJSON := `{"bar":"s3cr3t-new"}` //nolint:gosec // test fixture, not a real credential

	// The provider registers a throwaway resource name for this test.
	const resourceType = "rudderstack_account_testwh"
	const resourceName = resourceType + ".example"

	// --- Create mock: assert exact payload ---
	accounts.On("Create", mock.Anything, mock.MatchedBy(func(req *client.CreateAccountRequest) bool {
		return req.Name == "example" &&
			req.AccountDefinitionName == cm.APIType &&
			testutil.JSONEq(string(req.Options), createOptionsJSON) &&
			testutil.JSONEq(string(req.Secret), createSecretJSON)
	})).Return(&client.Account{
		ID:   "acct-001",
		Name: "example",
		Definition: struct {
			Name     string `json:"name"`
			Type     string `json:"type"`
			Category string `json:"category"`
		}{
			Name: cm.APIType,
			Type: cm.APIType,
		},
		Options:   json.RawMessage(createOptionsJSON),
		CreatedAt: testutil.TimePtr(time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt: testutil.TimePtr(time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)),
	}, nil)

	// --- Read mock (called after Create, and once during plan for next step = 3 total) ---
	accounts.On("Get", mock.Anything, "acct-001").Return(&client.Account{
		ID:   "acct-001",
		Name: "example",
		Definition: struct {
			Name     string `json:"name"`
			Type     string `json:"type"`
			Category string `json:"category"`
		}{
			Name: cm.APIType,
			Type: cm.APIType,
		},
		// Secret intentionally absent — API never returns it.
		Options:   json.RawMessage(createOptionsJSON),
		CreatedAt: testutil.TimePtr(time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt: testutil.TimePtr(time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)),
	}, nil).Times(3)

	// --- Update mock ---
	accounts.On("Update", mock.Anything, "acct-001", mock.MatchedBy(func(req *client.UpdateAccountRequest) bool {
		return req.Name == "example-updated" &&
			testutil.JSONEq(string(req.Options), updateOptionsJSON) &&
			testutil.JSONEq(string(req.Secret), updateSecretJSON)
	})).Return(&client.Account{
		ID:   "acct-001",
		Name: "example-updated",
		Definition: struct {
			Name     string `json:"name"`
			Type     string `json:"type"`
			Category string `json:"category"`
		}{
			Name: cm.APIType,
			Type: cm.APIType,
		},
		Options:   json.RawMessage(updateOptionsJSON),
		CreatedAt: testutil.TimePtr(time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt: testutil.TimePtr(time.Date(2024, 2, 3, 4, 5, 6, 0, time.UTC)),
	}, nil)

	// --- Read mock after Update (called twice: post-update refresh + destroy plan check) ---
	accounts.On("Get", mock.Anything, "acct-001").Return(&client.Account{
		ID:   "acct-001",
		Name: "example-updated",
		Definition: struct {
			Name     string `json:"name"`
			Type     string `json:"type"`
			Category string `json:"category"`
		}{
			Name: cm.APIType,
			Type: cm.APIType,
		},
		Options:   json.RawMessage(updateOptionsJSON),
		CreatedAt: testutil.TimePtr(time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt: testutil.TimePtr(time.Date(2024, 2, 3, 4, 5, 6, 0, time.UTC)),
	}, nil).Twice()

	// --- Delete mock ---
	accounts.On("Delete", mock.Anything, "acct-001").Return(nil)

	// The provider uses a hand-built schema.Provider so we can control both the
	// ResourcesMap (registering just our test resource type) and the injected Client.
	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"rudderstack": func() (*schema.Provider, error) {
				return &schema.Provider{
					Schema: map[string]*schema.Schema{
						"access_token": {Type: schema.TypeString, Optional: true},
					},
					ConfigureContextFunc: func(_ context.Context, _ *schema.ResourceData) (interface{}, diag.Diagnostics) {
						return &Client{Accounts: accounts}, diag.Diagnostics{}
					},
					ResourcesMap: map[string]*schema.Resource{
						resourceType: resourceAccount(cm),
					},
				}, nil
			},
		},
		Steps: []resource.TestStep{
			{
				// Step 1: create with foo=hello, bar=s3cr3t.
				Config: fmt.Sprintf(`
					provider "rudderstack" {
						access_token = "test-token"
					}

					resource %q "example" {
						name = "example"
						config {
							foo = "hello"
							bar = "s3cr3t"
						}
					}
				`, resourceType),
				Check: resource.ComposeTestCheckFunc(
					// ID is set.
					resource.TestCheckResourceAttr(resourceName, "id", "acct-001"),
					// name is reflected in state.
					resource.TestCheckResourceAttr(resourceName, "name", "example"),
					// created_at and updated_at are set correctly.
					resource.TestCheckResourceAttr(resourceName, "created_at", "2024-01-02T03:04:05Z"),
					resource.TestCheckResourceAttr(resourceName, "updated_at", "2024-01-02T03:04:05Z"),
					// options field foo is present in state.
					resource.TestCheckResourceAttr(resourceName, "config.0.foo", "hello"),
					// bar (sensitive secret) retains the user-configured value in state;
					// the API never echoes it back, so the read does not overwrite it.
					// A non-empty value here is expected and correct (no perpetual diff).
					resource.TestCheckResourceAttr(resourceName, "config.0.bar", "s3cr3t"),
				),
			},
			{
				// Step 2: update name + foo + bar.
				Config: fmt.Sprintf(`
					provider "rudderstack" {
						access_token = "test-token"
					}

					resource %q "example" {
						name = "example-updated"
						config {
							foo = "world"
							bar = "s3cr3t-new"
						}
					}
				`, resourceType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", "acct-001"),
					resource.TestCheckResourceAttr(resourceName, "name", "example-updated"),
					resource.TestCheckResourceAttr(resourceName, "created_at", "2024-01-02T03:04:05Z"),
					resource.TestCheckResourceAttr(resourceName, "updated_at", "2024-02-03T04:05:06Z"),
					resource.TestCheckResourceAttr(resourceName, "config.0.foo", "world"),
					// Secret value updated — preserved in state (API never returns it).
					resource.TestCheckResourceAttr(resourceName, "config.0.bar", "s3cr3t-new"),
				),
			},
		},
	})
}

// TestResourceAccountRead_404ClearsState verifies that when the accounts API returns
// a 404 (out-of-band deletion / import of a missing ID), resourceAccountRead clears
// the resource ID and returns no Terraform error — matching the drift-detection
// pattern used in rudderstack/retl/common.go.
func TestResourceAccountRead_404ClearsState(t *testing.T) {
	cm := testAccountConfigMeta()
	accounts := &mockAccountsService{}

	accounts.On("Get", mock.Anything, "acct-gone").
		Return(nil, &client.APIError{HTTPStatusCode: 404}).Once()

	res := resourceAccount(cm)
	d := res.Data(nil)
	d.SetId("acct-gone")

	c := &Client{Accounts: accounts}
	diags := resourceAccountRead(cm)(context.Background(), d, c)
	if diags.HasError() {
		t.Fatalf("expected no error on 404, got %+v", diags)
	}
	if d.Id() != "" {
		t.Fatalf("expected ID cleared on 404, got %q", d.Id())
	}

	accounts.AssertExpectations(t)
}

// TestResourceAccountRead_NonNotFoundErrorReturnsError verifies that non-404 API errors
// are surfaced as Terraform diagnostics (not silently swallowed).
func TestResourceAccountRead_NonNotFoundErrorReturnsError(t *testing.T) {
	cm := testAccountConfigMeta()
	accounts := &mockAccountsService{}

	accounts.On("Get", mock.Anything, "acct-err").
		Return(nil, &client.APIError{HTTPStatusCode: 500, Message: "internal server error"}).Once()

	res := resourceAccount(cm)
	d := res.Data(nil)
	d.SetId("acct-err")

	c := &Client{Accounts: accounts}
	diags := resourceAccountRead(cm)(context.Background(), d, c)
	if !diags.HasError() {
		t.Fatal("expected error on 500, got none")
	}
	// ID must NOT be cleared on a non-404 error.
	if d.Id() == "" {
		t.Fatal("expected ID to remain set on non-404 error")
	}

	accounts.AssertExpectations(t)
}

// TestResourceAccountSchemaNoSecretInState verifies that after a Read whose mocked
// Get returns no secret, the bar field is absent from state. This test exercises the
// schema + read function directly without needing a full provider lifecycle.
func TestResourceAccountSchemaNoSecretInState(t *testing.T) {
	cm := testAccountConfigMeta()
	accounts := &mockAccountsService{}

	accounts.On("Get", mock.Anything, "acct-002").Return(&client.Account{
		ID:   "acct-002",
		Name: "test-acct",
		Definition: struct {
			Name     string `json:"name"`
			Type     string `json:"type"`
			Category string `json:"category"`
		}{
			Name: cm.APIType,
			Type: cm.APIType,
		},
		// Options returned, secret NOT returned.
		Options:   json.RawMessage(`{"foo":"val1"}`),
		CreatedAt: testutil.TimePtr(time.Date(2024, 3, 4, 5, 6, 7, 0, time.UTC)),
		UpdatedAt: testutil.TimePtr(time.Date(2024, 3, 4, 5, 6, 7, 0, time.UTC)),
	}, nil)

	res := resourceAccount(cm)
	d := res.Data(nil)
	d.SetId("acct-002")

	c := &Client{Accounts: accounts}
	diags := resourceAccountRead(cm)(context.Background(), d, c)
	if diags.HasError() {
		t.Fatalf("unexpected error from Read: %v", diags)
	}

	// foo (option) must be in state.
	config := d.Get("config").([]interface{})
	if len(config) == 0 {
		t.Fatal("expected config block in state, got empty")
	}
	props, ok := config[0].(map[string]interface{})
	if !ok {
		t.Fatal("config[0] is not a map")
	}
	if v, exists := props["foo"]; !exists || v != "val1" {
		t.Errorf("expected config.foo=val1, got %v (exists=%v)", v, exists)
	}
	// bar (secret) must NOT be in state — API returned no secret.
	if v, exists := props["bar"]; exists && v != "" {
		t.Errorf("expected config.bar to be absent or empty after read with no secret, got %q", v)
	}

	// Timestamps are set.
	if got := d.Get("created_at").(string); got != "2024-03-04T05:06:07Z" {
		t.Errorf("expected created_at=2024-03-04T05:06:07Z, got %q", got)
	}
	if got := d.Get("updated_at").(string); got != "2024-03-04T05:06:07Z" {
		t.Errorf("expected updated_at=2024-03-04T05:06:07Z, got %q", got)
	}

	accounts.AssertExpectations(t)
}
