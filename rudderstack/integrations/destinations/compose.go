package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

// Delta describes how a new destination version's ConfigMeta differs from a
// base ConfigMeta (typically v1). It's consumed by ComposeConfigMeta, which
// clones the base and applies the delta, so a new version isn't a full
// copy-paste of the base's schema.
type Delta struct {
	// Renamed maps an existing ConfigSchema key to its new key. The
	// *schema.Schema value is carried over unchanged under the new key.
	Renamed map[string]string
	// Added introduces new top-level ConfigSchema fields.
	Added map[string]*schema.Schema
	// Removed drops existing ConfigSchema keys (and, if present, their renamed
	// equivalents are unaffected — Removed is matched against the base's
	// original keys, applied before Renamed... see ComposeConfigMeta ordering).
	Removed []string

	// AddedProperties is appended to the cloned Properties slice.
	AddedProperties []c.ConfigProperty
	// RemovedProperties drops entries from the cloned Properties slice whose
	// ConfigProperty.Name matches. ConfigProperty has no identity beyond an
	// optional Name (see configs.ConfigProperty): properties built by today's
	// constructors (Simple, Conditional, ArrayWithObjects, ...) don't set
	// Name, so they can't be targeted for removal until the base explicitly
	// tags them with a Name.
	RemovedProperties []string
}

// ComposeConfigMeta builds a new version's ConfigMeta from a base ConfigMeta
// (typically v1) and a Delta describing what changed, so a new version isn't a
// full copy-paste of the base's schema. version is stamped onto the result.
//
// The clone is shallow: ConfigSchema and Properties get new backing map/slice,
// but *schema.Schema values for keys that are neither renamed, removed, nor
// added are shared pointers with base. Callers must not mutate a *schema.Schema
// obtained from base.ConfigSchema in place (e.g. flipping Required/Optional on
// a renamed field) — build a new *schema.Schema and put it in delta.Added (or
// Renamed + a follow-up overwrite) instead. Mutating a shared *schema.Schema
// in place would leak the change back into base's ConfigMeta.
func ComposeConfigMeta(base c.ConfigMeta, delta Delta, version int) c.ConfigMeta {
	schemaClone := make(map[string]*schema.Schema, len(base.ConfigSchema)+len(delta.Added))
	for k, v := range base.ConfigSchema {
		schemaClone[k] = v
	}

	for _, key := range delta.Removed {
		delete(schemaClone, key)
	}
	for oldKey, newKey := range delta.Renamed {
		if s, ok := schemaClone[oldKey]; ok {
			delete(schemaClone, oldKey)
			schemaClone[newKey] = s
		}
	}
	for key, s := range delta.Added {
		schemaClone[key] = s
	}

	removedProps := make(map[string]bool, len(delta.RemovedProperties))
	for _, name := range delta.RemovedProperties {
		removedProps[name] = true
	}

	propsClone := make([]c.ConfigProperty, 0, len(base.Properties)+len(delta.AddedProperties))
	for _, p := range base.Properties {
		if p.Name != "" && removedProps[p.Name] {
			continue
		}
		propsClone = append(propsClone, p)
	}
	propsClone = append(propsClone, delta.AddedProperties...)

	return c.ConfigMeta{
		APIType:        base.APIType,
		Version:        version,
		SkipConfig:     base.SkipConfig,
		ConfigSchema:   schemaClone,
		Properties:     propsClone,
		SettingsSchema: base.SettingsSchema,
		Settings:       base.Settings,
	}
}
