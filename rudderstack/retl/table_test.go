package retl_test

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/mock"

	"github.com/rudderlabs/rudder-iac/api/client"
	iacretl "github.com/rudderlabs/rudder-iac/api/client/retl"
	"github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/retl"
)

func TestResourceTable(t *testing.T) {
	svc := &mockService{}

	createReq := &iacretl.RETLSourceCreateRequest{
		Name:                 "users-table",
		Config:               mustJSON(t, iacretl.RETLTableConfig{PrimaryKey: "id", Schema: "public", Table: "users"}),
		SourceType:           iacretl.TableSourceType,
		SourceDefinitionName: "snowflake",
		AccountID:            "acc-1",
		Enabled:              true,
	}
	created := &iacretl.RETLSource{
		ID:                   "tbl-1",
		Name:                 "users-table",
		Config:               createReq.Config,
		IsEnabled:            true,
		SourceType:           iacretl.TableSourceType,
		SourceDefinitionName: "snowflake",
		AccountID:            "acc-1",
		CreatedAt:            testutil.TimePtr(time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)),
		UpdatedAt:            testutil.TimePtr(time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)),
	}
	svc.On("CreateRetlSource", mock.Anything, createReq).Return(created, nil).Once()
	svc.On("GetRetlSource", mock.Anything, "tbl-1").Return(created, nil)
	svc.On("DeleteRetlSource", mock.Anything, "tbl-1").Return(nil).Once()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories(svc),
		Steps: []resource.TestStep{
			{
				Config: `
					provider "rudderstack" {
						access_token = "tok"
					}
					resource "rudderstack_retl_source_table" "example" {
						name                   = "users-table"
						source_definition_name = "snowflake"
						account_id             = "acc-1"
						config {
							primary_key = "id"
							schema      = "public"
							table       = "users"
						}
					}
				`,
				Check: func(s *terraform.State) error {
					attrs, err := resourceAttrs(s, "rudderstack_retl_source_table.example")
					if err != nil {
						return err
					}
					return checkAll(map[string]string{
						"id":                     "tbl-1",
						"source_definition_name": "snowflake",
						"config.0.schema":        "public",
						"config.0.table":         "users",
						"config.0.primary_key":   "id",
					}, attrs)
				},
			},
		},
	})

	svc.AssertExpectations(t)
}

func TestResourceTable_404RemovesFromState(t *testing.T) {
	svc := &mockService{}
	svc.On("GetRetlSource", mock.Anything, "tbl-gone").
		Return(nil, &client.APIError{HTTPStatusCode: 404}).Once()

	r := retl.ResourceTable()
	d := r.TestResourceData()
	d.SetId("tbl-gone")

	diags := r.ReadContext(context.Background(), d, &rudderstack.Client{RETLSources: svc})
	if diags.HasError() {
		t.Fatalf("expected no error, got %+v", diags)
	}
	if d.Id() != "" {
		t.Fatalf("expected ID cleared on 404, got %q", d.Id())
	}
	svc.AssertExpectations(t)
}
