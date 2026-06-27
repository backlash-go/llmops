package v1

import (
	"time"

	"llmops/internal/pkg/model"
)

type CreateUserRequest struct {
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8,max=32"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar"`
}

type UpdateUserRequest struct{}

type GetUserRequest struct {
}

type DeleteCollectionUserRequest struct {
}

type ListUserRequest struct {
}

type ChangeUserPasswordRequest struct {
}

// OAuthLoginRequest is the request used after an OAuth id_token has been validated.
type OAuthLoginRequest struct {
	Provider          string `json:"provider" binding:"required"`
	Issuer            string `json:"issuer" binding:"required"`
	Subject           string `json:"subject" binding:"required"`
	PreferredUsername string `json:"preferred_username"`
	Email             string `json:"email"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	DisplayName       string `json:"display_name"`
	Avatar            string `json:"avatar"`
}

// OAuthLoginResponse is the response returned after local OAuth login binding.
type OAuthLoginResponse struct {
	UserID        uint64     `json:"user_id"`
	IdentityID    uint64     `json:"identity_id"`
	Username      string     `json:"username"`
	Email         string     `json:"email"`
	Provider      string     `json:"provider"`
	Issuer        string     `json:"issuer"`
	Subject       string     `json:"subject"`
	Created       bool       `json:"created"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
	SessionID     string     `json:"-"`
	SessionMaxAge int        `json:"-"`
}

type ListUserResponse struct {
	TotalCount int64         `json:"totalCount"`
	Users      []*model.User `json:"users"`
}
