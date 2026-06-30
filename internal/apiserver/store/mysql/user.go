// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mysql

import (
	"context"

	"llmops/internal/pkg/model"

	gorm "gorm.io/gorm"
)

// UserStore defines the user storage interface.
type UserStore interface {
	Create(ctx context.Context, user *model.User) error
	Get(ctx context.Context, user *model.User) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
}

type user struct {
	db *gorm.DB
}

func newUser(ds *datastore) *user {
	return &user{ds.db}
}

// Create creates a new user account.
func (u *user) Create(ctx context.Context, user *model.User) error {
	return u.db.WithContext(ctx).Create(user).Error
}

// Get returns a user by query conditions.
func (u *user) Get(ctx context.Context, user *model.User) (*model.User, error) {
	ret := &model.User{}
	err := u.db.WithContext(ctx).Where(user).First(ret).Error

	return ret, err
}

// Update updates a user account.
func (u *user) Update(ctx context.Context, user *model.User) error {
	return u.db.WithContext(ctx).
		Updates(user).Error
}
