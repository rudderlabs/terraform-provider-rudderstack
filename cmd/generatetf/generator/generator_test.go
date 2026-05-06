package generator_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rudderlabs/rudder-iac/api/client"
	"github.com/rudderlabs/rudder-iac/api/client/retl"
	"github.com/rudderlabs/terraform-provider-rudderstack/cmd/generatetf/generator"
)

func intPtr(v int) *int    { return &v }
func boolPtr(v bool) *bool { return &v }

func TestGeneratorTerraform(t *testing.T) {
	sources := []client.Source{
		{
			ID:   "id-javascript",
			Name: "source-1",
			Type: "Javascript",
		},
		{
			ID:   "id-http",
			Name: "source-2",
			Type: "HTTP",
		},
		{
			ID:   "unknown",
			Name: "source-unknown",
			Type: "Unknown",
		},
	}

	destinations := []client.Destination{
		{
			ID:   "id-redshift",
			Name: "name-redshift",
			Type: "RS",
			Config: json.RawMessage(`{
				"host": "example.com",
				"port": "5439",
				"user": "example-user",
				"password": "example-password",
				"database": "example-database",
				"namespace": "example-namespace",
				"enableSSE": true,
				"useRudderStorage": false,
				"unknown": "unknown value",
				"whiteListedEvents": [],
				"blacklistedEvents": []
			}`),
		},
		{
			ID:   "id-facebook-pixel",
			Name: "name-facebook-pixel",
			Type: "FACEBOOK_PIXEL",
			Config: json.RawMessage(`{
				"pixelId": "facebook pixel id",
				"accessToken": "facebook access token",
				"standardPageCall": true,
				"valueFieldIdentifier": "properties.price",
				"advancedMapping": true,
				"testDestination": true,
				"testEventCode": "...",
				"eventsToEvents": [
				  { "from": "a1", "to": "b1" },
				  { "from": "a2", "to": "b2" }
				],
				"legacyConversionPixelId": { "from": "from", "to": "to" },
				"useNativeSDK": { "web": true },
				"consentManagement": {
					"web": [
						{
							"provider": "oneTrust",
							"consents": [
								{
									"consent": "one_web"
								},
								{
									"consent": "two_web"
								},
								{
									"consent": "three_web"
								}
							]
						},
						{
							"provider": "ketch",
							"consents": [
								{
									"consent": "one_web"
								},
								{
									"consent": "two_web"
								},
								{
									"consent": "three_web"
								}
							]
						},
						{
							"provider": "custom",
							"resolutionStrategy": "and",
							"consents": [
								{
									"consent": "one_web"
								},
								{
									"consent": "two_web"
								},
								{
									"consent": "three_web"
								}
							]
						}
					]
				},
				"blacklistedEvents": [
				  { "eventName": "one" },
				  { "eventName": "two" },
				  { "eventName": "three" }
				],
				"eventFilteringOption": "blacklistedEvents",
				"blacklistPiiProperties": [],
				"whiteListedEvents": []
			}`),
		},
		{
			ID:     "unknown",
			Name:   "destination-unknown",
			Type:   "Unknown",
			Config: json.RawMessage("{}"),
		},
	}

	connections := []client.Connection{
		{
			ID:            "id-connection-1",
			SourceID:      "id-javascript",
			DestinationID: "id-redshift",
		},
		{
			ID:            "id-connection-2",
			SourceID:      "id-http",
			DestinationID: "id-facebook-pixel",
		},
	}

	expected := `
resource "rudderstack_source_javascript" "src_id-javascript" {
  name = "source-1"
}

resource "rudderstack_source_http" "src_id-http" {
  name = "source-2"
}

resource "rudderstack_destination_redshift" "dst_id-redshift" {
  name = "name-redshift"
  config {
    database           = "example-database"
    enable_sse         = true
    host               = "example.com"
    namespace          = "example-namespace"
    password           = "example-password"
    port               = "5439"
    use_rudder_storage = false
    user               = "example-user"
  }
}

resource "rudderstack_destination_facebook_pixel" "dst_id-facebook-pixel" {
  name = "name-facebook-pixel"
  config {
    access_token             = "facebook access token"
    advanced_mapping         = true
    blacklist_pii_properties = []
    consent_management {
      web = [{
        consents            = ["one_web", "two_web", "three_web"]
        provider            = "oneTrust"
        resolution_strategy = ""
        }, {
        consents            = ["one_web", "two_web", "three_web"]
        provider            = "ketch"
        resolution_strategy = ""
        }, {
        consents            = ["one_web", "two_web", "three_web"]
        provider            = "custom"
        resolution_strategy = "and"
      }]
    }
    event_filtering {
      blacklist = ["one", "two", "three"]
    }
    events_to_events = [{
      from = "a1"
      to   = "b1"
      }, {
      from = "a2"
      to   = "b2"
    }]
    legacy_conversion_pixel_id {
      from = "from"
      to   = "to"
    }
    pixel_id           = "facebook pixel id"
    standard_page_call = true
    test_destination   = true
    test_event_code    = "..."
    use_native_sdk {
      web = true
    }
    value_field_identifier = "properties.price"
  }
}

resource "rudderstack_connection" "cnxn_id-connection-1" {
  source_id      = rudderstack_source_javascript.src_id-javascript.id
  destination_id = rudderstack_destination_redshift.dst_id-redshift.id
}

resource "rudderstack_connection" "cnxn_id-connection-2" {
  source_id      = rudderstack_source_http.src_id-http.id
  destination_id = rudderstack_destination_facebook_pixel.dst_id-facebook-pixel.id
}
`

	// trim new lines in expected var and add newlines at the end, generated by the Generator
	expected = fmt.Sprintf("%s\n\n", strings.Trim(expected, "\n"))

	data, err := generator.GenerateTerraform(sources, destinations, connections, nil, nil)
	require.NoError(t, err)
	assert.Equal(t, expected, string(data))
}

