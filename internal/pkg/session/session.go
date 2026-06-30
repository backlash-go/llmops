package session

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"
)

const (
	// CookieName is the browser cookie name used to carry the session ID.
	CookieName = "llmops_session"
	// CookiePath limits the browser session cookie to llmops pages and APIs.
	CookiePath = "/ops"
	// FrontendPath is the page users land on after OAuth login succeeds.
	FrontendPath = "/ops/portal"
	// KeyPrefix is the Redis key prefix for browser sessions.
	KeyPrefix = "llmops:session:"
	// DefaultTTL is the default browser session lifetime.
	DefaultTTL = 24 * time.Hour
)

// Data is the server-side session payload stored in Redis.
type Data struct {
	UserID      uint64   `json:"user_id"`
	IdentityID  uint64   `json:"identity_id"`
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	Provider    string   `json:"provider"`
	Issuer      string   `json:"issuer"`
	Subject     string   `json:"subject"`
	Roles       []string `json:"roles"`
	CreatedAt   int64    `json:"created_at"`
	ExpiresAt   int64    `json:"expires_at"`
	DisplayName string   `json:"display_name"`
}

// Key returns the Redis key for a session ID.
func Key(sessionID string) string {
	return KeyPrefix + sessionID
}

// NewID generates a cryptographically random browser session ID.
func NewID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate session id: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
