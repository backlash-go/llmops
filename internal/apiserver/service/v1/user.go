// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package v1

import (
	"context"
	"regexp"

	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
	"github.com/marmotedu/errors"

	"llmops/internal/apiserver/store"
	"llmops/internal/pkg/code"
	"llmops/internal/pkg/model"
	"llmops/pkg/log"
)

// UserSrv defines functions used to handle user request.
type UserSrv interface {
	Create(ctx context.Context, user *model.User, opts metav1.CreateOptions) error
	Update(ctx context.Context, user *model.User, opts metav1.UpdateOptions) error
	Delete(ctx context.Context, username string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, usernames []string, opts metav1.DeleteOptions) error
	Get(ctx context.Context, username string, opts metav1.GetOptions) (*model.User, error)
	List(ctx context.Context, opts metav1.ListOptions) (*model.UserList, error)
	ListWithBadPerformance(ctx context.Context, opts metav1.ListOptions) (*model.UserList, error)
	ChangePassword(ctx context.Context, user *model.User) error
}

type userService struct {
	store store.Factory
}

var _ UserSrv = (*userService)(nil)

func newUsers(srv *service) *userService {
	return &userService{store: srv.store}
}

// List returns user list in the storage.
func (u *userService) List(ctx context.Context, opts metav1.ListOptions) (*model.UserList, error) {
	users, err := u.store.Users().List(ctx, opts)
	if err != nil {
		log.L(ctx).Errorf("list users from storage failed: %s", err.Error())

		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	log.L(ctx).Debugf("get %d users from backend storage.", len(users.Items))

	return users, nil
}

// ListWithBadPerformance is kept for compatibility with older callers.
func (u *userService) ListWithBadPerformance(ctx context.Context, opts metav1.ListOptions) (*model.UserList, error) {
	return u.List(ctx, opts)
}

func (u *userService) Create(ctx context.Context, user *model.User, opts metav1.CreateOptions) error {
	if err := u.store.Users().Create(ctx, user, opts); err != nil {
		if match, _ := regexp.MatchString("Duplicate entry '.*' for key 'uk_(username|email)'", err.Error()); match {
			return errors.WithCode(code.ErrUserAlreadyExist, err.Error())
		}

		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return nil
}

func (u *userService) DeleteCollection(ctx context.Context, usernames []string, opts metav1.DeleteOptions) error {
	if err := u.store.Users().DeleteCollection(ctx, usernames, opts); err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return nil
}

func (u *userService) Delete(ctx context.Context, username string, opts metav1.DeleteOptions) error {
	if err := u.store.Users().Delete(ctx, username, opts); err != nil {
		return err
	}

	return nil
}

func (u *userService) Get(ctx context.Context, username string, opts metav1.GetOptions) (*model.User, error) {
	user, err := u.store.Users().Get(ctx, username, opts)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userService) Update(ctx context.Context, user *model.User, opts metav1.UpdateOptions) error {
	if err := u.store.Users().Update(ctx, user, opts); err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return nil
}

// ChangePassword is kept for compatibility; password is not part of model.User.
func (u *userService) ChangePassword(ctx context.Context, user *model.User) error {
	return u.Update(ctx, user, metav1.UpdateOptions{})
}
