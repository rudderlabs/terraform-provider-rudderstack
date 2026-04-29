package retl_test

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/rudderlabs/rudder-iac/api/client"
	iacretl "github.com/rudderlabs/rudder-iac/api/client/retl"
	"github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/retl"
)

func regexpMatches(pattern string) *regexp.Regexp { return regexp.MustCompile(pattern) }

// jsonMapperConnection returns a fully-populated JSON Mapper connection
// shared by several tests below.
func jsonMapperConnection(id string) *iacretl.RETLConnection {
	enabled := true
	every := 60
	return &iacretl.RETLConnection{
		ID:            id,
		SourceID:      "retl-src-1",
		DestinationID: "dest-1",
		Enabled:       enabled,
		Schedule:      iacretl.Schedule{Type: iacretl.ScheduleTypeBasic, EveryMinutes: &every},
		SyncBehaviour: iacretl.SyncBehaviourUpsert,
		Identifiers:   []iacretl.Mapping{{From: "email", To: "user_id"}},
		Mappings:      []iacretl.Mapping{{From: "name", To: "first_name"}},
		Event:         &iacretl.Event{Type: iacretl.EventTypeIdentify},
		CreatedAt:     testutil.TimePtr(time.Date(2026, 4, 20, 12, 0, 0, 0, time.UTC)),
		UpdatedAt:     testutil.TimePtr(time.Date(2026, 4, 20, 12, 0, 0, 0, time.UTC)),
	}
}

func TestResourceConnection_JSONMapper_CreateReadUpdateDelete(t *testing.T) {
	svc := &mockService{}
	enabled := true
	every60 := 60

	createReq := &iacretl.CreateRETLConnectionRequest{
		SourceID:      "retl-src-1",
		DestinationID: "dest-1",
		Enabled:       &enabled,
		Schedule:      iacretl.Schedule{Type: iacretl.ScheduleTypeBasic, EveryMinutes: &every60},
		SyncBehaviour: iacretl.SyncBehaviourUpsert,
		Identifiers:   []iacretl.Mapping{{From: "email", To: "user_id"}},
		Mappings:      []iacretl.Mapping{{From: "name", To: "first_name"}},
		Event:         &iacretl.Event{Type: iacretl.EventTypeIdentify},
	}
	created := jsonMapperConnection("conn-1")
	svc.On("CreateConnection", mock.Anything, createReq).Return(created, nil).Once()
	svc.On("GetConnection", mock.Anything, "conn-1").Return(created, nil).Times(3)

	every120 := 120
	updateReq := &iacretl.UpdateRETLConnectionRequest{
		Enabled:  &enabled,
		Schedule: iacretl.Schedule{Type: iacretl.ScheduleTypeBasic, EveryMinutes: &every120},
		Mappings: &[]iacretl.Mapping{
			{From: "name", To: "first_name"},
			{From: "phone", To: "phone_number"},
		},
		Constants: &[]iacretl.Constant{{Key: "properties.source", Value: "warehouse"}},
	}
	updated := *created
	updated.Schedule.EveryMinutes = &every120
	updated.Mappings = *updateReq.Mappings
	updated.Constants = *updateReq.Constants
	updated.UpdatedAt = testutil.TimePtr(time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC))
	svc.On("UpdateConnection", mock.Anything, "conn-1", updateReq).Return(&updated, nil).Once()
	svc.On("GetConnection", mock.Anything, "conn-1").Return(&updated, nil).Times(2)
	svc.On("DeleteConnection", mock.Anything, "conn-1").Return(nil).Once()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories(svc),
		Steps: []resource.TestStep{
			{
				Config: `
					provider "rudderstack" { access_token = "tok" }
					resource "rudderstack_retl_connection" "example" {
						source_id      = "retl-src-1"
						destination_id = "dest-1"
						sync_behaviour = "upsert"
						schedule {
							type          = "basic"
							every_minutes = 60
						}
						identifiers {
							from = "email"
							to   = "user_id"
						}
						mappings {
							from = "name"
							to   = "first_name"
						}
						event {
							type = "identify"
						}
					}
				`,
				Check: func(s *terraform.State) error {
					attrs, err := resourceAttrs(s, "rudderstack_retl_connection.example")
					if err != nil {
						return err
					}
					return checkAll(map[string]string{
						"id":                       "conn-1",
						"source_id":                "retl-src-1",
						"destination_id":           "dest-1",
						"enabled":                  "true",
						"sync_behaviour":           "upsert",
						"schedule.0.type":          "basic",
						"schedule.0.every_minutes": "60",
						"identifiers.0.from":       "email",
						"identifiers.0.to":         "user_id",
						"mappings.0.from":          "name",
						"mappings.0.to":            "first_name",
						"event.0.type":             "identify",
						"created_at":               "2026-04-20T12:00:00Z",
					}, attrs)
				},
			},
			{
				Config: `
					provider "rudderstack" { access_token = "tok" }
					resource "rudderstack_retl_connection" "example" {
						source_id      = "retl-src-1"
						destination_id = "dest-1"
						sync_behaviour = "upsert"
						schedule {
							type          = "basic"
							every_minutes = 120
						}
						identifiers {
							from = "email"
							to   = "user_id"
						}
						mappings {
							from = "name"
							to   = "first_name"
						}
						mappings {
							from = "phone"
							to   = "phone_number"
						}
						event {
							type = "identify"
						}
						constants {
							key   = "properties.source"
							value = "warehouse"
						}
					}
				`,
				Check: func(s *terraform.State) error {
					attrs, err := resourceAttrs(s, "rudderstack_retl_connection.example")
					if err != nil {
						return err
					}
					return checkAll(map[string]string{
						"schedule.0.every_minutes": "120",
						"mappings.#":               "2",
						"mappings.1.from":          "phone",
						"constants.0.key":          "properties.source",
						"constants.0.value":        "warehouse",
						"updated_at":               "2026-04-21T12:00:00Z",
					}, attrs)
				},
			},
		},
	})

	svc.AssertExpectations(t)
}

