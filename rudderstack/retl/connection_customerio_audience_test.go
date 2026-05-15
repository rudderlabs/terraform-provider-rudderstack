package retl_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	iacretl "github.com/rudderlabs/rudder-iac/api/client/retl"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/retl"
)

// The typed Customer.io Audience resource exposes audience_id as a top-level
// integer; internally it round-trips through destinationConfig as
// {"audienceId": N} so the API sees the same shape the generic resource used
// to produce.
func TestResourceConnectionCustomerIOAudience_CreateRead(t *testing.T) {
	svc := &mockService{}
	enabled := true

	createReq := &iacretl.CreateRETLConnectionRequest{
		SourceID:          "retl-src-1",
		DestinationID:     "dest-cio",
		Enabled:           &enabled,
		Schedule:          iacretl.Schedule{Type: iacretl.ScheduleTypeManual},
		SyncBehaviour:     iacretl.SyncBehaviourMirror,
		Identifiers:       []iacretl.Mapping{{From: "email", To: "email"}},
		DestinationConfig: json.RawMessage(`{"audienceId":42}`),
	}
	created := &iacretl.RETLConnection{
		ID:                "conn-cio",
		SourceID:          "retl-src-1",
		DestinationID:     "dest-cio",
		Enabled:           true,
		Schedule:          iacretl.Schedule{Type: iacretl.ScheduleTypeManual},
		SyncBehaviour:     iacretl.SyncBehaviourMirror,
		Identifiers:       []iacretl.Mapping{{From: "email", To: "email"}},
		DestinationConfig: json.RawMessage(`{"audienceId":42}`),
	}
	svc.On("CreateConnection", mock.Anything, createReq).Return(created, nil).Once()
	svc.On("GetConnection", mock.Anything, "conn-cio").Return(created, nil).Times(2)
	svc.On("DeleteConnection", mock.Anything, "conn-cio").Return(nil).Once()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories(svc),
		Steps: []resource.TestStep{
			{
				Config: `
					provider "rudderstack" { access_token = "tok" }
					resource "rudderstack_retl_connection_customerio_audience" "example" {
						source_id      = "retl-src-1"
						destination_id = "dest-cio"
						sync_behaviour = "mirror"
						schedule { type = "manual" }
						identifiers {
							from = "email"
							to   = "email"
						}
						audience_id = 42
					}
				`,
				Check: func(s *terraform.State) error {
					attrs, err := resourceAttrs(s, "rudderstack_retl_connection_customerio_audience.example")
					if err != nil {
						return err
					}
					return checkAll(map[string]string{
						"id":          "conn-cio",
						"audience_id": "42",
					}, attrs)
				},
			},
		},
	})

	svc.AssertExpectations(t)
}

