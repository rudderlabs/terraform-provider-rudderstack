package rudderstack

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/mock"

	"github.com/rudderlabs/rudder-iac/api/client"
)

// TestDataSourceAccountRead exercises dataSourceAccountRead using a mock
// AccountsService injected directly into *Client. It verifies that all
// non-secret attributes (name, definition_name, type, options) are populated
// from the mocked *Account, and that no "secret" attribute exists in the schema.
func TestDataSourceAccountRead(t *testing.T) {
	accounts := &mockAccountsService{}

	optionsJSON := `{"warehouse":"us-east-1","project":"my-project"}`

	accounts.On("Get", mock.Anything, "acc_1").Return(&client.Account{
		ID:   "acc_1",
		Name: "my-bigquery-account",
		Definition: struct {
			Name     string `json:"name"`
			Type     string `json:"type"`
			Category string `json:"category"`
		}{
			Name:     "SOURCE_BIGQUERY",
			Type:     "bigquery",
			Category: "warehouse",
		},
		Options: json.RawMessage(optionsJSON),
	}, nil)

	ds := dataSourceAccount()

	// Verify there is no "secret" attribute in the schema (CONTRACT §3).
	if _, exists := ds.Schema["secret"]; exists {
		t.Fatal("data source schema must not contain a 'secret' attribute")
	}

	d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
		"id": "acc_1",
	})

	c := &Client{Accounts: accounts}
	diags := dataSourceAccountRead(context.Background(), d, c)
	if diags.HasError() {
		t.Fatalf("unexpected error from dataSourceAccountRead: %v", diags)
	}

	// ID must be set from the returned account.
	if got := d.Id(); got != "acc_1" {
		t.Errorf("expected id=acc_1, got %q", got)
	}

	// name
	if got := d.Get("name").(string); got != "my-bigquery-account" {
		t.Errorf("expected name=my-bigquery-account, got %q", got)
	}

	// definition_name
	if got := d.Get("definition_name").(string); got != "SOURCE_BIGQUERY" {
		t.Errorf("expected definition_name=SOURCE_BIGQUERY, got %q", got)
	}

	// type
	if got := d.Get("type").(string); got != "bigquery" {
		t.Errorf("expected type=bigquery, got %q", got)
	}

	// options — must be a map[string]interface{} with string values.
	rawOpts := d.Get("options")
	optsMap, ok := rawOpts.(map[string]interface{})
	if !ok {
		t.Fatalf("expected options to be map[string]interface{}, got %T", rawOpts)
	}
	checkOpt := func(key, want string) {
		t.Helper()
		v, exists := optsMap[key]
		if !exists {
			t.Errorf("options[%q] not found", key)
			return
		}
		if got, _ := v.(string); got != want {
			t.Errorf("options[%q]: expected %q, got %q", key, want, got)
		}
	}
	checkOpt("warehouse", "us-east-1")
	checkOpt("project", "my-project")

	accounts.AssertExpectations(t)
}

// TestDataSourceAccountReadNilOptions verifies that a nil/empty options payload
// results in an empty map rather than an error.
func TestDataSourceAccountReadNilOptions(t *testing.T) {
	accounts := &mockAccountsService{}

	accounts.On("Get", mock.Anything, "acc_2").Return(&client.Account{
		ID:   "acc_2",
		Name: "no-options-account",
		Definition: struct {
			Name     string `json:"name"`
			Type     string `json:"type"`
			Category string `json:"category"`
		}{
			Name: "SOURCE_SNOWFLAKE",
			Type: "snowflake",
		},
		Options: nil,
	}, nil)

	ds := dataSourceAccount()
	d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
		"id": "acc_2",
	})

	c := &Client{Accounts: accounts}
	diags := dataSourceAccountRead(context.Background(), d, c)
	if diags.HasError() {
		t.Fatalf("unexpected error from dataSourceAccountRead with nil options: %v", diags)
	}

	rawOpts := d.Get("options")
	optsMap, ok := rawOpts.(map[string]interface{})
	if !ok {
		t.Fatalf("expected options to be map[string]interface{}, got %T", rawOpts)
	}
	if len(optsMap) != 0 {
		t.Errorf("expected empty options map, got %v", optsMap)
	}

	accounts.AssertExpectations(t)
}
