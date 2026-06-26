// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mysql

import (
	"context"

	"gorm.io/gorm"

	"llmops/internal/pkg/model"
)

// UserIdentityStore defines the user identity storage interface.
type UserIdentityStore interface {
	Create(ctx context.Context, identity *model.UserIdentity) error
	Get(ctx context.Context, identity *model.UserIdentity) (*model.UserIdentity, error)
	Update(ctx context.Context, identity *model.UserIdentity) error
	Delete(ctx context.Context, identity *model.UserIdentity) error
	List(ctx context.Context, identity *model.UserIdentity, limit int) ([]*model.UserIdentity, int, error)
}

type userIdentity struct {
	db *gorm.DB
}

func newUserIdentity(ds *datastore) *userIdentity {
	return &userIdentity{ds.db}
}

// Create creates a new user identity.
func (u *userIdentity) Create(ctx context.Context, identity *model.UserIdentity) error {
	return u.db.WithContext(ctx).Create(identity).Error
}

// Get returns a user identity by query conditions.
func (u *userIdentity) Get(ctx context.Context, identity *model.UserIdentity) (*model.UserIdentity, error) {
	ret := &model.UserIdentity{}
	err := u.db.WithContext(ctx).Where(identity).First(ret).Error

	return ret, err
}

// Update updates a user identity.
func (u *userIdentity) Update(ctx context.Context, identity *model.UserIdentity) error {
	return u.db.WithContext(ctx).Save(identity).Error
}

// Delete deletes a user identity.
func (u *userIdentity) Delete(ctx context.Context, identity *model.UserIdentity) error {
	return u.db.WithContext(ctx).Delete(identity).Error
}

// List returns user identities by query conditions.
func (u *userIdentity) List(ctx context.Context, identity *model.UserIdentity, limit int) ([]*model.UserIdentity, int, error) {
	var total int64
	items := make([]*model.UserIdentity, 0)

	db := u.db.WithContext(ctx).Model(&model.UserIdentity{}).Where(identity)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if limit > 0 {
		db = db.Limit(limit)
	}
	if err := db.Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, int(total), nil
}
