package v1

import "time"

// CreateUserIdentityRequest is the request used to create a user identity binding.
type CreateUserIdentityRequest struct {
	UserID   uint64 `json:"user_id" binding:"required"`
	Provider string `json:"provider" binding:"required"`
	Issuer   string `json:"issuer" binding:"required"`
	Subject  string `json:"subject" binding:"required"`
}

// GetUserIdentityRequest is the request used to query a user identity binding.
type GetUserIdentityRequest struct {
	ID       uint64 `json:"id" form:"id" uri:"id"`
	UserID   uint64 `json:"user_id" form:"user_id"`
	Provider string `json:"provider" form:"provider"`
	Issuer   string `json:"issuer" form:"issuer"`
	Subject  string `json:"subject" form:"subject"`
}

// UpdateUserIdentityRequest is the request used to update a user identity binding.
type UpdateUserIdentityRequest struct {
	ID       uint64 `json:"id" uri:"id"`
	UserID   uint64 `json:"user_id"`
	Provider string `json:"provider"`
	Issuer   string `json:"issuer"`
	Subject  string `json:"subject"`
}

// DeleteUserIdentityRequest is the request used to delete a user identity binding.
type DeleteUserIdentityRequest struct {
	ID       uint64 `json:"id" form:"id" uri:"id"`
	UserID   uint64 `json:"user_id" form:"user_id"`
	Provider string `json:"provider" form:"provider"`
	Issuer   string `json:"issuer" form:"issuer"`
	Subject  string `json:"subject" form:"subject"`
}

// ListUserIdentityRequest is the request used to list user identity bindings.
type ListUserIdentityRequest struct {
	UserID   uint64 `json:"user_id" form:"user_id"`
	Provider string `json:"provider" form:"provider"`
	Issuer   string `json:"issuer" form:"issuer"`
	Subject  string `json:"subject" form:"subject"`
	Limit    int    `json:"limit" form:"limit"`
}

// UserIdentityResponse is the response shape for a user identity binding.
type UserIdentityResponse struct {
	ID        uint64    `json:"id"`
	UserID    uint64    `json:"user_id"`
	Provider  string    `json:"provider"`
	Issuer    string    `json:"issuer"`
	Subject   string    `json:"subject"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ListUserIdentityResponse is the response shape for listing user identity bindings.
type ListUserIdentityResponse struct {
	TotalCount int                     `json:"total_count"`
	Items      []*UserIdentityResponse `json:"items"`
}
