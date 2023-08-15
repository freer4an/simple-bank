package token

import "time"

// Token manager
type Maker interface {
	// Creates token for a specific username
	CreateToken(username string, duration time.Duration) (string, error)

	// Check if the token is valid
	VerifyToken(token string) (*Payload, error)
}
