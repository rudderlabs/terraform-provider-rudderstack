package configs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// These tests exercise the versioned-registry guards (Version presence and
// (APIType, Version) uniqueness) via the unexported `versioned` field on a
// throwaway Registry, rather than mutating the shared, process-global
// Destinations registry (which every destination's init() registers into and
// must not be reset or polluted from tests).

func newVersionedRegistryForTest(name string) *Registry {
	return &Registry{name: name, versioned: true}
}

func TestVersionedRegistry_RequiresExplicitVersion(t *testing.T) {
	r := newVersionedRegistryForTest("destinations under test")

	assert.PanicsWithError(t, "'test' must set ConfigMeta.Version when registering with destinations under test", func() {
		r.Register("test", ConfigMeta{APIType: "APIType"})
	})
}

func TestVersionedRegistry_AllowsDistinctVersionsOfSameAPIType(t *testing.T) {
	r := newVersionedRegistryForTest("destinations under test")

	assert.NotPanics(t, func() {
		r.Register("braze", ConfigMeta{APIType: "BRAZE", Version: 1})
		r.Register("braze_v2", ConfigMeta{APIType: "BRAZE", Version: 2})
	})

	e := r.Entries()
	assert.Len(t, e, 2)
	assert.Equal(t, 1, e["braze"].Version)
	assert.Equal(t, 2, e["braze_v2"].Version)
}

func TestVersionedRegistry_RejectsDuplicateAPITypeVersion(t *testing.T) {
	r := newVersionedRegistryForTest("destinations under test")
	r.Register("braze", ConfigMeta{APIType: "BRAZE", Version: 1})

	assert.PanicsWithError(t, "apiType 'BRAZE' version 1 is already registered as 'braze' in destinations under test", func() {
		r.Register("braze_dup", ConfigMeta{APIType: "BRAZE", Version: 1})
	})
}

func TestVersionedRegistry_RejectsDuplicateNameRegardlessOfVersion(t *testing.T) {
	// Name-uniqueness (pre-existing behavior) still applies independently of
	// the (APIType, Version) guard.
	r := newVersionedRegistryForTest("destinations under test")
	r.Register("braze", ConfigMeta{APIType: "BRAZE", Version: 1})

	assert.PanicsWithError(t, "name 'braze' is already registered with destinations under test", func() {
		r.Register("braze", ConfigMeta{APIType: "BRAZE_OTHER", Version: 2})
	})
}
