package retl_test

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

	"github.com/rudderlabs/rudder-iac/api/client"
	iacretl "github.com/rudderlabs/rudder-iac/api/client/retl"
	"github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/retl"
)

func TestResourceModel(t *testing.T) {
	svc := &mockService{}

	createReq := &iacretl.RETLSourceCreateRequest{
		Name:                 "my-model",
		Config:               mustJSON(t, iacretl.RETLSQLModelConfig{PrimaryKey: "id", Sql: "select * from users"}),
		SourceType:           iacretl.ModelSourceType,
		SourceDefinitionName: "snowflake",
		AccountID:            "acc-1",
		Enabled:              true,
	}
	created := &iacretl.RETLSource{
		ID:                   "src-1",
		Name:                 "my-model",
		Config:               createReq.Config,
		IsEnabled:            true,
		SourceType:           iacretl.ModelSourceType,
		SourceDefinitionName: "snowflake",
		AccountID:            "acc-1",
		CreatedAt:            testutil.TimePtr(time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)),
		UpdatedAt:            testutil.TimePtr(time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)),
	}
	svc.On("CreateRetlSource", mock.Anything, createReq).Return(created, nil).Once()
	svc.On("GetRetlSource", mock.Anything, "src-1").Return(created, nil).Times(3)

	updatedConfig := mustJSON(t, iacretl.RETLSQLModelConfig{PrimaryKey: "id", Sql: "select id, name from users", Description: "v2"})
	updateReq := &iacretl.RETLSourceUpdateRequest{
		Name:      "my-model-v2",
		Config:    updatedConfig,
		IsEnabled: false,
		AccountID: "acc-1",
	}
	updated := &iacretl.RETLSource{
		ID:                   "src-1",
		Name:                 "my-model-v2",
		Config:               updatedConfig,
		IsEnabled:            false,
		SourceType:           iacretl.ModelSourceType,
		SourceDefinitionName: "snowflake",
		AccountID:            "acc-1",
		CreatedAt:            created.CreatedAt,
		UpdatedAt:            testutil.TimePtr(time.Date(2024, 2, 2, 3, 4, 5, 0, time.UTC)),
	}
	svc.On("UpdateRetlSource", mock.Anything, "src-1", updateReq).Return(updated, nil).Once()
	svc.On("GetRetlSource", mock.Anything, "src-1").Return(updated, nil).Times(2)
	svc.On("DeleteRetlSource", mock.Anything, "src-1").Return(nil).Once()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories(svc),
		Steps: []resource.TestStep{
			{
				Config: `
					provider "rudderstack" {
						access_token = "tok"
					}
					resource "rudderstack_retl_source_model" "example" {
						name                   = "my-model"
						source_definition_name = "snowflake"
						account_id             = "acc-1"
						config {
							primary_key = "id"
							sql         = "select * from users"
						}
					}
				`,
				Check: func(s *terraform.State) error {
					attrs, err := resourceAttrs(s, "rudderstack_retl_source_model.example")
					if err != nil {
						return err
					}
					return checkAll(map[string]string{
						"id":                     "src-1",
						"name":                   "my-model",
						"source_definition_name": "snowflake",
						"account_id":             "acc-1",
						"enabled":                "true",
						"config.0.primary_key":   "id",
						"config.0.sql":           "select * from users",
						"created_at":             "2024-01-02T03:04:05Z",
					}, attrs)
				},
			},
			{
				Config: `
					provider "rudderstack" {
						access_token = "tok"
					}
					resource "rudderstack_retl_source_model" "example" {
						name                   = "my-model-v2"
						source_definition_name = "snowflake"
						account_id             = "acc-1"
						enabled                = false
						config {
							primary_key = "id"
							sql         = "select id, name from users"
							description = "v2"
						}
					}
				`,
				Check: func(s *terraform.State) error {
					attrs, err := resourceAttrs(s, "rudderstack_retl_source_model.example")
					if err != nil {
						return err
					}
					return checkAll(map[string]string{
						"name":                 "my-model-v2",
						"enabled":              "false",
						"config.0.description": "v2",
						"config.0.sql":         "select id, name from users",
						"updated_at":           "2024-02-02T03:04:05Z",
					}, attrs)
				},
			},
		},
	})

	svc.AssertExpectations(t)
}

func TestResourceModel_404RemovesFromState(t *testing.T) {
	svc := &mockService{}
	svc.On("GetRetlSource", mock.Anything, "src-gone").
		Return(nil, &client.APIError{HTTPStatusCode: 404}).Once()

	r := retl.ResourceModel()
	d := r.TestResourceData()
	d.SetId("src-gone")

	diags := r.ReadContext(context.Background(), d, &rudderstack.Client{RETLSources: svc})
	if diags.HasError() {
		t.Fatalf("expected no error, got %+v", diags)
	}
	if d.Id() != "" {
		t.Fatalf("expected ID cleared on 404, got %q", d.Id())
	}
	svc.AssertExpectations(t)
}

// --- shared test helpers ---

func providerFactories(svc retl.Service) map[string]func() (*schema.Provider, error) {
	return map[string]func() (*schema.Provider, error){
		"rudderstack": func() (*schema.Provider, error) {
			return rudderstack.NewWithConfigureClientFunc(func(_ context.Context, _ *schema.ResourceData) (*rudderstack.Client, diag.Diagnostics) {
				return &rudderstack.Client{RETLSources: svc}, diag.Diagnostics{}
			}), nil
		},
	}
}

func resourceAttrs(state *terraform.State, name string) (map[string]string, error) {
	r, ok := state.RootModule().Resources[name]
	if !ok {
		return nil, fmt.Errorf("resource %q not found in state", name)
	}
	return r.Primary.Attributes, nil
}

func checkAll(want, got map[string]string) error {
	for k, v := range want {
		if g, ok := got[k]; !ok || g != v {
			return fmt.Errorf("attr %q: want %q, got %q", k, v, g)
		}
	}
	return nil
}

func mustJSON(t *testing.T, v interface{}) json.RawMessage {
	t.Helper()
	raw, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return raw
}
