package testutil

import "time"

// TimePtr converts a time value to a pointer. Useful for creating *time.Time values in tests.
func TimePtr(t time.Time) *time.Time {
	return &t
}
