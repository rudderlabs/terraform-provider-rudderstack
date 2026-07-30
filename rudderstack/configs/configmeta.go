package configs

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ConfigMeta struct {
	APIType string
	// Version is the major version of the underlying integration definition
	// (e.g. 1, 2, ...). It is only meaningful for destinations today; sources
	// and accounts leave it at the zero value. See rudderstack/configs/registries.go
	// for the registration-time guards enforced on this field for destinations.
	Version        int
	SkipConfig     bool
	ConfigSchema   map[string]*schema.Schema
	Properties     []ConfigProperty
	SettingsSchema map[string]*schema.Schema
	Settings       []ConfigProperty
}

// SensitiveConfigPaths returns the terraform state paths of every Sensitive
// (secret) field in the config schema, descending into nested blocks — e.g.
// "api_secret", "s3.0.access_key", "headers". Every secret is a leaf/field
// path. Used for per-secret read-preserve and redacted-API-key derivation.
func (cm *ConfigMeta) SensitiveConfigPaths() []string {
	paths, _ := collectSensitivePaths(cm.ConfigSchema, "", false)
	return paths
}

// SensitiveImportIgnorePaths is like SensitiveConfigPaths but collapses a block
// whose fields are ALL Sensitive to the block path itself. Such a block reads
// back empty and its .# / .0.% structural attributes differ on ImportStateVerify,
// so it must be ignored wholesale; blocks with non-secret fields still return only
// their sensitive leaves so the rest stays verified.
func (cm *ConfigMeta) SensitiveImportIgnorePaths() []string {
	paths, _ := collectSensitivePaths(cm.ConfigSchema, "", true)
	return paths
}

// collectSensitivePaths walks the schema returning Sensitive field paths. When
// collapseBlocks is true, a block whose every field is Sensitive is returned as
// the block path instead of its leaves. The second return reports whether the
// given schema level is entirely Sensitive (used for the collapse decision).
func collectSensitivePaths(sch map[string]*schema.Schema, prefix string, collapseBlocks bool) (paths []string, allSensitive bool) {
	allSensitive = len(sch) > 0
	for key, s := range sch {
		if s == nil {
			allSensitive = false
			continue
		}
		path := key
		if prefix != "" {
			path = prefix + "." + key
		}
		if s.Sensitive {
			paths = append(paths, path)
			continue
		}
		res, ok := s.Elem.(*schema.Resource)
		if !ok || res == nil {
			allSensitive = false
			continue
		}
		sub, subAll := collectSensitivePaths(res.Schema, path+".0", collapseBlocks)
		if collapseBlocks && subAll {
			paths = append(paths, path)
		} else {
			paths = append(paths, sub...)
			allSensitive = allSensitive && subAll
		}
	}
	return paths, allSensitive
}

func (cm *ConfigMeta) StateToAPI(state string) (string, error) {
	api := "{}"

	for _, p := range cm.Properties {
		r, err := p.FromStateFunc(api, state)
		if err != nil {
			return api, err
		}
		api = r
	}

	return api, nil
}

func (cm *ConfigMeta) APIToState(api string) (string, error) {
	state := "{}"
	for _, p := range cm.Properties {
		s, err := p.ToStateFunc(state, api)
		if err != nil {
			return state, err
		}
		state = s
	}

	return state, nil
}

func (cm *ConfigMeta) SettingsStateToAPI(state string) (string, error) {
	api := "{}"
	for _, p := range cm.Settings {
		r, err := p.FromStateFunc(api, state)
		if err != nil {
			return api, err
		}
		api = r
	}
	return api, nil
}

func (cm *ConfigMeta) SettingsAPIToState(api string) (string, error) {
	state := "{}"
	for _, p := range cm.Settings {
		s, err := p.ToStateFunc(state, api)
		if err != nil {
			return state, err
		}
		state = s
	}
	return state, nil
}
