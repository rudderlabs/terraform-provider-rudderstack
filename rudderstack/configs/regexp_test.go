package configs_test

import (
	"testing"

	"github.com/hashicorp/go-cty/cty"
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
