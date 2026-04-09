package acc

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
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
		return fmt.Errorf("API config verification failed:\n%s", strings.Join(mismatches, "\n"))
	}
	return nil
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

		switch ev := expectedVal.(type) {
		case map[string]any:
			if av, ok := actualVal.(map[string]any); ok {
				compareFields(path, ev, av, mismatches)
			} else {
				*mismatches = append(*mismatches, fmt.Sprintf("  field %q: expected object, got %T", path, actualVal))
			}
		default:
			if !reflect.DeepEqual(expectedVal, actualVal) {
				*mismatches = append(*mismatches, fmt.Sprintf("  field %q: expected %v (%T), got %v (%T)", path, expectedVal, expectedVal, actualVal, actualVal))
			}
		}
	}
}