// TestGeneratorTerraformNestedListInListOfObjects tests ctyValue handling for a list of objects
// where each object contains a nested list attribute. Three scenarios are covered:
//   - all items have an empty nested list
//   - all items have a non-empty nested list
//   - some items have an empty nested list and some have a non-empty nested list (the mixed case)
//
// The mixed case previously caused a panic ("inconsistent list element types") because empty
// lists were represented as List(String) while non-empty lists had a concrete element type,
// making it impossible to construct a uniform cty.ListVal.
func TestGeneratorTerraformNestedListInListOfObjects(t *testing.T) {
	// HS destination is used because it has hubspot_events (list of objects) where each
	// object contains event_properties (nested list), which is the structure that exercises
	// the ctyValue fix. The test is about the generator behavior, not the HS destination itself.
	sources := []client.Source{
		{ID: "id-js", Name: "source-js", Type: "Javascript"},
	}

	t.Run("all items have empty nested list", func(t *testing.T) {
		destinations := []client.Destination{
			{
				ID:   "id-dst-1",
				Name: "dest-1",
				Type: "HS",
				Config: json.RawMessage(`{
					"authorizationType": "newPrivateAppApi",
					"apiVersion": "newApi",
					"hubspotEvents": [
						{"rsEventName": "e1", "hubspotEventName": "pe_e1", "eventProperties": []},
						{"rsEventName": "e2", "hubspotEventName": "pe_e2", "eventProperties": []}
					]
				}`),
			},
		}
		data, err := generator.GenerateTerraform(sources, destinations, nil, nil, nil)
		require.NoError(t, err)
		output := string(data)
		assert.Contains(t, output, `hubspot_event_name = "pe_e1"`)
		assert.Contains(t, output, `hubspot_event_name = "pe_e2"`)
	})

	t.Run("all items have non-empty nested list", func(t *testing.T) {
		destinations := []client.Destination{
			{
				ID:   "id-dst-2",
				Name: "dest-2",
				Type: "HS",
				Config: json.RawMessage(`{
					"authorizationType": "newPrivateAppApi",
					"apiVersion": "newApi",
					"hubspotEvents": [
						{
							"rsEventName": "e1", "hubspotEventName": "pe_e1",
							"eventProperties": [{"from": "f1", "to": "t1"}]
						},
						{
							"rsEventName": "e2", "hubspotEventName": "pe_e2",
							"eventProperties": [{"from": "f2", "to": "t2"}]
						}
					]
				}`),
			},
		}
		data, err := generator.GenerateTerraform(sources, destinations, nil, nil, nil)
		require.NoError(t, err)
		output := string(data)
		assert.Contains(t, output, `from = "f1"`)
		assert.Contains(t, output, `from = "f2"`)
	})

	t.Run("mixed: some items have empty nested list and some have non-empty", func(t *testing.T) {
		destinations := []client.Destination{
			{
				ID:   "id-dst-3",
				Name: "dest-3",
				Type: "HS",
				Config: json.RawMessage(`{
					"authorizationType": "newPrivateAppApi",
					"apiVersion": "newApi",
					"hubspotEvents": [
						{"rsEventName": "e1", "hubspotEventName": "pe_e1", "eventProperties": []},
						{
							"rsEventName": "e2", "hubspotEventName": "pe_e2",
							"eventProperties": [{"from": "f2", "to": "t2"}]
						}
					]
				}`),
			},
		}
		data, err := generator.GenerateTerraform(sources, destinations, nil, nil, nil)
		require.NoError(t, err)
		output := string(data)
		assert.Contains(t, output, `hubspot_event_name = "pe_e1"`)
		assert.Contains(t, output, `hubspot_event_name = "pe_e2"`)
		assert.Contains(t, output, `from = "f2"`)
	})
}

