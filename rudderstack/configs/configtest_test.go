package configs_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

// The response expectation is the request minus the keys the backend prunes;
// the request JSON itself must stay untouched.
func TestTestConfigAPIResponsePrunesDeclaredKeys(t *testing.T) {
	tc := configs.TestConfig{
		APICreate:           `{"apiKey":"abc","additionalProperties":true}`,
		APICreatePrunedKeys: []string{"additionalProperties"},
		APIUpdate:           `{"apiKey":"xyz","collectContext":true}`,
		APIUpdatePrunedKeys: []string{"collectContext"},
	}

	assert.JSONEq(t, `{"apiKey":"abc"}`, tc.APICreateResponse())
	assert.JSONEq(t, `{"apiKey":"xyz"}`, tc.APIUpdateResponse())
	assert.JSONEq(t, `{"apiKey":"abc","additionalProperties":true}`, tc.APICreate, "request expectation stays strict")
}

func TestTestConfigAPIResponseWithoutPrunedKeys(t *testing.T) {
	tc := configs.TestConfig{APICreate: `{"apiKey":"abc"}`, APIUpdate: `{"apiKey":"xyz"}`}

	assert.Equal(t, tc.APICreate, tc.APICreateResponse())
	assert.Equal(t, tc.APIUpdate, tc.APIUpdateResponse())
}
