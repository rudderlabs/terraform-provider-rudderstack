package retl_test

import (
	"context"
	"fmt"
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

// JSON Mapper (no object, no customerio_audience_config) restricts
// identifiers[*].to to user_id / anonymous_id. The server enforces the rule
// with a generic apply-time error; customizeConnectionDiff surfaces it at plan
// time so users never hit the API for an obviously-invalid value.
func TestResourceConnection_JSONMapperRejectsInvalidIdentifierTarget(t *testing.T) {
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
						schedule { type = "manual" }
						identifiers {
							from = "email"
							to   = "email"
						}
					}
				`,
				ExpectError: regexpMatches(
					`identifiers\[0\]\.to must be one of \[user_id, anonymous_id\] for JSON Mapper flow \(got "email"\)`,
				),
			},
		},
	})
	svc.AssertNotCalled(t, "CreateConnection", mock.Anything, mock.Anything)
}

// Object Mapping (object set) accepts arbitrary identifier targets — the
// JSON Mapper rule must not apply to this flow.
func TestResourceConnection_ObjectMappingAcceptsArbitraryIdentifierTarget(t *testing.T) {
	svc := &mockService{}
	enabled := true
	createReq := &iacretl.CreateRETLConnectionRequest{
		SourceID:      "retl-src-1",
		DestinationID: "dest-1",
		Enabled:       &enabled,
		Schedule:      iacretl.Schedule{Type: iacretl.ScheduleTypeManual},
		SyncBehaviour: iacretl.SyncBehaviourUpsert,
		Identifiers:   []iacretl.Mapping{{From: "email", To: "email"}},
		Object:        "Contact",
	}
	created := &iacretl.RETLConnection{
		ID:            "conn-om",
		SourceID:      "retl-src-1",
		DestinationID: "dest-1",
		Enabled:       true,
		Schedule:      iacretl.Schedule{Type: iacretl.ScheduleTypeManual},
		SyncBehaviour: iacretl.SyncBehaviourUpsert,
		Identifiers:   []iacretl.Mapping{{From: "email", To: "email"}},
		Object:        "Contact",
	}
	svc.On("CreateConnection", mock.Anything, createReq).Return(created, nil).Once()
	svc.On("GetConnection", mock.Anything, "conn-om").Return(created, nil)
	svc.On("DeleteConnection", mock.Anything, "conn-om").Return(nil).Once()

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
						object         = "Contact"
						schedule { type = "manual" }
						identifiers {
							from = "email"
							to   = "email"
						}
					}
				`,
			},
		},
	})
	svc.AssertExpectations(t)
}

