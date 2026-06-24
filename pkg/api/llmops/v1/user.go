package v1

import "time"

type CreateUserInfoRequest struct {
	Username    string     `json:"username" validate:"required"`
	Email       string     `json:"email" validate:"required"`
	Password    string     `json:"password" validate:"required,min=8,max=32"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Avatar      string     `json:"avatar"`
	Status      uint8      `json:"status"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
}

