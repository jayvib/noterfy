package ptrconv

import "time"

// StringPointer converts s string into a pointer string.
func StringPointer(s string) *string {
	return &s
}

// TimePointer converts t time into a pointer time.
func TimePointer(t time.Time) *time.Time {
	return &t
}

// BoolPointer converts b boolean into a pointer boolean.
func BoolPointer(b bool) *bool {
	return &b
}

// BoolValue converts the pointer boolean b into a boolean value.
func BoolValue(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// StringValue converts the pointer string s into a string value
func StringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// TimeValue converts ponter time t into a time value
func TimeValue(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}
