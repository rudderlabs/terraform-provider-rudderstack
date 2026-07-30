package acc

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

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

// compareDestinationConfig verifies destination response config with the same
// subset semantics as compareConfig, but first removes expected fields that map
// back to sensitive Terraform schema values. Destination APIs intentionally omit
// or mask secrets on GET after create/update, while create/update request
// fixtures must remain strict.
func compareDestinationConfig(actualRaw json.RawMessage, expectedJSON string, cm configs.ConfigMeta) error {
	expectedJSON = strings.TrimSpace(expectedJSON)
	if expectedJSON == "" || expectedJSON == "{}" {
		return nil
	}

	sensitivePaths, err := destinationSensitiveAPIPaths(expectedJSON, cm)
	if err != nil {
		return err
	}
	if len(sensitivePaths) > 0 {
		expectedJSON = pruneUnavailableSensitiveJSONPaths(expectedJSON, actualRaw, sensitivePaths...)
	}

	return compareConfig(actualRaw, expectedJSON)
}

func destinationSensitiveAPIPaths(expectedJSON string, cm configs.ConfigMeta) ([]string, error) {
	if len(cm.ConfigSchema) == 0 {
		return nil, nil
	}

	stateJSON, err := cm.APIToState(expectedJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to map expected destination config to state: %w", err)
	}

	var state map[string]any
	if err := json.Unmarshal([]byte(stateJSON), &state); err != nil {
		return nil, fmt.Errorf("failed to unmarshal expected destination state: %w", err)
	}

	sensitiveState := filterSensitiveState(state, cm.ConfigSchema)
	if len(sensitiveState) == 0 {
		return nil, nil
	}

	sensitiveStateJSON, err := json.Marshal(sensitiveState)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal sensitive destination state: %w", err)
	}

	sensitiveAPIJSON, err := cm.StateToAPI(string(sensitiveStateJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to map sensitive destination state to API config: %w", err)
	}

	var sensitiveAPI map[string]any
	if err := json.Unmarshal([]byte(sensitiveAPIJSON), &sensitiveAPI); err != nil {
		return nil, fmt.Errorf("failed to unmarshal sensitive destination API config: %w", err)
	}

	var paths []string
	collectLeafPaths("", sensitiveAPI, &paths)
	return paths, nil
}

func filterSensitiveState(state map[string]any, schemaMap map[string]*schema.Schema) map[string]any {
	result := map[string]any{}
	for key, value := range state {
		fieldSchema, ok := schemaMap[key]
		if !ok {
			continue
		}
		if fieldSchema.Sensitive {
			result[key] = value
			continue
		}
		if nested := filterSensitiveNestedState(value, fieldSchema); nested != nil {
			result[key] = nested
		}
	}
	return result
}

func filterSensitiveNestedState(value any, fieldSchema *schema.Schema) any {
	resource, ok := fieldSchema.Elem.(*schema.Resource)
	if !ok {
		return nil
	}

	switch v := value.(type) {
	case []any:
		result := make([]any, 0, len(v))
		for _, item := range v {
			itemMap, ok := item.(map[string]any)
			if !ok {
				continue
			}
			nested := filterSensitiveState(itemMap, resource.Schema)
			if len(nested) > 0 {
				result = append(result, nested)
			}
		}
		if len(result) > 0 {
			return result
		}
	case map[string]any:
		nested := filterSensitiveState(v, resource.Schema)
		if len(nested) > 0 {
			return nested
		}
	}

	return nil
}

func collectLeafPaths(prefix string, value any, paths *[]string) {
	switch v := value.(type) {
	case map[string]any:
		for key, nested := range v {
			path := key
			if prefix != "" {
				path = prefix + "." + key
			}
			collectLeafPaths(path, nested, paths)
		}
	case []any:
		for i, nested := range v {
			collectLeafPaths(fmt.Sprintf("%s.%d", prefix, i), nested, paths)
		}
	default:
		if prefix != "" {
			*paths = append(*paths, prefix)
		}
	}
}

func pruneUnavailableSensitiveJSONPaths(jsonString string, actualRaw json.RawMessage, paths ...string) string {
	result := jsonString
	for _, path := range paths {
		actual := gjson.GetBytes(actualRaw, path)
		if actual.Exists() && actual.Value() != "" {
			continue
		}
		pruned, err := sjson.Delete(result, path)
		if err != nil {
			continue
		}
		result = pruned
	}
	return result
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
