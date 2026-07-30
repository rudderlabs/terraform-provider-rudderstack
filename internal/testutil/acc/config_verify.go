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
func compareConfig(actualRaw json.RawMessage, expectedJSON string, ignoredPaths ...string) error {
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
	ignoredPathSet := make(map[string]struct{}, len(ignoredPaths))
	for _, path := range ignoredPaths {
		ignoredPathSet[path] = struct{}{}
	}
	compareFields("", expected, actual, ignoredPathSet, &mismatches)

	if len(mismatches) > 0 {
		expectedPretty, _ := json.MarshalIndent(expected, "", "  ")
		actualPretty, _ := json.MarshalIndent(actual, "", "  ")
		fmt.Printf("\n=== expected config ===\n%s\n=== actual config ===\n%s\n===\n", expectedPretty, actualPretty)
		return fmt.Errorf("API config verification failed:\n%s", strings.Join(mismatches, "\n"))
	}
	return nil
}

// compareFields recursively checks that every key in expected exists in actual with the
// correct value. It collects all mismatches rather than failing on the first one.
func compareFields(prefix string, expected, actual map[string]any, ignoredPaths map[string]struct{}, mismatches *[]string) {
	for key, expectedVal := range expected {
		path := key
		if prefix != "" {
			path = prefix + "." + key
		}

		if ignoredPathCoversValue(path, expectedVal, ignoredPaths) {
			continue
		}

		actualVal, exists := actual[key]
		if !exists {
			*mismatches = append(*mismatches, fmt.Sprintf("  missing field %q: expected %v", path, expectedVal))
			continue
		}

		compareValue(path, expectedVal, actualVal, ignoredPaths, mismatches)
	}
}

// compareValue recursively compares JSON values using subset semantics for objects and arrays:
//   - objects: all expected keys must exist in actual, but extra actual keys are allowed
//   - arrays: all expected elements must exist in actual at the same indexes, but extra actual
//     elements are allowed; objects within arrays also use subset semantics
func compareValue(path string, expectedVal, actualVal any, ignoredPaths map[string]struct{}, mismatches *[]string) {
	switch ev := expectedVal.(type) {
	case map[string]any:
		if av, ok := actualVal.(map[string]any); ok {
			compareFields(path, ev, av, ignoredPaths, mismatches)
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
			ignoredExpectedValues := 0
			for i := len(av); i < len(ev); i++ {
				if ignoredPathCoversValue(fmt.Sprintf("%s[%d]", path, i), ev[i], ignoredPaths) {
					ignoredExpectedValues++
				}
			}
			if len(av)+ignoredExpectedValues < len(ev) {
				*mismatches = append(*mismatches, fmt.Sprintf("  field %q: expected array length >= %d, got %d", path, len(ev), len(av)))
			}
		}
		for i := 0; i < len(ev) && i < len(av); i++ {
			compareValue(fmt.Sprintf("%s[%d]", path, i), ev[i], av[i], ignoredPaths, mismatches)
		}
	default:
		if !reflect.DeepEqual(expectedVal, actualVal) {
			*mismatches = append(*mismatches, fmt.Sprintf("  field %q: expected %v (%T), got %v (%T)", path, expectedVal, expectedVal, actualVal, actualVal))
		}
	}
}

func ignoredPathCoversValue(path string, expectedVal any, ignoredPaths map[string]struct{}) bool {
	if _, ok := ignoredPaths[path]; ok {
		return true
	}

	switch ev := expectedVal.(type) {
	case map[string]any:
		if len(ev) == 0 {
			return false
		}
		for key, value := range ev {
			if !ignoredPathCoversValue(path+"."+key, value, ignoredPaths) {
				return false
			}
		}
		return true
	case []any:
		if len(ev) == 0 {
			return false
		}
		for i, value := range ev {
			if !ignoredPathCoversValue(fmt.Sprintf("%s[%d]", path, i), value, ignoredPaths) {
				return false
			}
		}
		return true
	default:
		return false
	}
}
