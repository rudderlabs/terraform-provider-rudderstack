package retl_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	iacretl "github.com/rudderlabs/rudder-iac/api/client/retl"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/retl"
)

// The typed Customer.io (VDM v2) resource exposes `object` as a top-level
// string; internally it round-trips through destinationConfig as
// {"object": "..."}. identifiers flow through the base schema; VDM v2 does NOT
// support field mappings, so this resource has no `mappings`.
func TestResourceConnectionCustomerIO_CreateRead(t *testing.T) {
	svc := &mockService{}
	enabled := true

	createReq := &iacretl.CreateRETLConnectionRequest{
		SourceID:          "retl-src-1",
		DestinationID:     "dest-cio",
		Enabled:           &enabled,
		Schedule:          iacretl.Schedule{Type: iacretl.ScheduleTypeManual},
		SyncBehaviour:     iacretl.SyncBehaviourUpsert,
		Identifiers:       []iacretl.Mapping{{From: "email", To: "email"}},
		DestinationConfig: json.RawMessage(`{"object":"person"}`),
	}
	created := &iacretl.RETLConnection{
		ID:                "conn-cio",
		SourceID:          "retl-src-1",
		DestinationID:     "dest-cio",
		Enabled:           true,
		Schedule:          iacretl.Schedule{Type: iacretl.ScheduleTypeManual},
		SyncBehaviour:     iacretl.SyncBehaviourUpsert,
		Identifiers:       []iacretl.Mapping{{From: "email", To: "email"}},
		DestinationConfig: json.RawMessage(`{"object":"person"}`),
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
					resource "rudderstack_retl_connection_customerio" "example" {
						source_id      = "retl-src-1"
						destination_id = "dest-cio"
						sync_behaviour = "upsert"
						object         = "person"
						schedule { type = "manual" }
						identifiers {
							from = "email"
							to   = "email"
						}
					}
				`,
				Check: func(s *terraform.State) error {
					attrs, err := resourceAttrs(s, "rudderstack_retl_connection_customerio.example")
					if err != nil {
						return err
					}
					return checkAll(map[string]string{
						"id":     "conn-cio",
						"object": "person",
					}, attrs)
				},
			},
		},
	})

	svc.AssertExpectations(t)
}

// VDM v2 does not support field mappings — the resource must reject a
// `mappings` block at plan time (unknown argument) rather than silently
// accept it.
func TestResourceConnectionCustomerIO_RejectsMappings(t *testing.T) {
	svc := &mockService{}
	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories(svc),
		Steps: []resource.TestStep{
			{
				Config: `
					provider "rudderstack" { access_token = "tok" }
					resource "rudderstack_retl_connection_customerio" "example" {
						source_id      = "retl-src-1"
						destination_id = "dest-cio"
						sync_behaviour = "upsert"
						object         = "person"
						schedule { type = "manual" }
						identifiers {
							from = "email"
							to   = "email"
						}
						mappings {
							from = "name"
							to   = "plan"
						}
					}
				`,
				ExpectError: regexpMatches(`[Uu]nsupported argument|not expected here`),
			},
		},
	})
	svc.AssertNotCalled(t, "CreateConnection", mock.Anything, mock.Anything)
}

// CustomerIO supports exactly one object — `person` (the listObjects value).
// The resource must reject any other object value at plan time rather than
// letting the server fail on apply.
func TestResourceConnectionCustomerIO_RejectsUnknownObject(t *testing.T) {
	svc := &mockService{}
	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories(svc),
		Steps: []resource.TestStep{
			{
				Config: `
					provider "rudderstack" { access_token = "tok" }
					resource "rudderstack_retl_connection_customerio" "example" {
						source_id      = "retl-src-1"
						destination_id = "dest-cio"
						sync_behaviour = "upsert"
						object         = "accounts"
						schedule { type = "manual" }
						identifiers {
							from = "email"
							to   = "email"
						}
					}
				`,
				ExpectError: regexpMatches(`expected object to be one of`),
			},
		},
	})
	svc.AssertNotCalled(t, "CreateConnection", mock.Anything, mock.Anything)
}

// VDM v2 supports only upsert and mirror sync modes; the typed resource must
// reject `full` at plan time rather than letting the API reject it on apply.
func TestResourceConnectionCustomerIO_RejectsFullSyncBehaviour(t *testing.T) {
	svc := &mockService{}
	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories(svc),
		Steps: []resource.TestStep{
			{
				Config: `
					provider "rudderstack" { access_token = "tok" }
					resource "rudderstack_retl_connection_customerio" "example" {
						source_id      = "retl-src-1"
						destination_id = "dest-cio"
						sync_behaviour = "full"
						object         = "person"
						schedule { type = "manual" }
						identifiers {
							from = "email"
							to   = "email"
						}
					}
				`,
				ExpectError: regexpMatches(`expected sync_behaviour to be one of`),
			},
		},
	})
	svc.AssertNotCalled(t, "CreateConnection", mock.Anything, mock.Anything)
}

// cursor_column is a generic source-side field (incremental watermark) sent
// as a top-level request field — NOT inside destinationConfig. config-be
// stores it on config.source. It is only valid with sync_behaviour="upsert".
func TestResourceConnectionCustomerIO_CursorColumnCreateRead(t *testing.T) {
	svc := &mockService{}
	enabled := true

	createReq := &iacretl.CreateRETLConnectionRequest{
		SourceID:          "retl-src-1",
		DestinationID:     "dest-cio",
		Enabled:           &enabled,
		Schedule:          iacretl.Schedule{Type: iacretl.ScheduleTypeManual},
		SyncBehaviour:     iacretl.SyncBehaviourUpsert,
		Identifiers:       []iacretl.Mapping{{From: "email", To: "email"}},
		CursorColumn:      "updated_at",
		DestinationConfig: json.RawMessage(`{"object":"person"}`),
	}
	created := &iacretl.RETLConnection{
		ID:                "conn-cio",
		SourceID:          "retl-src-1",
		DestinationID:     "dest-cio",
		Enabled:           true,
		Schedule:          iacretl.Schedule{Type: iacretl.ScheduleTypeManual},
		SyncBehaviour:     iacretl.SyncBehaviourUpsert,
		Identifiers:       []iacretl.Mapping{{From: "email", To: "email"}},
		CursorColumn:      "updated_at",
		DestinationConfig: json.RawMessage(`{"object":"person"}`),
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
					resource "rudderstack_retl_connection_customerio" "example" {
						source_id      = "retl-src-1"
						destination_id = "dest-cio"
						sync_behaviour = "upsert"
						object         = "person"
						cursor_column  = "updated_at"
						schedule { type = "manual" }
						identifiers {
							from = "email"
							to   = "email"
						}
					}
				`,
				Check: func(s *terraform.State) error {
					attrs, err := resourceAttrs(s, "rudderstack_retl_connection_customerio.example")
					if err != nil {
						return err
					}
					return checkAll(map[string]string{
						"id":            "conn-cio",
						"object":        "person",
						"cursor_column": "updated_at",
					}, attrs)
				},
			},
		},
	})

	svc.AssertExpectations(t)
}

