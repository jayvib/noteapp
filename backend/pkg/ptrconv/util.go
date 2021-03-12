package ptrconv

import "time"

func StringPointer(s string) *string {
	return &s
}

func TimePointer(t time.Time) *time.Time {
	return &t
}

func BoolPointer(b bool) *bool {
	return &b
}

func BoolValue(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func StringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func TimeValue(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}
