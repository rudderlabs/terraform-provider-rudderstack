package configs

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
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

	return filterJSONToSchema(state, cm.ConfigSchema)
}

func (cm *ConfigMeta) APIToStatePreservingWriteOnly(api, priorState string) (string, error) {
	state, err := cm.APIToState(api)
	if err != nil {
		return state, err
	}

	for _, path := range WriteOnlyStatePaths(cm.ConfigSchema) {
		prior := gjson.Get(priorState, path)
		if !prior.Exists() {
			continue
		}

		state, err = sjson.Set(state, path, prior.Value())
		if err != nil {
			return state, err
		}
	}

	return filterJSONToSchema(state, cm.ConfigSchema)
}

func WriteOnlyStatePaths(configSchema map[string]*schema.Schema) []string {
	pathSet := map[string]struct{}{}
	collectWriteOnlyStatePaths("", configSchema, pathSet)

	paths := make([]string, 0, len(pathSet))
	for path := range pathSet {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return paths
}

func collectWriteOnlyStatePaths(prefix string, configSchema map[string]*schema.Schema, paths map[string]struct{}) {
	for key, fieldSchema := range configSchema {
		if fieldSchema == nil {
			continue
		}

		path := key
		if prefix != "" {
			path = prefix + "." + key
		}

		if fieldSchema.Sensitive || schemaKeyLooksWriteOnly(key) {
			paths[path] = struct{}{}
			continue
		}

		resource, ok := fieldSchema.Elem.(*schema.Resource)
		if !ok {
			continue
		}

		switch fieldSchema.Type {
		case schema.TypeList, schema.TypeSet:
			collectWriteOnlyStatePaths(path+".0", resource.Schema, paths)
		case schema.TypeMap:
			collectWriteOnlyStatePaths(path, resource.Schema, paths)
		default:
			collectWriteOnlyStatePaths(path, resource.Schema, paths)
		}
	}
}

func schemaKeyLooksWriteOnly(key string) bool {
	normalized := strings.NewReplacer("_", "", "-", "").Replace(strings.ToLower(key))
	writeOnlyTokens := []string{
		"apikey",
		"apisecret",
		"apitoken",
		"eventkey",
		"accesstoken",
		"refreshtoken",
		"accesskey",
		"privatekey",
		"password",
		"secret",
		"token",
		"credential",
		"certificate",
	}
	for _, token := range writeOnlyTokens {
		if token == "eventkey" {
			if normalized == token {
				return true
			}
			continue
		}
		if strings.Contains(normalized, token) {
			return true
		}
	}
	return false
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
	return filterJSONToSchema(state, cm.SettingsSchema)
}

func filterJSONToSchema(state string, stateSchema map[string]*schema.Schema) (string, error) {
	if stateSchema == nil {
		return state, nil
	}

	var value interface{}
	if err := json.Unmarshal([]byte(state), &value); err != nil {
		return state, err
	}

	filtered, ok := filterValueToSchema(value, stateSchema)
	if !ok {
		return "{}", nil
	}

	filteredBytes, err := json.Marshal(filtered)
	if err != nil {
		return state, err
	}
	return string(filteredBytes), nil
}

func filterValueToSchema(value interface{}, stateSchema map[string]*schema.Schema) (interface{}, bool) {
	valueMap, ok := value.(map[string]interface{})
	if !ok {
		return nil, false
	}

	filtered := map[string]interface{}{}
	for key, rawValue := range valueMap {
		s, ok := stateSchema[key]
		if !ok {
			continue
		}

		filteredValue, keep := filterSchemaValue(rawValue, s)
		if keep {
			filtered[key] = filteredValue
		}
	}

	return filtered, len(filtered) > 0
}

func filterSchemaValue(value interface{}, s *schema.Schema) (interface{}, bool) {
	resource, ok := s.Elem.(*schema.Resource)
	if !ok {
		return value, true
	}

	switch s.Type {
	case schema.TypeList, schema.TypeSet:
		values, ok := value.([]interface{})
		if !ok {
			return nil, false
		}

		if len(values) == 0 {
			return values, true
		}

		filteredValues := make([]interface{}, 0, len(values))
		for _, item := range values {
			filteredItem, keep := filterValueToSchema(item, resource.Schema)
			if keep {
				filteredValues = append(filteredValues, filteredItem)
			}
		}
		return filteredValues, len(filteredValues) > 0
	case schema.TypeMap:
		filteredValue, keep := filterValueToSchema(value, resource.Schema)
		return filteredValue, keep
	default:
		return value, true
	}
}
