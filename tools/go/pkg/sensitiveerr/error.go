package sensitiveerr

// withheld retains an error cause and its public diagnostic separately.
type withheld struct {
	// message contains only information safe for command output.
	message string
	// cause can include authentication parameters, tokens, or response data.
	cause error
}

// Error reports the public diagnostic supplied by the caller.
func (w withheld) Error() string { return w.message }

// Unwrap allows trusted callers to match the original failure.
func (w withheld) Unwrap() error { return w.cause }

// WithContext wraps a non-nil cause with a diagnostic safe for command output.
// The message must omit credentials and untrusted request or response data.
func WithContext(message string, cause error) error {
	return withheld{message: message, cause: cause}
}
