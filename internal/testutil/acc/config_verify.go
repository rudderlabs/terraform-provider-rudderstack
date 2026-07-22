package acc

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

// compareConfig verifies that actualRaw contains all fields specified in expectedJSON.
// Extra fields in the actual config are allowed (the API may add defaults).
// Returns nil if expectedJSON is empty (nothing to verify).
func compareConfig(actualRaw json.RawMessage, expectedJSON string) error {
	expectedJSON = strings.TrimSpace(expectedJSON)
	if expectedJSON == "" || expectedJSON == "{}" {
		return nil
	}

	var actual map[string]any
	if err := json.Unmarshal(actualRaw, &actual); err != nil {
		return fmt.Errorf("failed to unmarshal actual API config: %w", err)
	}

	var expected map[string]any
	if err := json.Unmarshal([]byte(expectedJSON), &expected); err != nil {
		return fmt.Errorf("failed to unmarshal expected config JSON: %w", err)
	}

	var mismatches []string
	compareFields("", expected, actual, &mismatches)

	if len(mismatches) > 0 {
		expectedPretty, _ := json.MarshalIndent(expected, "", "  ")
		actualPretty, _ := json.MarshalIndent(actual, "", "  ")
		fmt.Printf("\n=== expected config ===\n%s\n=== actual config ===\n%s\n===\n", expectedPretty, actualPretty)
		return fmt.Errorf("API config verification failed:\n%s", strings.Join(mismatches, "\n"))
	}
	return nil
}

func compareDestinationReadConfig(cm configs.ConfigMeta, actualRaw json.RawMessage, expectedJSON string) error {
	expectedJSON = strings.TrimSpace(expectedJSON)
	if expectedJSON == "" || expectedJSON == "{}" {
		return nil
	}

	expectedState, err := cm.APIToState(expectedJSON)
	if err != nil {
		return fmt.Errorf("failed to convert expected config JSON to state: %w", err)
	}

	expectedStateProps := map[string]interface{}{}
	if err := json.Unmarshal([]byte(expectedState), &expectedStateProps); err != nil {
		return fmt.Errorf("failed to unmarshal expected state JSON: %w", err)
	}

	filteredStateProps := removeSensitiveStateValues(cm.ConfigSchema, expectedStateProps)
	if len(filteredStateProps) == 0 {
		return nil
	}

	filteredState, err := json.Marshal(filteredStateProps)
	if err != nil {
		return fmt.Errorf("failed to marshal expected state JSON: %w", err)
	}

	filteredExpected, err := cm.StateToAPI(string(filteredState))
	if err != nil {
		return fmt.Errorf("failed to convert expected state JSON to config: %w", err)
	}

	return compareConfig(actualRaw, filteredExpected)
}

func removeSensitiveStateValues(configSchema map[string]*schema.Schema, props map[string]interface{}) map[string]interface{} {
	filtered := make(map[string]interface{}, len(props))
	for key, value := range props {
		sch, exists := configSchema[key]
		if !exists {
			filtered[key] = value
			continue
		}
		if sch.Sensitive {
			continue
		}

		nestedSchema := nestedSchemaForTest(sch)
		if nestedSchema == nil {
			filtered[key] = value
			continue
		}

		list, ok := value.([]interface{})
		if !ok {
			filtered[key] = value
			continue
		}

		filteredList := make([]interface{}, 0, len(list))
		for _, item := range list {
			itemMap, ok := item.(map[string]interface{})
			if !ok {
				filteredList = append(filteredList, item)
				continue
			}
			filteredList = append(filteredList, removeSensitiveStateValues(nestedSchema, itemMap))
		}
		filtered[key] = filteredList
	}
	return filtered
}

func nestedSchemaForTest(sch *schema.Schema) map[string]*schema.Schema {
	resource, ok := sch.Elem.(*schema.Resource)
	if !ok {
		return nil
	}
	return resource.Schema
}

// compareFields recursively checks that every key in expected exists in actual with the
// correct value. It collects all mismatches rather than failing on the first one.
func compareFields(prefix string, expected, actual map[string]any, mismatches *[]string) {
	for key, expectedVal := range expected {
		path := key
		if prefix != "" {
			path = prefix + "." + key
		}

		actualVal, exists := actual[key]
		if !exists {
			*mismatches = append(*mismatches, fmt.Sprintf("  missing field %q: expected %v", path, expectedVal))
			continue
		}

		compareValue(path, expectedVal, actualVal, mismatches)
	}
}

// compareValue recursively compares JSON values using subset semantics for objects and arrays:
//   - objects: all expected keys must exist in actual, but extra actual keys are allowed
//   - arrays: all expected elements must exist in actual at the same indexes, but extra actual
//     elements are allowed; objects within arrays also use subset semantics
func compareValue(path string, expectedVal, actualVal any, mismatches *[]string) {
	switch ev := expectedVal.(type) {
	case map[string]any:
		if av, ok := actualVal.(map[string]any); ok {
			compareFields(path, ev, av, mismatches)
		} else {
			*mismatches = append(*mismatches, fmt.Sprintf("  field %q: expected object, got %T", path, actualVal))
		}
	case []any:
		av, ok := actualVal.([]any)
		if !ok {
			*mismatches = append(*mismatches, fmt.Sprintf("  field %q: expected array, got %T", path, actualVal))
			return
		}
		if len(av) < len(ev) {
			*mismatches = append(*mismatches, fmt.Sprintf("  field %q: expected array length >= %d, got %d", path, len(ev), len(av)))
		}
		for i := 0; i < len(ev) && i < len(av); i++ {
			compareValue(fmt.Sprintf("%s[%d]", path, i), ev[i], av[i], mismatches)
		}
	default:
		if !reflect.DeepEqual(expectedVal, actualVal) {
			*mismatches = append(*mismatches, fmt.Sprintf("  field %q: expected %v (%T), got %v (%T)", path, expectedVal, expectedVal, actualVal, actualVal))
		}
	}
}
