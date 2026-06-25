package retl_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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

// Identifiers are mutable: changing them updates the connection in place
// (same ID, an UpdateConnection call) rather than recreating it. The mock
// registers exactly one CreateConnection and one UpdateConnection — no second
// Create and no Delete — so a destroy+create would fail the expectations.
func TestResourceConnectionCustomerIOAudience_IdentifiersAreMutable(t *testing.T) {
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
		ID:                "conn-cio-ids",
		SourceID:          "retl-src-1",
		DestinationID:     "dest-cio",
		Enabled:           true,
		Schedule:          iacretl.Schedule{Type: iacretl.ScheduleTypeManual},
		SyncBehaviour:     iacretl.SyncBehaviourMirror,
		Identifiers:       []iacretl.Mapping{{From: "email", To: "email"}},
		DestinationConfig: json.RawMessage(`{"audienceId":42}`),
	}
	updated := &iacretl.RETLConnection{
		ID:                "conn-cio-ids",
		SourceID:          "retl-src-1",
		DestinationID:     "dest-cio",
		Enabled:           true,
		Schedule:          iacretl.Schedule{Type: iacretl.ScheduleTypeManual},
		SyncBehaviour:     iacretl.SyncBehaviourMirror,
		Identifiers:       []iacretl.Mapping{{From: "user_id", To: "external_id"}},
		DestinationConfig: json.RawMessage(`{"audienceId":42}`),
	}
	svc.On("CreateConnection", mock.Anything, createReq).Return(created, nil).Once()
	// Reads before the update (apply read-back + step-1 refresh) return the
	// original; reads after the update return the new identifiers. Sequencing
	// the GetConnection stubs around the UpdateConnection call keeps step 1's
	// post-apply refresh stable (no spurious diff) while step 2 sees the change.
	getBeforeUpdate := svc.On("GetConnection", mock.Anything, "conn-cio-ids").Return(created, nil)
	// Update forwards the changed identifiers; audience_id is ForceNew but
	// unchanged here, so this is an in-place update.
	svc.On("UpdateConnection", mock.Anything, "conn-cio-ids", mock.MatchedBy(func(r *iacretl.UpdateRETLConnectionRequest) bool {
		return len(r.Identifiers) == 1 && r.Identifiers[0].From == "user_id" && r.Identifiers[0].To == "external_id"
	})).Return(updated, nil).Once().Run(func(mock.Arguments) {
		// Once the update lands, subsequent reads reflect the new identifiers.
		getBeforeUpdate.Return(updated, nil)
	})
	svc.On("DeleteConnection", mock.Anything, "conn-cio-ids").Return(nil).Once()

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
						"id":                 "conn-cio-ids",
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
					// Same ID proves in-place update (UpdateConnection), not
					// destroy + create.
					return checkAll(map[string]string{
						"id":                 "conn-cio-ids",
						"identifiers.0.from": "user_id",
						"identifiers.0.to":   "external_id",
					}, attrs)
				},
			},
		},
	})

	svc.AssertExpectations(t)
}

// A JSON-`null` destinationConfig is a server-side inconsistency for a
// Customer.io Audience connection (audienceId is mandatory). Refresh must
// emit a Warning diagnostic and NOT zero audience_id — zeroing would
// produce a never-reconcilable plan because audience_id is ForceNew with
// IntAtLeast(1) validation. Same expectation for an empty payload.
func TestResourceConnectionCustomerIOAudience_NullDestinationConfigOnRead(t *testing.T) {
	cases := []struct {
		name              string
		destinationConfig json.RawMessage
	}{
		{name: "json null", destinationConfig: json.RawMessage(`null`)},
		{name: "empty payload", destinationConfig: nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc := &mockService{}
			conn := &iacretl.RETLConnection{
				ID:                "conn-cio-null",
				SourceID:          "retl-src-1",
				DestinationID:     "dest-cio",
				Enabled:           true,
				Schedule:          iacretl.Schedule{Type: iacretl.ScheduleTypeManual},
				SyncBehaviour:     iacretl.SyncBehaviourMirror,
				Identifiers:       []iacretl.Mapping{{From: "email", To: "email"}},
				DestinationConfig: tc.destinationConfig,
			}
			svc.On("GetConnection", mock.Anything, "conn-cio-null").Return(conn, nil).Once()

			r := retl.ResourceConnectionCustomerIOAudience()
			d := r.TestResourceData()
			d.SetId("conn-cio-null")
			// Seed prior audience_id so we can prove refresh leaves it untouched.
			require.NoError(t, d.Set("audience_id", 42))

			diags := r.ReadContext(context.Background(), d, &rudderstack.Client{RETLSources: svc})
			require.False(t, diags.HasError(), "diags=%v", diags)
			require.Len(t, diags, 1, "expected one warning diagnostic")
			require.Equal(t, diag.Warning, diags[0].Severity)
			require.Regexp(t, `missing audienceId`, diags[0].Summary)
			require.Equal(t, 42, d.Get("audience_id"),
				"prior audience_id must be preserved to avoid a never-reconcilable plan")
			svc.AssertExpectations(t)
		})
	}
}

// cursor_column is a base-schema field, so the audience resource inherits it.
// Verify it round-trips on a config-less refresh.
func TestResourceConnectionCustomerIOAudience_CursorColumnSurvivesRefresh(t *testing.T) {
	svc := &mockService{}
	conn := &iacretl.RETLConnection{
		ID:                "conn-aud",
		SourceID:          "retl-src-1",
		DestinationID:     "dest-cio",
		Enabled:           true,
		Schedule:          iacretl.Schedule{Type: iacretl.ScheduleTypeManual},
		SyncBehaviour:     iacretl.SyncBehaviourUpsert,
		Identifiers:       []iacretl.Mapping{{From: "email", To: "email"}},
		CursorColumn:      "updated_at",
		DestinationConfig: json.RawMessage(`{"audienceId":42}`),
	}
	svc.On("GetConnection", mock.Anything, "conn-aud").Return(conn, nil).Once()

	r := retl.ResourceConnectionCustomerIOAudience()
	d := r.TestResourceData()
	d.SetId("conn-aud")

	diags := r.ReadContext(context.Background(), d, &rudderstack.Client{RETLSources: svc})
	require.False(t, diags.HasError(), "diags=%v", diags)
	require.Equal(t, "updated_at", d.Get("cursor_column"))
	require.Equal(t, 42, d.Get("audience_id"))
	svc.AssertExpectations(t)
}

// The audience resource enforces the shared cursor+upsert rule at plan time.
func TestResourceConnectionCustomerIOAudience_RejectsCursorColumnWithNonUpsert(t *testing.T) {
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
						audience_id   = 42
						cursor_column = "updated_at"
					}
				`,
				ExpectError: regexpMatches(`cursor_column is only valid when sync_behaviour is`),
			},
		},
	})
	svc.AssertNotCalled(t, "CreateConnection", mock.Anything, mock.Anything)
}
