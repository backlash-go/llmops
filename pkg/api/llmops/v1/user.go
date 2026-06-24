package v1

import "llmops/internal/pkg/model"

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

type ListUserResponse struct {
	TotalCount int64         `json:"totalCount"`
	Users      []*model.User `json:"users"`
}
