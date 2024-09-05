package configs_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestRegistries(t *testing.T) {
	r := &configs.Registry{}
	r.Register("test", configs.ConfigMeta{APIType: "APIType"})

	e := r.Entries()
	assert.Len(t, e, 1)
	assert.Equal(t, "APIType", e["test"].APIType)
}
