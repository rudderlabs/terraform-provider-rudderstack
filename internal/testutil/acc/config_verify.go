package acc

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// compareConfig verifies that actualRaw contains all fields specified in expectedJSON.
// Extra fields in the actual config are allowed (the API may add defaults).
// Returns nil if expectedJSON is empty (nothing to verify).
//
// redactedFields lists top-level API config keys whose values are never compared
// because the backend redacts them. optionalFields lists fields the backend may
// omit from read responses; when returned, their values are still compared.
func compareConfig(actualRaw json.RawMessage, expectedJSON string, redactedFields, optionalFields map[string]bool) error {
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
	compareFields("", expected, actual, redactedFields, optionalFields, &mismatches)

	if len(mismatches) > 0 {
		expectedPretty, _ := json.MarshalIndent(expected, "", "  ")
		actualPretty, _ := json.MarshalIndent(actual, "", "  ")
		fmt.Printf("\n=== expected config ===\n%s\n=== actual config ===\n%s\n===\n", expectedPretty, actualPretty)
		return fmt.Errorf("API config verification failed:\n%s", strings.Join(mismatches, "\n"))
	}
	return nil
}

// summarizeValidatedFields returns the sorted leaf field paths that compareConfig
// always asserts for expectedJSON, plus top-level redacted and response-optional
// keys. Used to log what a passing CRUD verification checked.
func summarizeValidatedFields(expectedJSON string, redactedFields, optionalFields map[string]bool) (validated, redacted, optional []string) {
	expectedJSON = strings.TrimSpace(expectedJSON)
	if expectedJSON == "" || expectedJSON == "{}" {
		return nil, nil, nil
	}
	var expected map[string]any
	if err := json.Unmarshal([]byte(expectedJSON), &expected); err != nil {
		return nil, nil, nil
	}
	for key, val := range expected {
		if redactedFields[key] {
			redacted = append(redacted, key)
			continue
		}
		if optionalFields[key] {
			optional = append(optional, key)
			continue
		}
		validated = append(validated, leafPaths(key, val)...)
	}
	sort.Strings(validated)
	sort.Strings(redacted)
	sort.Strings(optional)
	return validated, redacted, optional
}

func stringSet(values []string) map[string]bool {
	set := make(map[string]bool, len(values))
	for _, value := range values {
		set[value] = true
	}
	return set
}

// leafPaths returns the dotted paths of every scalar leaf under val at prefix,
// mirroring how compareConfig descends objects and arrays.
func leafPaths(prefix string, val any) []string {
	switch v := val.(type) {
	case map[string]any:
		// An empty object/array is still an asserted leaf (compareValue checks its
		// type), so emit its path rather than dropping it from the count.
		if len(v) == 0 {
			return []string{prefix}
		}
		var out []string
		for k, sub := range v {
			out = append(out, leafPaths(prefix+"."+k, sub)...)
		}
		return out
	case []any:
		if len(v) == 0 {
			return []string{prefix}
		}
		var out []string
		for i, sub := range v {
			out = append(out, leafPaths(fmt.Sprintf("%s[%d]", prefix, i), sub)...)
		}
		return out
	default:
		return []string{prefix}
	}
}

// compareFields recursively checks that every key in expected exists in actual with the
// correct value. It collects all mismatches rather than failing on the first one.
func compareFields(prefix string, expected, actual map[string]any, redactedFields, optionalFields map[string]bool, mismatches *[]string) {
	for key, expectedVal := range expected {
		path := key
		if prefix != "" {
			path = prefix + "." + key
		}

		// The backend redacts secret fields from API responses — either omitting
		// them or returning them blanked in place (e.g. a Sensitive list whose
		// values are emptied). Don't verify a redacted field at all, present or
		// absent. (redactedFields holds top-level API keys, which is where
		// redaction applies.)
		if redactedFields[path] {
			continue
		}

		actualVal, exists := actual[key]
		if !exists {
			if optionalFields[path] {
				continue
			}
			*mismatches = append(*mismatches, fmt.Sprintf("  missing field %q: expected %v", path, expectedVal))
			continue
		}

		compareValue(path, expectedVal, actualVal, redactedFields, optionalFields, mismatches)
	}
}

// compareValue recursively compares JSON values using subset semantics for objects and arrays:
//   - objects: all expected keys must exist in actual, but extra actual keys are allowed
//   - arrays: all expected elements must exist in actual at the same indexes, but extra actual
//     elements are allowed; objects within arrays also use subset semantics
func compareValue(path string, expectedVal, actualVal any, redactedFields, optionalFields map[string]bool, mismatches *[]string) {
	switch ev := expectedVal.(type) {
	case map[string]any:
		if av, ok := actualVal.(map[string]any); ok {
			compareFields(path, ev, av, redactedFields, optionalFields, mismatches)
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
			compareValue(fmt.Sprintf("%s[%d]", path, i), ev[i], av[i], redactedFields, optionalFields, mismatches)
		}
	default:
		if !reflect.DeepEqual(expectedVal, actualVal) {
			*mismatches = append(*mismatches, fmt.Sprintf("  field %q: expected %v (%T), got %v (%T)", path, expectedVal, expectedVal, actualVal, actualVal))
		}
	}
}
