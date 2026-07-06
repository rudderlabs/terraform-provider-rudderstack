package destinations

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func baseConfigMetaForComposeTest() c.ConfigMeta {
	return c.ConfigMeta{
		APIType: "BASE_TYPE",
		Version: 1,
		ConfigSchema: map[string]*schema.Schema{
			"kept_field":    {Type: schema.TypeString, Optional: true},
			"renamed_field": {Type: schema.TypeString, Optional: true},
			"removed_field": {Type: schema.TypeString, Optional: true},
		},
		Properties: []c.ConfigProperty{
			c.Simple("keptField", "kept_field"),
		},
	}
}

func TestComposeConfigMeta_SetsVersionAndAPIType(t *testing.T) {
	base := baseConfigMetaForComposeTest()

	v2 := ComposeConfigMeta(base, Delta{}, 2)

	assert.Equal(t, base.APIType, v2.APIType)
	assert.Equal(t, 2, v2.Version)
	assert.Equal(t, 1, base.Version, "base.Version must be untouched")
}

func TestComposeConfigMeta_AddsRenamesAndRemovesSchemaFields(t *testing.T) {
	base := baseConfigMetaForComposeTest()

	v2 := ComposeConfigMeta(base, Delta{
		Renamed: map[string]string{"renamed_field": "renamed_field_v2"},
		Removed: []string{"removed_field"},
		Added: map[string]*schema.Schema{
			"new_field": {Type: schema.TypeBool, Optional: true},
		},
	}, 2)

	require.Contains(t, v2.ConfigSchema, "kept_field")
	require.Contains(t, v2.ConfigSchema, "renamed_field_v2")
	require.Contains(t, v2.ConfigSchema, "new_field")
	assert.NotContains(t, v2.ConfigSchema, "renamed_field")
	assert.NotContains(t, v2.ConfigSchema, "removed_field")

	// base must be entirely unaffected by composing v2.
	require.Contains(t, base.ConfigSchema, "renamed_field")
	require.Contains(t, base.ConfigSchema, "removed_field")
	assert.NotContains(t, base.ConfigSchema, "renamed_field_v2")
	assert.NotContains(t, base.ConfigSchema, "new_field")
	assert.Len(t, base.ConfigSchema, 3)
}

func TestComposeConfigMeta_RenamedFieldKeepsSameSchemaPointer(t *testing.T) {
	base := baseConfigMetaForComposeTest()
	originalSchema := base.ConfigSchema["renamed_field"]

	v2 := ComposeConfigMeta(base, Delta{
		Renamed: map[string]string{"renamed_field": "renamed_field_v2"},
	}, 2)

	assert.Same(t, originalSchema, v2.ConfigSchema["renamed_field_v2"])
}

func TestComposeConfigMeta_UnrelatedFieldsShareSchemaPointerWithBase(t *testing.T) {
	// Documents the shallow-copy contract: fields untouched by the delta are
	// the *same* *schema.Schema pointer as base's. Callers must not mutate
	// v2.ConfigSchema["kept_field"] in place.
	base := baseConfigMetaForComposeTest()

	v2 := ComposeConfigMeta(base, Delta{}, 2)

	assert.Same(t, base.ConfigSchema["kept_field"], v2.ConfigSchema["kept_field"])
}

func TestComposeConfigMeta_AddsAndRemovesProperties(t *testing.T) {
	base := c.ConfigMeta{
		APIType: "BASE_TYPE",
		Version: 1,
		Properties: []c.ConfigProperty{
			namedProperty("keep_me", c.Simple("keptField", "kept_field")),
			namedProperty("drop_me", c.Simple("droppedField", "dropped_field")),
		},
	}

	added := c.Simple("newField", "new_field")
	v2 := ComposeConfigMeta(base, Delta{
		AddedProperties:   []c.ConfigProperty{added},
		RemovedProperties: []string{"drop_me"},
	}, 2)

	assert.Len(t, v2.Properties, 2)
	names := propertyNames(v2.Properties)
	assert.Contains(t, names, "keep_me")
	assert.NotContains(t, names, "drop_me")

	// base is untouched.
	assert.Len(t, base.Properties, 2)
}

func TestComposeConfigMeta_PropertiesWithoutNameCannotBeRemoved(t *testing.T) {
	// Known limitation (documented on Delta.RemovedProperties): properties
	// built by today's constructors (Simple, Conditional, ...) don't set
	// Name, so RemovedProperties can't target them.
	base := c.ConfigMeta{
		APIType: "BASE_TYPE",
		Version: 1,
		Properties: []c.ConfigProperty{
			c.Simple("unnamedField", "unnamed_field"),
		},
	}

	v2 := ComposeConfigMeta(base, Delta{
		RemovedProperties: []string{"unnamedField"},
	}, 2)

	assert.Len(t, v2.Properties, 1, "property without a Name is not removable and must be carried over unchanged")
}

func TestAllRegisteredDestinationsHaveExplicitVersion(t *testing.T) {
	// Fleet-wide backfill verification (see the story's Open Items): every
	// destination registered via this package's init() functions must carry a
	// non-zero Version. Register itself already enforces this (it panics
	// during init() otherwise, which would fail the whole test binary), so
	// this test mostly documents and locks in the invariant; it also gives a
	// single readable assertion point if that guard is ever relaxed.
	entries := c.Destinations.Entries()
	require.NotEmpty(t, entries, "expected destinations to have registered via init()")

	for name, cm := range entries {
		assert.NotZero(t, cm.Version, "destination '%s' must have a non-zero Version", name)
	}
}

func namedProperty(name string, p c.ConfigProperty) c.ConfigProperty {
	p.Name = name
	return p
}

func propertyNames(props []c.ConfigProperty) []string {
	names := make([]string, 0, len(props))
	for _, p := range props {
		if p.Name != "" {
			names = append(names, p.Name)
		}
	}
	return names
}