func TestGeneratorImportScript(t *testing.T) {
	sources := []client.Source{
		{
			ID:   "id-source-1",
			Name: "source-1",
			Type: "Javascript",
		},
		{
			ID:   "id-source-2",
			Name: "source-2",
			Type: "HTTP",
		},
		{
			ID:   "unknown",
			Name: "source-unknown",
			Type: "Unknown",
		},
	}

	destinations := []client.Destination{
		{
			ID:   "id-destination-1",
			Name: "name-redshift",
			Type: "RS",
		},
		{
			ID:   "id-destination-2",
			Name: "name-facebook-pixel",
			Type: "FACEBOOK_PIXEL",
		},
		{
			ID:   "unknown",
			Name: "destination-unknown",
			Type: "Unknown",
		},
	}

	connections := []client.Connection{
		{
			ID:            "id-connection-1",
			SourceID:      "id-source-1",
			DestinationID: "id-destination-1",
		},
		{
			ID:            "id-connection-2",
			SourceID:      "id-source-2",
			DestinationID: "id-destination-2",
		},
		{
			ID:            "id-connection-3",
			SourceID:      "unknown",
			DestinationID: "id-destination-2",
		},
		{
			ID:            "id-connection-4",
			SourceID:      "id-source-1",
			DestinationID: "non-existing",
		},
	}

	expected := `
terraform import "rudderstack_source_javascript.src_id-source-1" "id-source-1"
terraform import "rudderstack_source_http.src_id-source-2" "id-source-2"
terraform import "rudderstack_destination_redshift.dst_id-destination-1" "id-destination-1"
terraform import "rudderstack_destination_facebook_pixel.dst_id-destination-2" "id-destination-2"
terraform import "rudderstack_connection.cnxn_id-connection-1" "id-connection-1"
terraform import "rudderstack_connection.cnxn_id-connection-2" "id-connection-2"
`

	// trim new line at the start of expected var
	expected = strings.Trim(expected, "\n")

	data, err := generator.GenerateImportScript(sources, destinations, connections, nil, nil)
	require.NoError(t, err)
	assert.Equal(t, expected, string(data))
}

