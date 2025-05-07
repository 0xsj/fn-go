// pkg/ptr/ptr.go
package ptr

// String returns a pointer to the given string value.
// Returns nil if the string is empty.
func String(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// Bool returns a pointer to the given boolean value.
func Bool(b bool) *bool {
	return &b
}

// Int32 returns a pointer to the given int32 value.
func Int32(i int32) *int32 {
	return &i
}

// Int64 returns a pointer to the given int64 value.
func Int64(i int64) *int64 {
	return &i
}

// Float64 returns a pointer to the given float64 value.
func Float64(f float64) *float64 {
	return &f
}

func getValueOrEmpty(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}