// Changing event.type must force destroy+create. Parent-block ForceNew alone
// doesn't catch this in terraform-plugin-sdk v2 — each nested field needs its
// own ForceNew. Different IDs across the two steps prove the resource was
// recreated rather than updated in place.
func TestResourceConnection_EventNestedFieldsAreForceNew(t *testing.T) {
	svc := &mockService{}
	enabled := true

	firstReq := &iacretl.CreateRETLConnectionRequest{
		SourceID:      "retl-src-1",
		DestinationID: "dest-1",
		Enabled:       &enabled,
		Schedule:      iacretl.Schedule{Type: iacretl.ScheduleTypeManual},
		SyncBehaviour: iacretl.SyncBehaviourUpsert,
		Identifiers:   []iacretl.Mapping{{From: "email", To: "user_id"}},
		Event:         &iacretl.Event{Type: iacretl.EventTypeIdentify},
	}
	firstCreated := &iacretl.RETLConnection{
		ID:            "conn-ev-1",
		SourceID:      "retl-src-1",
		DestinationID: "dest-1",
		Enabled:       true,
		Schedule:      iacretl.Schedule{Type: iacretl.ScheduleTypeManual},
		SyncBehaviour: iacretl.SyncBehaviourUpsert,
		Identifiers:   []iacretl.Mapping{{From: "email", To: "user_id"}},
		Event:         &iacretl.Event{Type: iacretl.EventTypeIdentify},
	}
	secondReq := &iacretl.CreateRETLConnectionRequest{
		SourceID:      "retl-src-1",
		DestinationID: "dest-1",
		Enabled:       &enabled,
		Schedule:      iacretl.Schedule{Type: iacretl.ScheduleTypeManual},
		SyncBehaviour: iacretl.SyncBehaviourUpsert,
		Identifiers:   []iacretl.Mapping{{From: "email", To: "user_id"}},
		Event:         &iacretl.Event{Type: iacretl.EventTypeTrack, Name: "user_synced"},
	}
	secondCreated := &iacretl.RETLConnection{
		ID:            "conn-ev-2",
		SourceID:      "retl-src-1",
		DestinationID: "dest-1",
		Enabled:       true,
		Schedule:      iacretl.Schedule{Type: iacretl.ScheduleTypeManual},
		SyncBehaviour: iacretl.SyncBehaviourUpsert,
		Identifiers:   []iacretl.Mapping{{From: "email", To: "user_id"}},
		Event:         &iacretl.Event{Type: iacretl.EventTypeTrack, Name: "user_synced"},
	}
	svc.On("CreateConnection", mock.Anything, firstReq).Return(firstCreated, nil).Once()
	svc.On("CreateConnection", mock.Anything, secondReq).Return(secondCreated, nil).Once()
	svc.On("GetConnection", mock.Anything, "conn-ev-1").Return(firstCreated, nil)
	svc.On("GetConnection", mock.Anything, "conn-ev-2").Return(secondCreated, nil)
	svc.On("DeleteConnection", mock.Anything, "conn-ev-1").Return(nil).Once()
	svc.On("DeleteConnection", mock.Anything, "conn-ev-2").Return(nil).Once()

	cfg := func(eventBlock string) string {
		return fmt.Sprintf(`
			provider "rudderstack" { access_token = "tok" }
			resource "rudderstack_retl_connection" "example" {
				source_id      = "retl-src-1"
				destination_id = "dest-1"
				sync_behaviour = "upsert"
				schedule { type = "manual" }
				identifiers {
					from = "email"
					to   = "user_id"
				}
				%s
			}
		`, eventBlock)
	}

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories(svc),
		Steps: []resource.TestStep{
			{
				Config: cfg("event {\n  type = \"identify\"\n}"),
				Check: func(s *terraform.State) error {
					attrs, err := resourceAttrs(s, "rudderstack_retl_connection.example")
					if err != nil {
						return err
					}
					return checkAll(map[string]string{
						"id":           "conn-ev-1",
						"event.0.type": "identify",
					}, attrs)
				},
			},
			{
				Config: cfg("event {\n  type = \"track\"\n  name = \"user_synced\"\n}"),
				Check: func(s *terraform.State) error {
					attrs, err := resourceAttrs(s, "rudderstack_retl_connection.example")
					if err != nil {
						return err
					}
					// New ID proves destroy+create: nested-field ForceNew
					// promoted the in-place update to a replacement.
					return checkAll(map[string]string{
						"id":           "conn-ev-2",
						"event.0.type": "track",
						"event.0.name": "user_synced",
					}, attrs)
				},
			},
		},
	})
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

// Refresh against an API response that has destination-specific config must
// fail loudly on the generic resource — otherwise users with an existing
// customerio_audience connection imported under the generic type would silently
// lose their audience config from state.
func TestResourceConnection_RefusesDestinationSpecificConfig(t *testing.T) {
	svc := &mockService{}
	conn := &iacretl.RETLConnection{
		ID:                "conn-bad",
		SourceID:          "retl-src-1",
		DestinationID:     "dest-cio",
		Enabled:           true,
		Schedule:          iacretl.Schedule{Type: iacretl.ScheduleTypeManual},
		SyncBehaviour:     iacretl.SyncBehaviourMirror,
		Identifiers:       []iacretl.Mapping{{From: "email", To: "email"}},
		DestinationConfig: []byte(`{"audienceId":42}`),
	}
	svc.On("GetConnection", mock.Anything, "conn-bad").Return(conn, nil).Once()

	r := retl.ResourceConnection()
	d := r.TestResourceData()
	d.SetId("conn-bad")

	diags := r.ReadContext(context.Background(), d, &rudderstack.Client{RETLSources: svc})
	require.True(t, diags.HasError(), "expected refresh to error on destination-specific config")
	// Match the stable pattern in the message rather than a specific example
	// destination — the example roster will grow as we add typed resources.
	require.Regexp(t,
		`rudderstack_retl_connection_<destination>`,
		diags[0].Summary,
	)
	svc.AssertExpectations(t)
}