// Changing `audience_id` must recreate the connection. The Customer.io
// Audience API rejects destinationConfig changes on update, so audience_id
// carries `ForceNew: true`. A new resource ID after the value change is the
// proof — terraform would have called UpdateConnection (which we don't mock)
// if ForceNew hadn't fired.
func TestResourceConnectionCustomerIOAudience_AudienceIDIsForceNew(t *testing.T) {
	svc := &mockService{}
	enabled := true

	firstReq := &iacretl.CreateRETLConnectionRequest{
		SourceID:          "retl-src-1",
		DestinationID:     "dest-cio",
		Enabled:           &enabled,
		Schedule:          iacretl.Schedule{Type: iacretl.ScheduleTypeManual},
		SyncBehaviour:     iacretl.SyncBehaviourMirror,
		Identifiers:       []iacretl.Mapping{{From: "email", To: "email"}},
		DestinationConfig: json.RawMessage(`{"audienceId":42}`),
	}
	firstCreated := &iacretl.RETLConnection{
		ID:                "conn-cio-1",
		SourceID:          "retl-src-1",
		DestinationID:     "dest-cio",
		Enabled:           true,
		Schedule:          iacretl.Schedule{Type: iacretl.ScheduleTypeManual},
		SyncBehaviour:     iacretl.SyncBehaviourMirror,
		Identifiers:       []iacretl.Mapping{{From: "email", To: "email"}},
		DestinationConfig: json.RawMessage(`{"audienceId":42}`),
	}
	secondReq := &iacretl.CreateRETLConnectionRequest{
		SourceID:          "retl-src-1",
		DestinationID:     "dest-cio",
		Enabled:           &enabled,
		Schedule:          iacretl.Schedule{Type: iacretl.ScheduleTypeManual},
		SyncBehaviour:     iacretl.SyncBehaviourMirror,
		Identifiers:       []iacretl.Mapping{{From: "email", To: "email"}},
		DestinationConfig: json.RawMessage(`{"audienceId":99}`),
	}
	secondCreated := &iacretl.RETLConnection{
		ID:                "conn-cio-2",
		SourceID:          "retl-src-1",
		DestinationID:     "dest-cio",
		Enabled:           true,
		Schedule:          iacretl.Schedule{Type: iacretl.ScheduleTypeManual},
		SyncBehaviour:     iacretl.SyncBehaviourMirror,
		Identifiers:       []iacretl.Mapping{{From: "email", To: "email"}},
		DestinationConfig: json.RawMessage(`{"audienceId":99}`),
	}
	svc.On("CreateConnection", mock.Anything, firstReq).Return(firstCreated, nil).Once()
	svc.On("CreateConnection", mock.Anything, secondReq).Return(secondCreated, nil).Once()
	svc.On("GetConnection", mock.Anything, "conn-cio-1").Return(firstCreated, nil)
	svc.On("GetConnection", mock.Anything, "conn-cio-2").Return(secondCreated, nil)
	svc.On("DeleteConnection", mock.Anything, "conn-cio-1").Return(nil).Once()
	svc.On("DeleteConnection", mock.Anything, "conn-cio-2").Return(nil).Once()

	cfg := func(audienceID int) string {
		return fmt.Sprintf(`
			provider "rudderstack" { access_token = "tok" }
			resource "rudderstack_retl_connection_customerio_audience" "example" {
				source_id      = "retl-src-1"
				destination_id = "dest-cio"
				sync_behaviour = "mirror"
				schedule { type = "manual" }
				identifiers {
					from = "email"
					to   = "email"
				}
				audience_id = %d
			}
		`, audienceID)
	}

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories(svc),
		Steps: []resource.TestStep{
			{
				Config: cfg(42),
				Check: func(s *terraform.State) error {
					attrs, err := resourceAttrs(s, "rudderstack_retl_connection_customerio_audience.example")
					if err != nil {
						return err
					}
					return checkAll(map[string]string{
						"id":          "conn-cio-1",
						"audience_id": "42",
					}, attrs)
				},
			},
			{
				Config: cfg(99),
				Check: func(s *terraform.State) error {
					attrs, err := resourceAttrs(s, "rudderstack_retl_connection_customerio_audience.example")
					if err != nil {
						return err
					}
					// New ID proves the resource was destroyed and recreated
					// rather than updated in place — i.e. ForceNew kicked in.
					return checkAll(map[string]string{
						"id":          "conn-cio-2",
						"audience_id": "99",
					}, attrs)
				},
			},
		},
	})

	svc.AssertExpectations(t)
}

// Customer.io audience IDs are always positive — plan-time validation must
// reject values < 1 so users see a clear error in `terraform plan` rather
// than a generic API rejection at apply time.
func TestResourceConnectionCustomerIOAudience_RejectsNonPositiveID(t *testing.T) {
	svc := &mockService{}
	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories(svc),
		Steps: []resource.TestStep{
			{
				Config: `
					provider "rudderstack" { access_token = "tok" }
					resource "rudderstack_retl_connection_customerio_audience" "example" {
						source_id      = "retl-src-1"
						destination_id = "dest-cio"
						sync_behaviour = "mirror"
						schedule { type = "manual" }
						identifiers {
							from = "email"
							to   = "email"
						}
						audience_id = 0
					}
				`,
				ExpectError: regexpMatches(`expected audience_id to be at least \(1\)`),
			},
		},
	})
	svc.AssertNotCalled(t, "CreateConnection", mock.Anything, mock.Anything)
}