// retlFixtures returns a representative slice of RETL sources + connections
// used by both the HCL and import-script tests:
//   - one model source (snowflake)
//   - one warehouse table source (postgres)
//   - one S3 table source (must be skipped — no _s3_table generator)
//   - one source with an unsupported sourceType (must be skipped)
//   - one JSON-Mapper connection wired to the model source + redshift destination
//   - one destination-specific connection (with destination_config) wired to the
//     table source + facebook-pixel destination
//   - one connection whose RETL source is the unsupported one (must be skipped)
//   - one connection whose destination doesn't exist in the generated set (must be skipped)
func retlFixtures() ([]retl.RETLSource, []retl.RETLConnection) {
	sources := []retl.RETLSource{
		{
			ID:                   "src-model-1",
			Name:                 "users-model",
			IsEnabled:            true,
			SourceType:           retl.ModelSourceType,
			SourceDefinitionName: "snowflake",
			AccountID:            "acc-snow-1",
			Config: retl.RETLSQLModelConfig{
				PrimaryKey:  "user_id",
				Sql:         "SELECT * FROM users",
				Description: "All users",
			},
		},
		{
			ID:                   "src-table-1",
			Name:                 "events-table",
			IsEnabled:            true,
			SourceType:           retl.TableSourceType,
			SourceDefinitionName: "postgres",
			AccountID:            "acc-pg-1",
			Config: retl.RETLTableConfig{
				PrimaryKey: "event_id",
				Schema:     "public",
				Table:      "events",
			},
		},
		{
			// S3 table — must be skipped (no _s3_table generator in this PR).
			ID:                   "src-s3-1",
			Name:                 "s3-events",
			IsEnabled:            true,
			SourceType:           retl.TableSourceType,
			SourceDefinitionName: "s3",
		},
		{
			// Audience / profiles / etc. — unsupported sourceType, must be skipped.
			ID:         "src-unsupported-1",
			Name:       "audience-1",
			IsEnabled:  true,
			SourceType: retl.SourceType("audience"),
		},
	}

	connections := []retl.RETLConnection{
		{
			ID:            "cnxn-jm-1",
			SourceID:      "src-model-1",
			DestinationID: "id-redshift",
			Enabled:       true,
			Schedule: retl.Schedule{
				Type:         retl.ScheduleTypeBasic,
				EveryMinutes: intPtr(60),
			},
			SyncSettings: &retl.SyncSettings{
				SyncLogsConfig: &retl.SyncLogsConfig{
					Enabled:            boolPtr(true),
					LogRetentionInDays: intPtr(30),
					SnapshotsToRetain:  intPtr(5),
				},
				FailedKeysConfig: &retl.FailedKeysConfig{
					EnableFailedKeysRetry: boolPtr(true),
				},
			},
			SyncBehaviour: retl.SyncBehaviourUpsert,
			Identifiers: []retl.Mapping{
				{From: "email", To: "user_id"},
			},
			Mappings: []retl.Mapping{
				{From: "first_name", To: "fname"},
			},
			Event:        &retl.Event{Type: retl.EventTypeIdentify},
			Constants:    []retl.Constant{{Key: "source", Value: "rudderstack"}},
			CursorColumn: "updated_at",
		},
		{
			ID:            "cnxn-ds-1",
			SourceID:      "src-table-1",
			DestinationID: "id-facebook-pixel",
			Enabled:       true,
			Schedule: retl.Schedule{
				Type:         retl.ScheduleTypeBasic,
				EveryMinutes: intPtr(120),
			},
			SyncBehaviour: retl.SyncBehaviourMirror,
			Identifiers: []retl.Mapping{
				{From: "email_col", To: "EMAIL"},
			},
			DestinationConfig: json.RawMessage(`{"audienceId":"segment_abc123"}`),
		},
		{
			// Connection whose RETL source is the unsupported one — must be skipped.
			ID:            "cnxn-skip-source",
			SourceID:      "src-unsupported-1",
			DestinationID: "id-redshift",
			Schedule:      retl.Schedule{Type: retl.ScheduleTypeManual},
			SyncBehaviour: retl.SyncBehaviourUpsert,
			Identifiers:   []retl.Mapping{{From: "x", To: "y"}},
		},
		{
			// Connection referencing a non-existent destination — must be skipped.
			ID:            "cnxn-skip-dest",
			SourceID:      "src-model-1",
			DestinationID: "does-not-exist",
			Schedule:      retl.Schedule{Type: retl.ScheduleTypeManual},
			SyncBehaviour: retl.SyncBehaviourUpsert,
			Identifiers:   []retl.Mapping{{From: "x", To: "y"}},
		},
	}
	return sources, connections
}

