package configs

import (
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// ConfigProperty defines how a property in a API config object (e.g source/destination config)
// maps to terraform state and vice versa.
type ConfigProperty struct {
	ToStateFunc   ToStateFunc
	FromStateFunc FromStateFunc
}

// FromStateFunc modifies am API config json object using terraform state information
// provided by a ResourceData object. It returns the modified config and an optional error.
type FromStateFunc func(config string, state string) (string, error)

// ToStateFunc modifies a terraform state json object that represents a config object
// by extracting data from a provided API config json object. It returns the modified
// terraform state and an optional error.
type ToStateFunc func(state string, config string) (string, error)

// Simple returns a ConfigProperty that maps an API config key to a terraform config key
// and vice versa.
func Simple(apiKey, terraformKey string) ConfigProperty {
	return ConfigProperty{
		FromStateFunc: copyFromState(apiKey, terraformKey),
		ToStateFunc:   copyToState(apiKey, terraformKey),
	}
}

// conditional returns a ConfigProperty that maps an API config key to a terraform config key
// only if provided condition is satisfied for that API config.
func Conditional(apiKey, terraformKey string, condition ConfigConditionFunc) ConfigProperty {
	return ConfigProperty{
		FromStateFunc: copyFromState(apiKey, terraformKey),
		ToStateFunc:   copyToStateConditional(apiKey, terraformKey, condition),
	}
}

// ConfigConditionFunc is a function that checks a provided API config object for
// some condition and returns true if the condition is met.
type ConfigConditionFunc func(config string) bool

// equals returns a ConfigConditionFunc that is true if the API config contains
// the specified key and it has the specified value.
func Equals(key string, value string) ConfigConditionFunc {
	return func(config string) bool {
		r := gjson.Get(config, key)
		return r.Exists() && r.Value() == value
	}
}

// discriminator returns a ConfigProperty that is not stored directly in terraform state.
// The corresponding API config value is set based on the provided DiscriminatorValues.
// if a DiscriminatorValues key exists in terraform state, the corresponding value in the
// values object is used.
func Discriminator(apiKey string, values DiscriminatorValues) ConfigProperty {
	return ConfigProperty{
		FromStateFunc: discriminatorValue(apiKey, values),
		ToStateFunc:   func(state, config string) (string, error) { return state, nil },
	}
}

// DiscriminatorValues is a map of API config values for discriminator fields, mapped
// to a terraform state key of a config.
type DiscriminatorValues map[string]string

func copyFromState(apiKey, terraformKey string) FromStateFunc {
	return func(config string, state string) (string, error) {
		result := config
		v := gjson.Get(state, terraformKey)
		if v.Exists() && v.Value() != nil {
			sresult, err := sjson.Set(result, apiKey, v.Value())
			if err != nil {
				return result, err
			}
			result = sresult
		}

		return result, nil
	}
}

func copyToState(apiKey, terraformKey string) ToStateFunc {
	return func(state string, config string) (string, error) {
		r := gjson.Get(config, apiKey)
		if r.Exists() {
			s, err := sjson.Set(state, terraformKey, r.Value())
			if err != nil {
				return state, err
			}
			state = s
		}

		return state, nil
	}
}

func copyToStateConditional(apiKey, terraformKey string, condition ConfigConditionFunc) ToStateFunc {
	return func(state string, config string) (string, error) {
		if !condition(config) {
			return state, nil
		}

		return copyToState(apiKey, terraformKey)(state, config)
	}
}

func discriminatorValue(apiKey string, values DiscriminatorValues) FromStateFunc {
	return func(config string, state string) (string, error) {
		for k, v := range values {
			r := gjson.Get(state, k)

			if !r.Exists() {
				continue
			}

			// this is necessary to ignore empty state blocks
			if r.IsArray() && len(r.Value().([]interface{})) == 0 {
				continue
			}

			return sjson.Set(config, apiKey, v)
		}

		return config, nil
	}
}
