package helper

import "time"

// IntPtr returns a pointer to the given int
func IntPtr(v int) *int {
	return &v
}

// StringPtr returns a pointer to the given string
func StringPtr(v string) *string {
	return &v
}

// TimePtr returns a pointer to the given time.Time
func TimePtr(v time.Time) *time.Time {
	return &v
}

// Float64Ptr returns a pointer to the given float64
func Float64Ptr(v float64) *float64 {
	return &v
}