// retlEventStreamingFixtures returns the regular sources/destinations needed
// by retlFixtures (RETL connections reuse event-streaming destinations).
func retlEventStreamingFixtures() ([]client.Source, []client.Destination) {
	return []client.Source{},
		[]client.Destination{
			{
				ID:     "id-redshift",
				Name:   "name-redshift",
				Type:   "RS",
				Config: json.RawMessage(`{}`),
			},
			{
				ID:     "id-facebook-pixel",
				Name:   "name-facebook-pixel",
				Type:   "FACEBOOK_PIXEL",
				Config: json.RawMessage(`{}`),
			},
		}
}

func TestGeneratorTerraform_RETL(t *testing.T) {
	esSources, esDestinations := retlEventStreamingFixtures()
	retlSources, retlConnections := retlFixtures()

	data, err := generator.GenerateTerraform(esSources, esDestinations, nil, retlSources, retlConnections)
	require.NoError(t, err)
	output := string(data)

	// supported RETL sources are emitted with the expected resource type + name
	assert.Contains(t, output, `resource "rudderstack_retl_source_model" "retl_src_src-model-1"`)
	assert.Contains(t, output, `resource "rudderstack_retl_source_table" "retl_src_src-table-1"`)

	// model source has its config block populated
	assert.Contains(t, output, `primary_key = "user_id"`)
	assert.Contains(t, output, `sql         = "SELECT * FROM users"`)
	assert.Contains(t, output, `description = "All users"`)
	assert.Contains(t, output, `source_definition_name = "snowflake"`)
	assert.Contains(t, output, `account_id             = "acc-snow-1"`)

	// table source has its config block populated
	assert.Contains(t, output, `schema      = "public"`)
	assert.Contains(t, output, `table       = "events"`)

	// S3 and unsupported sources are skipped (no resource block, no config).
	assert.NotContains(t, output, `retl_src_src-s3-1`)
	assert.NotContains(t, output, `retl_src_src-unsupported-1`)
	assert.NotContains(t, output, `rudderstack_retl_source_s3_table`)

	// JSON-Mapper connection block is emitted
	assert.Contains(t, output, `resource "rudderstack_retl_connection" "retl_cnxn_cnxn-jm-1"`)
	// cross-references the RETL source via terraform reference (not the API ID)
	assert.Contains(t, output, `source_id      = rudderstack_retl_source_model.retl_src_src-model-1.id`)
	assert.Contains(t, output, `destination_id = rudderstack_destination_redshift.dst_id-redshift.id`)
	assert.Contains(t, output, `sync_behaviour = "upsert"`)
	// cursor_column is in its own attribute group (after blocks), so it
	// gets single-space alignment rather than aligning with sync_behaviour.
	assert.Contains(t, output, `cursor_column = "updated_at"`)

	// schedule + sync_settings + identifiers + mappings + constants + event blocks
	assert.Contains(t, output, `every_minutes = 60`)
	assert.Contains(t, output, `log_retention_in_days = 30`)
	assert.Contains(t, output, `enable_failed_keys_retry = true`)
	assert.Contains(t, output, `from = "email"`)
	assert.Contains(t, output, `to   = "user_id"`)
	assert.Contains(t, output, `from = "first_name"`)
	assert.Contains(t, output, `key   = "source"`)
	assert.Contains(t, output, `value = "rudderstack"`)
	assert.Contains(t, output, `type = "identify"`)

	// destination-specific connection: destination_config rendered as jsonencode(...)
	assert.Contains(t, output, `resource "rudderstack_retl_connection" "retl_cnxn_cnxn-ds-1"`)
	assert.Contains(t, output, `destination_id = rudderstack_destination_facebook_pixel.dst_id-facebook-pixel.id`)
	// destination_config sits in its own attribute group after the blocks, so
	// sync_behaviour is aligned with the top group (destination_id, 14 chars).
	assert.Contains(t, output, `sync_behaviour = "mirror"`)
	assert.Contains(t, output, `destination_config = jsonencode({`)
	assert.Contains(t, output, `audienceId = "segment_abc123"`)

	// skipped connections aren't emitted
	assert.NotContains(t, output, `retl_cnxn_cnxn-skip-source`)
	assert.NotContains(t, output, `retl_cnxn_cnxn-skip-dest`)
}

