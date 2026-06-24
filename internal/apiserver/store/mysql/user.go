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
}

type users struct {
	db *gorm.DB
}

func newUsers(ds *datastore) *users {
	return &users{ds.db}
}

// Create creates a new user account.
func (u *users) Create(ctx context.Context, user *model.User) error {
	return u.db.Create(&user).Error
}
