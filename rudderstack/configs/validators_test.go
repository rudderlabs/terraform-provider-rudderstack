package configs_test

import (
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
	"github.com/stretchr/testify/assert"
)

func TestStringMatchesRegexp(t *testing.T) {
	f := configs.StringMatchesRegexp("^start-[0-9]*-end$")

	// a map of test strings and wether matching the regexp returns an error or not
	tests := map[string]bool{
		"start-1234-end": false,
		"start--end":     false,
		"start-abcd-end": true,
		"not-true":       true,
	}

	for s, b := range tests {
		d := f(s, cty.GetAttrPath("some-path"))
		assert.Equal(t, b, d.HasError())
	}
}

func TestStringNotMatchesRegexp(t *testing.T) {
	f := configs.StringNotMatchesRegexp("^(pg_|PG_|pG_|PG_)")

	// a map of test strings and wether matching the regexp returns an error or not
	tests := map[string]bool{
		"pg_any":     true,
		"not_pg_any": false,
	}

	for s, b := range tests {
		d := f(s, cty.GetAttrPath("some-path"))
		assert.Equal(t, b, d.HasError())
	}
}

func TestValidateAll(t *testing.T) {
	f := configs.ValidateAll(
		validation.ToDiagFunc(validation.StringLenBetween(1, 3)),
		validation.ToDiagFunc(validation.StringDoesNotContainAny("foo")),
	)

	// a map of test strings and wether matching the regexp returns an error or not
	tests := map[string]bool{
		"bar":    false,
		"bar123": true,
		"foo":    true,
	}

	for s, b := range tests {
		d := f(s, cty.GetAttrPath("some-path"))
		assert.Equal(t, b, d.HasError())
	}
}
