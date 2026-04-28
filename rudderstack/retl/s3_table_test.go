package retl_test

import (
	"context"
	"fmt"
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

func TestResourceS3Table(t *testing.T) {
	svc := &mockService{}

	// S3 source: source_definition_name is hardcoded to "s3" by the provider
	// and is not part of the resource schema — verify both behaviors.
	createReq := &iacretl.RETLSourceCreateRequest{
		Name:                 "s3-events",
		Config:               iacretl.RETLS3TableConfig{BucketName: "my-bucket", ObjectPrefix: "events/"},
		SourceType:           iacretl.TableSourceType,
		SourceDefinitionName: "s3",
		AccountID:            "acc-1",
		Enabled:              true,
	}
	created := &iacretl.RETLSource{
		ID:                   "s3-1",
		Name:                 "s3-events",
		Config:               createReq.Config,
		IsEnabled:            true,
		SourceType:           iacretl.TableSourceType,
		SourceDefinitionName: "s3",
		AccountID:            "acc-1",
		CreatedAt:            testutil.TimePtr(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)),
		UpdatedAt:            testutil.TimePtr(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)),
	}
	svc.On("CreateRetlSource", mock.Anything, createReq).Return(created, nil).Once()
	svc.On("GetRetlSource", mock.Anything, "s3-1").Return(created, nil)
	svc.On("DeleteRetlSource", mock.Anything, "s3-1").Return(nil).Once()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories(svc),
		Steps: []resource.TestStep{
			{
				Config: `
					provider "rudderstack" {
						access_token = "tok"
					}
					resource "rudderstack_retl_source_s3_table" "example" {
						name       = "s3-events"
						account_id = "acc-1"
						config {
							bucket_name   = "my-bucket"
							object_prefix = "events/"
						}
					}
				`,
				Check: func(s *terraform.State) error {
					attrs, err := resourceAttrs(s, "rudderstack_retl_source_s3_table.example")
					if err != nil {
						return err
					}
					if _, ok := attrs["source_definition_name"]; ok {
						return fmt.Errorf("expected source_definition_name to not be a schema attribute on s3 resource")
					}
					return checkAll(map[string]string{
						"id":                     "s3-1",
						"config.0.bucket_name":   "my-bucket",
						"config.0.object_prefix": "events/",
					}, attrs)
				},
			},
		},
	})

	svc.AssertExpectations(t)
}

func TestResourceS3Table_404RemovesFromState(t *testing.T) {
	svc := &mockService{}
	svc.On("GetRetlSource", mock.Anything, "s3-gone").
		Return(nil, &client.APIError{HTTPStatusCode: 404}).Once()

	r := retl.ResourceS3Table()
	d := r.TestResourceData()
	d.SetId("s3-gone")

	diags := r.ReadContext(context.Background(), d, &rudderstack.Client{RETLSources: svc})
	if diags.HasError() {
		t.Fatalf("expected no error, got %+v", diags)
	}
	if d.Id() != "" {
		t.Fatalf("expected ID cleared on 404, got %q", d.Id())
	}
	svc.AssertExpectations(t)
}