func TestResourceConnection_404RemovesFromState(t *testing.T) {
	svc := &mockService{}
	svc.On("GetConnection", mock.Anything, "conn-gone").
		Return(nil, &client.APIError{HTTPStatusCode: 404}).Once()

	r := retl.ResourceConnection()
	d := r.TestResourceData()
	d.SetId("conn-gone")

	diags := r.ReadContext(context.Background(), d, &rudderstack.Client{RETLSources: svc})
	require.False(t, diags.HasError())
	require.Empty(t, d.Id())
	svc.AssertExpectations(t)
}

func TestResourceConnection_CursorColumnRequiresUpsert(t *testing.T) {
	svc := &mockService{}
	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories(svc),
		Steps: []resource.TestStep{
			{
				Config: `
					provider "rudderstack" { access_token = "tok" }
					resource "rudderstack_retl_connection" "example" {
						source_id      = "retl-src-1"
						destination_id = "dest-1"
						sync_behaviour = "mirror"
						cursor_column  = "updated_at"
						schedule {
							type = "manual"
						}
						identifiers {
							from = "email"
							to   = "user_id"
						}
					}
				`,
				ExpectError: regexpMatches(`cursor_column is only valid when sync_behaviour is "upsert"`),
			},
		},
	})
	svc.AssertNotCalled(t, "CreateConnection", mock.Anything, mock.Anything)
}

func TestResourceConnection_CronScheduleRequiresExpression(t *testing.T) {
	svc := &mockService{}
	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories(svc),
		Steps: []resource.TestStep{
			{
				Config: `
					provider "rudderstack" { access_token = "tok" }
					resource "rudderstack_retl_connection" "example" {
						source_id      = "retl-src-1"
						destination_id = "dest-1"
						sync_behaviour = "upsert"
						schedule {
							type = "cron"
						}
						identifiers {
							from = "email"
							to   = "user_id"
						}
					}
				`,
				ExpectError: regexpMatches(`schedule\.cron_expression is required when schedule\.type is "cron"`),
			},
		},
	})
	svc.AssertNotCalled(t, "CreateConnection", mock.Anything, mock.Anything)
}

func TestResourceConnection_BasicScheduleRequiresEveryMinutes(t *testing.T) {
	svc := &mockService{}
	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories(svc),
		Steps: []resource.TestStep{
			{
				Config: `
					provider "rudderstack" { access_token = "tok" }
					resource "rudderstack_retl_connection" "example" {
						source_id      = "retl-src-1"
						destination_id = "dest-1"
						sync_behaviour = "upsert"
						schedule {
							type = "basic"
						}
						identifiers {
							from = "email"
							to   = "user_id"
						}
					}
				`,
				ExpectError: regexpMatches(`schedule\.every_minutes is required when schedule\.type is "basic"`),
			},
		},
	})
	svc.AssertNotCalled(t, "CreateConnection", mock.Anything, mock.Anything)
}

func TestResourceConnection_Delete404IsTreatedAsAlreadyGone(t *testing.T) {
	svc := &mockService{}
	svc.On("DeleteConnection", mock.Anything, "conn-gone").
		Return(&client.APIError{HTTPStatusCode: 404}).Once()

	r := retl.ResourceConnection()
	d := r.TestResourceData()
	d.SetId("conn-gone")

	diags := r.DeleteContext(context.Background(), d, &rudderstack.Client{RETLSources: svc})
	require.False(t, diags.HasError())
	require.Empty(t, d.Id())
	svc.AssertExpectations(t)
}

