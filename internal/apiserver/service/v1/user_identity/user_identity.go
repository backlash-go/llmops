// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// Package
package useridentity

import (
	"context"
	stderrors "errors"

	"github.com/jinzhu/copier"
	"github.com/marmotedu/errors"
	"gorm.io/gorm"

	"llmops/internal/apiserver/store/mysql"
	"llmops/internal/pkg/code"
	"llmops/internal/pkg/model"
	apiv1 "llmops/pkg/api/llmops/v1"
)

// UserIdentitySrv defines functions used to handle user identity requests.
type UserIdentitySrv interface {
	Create(ctx context.Context, r *apiv1.CreateUserIdentityRequest) error
	Get(ctx context.Context, r *apiv1.GetUserIdentityRequest) (*apiv1.UserIdentityResponse, error)
	Update(ctx context.Context, r *apiv1.UpdateUserIdentityRequest) error
	Delete(ctx context.Context, r *apiv1.DeleteUserIdentityRequest) error
	List(ctx context.Context, r *apiv1.ListUserIdentityRequest) (*apiv1.ListUserIdentityResponse, error)
}

type userIdentityService struct {
	store mysql.Factory
}

var _ UserIdentitySrv = (*userIdentityService)(nil)

// NewUserIdentity creates a user identity service.
func NewUserIdentity(store mysql.Factory) UserIdentitySrv {
	return &userIdentityService{store: store}
}

// Create creates a user identity binding.
func (u *userIdentityService) Create(ctx context.Context, r *apiv1.CreateUserIdentityRequest) error {
	var identity model.UserIdentity

	_ = copier.Copy(&identity, r)

	if err := u.store.UserIdentity().Create(ctx, &identity); err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return nil
}

// Get returns a user identity binding.
func (u *userIdentityService) Get(ctx context.Context, r *apiv1.GetUserIdentityRequest) (*apiv1.UserIdentityResponse, error) {
	var identity model.UserIdentity

	_ = copier.Copy(&identity, r)

	ret, err := u.store.UserIdentity().Get(ctx, &identity)
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrUserNotFound, err.Error())
		}

		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return userIdentityResponseFromModel(ret), nil
}

// Update updates a user identity binding.
func (u *userIdentityService) Update(ctx context.Context, r *apiv1.UpdateUserIdentityRequest) error {
	var identity model.UserIdentity

	_ = copier.Copy(&identity, r)

	if err := u.store.UserIdentity().Update(ctx, &identity); err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return nil
}

// Delete deletes a user identity binding.
func (u *userIdentityService) Delete(ctx context.Context, r *apiv1.DeleteUserIdentityRequest) error {
	var identity model.UserIdentity

	_ = copier.Copy(&identity, r)

	if err := u.store.UserIdentity().Delete(ctx, &identity); err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return nil
}

// List returns user identity bindings.
func (u *userIdentityService) List(ctx context.Context, r *apiv1.ListUserIdentityRequest) (*apiv1.ListUserIdentityResponse, error) {
	var identity model.UserIdentity

	_ = copier.Copy(&identity, r)

	identities, total, err := u.store.UserIdentity().List(ctx, &identity, r.Limit)
	if err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	items := make([]*apiv1.UserIdentityResponse, 0, len(identities))
	for _, item := range identities {
		items = append(items, userIdentityResponseFromModel(item))
	}

	return &apiv1.ListUserIdentityResponse{
		TotalCount: total,
		Items:      items,
	}, nil
}

func userIdentityResponseFromModel(identity *model.UserIdentity) *apiv1.UserIdentityResponse {
	if identity == nil {
		return nil
	}

	return &apiv1.UserIdentityResponse{
		ID:        identity.ID,
		UserID:    identity.UserID,
		Provider:  identity.Provider,
		Issuer:    identity.Issuer,
		Subject:   identity.Subject,
		CreatedAt: identity.CreatedAt,
		UpdatedAt: identity.UpdatedAt,
	}
}
