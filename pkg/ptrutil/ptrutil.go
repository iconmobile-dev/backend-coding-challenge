package ptrutil

import "time"

// Int returns a pointer to the given int value
func Int(i int) *int {
	return &i
}

// String returns a pointer to the given string value
func String(i string) *string {
	return &i
}

// Time returns a pointer to the given time value
func Time(i time.Time) *time.Time {
	return &i
}

// Bool returns a pointer to the given bool value
func Bool(i bool) *bool {
	return &i
}
