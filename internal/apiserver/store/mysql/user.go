// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mysql

import (
	"context"

	v1 "github.com/marmotedu/api/apiserver/v1"
	"github.com/marmotedu/component-base/pkg/fields"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
	"github.com/marmotedu/errors"
	gorm "gorm.io/gorm"

	"llmops/internal/pkg/code"
	"llmops/internal/pkg/util/gormutil"
)




// UserStore defines the user storage interface.
type UserStore interface {
	Create(ctx context.Context, user *, opts metav1.CreateOptions) error
}




type users struct {
	db *gorm.DB
}

func newUsers(ds *datastore) *users {
	return &users{ds.db}
}

// Create creates a new user account.
func (u *users) Create(ctx context.Context, user *model.) error {
	return u.db.Create(&user).Error
}
