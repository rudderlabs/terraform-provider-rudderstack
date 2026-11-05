package destinations_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil"
	acc "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/acc"
	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

var httpV2TestConfigs = []c.TestConfig{
	{
		TerraformCreate: `
			api_url = "https://example.com/base"
			request_body_template = "[$$.events.{\"customer\": {\"id\": userId}, \"items\": properties.items}]"
			supported_message_types = ["track", "identify"]
			max_batch_size = "1"
			headers = [{ to = "x-customer", from = "$$.message.userId" }]
			query_params = [{ to = "environment", from = "production" }]
			path_params = [{ path = "$$.message.properties.region" }]
		`,
		APICreate: `{
			"apiUrl": "https://example.com/base",
			"auth": "noAuth", "method": "POST", "format": "JSON",
			"requestBodyTemplate": "[$$.events.{\"customer\": {\"id\": userId}, \"items\": properties.items}]",
			"supportedMessageTypes": ["track", "identify"],
			"maxBatchSize": "1",
			"headers": [{"to": "x-customer", "from": "$$.message.userId"}],
			"queryParams": [{"to": "environment", "from": "production"}],
			"pathParams": [{"path": "$$.message.properties.region"}]
		}`,
		TerraformUpdate: `
			api_url = "https://example.com/base"
			auth = "bearerTokenAuth"
			bearer_token = "my-token"
			method = "PATCH"
			request_body_template = "[$$.events.{\"profile\": {\"id\": userId}, \"orders\": properties.items}]"
			supported_message_types = ["track"]
			max_batch_size = "5000"
			headers = [{ to = "x-customer", from = "$$.message.userId" }]
			query_params = [{ to = "environment", from = "production" }]
			path_params = [{ path = "$$.message.properties.region" }]
		`,
		APIUpdate: `{
			"apiUrl": "https://example.com/base",
			"auth": "bearerTokenAuth", "bearerToken": "my-token",
			"method": "PATCH", "format": "JSON",
			"requestBodyTemplate": "[$$.events.{\"profile\": {\"id\": userId}, \"orders\": properties.items}]",
			"supportedMessageTypes": ["track"],
			"maxBatchSize": "5000",
			"headers": [{"to": "x-customer", "from": "$$.message.userId"}],
			"queryParams": [{"to": "environment", "from": "production"}],
			"pathParams": [{"path": "$$.message.properties.region"}]
		}`,
	},
}

func TestDestinationResourceHttpV2(t *testing.T) {
	cmt.AssertDestination(t, "http_v2", httpV2TestConfigs)
}

func TestAccDestinationHttpV2(t *testing.T) {
	acc.AccAssertDestination(t, "http_v2", httpV2TestConfigs)
}

func TestDestinationResourceHttpV2OptionalFields(t *testing.T) {
	cm := c.Destinations.Entries()["http_v2"]
	got, err := cm.StateToAPI(`{"api_url":"https://example.com/base","auth":"noAuth","method":"POST","format":"JSON","request_body_template":"[$$.events]","max_batch_size":"1"}`)
	require.NoError(t, err)
	assert.True(t, testutil.JSONEq(got, `{"apiUrl":"https://example.com/base","auth":"noAuth","method":"POST","format":"JSON","requestBodyTemplate":"[$$.events]","maxBatchSize":"1"}`), got)

	got, err = cm.APIToState(`{"apiUrl":"https://example.com/base","auth":"noAuth","method":"POST","format":"JSON","requestBodyTemplate":"[$$.events]","maxBatchSize":"1"}`)
	require.NoError(t, err)
	assert.True(t, testutil.JSONEq(got, `{"api_url":"https://example.com/base","auth":"noAuth","method":"POST","format":"JSON","request_body_template":"[$$.events]","max_batch_size":"1"}`), got)

	for _, key := range []string{"properties_mapping", "is_default_mapping", "is_batching_enabled", "event_filtering"} {
		assert.NotContains(t, cm.ConfigSchema, key)
	}
	assert.Equal(t, "HTTP", cm.APIType)
	assert.Equal(t, 2, cm.Version)
	assert.Equal(t, 1, c.Destinations.Entries()["http"].Version)
}

func TestDestinationResourceHttpV2BodyAndMessageValidation(t *testing.T) {
	config := c.Destinations.Entries()["http_v2"].ConfigSchema
	body := config["request_body_template"]
	assert.True(t, body.Required)
	for _, value := range []string{"", "  \n  "} {
		assert.True(t, body.ValidateDiagFunc(value, cty.Path{}.IndexString("request_body_template")).HasError())
	}
	assert.False(t, body.ValidateDiagFunc("[$$.events]", cty.Path{}.IndexString("request_body_template")).HasError())

	messageType := config["supported_message_types"].Elem.(*schema.Schema).ValidateDiagFunc
	for _, value := range []string{"track", "identify", "page", "screen", "group", "alias"} {
		assert.False(t, messageType(value, cty.Path{}.IndexString("supported_message_types")).HasError(), value)
	}
	for _, value := range []string{"Track", "unknown", ""} {
		assert.True(t, messageType(value, cty.Path{}.IndexString("supported_message_types")).HasError(), value)
	}
}