// Non-null, non-empty destinationConfig must never silently slip through the
// generic resource. Any payload that decodes to a non-nil value is a
// destination-specific flow that this resource can't represent; any payload
// that fails to decode is a malformed signal worth surfacing loudly. Both
// cases produce an error diagnostic, but with different shapes so the user
// sees the underlying cause.
func TestResourceConnection_DestinationConfigOnRead_NonNullAndMalformed(t *testing.T) {
	cases := []struct {
		name              string
		destinationConfig []byte
		wantRegex         string
	}{
		// Malformed payloads — surface a decode error.
		{name: "truncated object", destinationConfig: []byte(`{"broken`), wantRegex: `decode destinationConfig`},
		{name: "whitespace only", destinationConfig: []byte(`   `), wantRegex: `decode destinationConfig`},
		// Well-formed but non-null payloads of any shape — typed-resource signal.
		{name: "empty object", destinationConfig: []byte(`{}`), wantRegex: `rudderstack_retl_connection_<destination>`},
		{name: "array", destinationConfig: []byte(`[]`), wantRegex: `rudderstack_retl_connection_<destination>`},
		{name: "scalar number", destinationConfig: []byte(`42`), wantRegex: `rudderstack_retl_connection_<destination>`},
		{name: "scalar string", destinationConfig: []byte(`"hello"`), wantRegex: `rudderstack_retl_connection_<destination>`},
		{name: "scalar bool", destinationConfig: []byte(`true`), wantRegex: `rudderstack_retl_connection_<destination>`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc := &mockService{}
			conn := &iacretl.RETLConnection{
				ID:                "conn-x",
				SourceID:          "retl-src-1",
				DestinationID:     "dest-1",
				Enabled:           true,
				Schedule:          iacretl.Schedule{Type: iacretl.ScheduleTypeManual},
				SyncBehaviour:     iacretl.SyncBehaviourFull,
				Identifiers:       []iacretl.Mapping{{From: "email", To: "user_id"}},
				DestinationConfig: tc.destinationConfig,
			}
			svc.On("GetConnection", mock.Anything, "conn-x").Return(conn, nil).Once()

			r := retl.ResourceConnection()
			d := r.TestResourceData()
			d.SetId("conn-x")

			diags := r.ReadContext(context.Background(), d, &rudderstack.Client{RETLSources: svc})
			require.True(t, diags.HasError(), "expected refresh to error on destinationConfig=%s", string(tc.destinationConfig))
			require.Regexp(t, tc.wantRegex, diags[0].Summary)
			svc.AssertExpectations(t)
		})
	}
}

// Regression: JSON-`null` destinationConfig in the API response must NOT error
// on the generic resource — `null` is the server's "no destination-specific
// config" signal and should be treated as a no-op.
func TestResourceConnection_NullDestinationConfigOnRead(t *testing.T) {
	svc := &mockService{}
	conn := &iacretl.RETLConnection{
		ID:                "conn-null",
		SourceID:          "retl-src-1",
		DestinationID:     "dest-1",
		Enabled:           true,
		Schedule:          iacretl.Schedule{Type: iacretl.ScheduleTypeManual},
		SyncBehaviour:     iacretl.SyncBehaviourFull,
		Identifiers:       []iacretl.Mapping{{From: "email", To: "user_id"}},
		DestinationConfig: []byte(`null`),
	}
	svc.On("GetConnection", mock.Anything, "conn-null").Return(conn, nil).Once()

	r := retl.ResourceConnection()
	d := r.TestResourceData()
	d.SetId("conn-null")

	diags := r.ReadContext(context.Background(), d, &rudderstack.Client{RETLSources: svc})
	require.False(t, diags.HasError(), "diags=%v", diags)
	svc.AssertExpectations(t)
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
