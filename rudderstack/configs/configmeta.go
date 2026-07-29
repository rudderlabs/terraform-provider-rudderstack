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
// "api_secret", "s3.0.access_key", "headers". A Sensitive field terminates the
// descent: the whole subtree is treated as the secret. Used to handle the
// backend redacting secrets from API responses (read/verify/import).
func (cm *ConfigMeta) SensitiveConfigPaths() []string {
	return sensitivePaths(cm.ConfigSchema, "")
}

func sensitivePaths(sch map[string]*schema.Schema, prefix string) []string {
	var out []string
	for key, s := range sch {
		if s == nil {
			continue
		}
		path := key
		if prefix != "" {
			path = prefix + "." + key
		}
		if s.Sensitive {
			out = append(out, path)
			continue
		}
		if res, ok := s.Elem.(*schema.Resource); ok && res != nil {
			out = append(out, sensitivePaths(res.Schema, path+".0")...)
		}
	}
	return out
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