func TestGeneratorImportScript_RETL(t *testing.T) {
	esSources, esDestinations := retlEventStreamingFixtures()
	retlSources, retlConnections := retlFixtures()

	// Full golden-string assertion (matching the style of
	// TestGeneratorImportScript): order is event-streaming first, then RETL
	// sources, then RETL connections. Skipped resources (s3 source,
	// unsupported source-type, connections referencing them) must not appear.
	expected := strings.Join([]string{
		`terraform import "rudderstack_destination_redshift.dst_id-redshift" "id-redshift"`,
		`terraform import "rudderstack_destination_facebook_pixel.dst_id-facebook-pixel" "id-facebook-pixel"`,
		`terraform import "rudderstack_retl_source_model.retl_src_src-model-1" "src-model-1"`,
		`terraform import "rudderstack_retl_source_table.retl_src_src-table-1" "src-table-1"`,
		`terraform import "rudderstack_retl_connection.retl_cnxn_cnxn-jm-1" "cnxn-jm-1"`,
		`terraform import "rudderstack_retl_connection.retl_cnxn_cnxn-ds-1" "cnxn-ds-1"`,
	}, "\n")

	data, err := generator.GenerateImportScript(esSources, esDestinations, nil, retlSources, retlConnections)
	require.NoError(t, err)
	assert.Equal(t, expected, string(data))
}

// TestGeneratorTerraform_RETL_DestinationConfigEdgeCases verifies that
// connections whose `destinationConfig` is JSON null or a non-object (string,
// number, array) are still emitted but with the `destination_config` attribute
// suppressed — the rest of the connection is usable. This guards the loose
// API contract where the field is an opaque blob.
func TestGeneratorTerraform_RETL_DestinationConfigEdgeCases(t *testing.T) {
	_, esDestinations := retlEventStreamingFixtures()
	retlSources := []retl.RETLSource{
		{
			ID: "src-1", Name: "src-1", IsEnabled: true,
			SourceType: retl.ModelSourceType, SourceDefinitionName: "snowflake",
			AccountID: "acc-1",
			Config:    retl.RETLSQLModelConfig{PrimaryKey: "id", Sql: "SELECT 1"},
		},
	}
	mkConn := func(id string, destCfg json.RawMessage) retl.RETLConnection {
		return retl.RETLConnection{
			ID: id, SourceID: "src-1", DestinationID: "id-redshift", Enabled: true,
			Schedule:          retl.Schedule{Type: retl.ScheduleTypeManual},
			SyncBehaviour:     retl.SyncBehaviourUpsert,
			Identifiers:       []retl.Mapping{{From: "x", To: "y"}},
			DestinationConfig: destCfg,
		}
	}
	retlConnections := []retl.RETLConnection{
		mkConn("cnxn-null", json.RawMessage(`null`)),
		mkConn("cnxn-string", json.RawMessage(`"opaque-token"`)),
		mkConn("cnxn-array", json.RawMessage(`[1,2,3]`)),
	}

	data, err := generator.GenerateTerraform(nil, esDestinations, nil, retlSources, retlConnections)
	require.NoError(t, err)
	output := string(data)

	// All three connection blocks should still be present, with their other fields intact.
	assert.Contains(t, output, `"retl_cnxn_cnxn-null"`)
	assert.Contains(t, output, `"retl_cnxn_cnxn-string"`)
	assert.Contains(t, output, `"retl_cnxn_cnxn-array"`)
	// But destination_config must NOT have been emitted for any of them.
	assert.NotContains(t, output, "destination_config")
}
