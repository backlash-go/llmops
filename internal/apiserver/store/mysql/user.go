// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mysql

import (
	"context"

	"github.com/marmotedu/component-base/pkg/fields"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
	"github.com/marmotedu/errors"
	gorm "gorm.io/gorm"

	"llmops/internal/pkg/code"
	"llmops/internal/pkg/model"
	"llmops/internal/pkg/util/gormutil"
)

type users struct {
	db *gorm.DB
}

func newUsers(ds *datastore) *users {
	return &users{ds.db}
}

// Create creates a new user account.
func (u *users) Create(ctx context.Context, user *model.User, opts metav1.CreateOptions) error {
	return u.db.Create(user).Error
}

// Update updates an user account information.
func (u *users) Update(ctx context.Context, user *model.User, opts metav1.UpdateOptions) error {
	return u.db.Save(user).Error
}

// Delete deletes the user by the user identifier.
func (u *users) Delete(ctx context.Context, username string, opts metav1.DeleteOptions) error {
	// delete related policy first
	pol := newPolicies(&datastore{u.db})
	if err := pol.DeleteByUser(ctx, username, opts); err != nil {
		return err
	}

	if opts.Unscoped {
		u.db = u.db.Unscoped()
	}

	err := u.db.Where("username = ?", username).Delete(&model.User{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return nil
}

// DeleteCollection batch deletes the users.
func (u *users) DeleteCollection(ctx context.Context, usernames []string, opts metav1.DeleteOptions) error {
	// delete related policy first
	pol := newPolicies(&datastore{u.db})
	if err := pol.DeleteCollectionByUser(ctx, usernames, opts); err != nil {
		return err
	}

	if opts.Unscoped {
		u.db = u.db.Unscoped()
	}

	return u.db.Where("username in (?)", usernames).Delete(&model.User{}).Error
}

// Get return an user by the user identifier.
func (u *users) Get(ctx context.Context, username string, opts metav1.GetOptions) (*model.User, error) {
	user := &model.User{}
	err := u.db.Where("username = ? and status = 1", username).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrUserNotFound, err.Error())
		}

		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return user, nil
}

// List return all users.
func (u *users) List(ctx context.Context, opts metav1.ListOptions) (*model.UserList, error) {
	ret := &model.UserList{}
	ol := gormutil.Unpointer(opts.Offset, opts.Limit)

	selector, _ := fields.ParseSelector(opts.FieldSelector)
	username, _ := selector.RequiresExactMatch("username")
	d := u.db.Where("username like ? and status = 1", "%"+username+"%").
		Offset(ol.Offset).
		Limit(ol.Limit).
		Order("id desc").
		Find(&ret.Items).
		Offset(-1).
		Limit(-1).
		Count(&ret.TotalCount)

	return ret, d.Error
}

// ListOptional show a more graceful query method.
func (u *users) ListOptional(ctx context.Context, opts metav1.ListOptions) (*model.UserList, error) {
	ret := &model.UserList{}
	ol := gormutil.Unpointer(opts.Offset, opts.Limit)

	where := model.User{}
	selector, _ := fields.ParseSelector(opts.FieldSelector)
	username, found := selector.RequiresExactMatch("username")
	if found {
		where.Username = username
	}

	d := u.db.Where(where).
		Offset(ol.Offset).
		Limit(ol.Limit).
		Order("id desc").
		Find(&ret.Items).
		Offset(-1).
		Limit(-1).
		Count(&ret.TotalCount)

	return ret, d.Error
}