// Identifier changes must recreate the typed resource (ForceNew is set at
// both the top level and the nested from/to in baseConnectionSchema). A new
// ID after the value change proves destroy + create — terraform would have
// called UpdateConnection (which the mock does not register) if the schema
// had allowed in-place updates.
func TestResourceConnectionCustomerIOAudience_IdentifiersAreForceNew(t *testing.T) {
	svc := &mockService{}
	enabled := true

	firstReq := &iacretl.CreateRETLConnectionRequest{
		SourceID:          "retl-src-1",
		DestinationID:     "dest-cio",
		Enabled:           &enabled,
		Schedule:          iacretl.Schedule{Type: iacretl.ScheduleTypeManual},
		SyncBehaviour:     iacretl.SyncBehaviourMirror,
		Identifiers:       []iacretl.Mapping{{From: "email", To: "email"}},
		DestinationConfig: json.RawMessage(`{"audienceId":42}`),
	}
	firstCreated := &iacretl.RETLConnection{
		ID:                "conn-cio-ids-1",
		SourceID:          "retl-src-1",
		DestinationID:     "dest-cio",
		Enabled:           true,
		Schedule:          iacretl.Schedule{Type: iacretl.ScheduleTypeManual},
		SyncBehaviour:     iacretl.SyncBehaviourMirror,
		Identifiers:       []iacretl.Mapping{{From: "email", To: "email"}},
		DestinationConfig: json.RawMessage(`{"audienceId":42}`),
	}
	secondReq := &iacretl.CreateRETLConnectionRequest{
		SourceID:          "retl-src-1",
		DestinationID:     "dest-cio",
		Enabled:           &enabled,
		Schedule:          iacretl.Schedule{Type: iacretl.ScheduleTypeManual},
		SyncBehaviour:     iacretl.SyncBehaviourMirror,
		Identifiers:       []iacretl.Mapping{{From: "user_id", To: "external_id"}},
		DestinationConfig: json.RawMessage(`{"audienceId":42}`),
	}
	secondCreated := &iacretl.RETLConnection{
		ID:                "conn-cio-ids-2",
		SourceID:          "retl-src-1",
		DestinationID:     "dest-cio",
		Enabled:           true,
		Schedule:          iacretl.Schedule{Type: iacretl.ScheduleTypeManual},
		SyncBehaviour:     iacretl.SyncBehaviourMirror,
		Identifiers:       []iacretl.Mapping{{From: "user_id", To: "external_id"}},
		DestinationConfig: json.RawMessage(`{"audienceId":42}`),
	}
	svc.On("CreateConnection", mock.Anything, firstReq).Return(firstCreated, nil).Once()
	svc.On("CreateConnection", mock.Anything, secondReq).Return(secondCreated, nil).Once()
	svc.On("GetConnection", mock.Anything, "conn-cio-ids-1").Return(firstCreated, nil)
	svc.On("GetConnection", mock.Anything, "conn-cio-ids-2").Return(secondCreated, nil)
	svc.On("DeleteConnection", mock.Anything, "conn-cio-ids-1").Return(nil).Once()
	svc.On("DeleteConnection", mock.Anything, "conn-cio-ids-2").Return(nil).Once()

	cfg := func(from, to string) string {
		return fmt.Sprintf(`
			provider "rudderstack" { access_token = "tok" }
			resource "rudderstack_retl_connection_customerio_audience" "example" {
				source_id      = "retl-src-1"
				destination_id = "dest-cio"
				sync_behaviour = "mirror"
				schedule { type = "manual" }
				identifiers {
					from = "%s"
					to   = "%s"
				}
				audience_id = 42
			}
		`, from, to)
	}

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories(svc),
		Steps: []resource.TestStep{
			{
				Config: cfg("email", "email"),
				Check: func(s *terraform.State) error {
					attrs, err := resourceAttrs(s, "rudderstack_retl_connection_customerio_audience.example")
					if err != nil {
						return err
					}
					return checkAll(map[string]string{
						"id":                 "conn-cio-ids-1",
						"identifiers.0.from": "email",
						"identifiers.0.to":   "email",
					}, attrs)
				},
			},
			{
				Config: cfg("user_id", "external_id"),
				Check: func(s *terraform.State) error {
					attrs, err := resourceAttrs(s, "rudderstack_retl_connection_customerio_audience.example")
					if err != nil {
						return err
					}
					// New ID proves destroy + create — ForceNew fired. The
					// mock also expects two distinct CreateConnection calls
					// (one per identifier value).
					return checkAll(map[string]string{
						"id":                 "conn-cio-ids-2",
						"identifiers.0.from": "user_id",
						"identifiers.0.to":   "external_id",
					}, attrs)
				},
			},
		},
	})

	svc.AssertExpectations(t)
}

// Regression: a JSON-`null` destinationConfig in the API response must NOT
// error on the typed resource. The read path treats `null` as "no typed
// config" and clears audience_id from state instead of failing refresh.
func TestResourceConnectionCustomerIOAudience_NullDestinationConfigOnRead(t *testing.T) {
	svc := &mockService{}
	conn := &iacretl.RETLConnection{
		ID:                "conn-cio-null",
		SourceID:          "retl-src-1",
		DestinationID:     "dest-cio",
		Enabled:           true,
		Schedule:          iacretl.Schedule{Type: iacretl.ScheduleTypeManual},
		SyncBehaviour:     iacretl.SyncBehaviourMirror,
		Identifiers:       []iacretl.Mapping{{From: "email", To: "email"}},
		DestinationConfig: json.RawMessage(`null`),
	}
	svc.On("GetConnection", mock.Anything, "conn-cio-null").Return(conn, nil).Once()

	r := retl.ResourceConnectionCustomerIOAudience()
	d := r.TestResourceData()
	d.SetId("conn-cio-null")

	diags := r.ReadContext(context.Background(), d, &rudderstack.Client{RETLSources: svc})
	require.False(t, diags.HasError(), "diags=%v", diags)
	require.Equal(t, 0, d.Get("audience_id"))
	svc.AssertExpectations(t)
}
