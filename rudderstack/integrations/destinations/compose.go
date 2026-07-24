package destinations

import (
	"fmt"
	"sort"
	"strings"

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
// Invalid deltas panic with a deterministic message (programmer/config error),
// matching Destinations.Register guards: unknown Removed/Renamed/RemovedProperties
// keys, Removed/Renamed source overlap, chained renames, duplicate rename
// targets, and collisions with retained or Added keys.
//
// The clone is shallow: ConfigSchema and Properties get new backing map/slice,
// but *schema.Schema values for keys that are neither renamed, removed, nor
// added are shared pointers with base. SettingsSchema and Settings are shared
// with base (destinations do not use them today). Callers must not mutate a
// *schema.Schema obtained from base.ConfigSchema in place (e.g. flipping
// Required/Optional on a renamed field) — build a new *schema.Schema and put
// it in delta.Added (or Renamed + a follow-up overwrite) instead. Mutating a
// shared *schema.Schema in place would leak the change back into base's
// ConfigMeta.
func ComposeConfigMeta(base c.ConfigMeta, delta Delta, version int) c.ConfigMeta {
	validateDelta(base, delta)

	schemaClone := make(map[string]*schema.Schema, len(base.ConfigSchema)+len(delta.Added))
	for k, v := range base.ConfigSchema {
		schemaClone[k] = v
	}

	for _, key := range delta.Removed {
		delete(schemaClone, key)
	}

	// Apply renames from original base values in sorted old-key order so the
	// result never depends on map iteration order.
	oldKeys := make([]string, 0, len(delta.Renamed))
	for oldKey := range delta.Renamed {
		oldKeys = append(oldKeys, oldKey)
	}
	sort.Strings(oldKeys)
	for _, oldKey := range oldKeys {
		newKey := delta.Renamed[oldKey]
		delete(schemaClone, oldKey)
		schemaClone[newKey] = base.ConfigSchema[oldKey]
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

func validateDelta(base c.ConfigMeta, delta Delta) {
	var errs []string

	removed := make(map[string]bool, len(delta.Removed))
	for _, key := range delta.Removed {
		if removed[key] {
			continue
		}
		removed[key] = true
		if _, ok := base.ConfigSchema[key]; !ok {
			errs = append(errs, fmt.Sprintf("Removed schema key %q does not exist on base", key))
		}
	}

	renameOldKeys := make([]string, 0, len(delta.Renamed))
	for oldKey := range delta.Renamed {
		renameOldKeys = append(renameOldKeys, oldKey)
	}
	sort.Strings(renameOldKeys)

	renameSources := make(map[string]bool, len(delta.Renamed))
	renameTargets := make(map[string]string, len(delta.Renamed)) // newKey -> oldKey
	for _, oldKey := range renameOldKeys {
		newKey := delta.Renamed[oldKey]
		renameSources[oldKey] = true

		if _, ok := base.ConfigSchema[oldKey]; !ok {
			errs = append(errs, fmt.Sprintf("Renamed schema key %q does not exist on base", oldKey))
		}
		if removed[oldKey] {
			errs = append(errs, fmt.Sprintf("schema key %q cannot be both Removed and Renamed", oldKey))
		}
		if otherOld, ok := renameTargets[newKey]; ok {
			errs = append(errs, fmt.Sprintf("duplicate Renamed target %q from %q and %q", newKey, otherOld, oldKey))
		} else {
			renameTargets[newKey] = oldKey
		}
	}

	// Chained renames (A→B and B→C) are rejected: B would be both a source and a target.
	for _, oldKey := range renameOldKeys {
		newKey := delta.Renamed[oldKey]
		if renameSources[newKey] {
			errs = append(errs, fmt.Sprintf("chained Renamed is not supported: %q -> %q is also a Renamed source", oldKey, newKey))
		}
	}

	// Rename into a retained base key (present after removals, not itself renamed away).
	for _, oldKey := range renameOldKeys {
		newKey := delta.Renamed[oldKey]
		if _, exists := base.ConfigSchema[newKey]; !exists {
			continue
		}
		if removed[newKey] || renameSources[newKey] {
			continue
		}
		errs = append(errs, fmt.Sprintf("Renamed target %q collides with retained base schema key", newKey))
	}

	addedKeys := make([]string, 0, len(delta.Added))
	for key := range delta.Added {
		addedKeys = append(addedKeys, key)
	}
	sort.Strings(addedKeys)
	for _, key := range addedKeys {
		if removed[key] {
			continue
		}
		if renameSources[key] {
			continue
		}
		if _, ok := renameTargets[key]; ok {
			errs = append(errs, fmt.Sprintf("Added schema key %q collides with Renamed target", key))
			continue
		}
		if _, ok := base.ConfigSchema[key]; ok {
			errs = append(errs, fmt.Sprintf("Added schema key %q already exists on base", key))
		}
	}

	namedProps := make(map[string]bool)
	for _, p := range base.Properties {
		if p.Name != "" {
			namedProps[p.Name] = true
		}
	}
	removedPropNames := append([]string(nil), delta.RemovedProperties...)
	sort.Strings(removedPropNames)
	seenProp := make(map[string]bool, len(removedPropNames))
	for _, name := range removedPropNames {
		if seenProp[name] {
			continue
		}
		seenProp[name] = true
		if !namedProps[name] {
			errs = append(errs, fmt.Sprintf("RemovedProperties name %q does not match a named base property", name))
		}
	}

	if len(errs) == 0 {
		return
	}
	sort.Strings(errs)
	panic(fmt.Errorf("invalid ComposeConfigMeta delta: %s", strings.Join(errs, "; ")))
}
