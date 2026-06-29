package v1

import "time"

// CreateUserAIAPIKeyRequest is the request used to create a user AI API key.
type CreateUserAIAPIKeyRequest struct {
	UserID    uint64     `json:"-"`
	Name      string     `json:"name" binding:"required"`
	ExpiresAt *time.Time `json:"expires_at"`
}

// UserAIAPIKeyResponse is the response shape for a user AI API key.
type UserAIAPIKeyResponse struct {
	ID           uint64     `json:"id"`
	UserID       uint64     `json:"user_id"`
	Name         string     `json:"name"`
	APIKey       string     `json:"api_key,omitempty"`
	APIKeyPrefix string     `json:"api_key_prefix"`
	APIKeyLast4  string     `json:"api_key_last4"`
	Status       uint8      `json:"status"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
	LastUsedAt   *time.Time `json:"last_used_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// BatchDeleteUserAIAPIKeyRequest is the request used to delete user AI API keys.
type BatchDeleteUserAIAPIKeyRequest struct {
	UserID uint64   `json:"-"`
	IDs    []uint64 `json:"ids" binding:"required"`
}

// BatchDeleteUserAIAPIKeyResponse is the response returned after deleting user AI API keys.
type BatchDeleteUserAIAPIKeyResponse struct {
	DeletedCount int64 `json:"deleted_count"`
}
