package testutil

import (
	"encoding/json"
	"reflect"
	"time"
)

// TimePtr converts a time value to a pointer. Useful for creating *time.Time values in tests.
func TimePtr(t time.Time) *time.Time {
	return &t
}

// JSONEq returns true if JSON strings a & b are equal.
// Returns false if any of them is not a valid json
func JSONEq(a, b string) bool {
	var am, bm interface{}
	if err := json.Unmarshal([]byte(a), &am); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(b), &bm); err != nil {
		return false
	}

	return reflect.DeepEqual(am, bm)
}