// cursor_column must round-trip on a bare refresh (ReadContext with no config,
// only a prior Id) — the real terraform refresh scenario. This guards the
// regression where a config-dependent schema-presence check would silently
// drop cursor_column from state because config is absent during refresh,
// producing perpetual drift.
func TestResourceConnectionCustomerIO_CursorColumnSurvivesRefresh(t *testing.T) {
	svc := &mockService{}
	conn := &iacretl.RETLConnection{
		ID:                "conn-cio",
		SourceID:          "retl-src-1",
		DestinationID:     "dest-cio",
		Enabled:           true,
		Schedule:          iacretl.Schedule{Type: iacretl.ScheduleTypeManual},
		SyncBehaviour:     iacretl.SyncBehaviourUpsert,
		Identifiers:       []iacretl.Mapping{{From: "email", To: "email"}},
		CursorColumn:      "updated_at",
		DestinationConfig: json.RawMessage(`{"object":"person"}`),
	}
	svc.On("GetConnection", mock.Anything, "conn-cio").Return(conn, nil).Once()

	r := retl.ResourceConnectionCustomerIO()
	d := r.TestResourceData()
	d.SetId("conn-cio")

	diags := r.ReadContext(context.Background(), d, &rudderstack.Client{RETLSources: svc})
	require.False(t, diags.HasError(), "diags=%v", diags)
	require.Equal(t, "updated_at", d.Get("cursor_column"),
		"cursor_column must round-trip on refresh (no config present)")
	require.Equal(t, "person", d.Get("object"))
	svc.AssertExpectations(t)
}

// cursor_column is only meaningful for incremental upsert syncs. The resource
// must reject it at plan time when sync_behaviour is not "upsert" (e.g.
// "mirror"), mirroring the generic resource's CustomizeDiff check, rather than
// letting config-be reject it on apply.
func TestResourceConnectionCustomerIO_RejectsCursorColumnWithNonUpsert(t *testing.T) {
	svc := &mockService{}
	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories(svc),
		Steps: []resource.TestStep{
			{
				Config: `
					provider "rudderstack" { access_token = "tok" }
					resource "rudderstack_retl_connection_customerio" "example" {
						source_id      = "retl-src-1"
						destination_id = "dest-cio"
						sync_behaviour = "mirror"
						object         = "person"
						cursor_column  = "updated_at"
						schedule { type = "manual" }
						identifiers {
							from = "email"
							to   = "email"
						}
					}
				`,
				ExpectError: regexpMatches(`cursor_column is only valid when sync_behaviour is`),
			},
		},
	})
	svc.AssertNotCalled(t, "CreateConnection", mock.Anything, mock.Anything)
}

// A 200 GetConnection with no usable object (JSON `null` or empty payload) is
// a persistent server-side inconsistency, not a transient glitch the HTTP
// layer would retry. Refresh must hard-error so the broken connection surfaces
// instead of being masked by a warning that leaves the plan a silent no-op.
func TestResourceConnectionCustomerIO_MissingObjectOnReadErrors(t *testing.T) {
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
				SyncBehaviour:     iacretl.SyncBehaviourUpsert,
				Identifiers:       []iacretl.Mapping{{From: "email", To: "email"}},
				DestinationConfig: tc.destinationConfig,
			}
			svc.On("GetConnection", mock.Anything, "conn-cio-null").Return(conn, nil).Once()

			r := retl.ResourceConnectionCustomerIO()
			d := r.TestResourceData()
			d.SetId("conn-cio-null")

			diags := r.ReadContext(context.Background(), d, &rudderstack.Client{RETLSources: svc})
			require.True(t, diags.HasError(), "expected a hard error, got %v", diags)
			// Either surfacing is acceptable: a nil/empty payload fails at JSON
			// decode, a JSON `null` fails the missing-object check.
			require.Regexp(t, `no object|decode customerio destination config`, diags[0].Summary)
			svc.AssertExpectations(t)
		})
	}
}
