package base

type contextKey string

const (
	// RequestIPAddressContextKey holds the key used to store a user ip address in the context.
	RequestIPAddressContextKey contextKey = "RequestIPAddress"
)
