// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"time"

	"gorm.io/gorm"
)

// UserAIAPIKey maps to the user_ai_api_key table.
type UserAIAPIKey struct {
	ID           uint64         `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement" json:"id"`
	UserID       uint64         `gorm:"column:user_id;type:bigint unsigned;not null;index:idx_user_id" json:"user_id"`
	Name         string         `gorm:"column:name;type:varchar(64);not null" json:"name"`
	APIKeyHash   string         `gorm:"column:api_key_hash;type:varchar(255);not null;uniqueIndex:uk_api_key_hash" json:"-"`
	APIKeyPrefix string         `gorm:"column:api_key_prefix;type:varchar(32);not null;default:'';index:idx_api_key_prefix" json:"api_key_prefix"`
	APIKeyLast4  string         `gorm:"column:api_key_last4;type:varchar(16);not null;default:''" json:"api_key_last4"`
	Status       uint8          `gorm:"column:status;type:tinyint unsigned;not null;default:1" json:"status"`
	ExpiresAt    *time.Time     `gorm:"column:expires_at;type:datetime" json:"expires_at,omitempty"`
	LastUsedAt   *time.Time     `gorm:"column:last_used_at;type:datetime" json:"last_used_at,omitempty"`
	CreatedAt    time.Time      `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;type:datetime" json:"deleted_at,omitempty"`
}

// TableName returns the database table name.
func (UserAIAPIKey) TableName() string {
	return "user_ai_api_key"
}
