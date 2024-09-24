package configs

import (
	"fmt"
	"reflect"

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
type FromStateFunc func(config, state string) (string, error)

// ToStateFunc modifies a terraform state json object that represents a config object
// by extracting data from a provided API config json object. It returns the modified
// terraform state and an optional error.
type ToStateFunc func(state, config string) (string, error)

// Simple returns a ConfigProperty that maps an API config key to a terraform config key
// and vice versa. Additional ValueFilter filters can be applied to ignore a field in state
// depending on its value.
func Simple(apiKey, terraformKey string, filters ...ValueFilter) ConfigProperty {
	return ConfigProperty{
		FromStateFunc: copyFromState(apiKey, terraformKey, filters...),
		ToStateFunc:   copyToState(apiKey, terraformKey),
	}
}

type ValueFilter func(a interface{}) bool

type APINestedObject struct {
	TerraformKey string
	NestedKey    string
}

type terraformNestedObject struct {
	APIKey    string
	NestedKey string
}

// SkipZeroValue is a ValueFilter that returns true if the value is golang's zero value or an empty slice.
func SkipZeroValue(a interface{}) bool {
	switch v := a.(type) {
	case []interface{}:
		return len(v) == 0
	default:
		return reflect.ValueOf(a).IsZero()
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
func Equals(key, value string) ConfigConditionFunc {
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
type DiscriminatorValues map[string]interface{}

func ArrayWithStrings(rootAPIKey, nestedAPIField, terraformKey string) ConfigProperty {
	return ConfigProperty{
		FromStateFunc: func(config, state string) (string, error) {
			result := config
			v := gjson.Get(state, terraformKey)
			if v.Exists() && v.Value() != nil {
				switch a := v.Value().(type) {
				case []interface{}:
					contents := []interface{}{}
					for _, i := range a {
						contents = append(contents, map[string]interface{}{nestedAPIField: i})
					}

					if len(contents) > 0 {
						r, err := sjson.Set(result, rootAPIKey, contents)
						if err != nil {
							return result, err
						}
						result = r
					}
				default:
					return result, fmt.Errorf("provided value was not an array")
				}
			}
			return result, nil
		},
		ToStateFunc: func(state, config string) (string, error) {
			result := state

			r := gjson.Get(config, rootAPIKey)
			if r.Exists() && r.IsArray() {
				contents := []interface{}{}
				for _, i := range r.Value().([]interface{}) {
					if m, ok := i.(map[string]interface{}); ok {
						if v, ok := m[nestedAPIField]; ok {
							contents = append(contents, v)
						}
					}
				}
				s, err := sjson.Set(result, terraformKey, contents)
				if err != nil {
					return result, err
				}
				result = s
			}

			return result, nil
		},
	}
}

func GetInverseFields(fields map[string]interface{}) map[string]interface{} {
	inverseFields := map[string]interface{}{}
	for a, t := range fields {
		switch fieldVal := t.(type) {
		case string:
			inverseFields[fieldVal] = a
		case APINestedObject:
			tfKey := fieldVal.TerraformKey
			nestedKey := fieldVal.NestedKey
			inverseFields[tfKey] = terraformNestedObject{APIKey: a, NestedKey: nestedKey}
		}
	}
	return inverseFields
}

func ArrayWithObjects(rootAPIKey string, terraformKey string, fields map[string]interface{}) ConfigProperty {
	// we also need the inverse field map to convert terraform keys to api keys
	inverseFields := GetInverseFields(fields)

	return ConfigProperty{
		FromStateFunc: func(config, state string) (string, error) {
			result := config
			v := gjson.Get(state, terraformKey)
			if v.Exists() && v.Value() != nil {
				switch a := v.Value().(type) {
				case []interface{}:

					contents := GetConfigValue(a, inverseFields)
					if len(contents) > 0 {
						r, err := sjson.Set(result, rootAPIKey, contents)
						if err != nil {
							return result, err
						}
						result = r
					}
				default:
					return result, fmt.Errorf("provided value was not an array")
				}
			}
			return result, nil
		},
		ToStateFunc: func(state, config string) (string, error) {
			result := state

			r := gjson.Get(config, rootAPIKey)
			if r.Exists() && r.IsArray() {
				contents := GetTerraformValue(r.Value().([]interface{}), fields)
				s, err := sjson.Set(result, terraformKey, contents)
				if err != nil {
					return result, err
				}
				result = s
			}

			return result, nil
		},
	}
}

func GetTerraformValue(configValue []interface{}, fields map[string]interface{}) []interface{} {
	contents := []interface{}{}
	for _, i := range configValue {
		tv := map[string]interface{}{} // terraform value

		// iterate api values
		if av, ok := i.(map[string]interface{}); ok {
			// iterate terraform value fields
			for af, v := range av {
				if tf, ok := fields[af]; ok {
					switch fieldVal := tf.(type) {
					case string:
						tv[fieldVal] = v
					case APINestedObject:
						tfValues := []interface{}{}
						tfKey := fieldVal.TerraformKey
						nestedKey := fieldVal.NestedKey

						// iterate nested api values
						for _, nestedObj := range v.([]interface{}) {
							tfValues = append(tfValues, nestedObj.(map[string]interface{})[nestedKey])
						}
						tv[tfKey] = tfValues
					}
				}
			}
		}

		if len(tv) > 0 {
			contents = append(contents, tv)
		}
	}
	return contents
}

func GetConfigValue(stateValue []interface{}, fields map[string]interface{}) []interface{} {
	contents := []interface{}{}
	for _, i := range stateValue {
		av := map[string]interface{}{} // api value
		
		// iterate terraform values
		if tv, ok := i.(map[string]interface{}); ok {
			// iterate api value fields
			for tf, tfValue := range tv {
				switch fieldVal := fields[tf].(type) {
				case string:
					av[fieldVal] = tfValue
				case terraformNestedObject:
					tfValues := []interface{}{}
					af := fieldVal.APIKey
					nestedKey := fieldVal.NestedKey

					// iterate terraform values and convert to nested api values
					for _, nestedVal := range tfValue.([]interface{}) {
						nestedObj := map[string]interface{}{
							nestedKey: nestedVal,
						}
						tfValues = append(tfValues, nestedObj)
					}
					av[af] = tfValues
				}
			}
		}

		if len(av) > 0 {
			contents = append(contents, av)
		}
	}
	return contents
}

func applyFilters(a interface{}, filters []ValueFilter) bool {
	for _, o := range filters {
		if o(a) {
			return false
		}
	}

	return true
}

func copyFromState(apiKey, terraformKey string, options ...ValueFilter) FromStateFunc {
	return func(config, state string) (string, error) {
		result := config
		v := gjson.Get(state, terraformKey)
		if v.Exists() && v.Value() != nil && applyFilters(v.Value(), options) {
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
	return func(state, config string) (string, error) {
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
	return func(state, config string) (string, error) {
		if !condition(config) {
			return state, nil
		}

		return copyToState(apiKey, terraformKey)(state, config)
	}
}

func discriminatorValue(apiKey string, values DiscriminatorValues) FromStateFunc {
	return func(config, state string) (string, error) {
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
