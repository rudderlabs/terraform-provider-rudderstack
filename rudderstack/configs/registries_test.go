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

func TestRegistries_UnversionedAllowsZeroVersion(t *testing.T) {
	// Sources/Accounts-style registries aren't versioned: Version 0 (the zero
	// value) is allowed and there's no (APIType, Version) uniqueness guard.
	r := &configs.Registry{}
	assert.NotPanics(t, func() {
		r.Register("a", configs.ConfigMeta{APIType: "SAME"})
		r.Register("b", configs.ConfigMeta{APIType: "SAME"})
	})
}
