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

// The typed Customer.io (VDM v2) resource exposes `object` as a top-level
// string; internally it round-trips through destinationConfig as
// {"object": "..."}. identifiers and mappings flow through the base schema
// (VDM v2 identifierMappings / fieldMappings); config-be assembles the VDM v2
// connectionConfig from the destination definition.
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
		Mappings:          []iacretl.Mapping{{From: "name", To: "plan"}},
		DestinationConfig: json.RawMessage(`{"object":"customers"}`),
	}
	created := &iacretl.RETLConnection{
		ID:                "conn-cio",
		SourceID:          "retl-src-1",
		DestinationID:     "dest-cio",
		Enabled:           true,
		Schedule:          iacretl.Schedule{Type: iacretl.ScheduleTypeManual},
		SyncBehaviour:     iacretl.SyncBehaviourUpsert,
		Identifiers:       []iacretl.Mapping{{From: "email", To: "email"}},
		Mappings:          []iacretl.Mapping{{From: "name", To: "plan"}},
		DestinationConfig: json.RawMessage(`{"object":"customers"}`),
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
						object         = "customers"
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
				Check: func(s *terraform.State) error {
					attrs, err := resourceAttrs(s, "rudderstack_retl_connection_customerio.example")
					if err != nil {
						return err
					}
					return checkAll(map[string]string{
						"id":              "conn-cio",
						"object":          "customers",
						"mappings.0.from": "name",
						"mappings.0.to":   "plan",
					}, attrs)
				},
			},
		},
	})

	svc.AssertExpectations(t)
}

// Changing `object` must recreate the connection (ForceNew) — the
// destinationConfig shape is not mutable in place on this flow. A new
// resource ID after the value change proves destroy + create; terraform would
// have called UpdateConnection (which the mock does not register) otherwise.
func TestResourceConnectionCustomerIO_ObjectIsForceNew(t *testing.T) {
	svc := &mockService{}
	enabled := true

	mkReq := func(object string) *iacretl.CreateRETLConnectionRequest {
		return &iacretl.CreateRETLConnectionRequest{
			SourceID:          "retl-src-1",
			DestinationID:     "dest-cio",
			Enabled:           &enabled,
			Schedule:          iacretl.Schedule{Type: iacretl.ScheduleTypeManual},
			SyncBehaviour:     iacretl.SyncBehaviourUpsert,
			Identifiers:       []iacretl.Mapping{{From: "email", To: "email"}},
			DestinationConfig: json.RawMessage(fmt.Sprintf(`{"object":%q}`, object)),
		}
	}
	mkConn := func(id, object string) *iacretl.RETLConnection {
		return &iacretl.RETLConnection{
			ID:                id,
			SourceID:          "retl-src-1",
			DestinationID:     "dest-cio",
			Enabled:           true,
			Schedule:          iacretl.Schedule{Type: iacretl.ScheduleTypeManual},
			SyncBehaviour:     iacretl.SyncBehaviourUpsert,
			Identifiers:       []iacretl.Mapping{{From: "email", To: "email"}},
			DestinationConfig: json.RawMessage(fmt.Sprintf(`{"object":%q}`, object)),
		}
	}
	first := mkConn("conn-cio-1", "customers")
	second := mkConn("conn-cio-2", "accounts")
	svc.On("CreateConnection", mock.Anything, mkReq("customers")).Return(first, nil).Once()
	svc.On("CreateConnection", mock.Anything, mkReq("accounts")).Return(second, nil).Once()
	svc.On("GetConnection", mock.Anything, "conn-cio-1").Return(first, nil)
	svc.On("GetConnection", mock.Anything, "conn-cio-2").Return(second, nil)
	svc.On("DeleteConnection", mock.Anything, "conn-cio-1").Return(nil).Once()
	svc.On("DeleteConnection", mock.Anything, "conn-cio-2").Return(nil).Once()

	cfg := func(object string) string {
		return fmt.Sprintf(`
			provider "rudderstack" { access_token = "tok" }
			resource "rudderstack_retl_connection_customerio" "example" {
				source_id      = "retl-src-1"
				destination_id = "dest-cio"
				sync_behaviour = "upsert"
				object         = %q
				schedule { type = "manual" }
				identifiers {
					from = "email"
					to   = "email"
				}
			}
		`, object)
	}

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories(svc),
		Steps: []resource.TestStep{
			{
				Config: cfg("customers"),
				Check: func(s *terraform.State) error {
					attrs, err := resourceAttrs(s, "rudderstack_retl_connection_customerio.example")
					if err != nil {
						return err
					}
					return checkAll(map[string]string{"id": "conn-cio-1", "object": "customers"}, attrs)
				},
			},
			{
				Config: cfg("accounts"),
				Check: func(s *terraform.State) error {
					attrs, err := resourceAttrs(s, "rudderstack_retl_connection_customerio.example")
					if err != nil {
						return err
					}
					// New ID proves destroy + create — ForceNew fired.
					return checkAll(map[string]string{"id": "conn-cio-2", "object": "accounts"}, attrs)
				},
			},
		},
	})

	svc.AssertExpectations(t)
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
						object         = "customers"
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
		DestinationConfig: json.RawMessage(`{"object":"customers"}`),
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
		DestinationConfig: json.RawMessage(`{"object":"customers"}`),
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
						object         = "customers"
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
						"object":        "customers",
						"cursor_column": "updated_at",
					}, attrs)
				},
			},
		},
	})

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
						object         = "customers"
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

// A JSON-`null` or empty destinationConfig is a server-side inconsistency for
// a Customer.io VDM v2 connection (object is mandatory). Refresh must emit a
// Warning and NOT zero `object` — zeroing would produce a never-reconcilable
// plan because `object` is ForceNew.
func TestResourceConnectionCustomerIO_NullDestinationConfigOnRead(t *testing.T) {
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
			// Seed prior object so we can prove refresh leaves it untouched.
			require.NoError(t, d.Set("object", "customers"))

			diags := r.ReadContext(context.Background(), d, &rudderstack.Client{RETLSources: svc})
			require.False(t, diags.HasError(), "diags=%v", diags)
			require.Len(t, diags, 1, "expected one warning diagnostic")
			require.Equal(t, diag.Warning, diags[0].Severity)
			require.Regexp(t, `missing object`, diags[0].Summary)
			require.Equal(t, "customers", d.Get("object"),
				"prior object must be preserved to avoid a never-reconcilable plan")
			svc.AssertExpectations(t)
		})
	}
}
