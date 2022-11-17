package cfg

// Pointer returns a pointer to t.
func Pointer[T any](t T) *T {
	return &t
}
