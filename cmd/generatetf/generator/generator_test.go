package generator_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rudderlabs/rudder-api-go/client"
	"github.com/rudderlabs/terraform-provider-rudderstack/cmd/generatetf/generator"
)

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

	data, err := generator.GenerateTerraform(sources, destinations, connections)
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
		data, err := generator.GenerateTerraform(sources, destinations, nil)
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
		data, err := generator.GenerateTerraform(sources, destinations, nil)
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
		data, err := generator.GenerateTerraform(sources, destinations, nil)
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

	data, err := generator.GenerateImportScript(sources, destinations, connections)
	fmt.Println(string(data))
	require.NoError(t, err)
	assert.Equal(t, expected, string(data))
}
