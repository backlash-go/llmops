// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mysql

import (
	"context"

	"gorm.io/gorm"

	"llmops/internal/pkg/model"
)

// UserAIAPIKeyStore defines the user AI API key storage interface.
type UserAIAPIKeyStore interface {
	Create(ctx context.Context, key *model.UserAIAPIKey) error
	GetByHash(ctx context.Context, apiKeyHash string) (*model.UserAIAPIKey, error)
	BatchDelete(ctx context.Context, userID uint64, ids []uint64) (int64, error)
}

type userAIAPIKey struct {
	db *gorm.DB
}

func newUserAIAPIKey(ds *datastore) *userAIAPIKey {
	return &userAIAPIKey{ds.db}
}

// Create creates a new user AI API key.
func (u *userAIAPIKey) Create(ctx context.Context, key *model.UserAIAPIKey) error {
	return u.db.WithContext(ctx).Create(key).Error
}

// GetByHash returns a user AI API key by API key hash.
func (u *userAIAPIKey) GetByHash(ctx context.Context, apiKeyHash string) (*model.UserAIAPIKey, error) {
	ret := &model.UserAIAPIKey{}
	err := u.db.WithContext(ctx).Where("api_key_hash = ?", apiKeyHash).First(ret).Error

	return ret, err
}

// BatchDelete deletes user AI API keys by user ID and key IDs.
func (u *userAIAPIKey) BatchDelete(ctx context.Context, userID uint64, ids []uint64) (int64, error) {
	db := u.db.WithContext(ctx).
		Where("user_id = ? AND id IN ?", userID, ids).
		Delete(&model.UserAIAPIKey{})

	return db.RowsAffected, db.Error
}
