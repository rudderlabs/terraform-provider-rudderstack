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
	// Name, so RemovedProperties can't target them. Targeting a missing Name
	// is a hard validation error rather than a silent no-op.
	base := c.ConfigMeta{
		APIType: "BASE_TYPE",
		Version: 1,
		Properties: []c.ConfigProperty{
			c.Simple("unnamedField", "unnamed_field"),
		},
	}

	assert.Panics(t, func() {
		ComposeConfigMeta(base, Delta{
			RemovedProperties: []string{"unnamedField"},
		}, 2)
	})

	v2 := ComposeConfigMeta(base, Delta{}, 2)
	assert.Len(t, v2.Properties, 1)
}

func TestComposeConfigMeta_InvalidDeltaPanics(t *testing.T) {
	base := baseConfigMetaForComposeTest()

	tests := []struct {
		name  string
		delta Delta
		msg   string
	}{
		{
			name:  "unknown removed key",
			delta: Delta{Removed: []string{"missing_field"}},
			msg:   `Removed schema key "missing_field" does not exist on base`,
		},
		{
			name:  "unknown renamed key",
			delta: Delta{Renamed: map[string]string{"missing_field": "x"}},
			msg:   `Renamed schema key "missing_field" does not exist on base`,
		},
		{
			name: "removed and renamed overlap",
			delta: Delta{
				Removed: []string{"renamed_field"},
				Renamed: map[string]string{"renamed_field": "renamed_field_v2"},
			},
			msg: `schema key "renamed_field" cannot be both Removed and Renamed`,
		},
		{
			name: "chained renames",
			delta: Delta{
				Renamed: map[string]string{
					"kept_field":    "renamed_field",
					"renamed_field": "renamed_field_v2",
				},
			},
			msg: "chained Renamed is not supported",
		},
		{
			name: "duplicate rename targets",
			delta: Delta{
				Renamed: map[string]string{
					"kept_field":    "shared_target",
					"renamed_field": "shared_target",
				},
			},
			msg: `duplicate Renamed target "shared_target"`,
		},
		{
			name: "rename into retained key",
			delta: Delta{
				Renamed: map[string]string{"renamed_field": "kept_field"},
			},
			msg: `Renamed target "kept_field" collides with retained base schema key`,
		},
		{
			name: "added key already exists",
			delta: Delta{
				Added: map[string]*schema.Schema{
					"kept_field": {Type: schema.TypeBool, Optional: true},
				},
			},
			msg: `Added schema key "kept_field" already exists on base`,
		},
		{
			name: "added key collides with rename target",
			delta: Delta{
				Renamed: map[string]string{"renamed_field": "renamed_field_v2"},
				Added: map[string]*schema.Schema{
					"renamed_field_v2": {Type: schema.TypeBool, Optional: true},
				},
			},
			msg: `Added schema key "renamed_field_v2" collides with Renamed target`,
		},
		{
			name:  "unknown removed property",
			delta: Delta{RemovedProperties: []string{"no_such_prop"}},
			msg:   `RemovedProperties name "no_such_prop" does not match a named base property`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var panicked any
			func() {
				defer func() { panicked = recover() }()
				ComposeConfigMeta(base, tt.delta, 2)
			}()
			require.NotNil(t, panicked, "expected panic")
			assert.Contains(t, fmtPanic(panicked), tt.msg)
		})
	}
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

func fmtPanic(p any) string {
	switch v := p.(type) {
	case error:
		return v.Error()
	case string:
		return v
	default:
		return ""
	}
}
