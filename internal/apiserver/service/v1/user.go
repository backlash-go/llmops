// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package v1

import (
	"context"
	v1 "llmops/pkg/api/llmops/v1"
	"regexp"

	"github.com/marmotedu/errors"

	"github.com/jinzhu/copier"
	"llmops/internal/apiserver/store"
	"llmops/internal/pkg/code"
	"llmops/internal/pkg/model"
)

// UserSrv defines functions used to handle user request.
type UserSrv interface {
	Create(ctx context.Context, r *v1.CreateUserRequest) error
	//Update(ctx context.Context, r *v1.UpdateUserRequest) error
	//DeleteCollection(ctx context.Context, r *v1.DeleteCollectionUserRequest) error
	//Get(ctx context.Context, r *v1.GetUserRequest) (*model.User, error)
	//List(ctx context.Context, r *v1.ListUserRequest) (*model.UserList, error)
	//ChangePassword(ctx context.Context, r *v1.ChangeUserPasswordRequest) error
}

type userService struct {
	store store.Factory
}

var _ UserSrv = (*userService)(nil)

func newUsers(srv *service) *userService {
	return &userService{store: srv.store}
}

// List returns user list in the storage.

// ListWithBadPerformance is kept for compatibility with older callers.

func (u *userService) Create(ctx context.Context, r *v1.CreateUserRequest) error {
	var user model.User

	_ = copier.Copy(&user, r)

	if err := u.store.Users().Create(ctx, &user); err != nil {
		if match, _ := regexp.MatchString("Duplicate entry '.*' for key 'uk_(username|email)'", err.Error()); match {
			return errors.WithCode(code.ErrUserAlreadyExist, err.Error())
		}

		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return nil
}
