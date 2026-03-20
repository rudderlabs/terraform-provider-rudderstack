package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"
)

func TestExtractResourceType(t *testing.T) {
	tests := []struct {
		input           string
		expectedKind    string
		expectedType    string
		expectErr       bool
	}{
		{"rudderstack_destination_webhook", "destination", "webhook", false},
		{"rudderstack_destination_google_analytics", "destination", "google_analytics", false},
		{"rudderstack_source_shopify", "source", "shopify", false},
		{"rudderstack_source_http", "source", "http", false},
		{"rudderstack_connection_foo", "", "", true},
		{"aws_s3_bucket", "", "", true},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			kind, intType, err := ExtractResourceType(tc.input)
			if tc.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedKind, kind)
			assert.Equal(t, tc.expectedType, intType)
		})
	}
}

func TestCtyToGo(t *testing.T) {
	tests := []struct {
		name     string
		val      cty.Value
		sch      *schema.Schema
		expected interface{}
	}{
		{
			name:     "string",
			val:      cty.StringVal("hello"),
			sch:      &schema.Schema{Type: schema.TypeString},
			expected: "hello",
		},
		{
			name:     "bool true",
			val:      cty.True,
			sch:      &schema.Schema{Type: schema.TypeBool},
			expected: true,
		},
		{
			name:     "bool false",
			val:      cty.False,
			sch:      &schema.Schema{Type: schema.TypeBool},
			expected: false,
		},
		{
			name:     "integer",
			val:      cty.NumberIntVal(42),
			sch:      &schema.Schema{Type: schema.TypeInt},
			expected: int64(42),
		},
		{
			name:     "float",
			val:      cty.NumberFloatVal(3.14),
			sch:      &schema.Schema{Type: schema.TypeFloat},
			expected: 3.14,
		},
		{
			name:     "null value",
			val:      cty.NullVal(cty.String),
			sch:      &schema.Schema{Type: schema.TypeString},
			expected: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ctyToGo(tc.val, tc.sch)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestCtyListToGo_Primitives(t *testing.T) {
	val := cty.ListVal([]cty.Value{
		cty.StringVal("a"),
		cty.StringVal("b"),
		cty.StringVal("c"),
	})
	sch := &schema.Schema{
		Type: schema.TypeList,
		Elem: &schema.Schema{Type: schema.TypeString},
	}

	result, err := ctyListToGo(val, sch)
	require.NoError(t, err)

	items, ok := result.([]interface{})
	require.True(t, ok)
	assert.Equal(t, []interface{}{"a", "b", "c"}, items)
}

func TestCtyListToGo_Objects(t *testing.T) {
	val := cty.ListVal([]cty.Value{
		cty.ObjectVal(map[string]cty.Value{
			"key":   cty.StringVal("k1"),
			"value": cty.StringVal("v1"),
		}),
		cty.ObjectVal(map[string]cty.Value{
			"key":   cty.StringVal("k2"),
			"value": cty.StringVal("v2"),
		}),
	})
	sch := &schema.Schema{
		Type: schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"key":   {Type: schema.TypeString},
				"value": {Type: schema.TypeString},
			},
		},
	}

	result, err := ctyListToGo(val, sch)
	require.NoError(t, err)

	items, ok := result.([]interface{})
	require.True(t, ok)
	require.Len(t, items, 2)

	first := items[0].(map[string]interface{})
	assert.Equal(t, "k1", first["key"])
	assert.Equal(t, "v1", first["value"])

	second := items[1].(map[string]interface{})
	assert.Equal(t, "k2", second["key"])
	assert.Equal(t, "v2", second["value"])
}

func TestCtyListToGo_EmptyNonIterable(t *testing.T) {
	val := cty.NullVal(cty.List(cty.String))
	sch := &schema.Schema{
		Type: schema.TypeList,
		Elem: &schema.Schema{Type: schema.TypeString},
	}

	result, err := ctyListToGo(val, sch)
	require.NoError(t, err)
	assert.Equal(t, []interface{}{}, result)
}

func TestSubsetDiff(t *testing.T) {
	t.Run("matching subset", func(t *testing.T) {
		expected := map[string]interface{}{
			"url":    "https://example.com",
			"method": "POST",
		}
		actual := map[string]interface{}{
			"url":    "https://example.com",
			"method": "POST",
			"extra":  "ignored",
		}

		diffs := SubsetDiff(expected, actual, "")
		assert.Empty(t, diffs)
	})

	t.Run("missing key", func(t *testing.T) {
		expected := map[string]interface{}{
			"url":    "https://example.com",
			"apiKey": "secret",
		}
		actual := map[string]interface{}{
			"url": "https://example.com",
		}

		diffs := SubsetDiff(expected, actual, "")
		require.Len(t, diffs, 1)
		assert.Contains(t, diffs[0], "apiKey")
		assert.Contains(t, diffs[0], "missing")
	})

	t.Run("value mismatch", func(t *testing.T) {
		expected := map[string]interface{}{
			"url": "https://example.com",
		}
		actual := map[string]interface{}{
			"url": "https://other.com",
		}

		diffs := SubsetDiff(expected, actual, "")
		require.Len(t, diffs, 1)
		assert.Contains(t, diffs[0], "url")
	})

	t.Run("nested map match", func(t *testing.T) {
		expected := map[string]interface{}{
			"config": map[string]interface{}{
				"key": "value",
			},
		}
		actual := map[string]interface{}{
			"config": map[string]interface{}{
				"key":   "value",
				"extra": "ok",
			},
		}

		diffs := SubsetDiff(expected, actual, "")
		assert.Empty(t, diffs)
	})

	t.Run("nested map mismatch", func(t *testing.T) {
		expected := map[string]interface{}{
			"config": map[string]interface{}{
				"key": "value",
			},
		}
		actual := map[string]interface{}{
			"config": map[string]interface{}{
				"key": "wrong",
			},
		}

		diffs := SubsetDiff(expected, actual, "")
		require.Len(t, diffs, 1)
		assert.Contains(t, diffs[0], "config.key")
	})

	t.Run("with prefix", func(t *testing.T) {
		expected := map[string]interface{}{"a": "1"}
		actual := map[string]interface{}{"a": "2"}

		diffs := SubsetDiff(expected, actual, "root")
		require.Len(t, diffs, 1)
		assert.Contains(t, diffs[0], "root.a")
	})
}

func TestArrayDiff(t *testing.T) {
	t.Run("matching arrays", func(t *testing.T) {
		a := []interface{}{"x", "y"}
		b := []interface{}{"x", "y"}

		diffs := arrayDiff(a, b, "arr")
		assert.Empty(t, diffs)
	})

	t.Run("length mismatch", func(t *testing.T) {
		a := []interface{}{"x"}
		b := []interface{}{"x", "y"}

		diffs := arrayDiff(a, b, "arr")
		require.Len(t, diffs, 1)
		assert.Contains(t, diffs[0], "array length")
	})

	t.Run("element mismatch", func(t *testing.T) {
		a := []interface{}{"x", "y"}
		b := []interface{}{"x", "z"}

		diffs := arrayDiff(a, b, "arr")
		require.Len(t, diffs, 1)
		assert.Contains(t, diffs[0], "arr[1]")
	})

	t.Run("nested object arrays", func(t *testing.T) {
		a := []interface{}{
			map[string]interface{}{"k": "v1"},
		}
		b := []interface{}{
			map[string]interface{}{"k": "v2"},
		}

		diffs := arrayDiff(a, b, "items")
		require.Len(t, diffs, 1)
		assert.Contains(t, diffs[0], "items[0].k")
	})
}

func TestFormatResult(t *testing.T) {
	info := &IntegrationResource{
		ResourceType: "rudderstack_destination_webhook",
	}

	t.Run("pass", func(t *testing.T) {
		result := &VerifyResult{
			Match:    true,
			Expected: `{"url":"https://example.com"}`,
			Actual:   `{"url":"https://example.com"}`,
		}

		output := FormatResult(info, "abc12345678", result)
		assert.Contains(t, output, "PASS")
		assert.Contains(t, output, "abc12345...")
	})

	t.Run("fail with diffs", func(t *testing.T) {
		result := &VerifyResult{
			Match:       false,
			Expected:    `{"url":"https://a.com"}`,
			Actual:      `{"url":"https://b.com"}`,
			Differences: []string{"  url: expected https://a.com, got https://b.com"},
		}

		output := FormatResult(info, "short", result)
		assert.Contains(t, output, "FAIL")
		assert.Contains(t, output, "Differences")
		assert.Contains(t, output, "url")
	})

	t.Run("short id not truncated", func(t *testing.T) {
		result := &VerifyResult{Match: true, Expected: "{}", Actual: "{}"}
		output := FormatResult(info, "abcd", result)
		assert.Contains(t, output, "abcd")
		assert.NotContains(t, output, "...")
	})
}

func TestParseTFFile(t *testing.T) {
	t.Run("parses destination resource", func(t *testing.T) {
		tf := `
resource "rudderstack_destination_webhook" "test" {
  name = "my-webhook"

  config {
    webhook_url = "https://example.com"
  }
}
`
		dir := t.TempDir()
		path := filepath.Join(dir, "main.tf")
		require.NoError(t, os.WriteFile(path, []byte(tf), 0600))

		info, err := ParseTFFile(path, "")
		require.NoError(t, err)
		assert.Equal(t, "destination", info.Kind)
		assert.Equal(t, "webhook", info.IntegrationType)
		assert.Equal(t, "rudderstack_destination_webhook", info.ResourceType)
		assert.Equal(t, "test", info.Name)
		assert.Contains(t, info.ConfigState, "webhook_url")
	})

	t.Run("filters by target resource name", func(t *testing.T) {
		tf := `
resource "rudderstack_destination_webhook" "first" {
  name = "first"
  config {
    webhook_url = "https://first.com"
  }
}

resource "rudderstack_destination_webhook" "second" {
  name = "second"
  config {
    webhook_url = "https://second.com"
  }
}
`
		dir := t.TempDir()
		path := filepath.Join(dir, "main.tf")
		require.NoError(t, os.WriteFile(path, []byte(tf), 0600))

		info, err := ParseTFFile(path, "second")
		require.NoError(t, err)
		assert.Equal(t, "second", info.Name)
	})

	t.Run("returns error for no matching resource", func(t *testing.T) {
		tf := `
resource "aws_s3_bucket" "example" {
  bucket = "my-bucket"
}
`
		dir := t.TempDir()
		path := filepath.Join(dir, "main.tf")
		require.NoError(t, os.WriteFile(path, []byte(tf), 0600))

		_, err := ParseTFFile(path, "")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no matching")
	})

	t.Run("resource without config block returns empty JSON", func(t *testing.T) {
		tf := `
resource "rudderstack_source_http" "test" {
  name = "my-source"
}
`
		dir := t.TempDir()
		path := filepath.Join(dir, "main.tf")
		require.NoError(t, os.WriteFile(path, []byte(tf), 0600))

		info, err := ParseTFFile(path, "")
		require.NoError(t, err)
		assert.Equal(t, "source", info.Kind)
		assert.Equal(t, "{}", info.ConfigState)
	})
}
