package sharedkernel

// JWT tokens
type JWT string

// JWT type to string
func (j JWT) String() string {
	return string(j)
}

// Header names
type HeaderName string

const (
	// Authorization header
	Authorization HeaderName = "Authorization"
	// Custom refresh token header
	RefreshCustom HeaderName = "X-Refresh"
)

// Headername type to string
func (h HeaderName) String() string {
	return string(h)
}

// User claim decoded from JWT
type UserClaim struct {
	UserID       string
	TokenVersion int64
}

// Custom context key type for storing values in contexts
type CtxKey string

const (
	Claim CtxKey = "claim"
)