func TestDestinationResourceHttpV2XMLAndFORM(t *testing.T) {
	cmt.AssertDestination(t, "http_v2", []c.TestConfig{
		{
			TerraformCreate: `
				api_url = "https://example.com/xml"
				format = "XML"
				xml_root_key = "Envelope"
				max_batch_size = "10"
				request_body_template = "{\"event\": $$.events}"
			`,
			APICreate: `{"apiUrl":"https://example.com/xml","auth":"noAuth","method":"POST","format":"XML","xmlRootKey":"Envelope","maxBatchSize":"10","requestBodyTemplate":"{\"event\": $$.events}"}`,
			TerraformUpdate: `
				api_url = "https://example.com/form"
				format = "FORM"
				request_body_template = "{\"email\": $$.events[0].traits.email}"
			`,
			APIUpdate: `{"apiUrl":"https://example.com/form","auth":"noAuth","method":"POST","format":"FORM","requestBodyTemplate":"{\"email\": $$.events[0].traits.email}"}`,
		},
	})
}

func TestDestinationResourceHttpV2Formats(t *testing.T) {
	cm := c.Destinations.Entries()["http_v2"]
	format := cm.ConfigSchema["format"]
	assert.Equal(t, "JSON", format.Default)
	for _, value := range []string{"JSON", "XML", "FORM"} {
		assert.False(t, format.ValidateDiagFunc(value, cty.Path{}.IndexString("format")).HasError(), value)
	}
	assert.True(t, format.ValidateDiagFunc("SOAP", cty.Path{}.IndexString("format")).HasError())

	root := cm.ConfigSchema["xml_root_key"]
	assert.True(t, root.Optional)
	assert.Nil(t, root.Default)
	for _, value := range []string{"", "Envelope", strings.Repeat("x", 100)} {
		assert.False(t, root.ValidateDiagFunc(value, cty.Path{}.IndexString("xml_root_key")).HasError(), value)
	}
	assert.True(t, root.ValidateDiagFunc(strings.Repeat("x", 101), cty.Path{}.IndexString("xml_root_key")).HasError())

	batch := cm.ConfigSchema["max_batch_size"]
	assert.True(t, batch.Optional)
	assert.Nil(t, batch.Default)
}

func TestDestinationResourceHttpV2FormatRoundTrip(t *testing.T) {
	cm := c.Destinations.Entries()["http_v2"]
	for _, tc := range []struct {
		name  string
		state string
		api   string
	}{
		{"JSON", `{"format":"JSON","max_batch_size":"5000"}`, `{"format":"JSON","maxBatchSize":"5000"}`},
		{"XML", `{"format":"XML","xml_root_key":"Envelope","max_batch_size":"10"}`, `{"format":"XML","xmlRootKey":"Envelope","maxBatchSize":"10"}`},
		{"FORM", `{"format":"FORM"}`, `{"format":"FORM"}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			api, err := cm.StateToAPI(tc.state)
			require.NoError(t, err)
			assert.True(t, testutil.JSONEq(api, tc.api), api)
			state, err := cm.APIToState(tc.api)
			require.NoError(t, err)
			assert.True(t, testutil.JSONEq(state, tc.state), state)
		})
	}

	v1 := c.Destinations.Entries()["http"]
	assert.Equal(t, 1, v1.Version)
	assert.Contains(t, v1.ConfigSchema, "properties_mapping")
	assert.Contains(t, v1.ConfigSchema, "is_batching_enabled")
	assert.Equal(t, "JSON", v1.ConfigSchema["format"].Default)
}

func TestDestinationResourceHttpV2FieldNames(t *testing.T) {
	config := c.Destinations.Entries()["http_v2"].ConfigSchema
	for _, key := range []string{"headers", "query_params"} {
		validate := config[key].Elem.(*schema.Resource).Schema["to"].ValidateDiagFunc
		assert.True(t, validate("$.token", cty.Path{}.IndexString("to")).HasError(), key)
		assert.False(t, validate("x-token", cty.Path{}.IndexString("to")).HasError(), key)
	}
}

func TestDestinationResourceHttpV2BatchBounds(t *testing.T) {
	validate := c.Destinations.Entries()["http_v2"].ConfigSchema["max_batch_size"].ValidateDiagFunc
	for _, tc := range []struct {
		value string
		valid bool
	}{
		{"0", false}, {"1", true}, {"100", true}, {"101", true}, {"4999", true}, {"5000", true}, {"5001", false}, {"10000", false}, {"01", false}, {"1.5", false},
	} {
		t.Run(tc.value, func(t *testing.T) {
			assert.Equal(t, tc.valid, !validate(tc.value, cty.Path{}.IndexString("max_batch_size")).HasError())
		})
	}
}

func TestDestinationResourceHttpV2ExpressionFields(t *testing.T) {
	cm := c.Destinations.Entries()["http_v2"]
	for _, tc := range []struct {
		field  string
		nested string
	}{
		{"headers", "from"}, {"query_params", "from"}, {"path_params", "path"},
	} {
		t.Run(tc.field, func(t *testing.T) {
			field := cm.ConfigSchema[tc.field].Elem.(*schema.Resource).Schema[tc.nested]
			assert.Nil(t, field.ValidateDiagFunc, "JSONata must not be checked with a v1 JSONPath validator")
		})
	}
}
