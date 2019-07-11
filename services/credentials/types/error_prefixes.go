package types

// These constants are the prefixes of errors that should be exposed to the client.
// All error encoders should check for them.
const (
	NotFound        = "not found:"
	InvalidArgument = "invalid argument:"
)
